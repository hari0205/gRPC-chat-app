package main

import (
	"context"
	"log"
	"net"
	"sync"

	chatapp "github.com/hari0205/grpc-chat-app/proto"
	"google.golang.org/grpc"
)

type Connection struct {
	stream chatapp.Broadcast_CreateStreamServer
	id     string
	active bool
	error  chan error
}

type Server struct {
	Connection []*Connection
	chatapp.UnimplementedBroadcastServer
}

func (s *Server) CreateStream(conn *chatapp.Connect, stream chatapp.Broadcast_CreateStreamServer) error {
	c := &Connection{
		stream: stream,
		id:     conn.User.Id,
		active: true,
		error:  make(chan error),
	}

	s.Connection = append(s.Connection, c)
	return <-c.error
}

func (s *Server) BroadcastMessage(ctx context.Context, msg *chatapp.Message) (*chatapp.Close, error) {
	wait := sync.WaitGroup{}
	done := make(chan int)

	for _, conn := range s.Connection {
		wait.Add(1)

		go func(msg *chatapp.Message, conn *Connection) {
			defer wait.Done()

			if conn.active {
				err := conn.stream.Send(msg)

				if err != nil {
					conn.active = false
					conn.error <- err
				}
			}
		}(msg, conn)
	}
	go func() {
		wait.Wait()
		close(done)
	}()

	<-done
	return &chatapp.Close{}, nil
}

func main() {

	var Connections []*Connection
	server := &Server{Connections, chatapp.UnimplementedBroadcastServer{}}

	grpcserver := grpc.NewServer()
	lis, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		log.Fatalf("Error with grpcserver: %v", err)

	}

	chatapp.RegisterBroadcastServer(grpcserver, server)
	grpcserver.Serve(lis)

}
