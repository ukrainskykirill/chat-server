package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/fatih/color"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/timestamppb"

	desc "github.com/ukrainskykirill/chat-server/pkg/chat_v1"
)

const (
	address         = "localhost:50051"
	myUserID        = 133
	chatParticipant = 4344
)

func main() {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("failed to connect to server: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()
	client := desc.NewChatV1Client(conn)

	chatID, err := createChat(ctx, client)
	if err != nil {
		fmt.Println(fmt.Sprintf(color.RedString("chat creation error %v", err)))
	}

	log.Printf(fmt.Sprintf("%s: %s\n", color.GreenString("Chat created"), color.YellowString(string(chatID))))

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		err = connectChat(ctx, client, strconv.Itoa(int(chatID)), strconv.Itoa(int(myUserID)), time.Duration(time.Second*10))
		if err != nil {
			log.Fatalf("failed to connect chat: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		err = connectChat(ctx, client, strconv.Itoa(int(chatID)), strconv.Itoa(int(myUserID)), time.Duration(time.Second*5))
		if err != nil {
			log.Fatalf("failed to connect chat: %v", err)
		}
	}()

	wg.Wait()

}

func connectChat(ctx context.Context, client desc.ChatV1Client, chatID string, userID string, period time.Duration) error {
	stream, err := client.ConnectChat(ctx, &desc.ConnectChatRequest{
		ChatId: chatID,
		UserId: userID,
	})

	if err != nil {
		return nil
	}

	go func() {
		for {
			message, errRecv := stream.Recv()
			if errRecv == io.EOF {
				return
			}
			if errRecv != nil {
				log.Println("failed to receive message from stream: ", errRecv)
				return
			}

			log.Printf("[from: %s]: %s", message.GetUserId(), message.GetText())
		}
	}()

	intChatID, err := strconv.Atoi(chatID)
	if err != nil {
		fmt.Printf("cannot parse str chatID to int, %s", chatID)
	}
	intUserID, err := strconv.Atoi(userID)
	if err != nil {
		fmt.Printf("cannot parse str userID to int, %s", userID)
	}

	for {
		time.Sleep(period)
		text := gofakeit.BeerName()

		_, err := client.SendMessage(ctx, &desc.SendMessageRequest{
			ChatID:    int64(intChatID),
			UserId:    int64(intUserID),
			Text:      text,
			Timestamp: timestamppb.Now(),
		},
		)
		if err != nil {
			log.Println("failed to send message: ", err)
			return err
		}
	}
}

func createChat(ctx context.Context, client desc.ChatV1Client) (int64, error) {
	res, err := client.Create(ctx, &desc.CreateRequest{
		UserIDs: []int64{myUserID, chatParticipant},
	})
	if err != nil {
		return 0, err
	}

	return res.GetId(), nil
}
