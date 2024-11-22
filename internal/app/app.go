package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/fatih/color"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/ukrainskykirill/platform_common/pkg/closer"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"

	"github.com/ukrainskykirill/chat-server/internal/config"
	"github.com/ukrainskykirill/chat-server/internal/interceptor"
	"github.com/ukrainskykirill/chat-server/internal/metrics"
	gchat "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
)

type App struct {
	serviceProvider *serviceProvider
	grpcServer      *grpc.Server
}

func NewApp(ctx context.Context) (*App, error) {
	a := &App{}

	err := a.initDeps(ctx)
	if err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run(ctx context.Context) error {
	defer func() {
		closer.CloseAll()
		closer.Wait()
	}()

	ctx, cancel := context.WithCancel(ctx)

	err := metrics.Init(ctx)
	if err != nil {
		fmt.Println("metrics init error")
	}

	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		err := a.runGRPCServer()
		if err != nil {
			fmt.Printf("run grpc service: %s", err)
		}
	}()

	go func() {
		err := runPrometheus()
		if err != nil {
			fmt.Printf("run prometheus: %s", err)
		}
	}()

	gracefulShutdown(ctx, cancel, wg)
	return nil
}

func gracefulShutdown(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	select {
	case <-ctx.Done():
		log.Println("terminating: context cancelled")
	case <-waitSignal():
		log.Println("terminating: via signal")
	}

	cancel()
	if wg != nil {
		wg.Wait()
	}
}

func waitSignal() chan os.Signal {
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	return sigterm
}

func (a *App) initDeps(ctx context.Context) error {
	inits := []func(context.Context) error{
		a.initConfig,
		a.initServiceProvider,
		a.initGRPCServer,
	}

	for _, f := range inits {
		err := f(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *App) initConfig(_ context.Context) error {
	err := config.LoadConfig()
	if err != nil {
		return err
	}

	return nil
}

func (a *App) initServiceProvider(_ context.Context) error {
	a.serviceProvider = newServiceProvider()
	return nil
}

func (a *App) initGRPCServer(ctx context.Context) error {
	a.grpcServer = grpc.NewServer(
		grpc.Creds(insecure.NewCredentials()),
		grpc.ChainUnaryInterceptor(
			interceptor.MetricsInterceptor,
			interceptor.ValidationInterceptor,
		),
	)

	reflection.Register(a.grpcServer)

	gchat.RegisterChatV1Server(a.grpcServer, a.serviceProvider.ChatsAPI(ctx))

	return nil
}

func (a *App) runGRPCServer() error {
	lis, err := net.Listen("tcp", a.serviceProvider.GRPCConfig().Address())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(color.GreenString("run server at %s", lis.Addr()))
	if err = a.grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}

	return nil
}

func runPrometheus() error {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	prometheusServer := &http.Server{
		Addr:    "localhost:2112",
		Handler: mux,
	}

	log.Printf(color.BlueString("Prometheus server is running on %s", "localhost:2112"))

	err := prometheusServer.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
