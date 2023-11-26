package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	Node "github.com/SkarpKat/A-Distributed-Auction-System/Server/proto"
	"google.golang.org/grpc"
)

var (
	clientName = flag.String("name", "client", "The name of the client")
	// clientID   = flag.Int("id", 0, "The ID of the client")
	nodePorts = []string{"8080", "8081", "8082"}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.Parse()

	clientConnections := make([]Node.AuctionClient, len(nodePorts))

	for i, port := range nodePorts {
		conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Printf("Failed to dial to node with port: %v", err)
		}
		defer conn.Close()
		clientConnections[i] = Node.NewAuctionClient(conn)
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Commands:\n bid - Bid on the auction\n result - Get the result of the auction\n exit - Exit the client\n")
	for scanner.Scan() {
		command := scanner.Text()
		switch command {
		case "bid":
			fmt.Printf("Enter the amount you would like to bid: ")
			bid := 0
			if scanner.Scan() {
				input := scanner.Text()
				_, err := fmt.Sscanf(input, "%d", &bid)
				if err != nil {
					log.Printf("Failed to read bid: %v", err)
				}
			}
			// Read the bid from the user
			// _, err := fmt.Scan(&bid)

			// if err != nil {
			// 	log.Printf("Failed to read bid: %v", err)
			// }

			winner := ""
			winningBid := 0

			// Send the bid to all nodes
			for _, client := range clientConnections {
				rsp, err := client.Bid(ctx, &Node.BidRequest{Bidder: *clientName, Bid: int64(bid)})
				if err != nil {
					log.Printf("Failed to bid: %v", err)
					continue
				}
				winner = rsp.Bidder
				winningBid = int(rsp.Bid)
			}

			if winner == *clientName {
				fmt.Printf("You are winning with a bid of %d\n", winningBid)
			} else {
				fmt.Printf("Your bid of %d could not competet with %d from %v\n", bid, winningBid, winner)
			}
		case "result":
			winner := ""
			winningBid := 0
			for _, client := range clientConnections {
				rsp, err := client.Result(ctx, &Node.ResultRequest{Bidder: *clientName})
				if err != nil {
					log.Printf("Failed to get result: %v", err)
					continue
				}
				winner = rsp.Bidder
				winningBid = int(rsp.Bid)
			}

			if winner == *clientName {
				fmt.Printf("You won with a bid of %d\n", winningBid)
			} else {
				fmt.Printf("You lost with a bid of %d\n", winningBid)
			}
		case "exit":
			os.Exit(0)
		default:
			fmt.Printf("Invalid command\n Commands:\n bid - Bid on the auction\n status - Get the status of the auction\n Join - Join the auction\n")
		}

	}

}
