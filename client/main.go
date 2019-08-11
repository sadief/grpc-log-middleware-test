package main

import (
	"context"
	"log"
	"time"

	api "github.com/axiomzen/grpc-testing/api"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50051"
)

func main() {

	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewPingServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Ping(ctx, &api.PingRequest{Message: "Ping"})
	if err != nil {
		log.Fatalf("could not ping: %v", err)
	}
	log.Printf("Successfully Pinged: %s", r.Message)
}
