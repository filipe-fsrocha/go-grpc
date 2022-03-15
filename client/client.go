package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"pb/pb"
	"strconv"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	opts := grpc.WithTransportCredentials(insecure.NewCredentials())
	connection, err := grpc.Dial("localhost:50051", opts)

	if err != nil {
		log.Fatalf("Could not connect: %v", err)
	}

	defer connection.Close()

	client := pb.NewTestGRPCServiceClient(connection)
	stream, _ := client.StreamRPC(context.Background())
	waitcResponse := make(chan struct{})

	fmt.Println("#### API Unary ###")
	Hello(client)
	CreateUser(client)
	fmt.Println("#### API Stream Server ####")
	Fibonacci(client)

	fmt.Println("#### API Bidirectional streaming ###")
	go func() {
		for i := 0; i <= 5; i++ {
			time.Sleep(2 * time.Second)
			log.Println("Client request: " + "User: " + strconv.Itoa(i))
			stream.Send(&pb.StreamDataRequest{Msg: "User " + strconv.Itoa(i)})
		}
		stream.CloseSend()
	}()

	go func() {
		for {
			res, err := stream.Recv()
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Fatalf("Erro when receiving response: %v", err)
			}
			log.Println("Server response: ", res)
		}
		close(waitcResponse)
	}()
	<-waitcResponse
}

func Fibonacci(client pb.TestGRPCServiceClient) {
	request := &pb.FibonacciRequest{
		N: 5,
	}

	responseStream, err := client.Fibonacci(context.Background(), request)

	if err != nil {
		LogError(err)
	}

	for {
		stream, err := responseStream.Recv()
		if err == io.EOF {
			break
		}
		log.Printf("Fibonacci: %v", stream.GetResult())
	}
}

func Hello(client pb.TestGRPCServiceClient) {
	request := &pb.HelloRequest{
		Name: "Filipe Rocha",
	}

	res, err := client.Hello(context.Background(), request)

	if err != nil {
		LogError(err)
	}

	log.Println(res.Msg)
}

func CreateUser(client pb.TestGRPCServiceClient) {
	request := &pb.CreateUserRequest{
		Username: "user.test",
		FullName: "Usuário test",
		Email:    "user.test@email.com",
	}

	res, err := client.CreateUser(context.Background(), request)

	if err != nil {
		LogError(err)
	}

	log.Println(res)
}

func LogError(err error) {
	log.Fatalf("Error during the execution %v", err)
}
