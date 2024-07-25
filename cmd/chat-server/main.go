package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/jackc/pgx/v5"
	pool "github.com/jackc/pgx/v5/pgxpool"

	"github.com/fatih/color"
	"github.com/ukrainskykirill/chat-server/internal/config"
	gchat "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type server struct {
	gchat.UnimplementedChatV1Server
	pg *pool.Pool
}

func (s *server) Create(ctx context.Context, req *gchat.CreateRequest) (*gchat.CreateResponse, error) {
	tx, err := s.pg.Begin(ctx)
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Printf("rollback error: %v", err)
		}
	}()

	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var chatID int64
	err = tx.QueryRow(
		ctx,
		`INSERT INTO chats (is_regular) VALUES ($1) RETURNING id;`,
		len(req.Usernames) > 2,
	).Scan(&chatID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}
	rows := [][]any{}
	for _, username := range req.Usernames {
		newRow := []any{username, chatID}
		rows = append(rows, newRow)
	}
	_, err = tx.CopyFrom(
		ctx,
		pgx.Identifier{"chat_users"},
		[]string{"username", "chat_id"},
		pgx.CopyFromRows(rows),
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	fmt.Println(color.BlueString("Create chat with: name - %v", req.Usernames, ctx))
	return &gchat.CreateResponse{
		Id: chatID,
	}, nil
}

func (s *server) Delete(ctx context.Context, req *gchat.DeleteRequest) (*emptypb.Empty, error) {
	tx, err := s.pg.Begin(ctx)
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Printf("rollback error: %v", err)
		}
	}()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var chatUserID int64
	err = tx.QueryRow(
		ctx,
		`SELECT id FROM chat_users WHERE chat_id = $1;`,
		req.Id,
	).Scan(&chatUserID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	_, err = tx.Exec(
		ctx,
		`DELETE FROM messages WHERE chat_user_id = $1;`,
		chatUserID,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	_, err = tx.Exec(
		ctx,
		`DELETE FROM chat_users WHERE chat_id = $1;`,
		req.Id,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	_, err = tx.Exec(
		ctx,
		`DELETE FROM chats WHERE id = $1;`,
		req.Id,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	fmt.Println(color.RedString("Delete chat: id %d, with ctx: %v", req.Id, ctx))
	return &emptypb.Empty{}, nil
}

func (s *server) SendMessage(ctx context.Context, req *gchat.SendMessageRequest) (*emptypb.Empty, error) {
	timestamp := req.Timestamp.AsTime().Format(time.UnixDate)

	tx, err := s.pg.Begin(ctx)
	defer func() {
		if err := tx.Rollback(ctx); err != nil {
			log.Printf("rollback error: %v", err)
		}
	}()
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	var chatUserID int64
	err = tx.QueryRow(
		ctx,
		`SELECT id FROM chat_users WHERE username = $1;`,
		req.From,
	).Scan(&chatUserID)
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}

	_, err = tx.Exec(
		ctx,
		`INSERT INTO messages (chat_user_id, text, timestamp) VALUES ($1, $2, $3);`,
		chatUserID, req.Text, timestamp,
	)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	fmt.Println(color.WhiteString("Send message: %s to %s at %s, with ctx: %v", req.Text, req.From, timestamp, ctx))
	return &emptypb.Empty{}, nil
}

func main() {
	appConf, err := config.InitConfig()
	if err != nil {
		log.Fatal(err)
	}

	pgxPool, err := pool.New(context.Background(), appConf.DB.URL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pgxPool.Close()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", appConf.GRPC.Port))
	if err != nil {
		log.Fatal(err)
	}

	s := grpc.NewServer()
	reflection.Register(s)
	gchat.RegisterChatV1Server(s, &server{
		pg: pgxPool,
	})

	fmt.Println(color.GreenString("run server at %s", lis.Addr()))

	if err = s.Serve(lis); err != nil {
		log.Fatal(err)
	}
}
