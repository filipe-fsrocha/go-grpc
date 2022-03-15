package main

import (
	"context"
	"io"
	"log"
	"net"
	"time"

	"pb/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type server struct{}

func main() {
	grpcServer := grpc.NewServer()
	pb.RegisterTestGRPCServiceServer(grpcServer, &server{})
	listener, err := net.Listen("tcp", "0.0.0.0:50051")

	if err != nil {
		log.Fatalf("Failed to listener: %v", err)
	}

	reflection.Register(grpcServer)
	log.Println("Listening on tcp://localhost:50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}

func (s *server) Hello(ctx context.Context, request *pb.HelloRequest) (*pb.HelloResponse, error) {
	result := "Hello " + request.GetName()

	res := &pb.HelloResponse{
		Msg: result,
	}

	return res, nil
}

func (s *server) CreateUser(ctx context.Context, request *pb.CreateUserRequest) (*pb.CreateUserResponse, error) {
	log.Printf("Creating user: %v", request)

	res := &pb.CreateUserResponse{
		Username: request.Username,
		Msg:      "User created successfully",
	}

	return res, nil
}

func (s *server) Fibonacci(in *pb.FibonacciRequest, stream pb.TestGRPCService_FibonacciServer) error {
	n := in.GetN()
	var i int32

	for i = 1; i <= n; i++ {
		res := &pb.FibonacciReponse{
			Result: FibonacciResolver(i),
		}
		stream.Send(res)
		time.Sleep(time.Second * 2)
	}

	return nil
}

func (s *server) StreamRPC(stream pb.TestGRPCService_StreamRPCServer) error {
	log.Println("Started stream")
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}

		res := &pb.StreamDataResponse{
			Result: in.GetMsg(),
		}

		if err := stream.Send(res); err != nil {
			return err
		}
	}
}

func FibonacciResolver(n int32) int32 {
	if n <= 1 {
		return n
	}

	return FibonacciResolver(n-1) + FibonacciResolver(n-2)
}
