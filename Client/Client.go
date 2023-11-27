package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	Node "github.com/SkarpKat/A-Distributed-Auction-System/Server/proto"
	"google.golang.org/grpc"
)

var (
	clientName = flag.String("name", "client", "The name of the client")
	nodePorts  = []string{"8080", "8081", "8082"}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	flag.Parse()

	logPath := fmt.Sprintf("Client/Logs/%s.log", *clientName)

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.SetOutput(file)

	clientConnections := make([]Node.AuctionClient, len(nodePorts))

	for i, port := range nodePorts {
		conn, err := grpc.Dial("localhost:"+port, grpc.WithInsecure(), grpc.WithBlock(), grpc.WithTimeout(60000*time.Millisecond))
		if err != nil {
			log.Printf("Failed to dial to node with port: %v", err)
			continue
		}
		defer conn.Close()
		clientConnections[i] = Node.NewAuctionClient(conn)
	}

	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Commands:\n bid - Bid on the auction\n result - Get the result of the auction\n exit - Exit the client\n")
	for scanner.Scan() {
		command := scanner.Text()
		command = strings.ToLower(command)
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
				fmt.Printf("You are leading with a bid of %d kr.\n", winningBid)
			} else {
				fmt.Printf("Your bid of %d could not competet with %d kr. from %v\n", bid, winningBid, winner)
			}
		case "result":
			winner := ""
			winningBid := 0
			state := ""
			for _, client := range clientConnections {
				rsp, err := client.Result(ctx, &Node.ResultRequest{Bidder: *clientName})
				if err != nil {
					log.Printf("Failed to get result: %v", err)
					continue
				}
				winner = rsp.Bidder
				winningBid = int(rsp.Bid)
				state = rsp.Status
			}

			if state == "closed" {
				if winner == *clientName {
					fmt.Printf("You won with a bid of %d\n", winningBid)
				} else {
					fmt.Printf("You lost to a bid of %d\n", winningBid)
				}
			} else {
				fmt.Printf("The auction is still open\n")
			}

		case "exit":
			os.Exit(0)
		default:
			fmt.Printf("Invalid command\n Commands:\n bid - Bid on the auction\n result - Get's the result of the auction and tells if it's still going\n")
		}

	}

}
