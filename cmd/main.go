package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/fatih/color"
	gchat "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"google.golang.org/protobuf/types/known/emptypb"
)

const grpcPort = 50051

type server struct {
	gchat.UnimplementedChatV1Server
}

func (s *server) Create(ctx context.Context, req *gchat.CreateRequest) (*gchat.CreateResponse, error) {
	fmt.Println(color.BlueString("Create chat with: name - %v", req.Usernames, ctx))
	return &gchat.CreateResponse{
		Id: gofakeit.Int64(),
	}, nil
}

func (s *server) Delete(ctx context.Context, req *gchat.DeleteRequest) (*emptypb.Empty, error) {
	fmt.Println(color.RedString("Delete chat: id %d, with ctx: %v", req.Id, ctx))
	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *gchat.SendMessageRequest) (*emptypb.Empty, error) {
	timestamp := req.Timestamp.AsTime().Format(time.UnixDate)
	fmt.Println(color.WhiteString("Send message: %s to %s at %s, with ctx: %v", req.Text, req.From, timestamp, ctx))
	return &emptypb.Empty{}, nil
}

func main() {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", grpcPort))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	gchat.RegisterChatV1Server(s, &server{})

	fmt.Println(color.GreenString("run server at %s", lis.Addr()))

	if err = s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
