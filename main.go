package main

import (
	"context"
	"flag"
	"net"

	pb "github.com/jboykin-bread/mirrord-traffic-stealing-reproduction/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"log/slog"
)

// Server is used to implement the gRPC ColorService.
type Server struct {
	pb.UnimplementedColorServiceServer
	color string
}

// GetColor implements the GetColor RPC for the gRPC service.
func (s *Server) GetColor(ctx context.Context, in *pb.ColorRequest) (*pb.ColorResponse, error) {
	slog.Info("gRPC Request: GetColor")
	return &pb.ColorResponse{Color: s.color}, nil
}

func main() {
	// flag for color
	color := flag.String("color", "blue", "the color to return in the response")
	flag.Parse()

	// Set up the gRPC server on port 9000
	lis, err := net.Listen("tcp", ":9000")
	if err != nil {
		slog.Error("Failed to listen on port 9000", "error", err)
		panic(err)
	}

	grpcServer := grpc.NewServer()
	reflection.Register(grpcServer)
	s := &Server{color: *color}
	pb.RegisterColorServiceServer(grpcServer, s)

	slog.Info("Starting gRPC server on :9000")
	if err := grpcServer.Serve(lis); err != nil {
		slog.Error("Failed to serve gRPC", "error", err)
		panic(err)
	}
}
