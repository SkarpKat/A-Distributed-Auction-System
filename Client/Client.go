package main

import (
	"context"
	"flag"
	"log"

	Node "github.com/SkarpKat/A-Distributed-Auction-System/Server/proto"
	"google.golang.org/grpc"
)

var (
	clientName = flag.String("name", "client", "The name of the client")
	clientID   = flag.Int("id", 0, "The ID of the client")
	nodePorts  = []string{"8080", "8081", "8082", "8083", "8084", "8085", "8086", "8087", "8088", "8089"}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.Parse()

	clientConnections := make([]Node.AuctionClient, len(nodePorts))

	for i, port := range nodePorts {
		conn, err := grpc.DialContext(ctx, "localhost:"+port, grpc.WithInsecure())
		if err != nil {
			log.Printf("Failed to dial to node with port: %v", err)
		}
		defer conn.Close()
		clientConnections[i] = Node.NewAuctionClient(conn)
	}

}
