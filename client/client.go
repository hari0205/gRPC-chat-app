package main

import (
	"bufio"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	chatapp "github.com/hari0205/grpc-chat-app/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var client chatapp.BroadcastClient
var wait *sync.WaitGroup

func init() {
	wait = &sync.WaitGroup{}
}

func Connect(user *chatapp.User) error {
	var streamerror error

	stream, err := client.CreateStream(context.Background(), &chatapp.Connect{
		User:   user,
		Active: true,
	})

	if err != nil {
		log.Fatalf("connection failed: %v", err)
	}
	wait.Add(1)
	go func(str chatapp.Broadcast_CreateStreamClient) {
		defer wait.Done()
		for {
			msg, err := str.Recv()
			if err != nil {
				streamerror = fmt.Errorf("recv error: %v", err)
				break
			}
			fmt.Printf("%v sent %s\n", msg.Id, msg.Content)
		}
		str.Recv()
	}(stream)
	return streamerror
}

func main() {
	timestamp := time.Now()
	done := make(chan int)

	name := flag.String("name", "anonymous", "The name to use")
	flag.Parse()
	id := sha256.Sum256([]byte(timestamp.String() + *name))

	conn, err := grpc.Dial("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("Could not connnect to service: %v", err)
	}

	client = chatapp.NewBroadcastClient(conn)
	user := &chatapp.User{
		Id:   hex.EncodeToString(id[:]),
		Name: *name,
	}

	Connect(user)
	wait.Add(1)
	go func() {
		defer wait.Done()

		scanner := bufio.NewScanner(os.Stdin)
		for scanner.Scan() {
			msg := &chatapp.Message{
				Id:        user.Id,
				Content:   scanner.Text(),
				Timestamp: timestamp.String(),
			}

			_, err := client.BroadcastMessage(context.Background(), msg)
			if err != nil {
				fmt.Printf("Error while sending message: %v", err)
				break
			}

		}
	}()

	go func() {
		wait.Wait()
		close(done)
	}()
	<-done
}
