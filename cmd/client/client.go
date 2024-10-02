package main

import (
	"context"
	"log"
	"time"

	pb "github.com/jboykin-bread/mirrord-traffic-stealing-reproduction/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/status"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.NewClient("localhost:9000", grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("Failed to connect to server: %v", err)
	}
	defer conn.Close()

	client := pb.NewColorServiceClient(conn)

	for {
		// Create a context with a 5-second timeout for the request
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

		// Make the gRPC request
		resp, err := client.GetColor(ctx, &pb.ColorRequest{})
		if err != nil {
			st, _ := status.FromError(err)
			log.Printf("Error calling GetColor: %v", st.Message())
		} else {
			log.Printf("Color: %s", resp.Color)
		}

		// Explicitly cancel the context to avoid memory leaks
		cancel()

		// Wait for 2 seconds before the next request
		time.Sleep(2 * time.Second)
	}
}
