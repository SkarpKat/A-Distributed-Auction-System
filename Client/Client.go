package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sync"

	Node "github.com/SkarpKat/A-Distributed-Auction-System/Server/proto"
	"google.golang.org/grpc"
)

var (
	clientName = flag.String("name", "client", "The name of the client")
	// clientID   = flag.Int("id", 0, "The ID of the client")
	nodePorts = []string{"8080", "8081", "8082"}
	joined    = false
)

type resultResponse struct {
	bidder string
	bid    int64
	status string
}

type bidResponse struct {
	bidder string
	bid    int64
	status string
}

type statusResponse struct {
	bidder string
	bid    int64
	status string
	winner bool
}

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

	for scanner.Scan() {
		command := scanner.Text()
		switch command {
		case "bid":
			fmt.Printf("Enter the amount you would like to bid: ")
			bid := 0
			// Read the bid from the user
			_, err := fmt.Scanf("%d", &bid)
			if err != nil {
				log.Printf("Failed to read bid: %v", err)
			}
			responses := make([]bidResponse, len(clientConnections))

			// Send the bid to all nodes
			for _, client := range clientConnections {
				rsp, err := client.Bid(ctx, &Node.BidRequest{Bidder: *clientName, Bid: int64(bid)})
				if err != nil {
					log.Printf("Failed to bid: %v", err)
				}
				log.Printf("Bid response: %v", rsp.Bidder)
				responses = append(responses, bidResponse{bidder: rsp.Bidder, bid: rsp.Bid, status: rsp.Status})
			}

			// Check if all responses are the same and print the status
			for _, rsp := range responses {
				// Check if all responses are the same
				if rsp.bidder != responses[0].bidder || rsp.bid != responses[0].bid || rsp.status != responses[0].status {
					log.Printf("Error: Not all responses are the same")
				}

			}
			if responses[0].bidder == *clientName {
				fmt.Printf("You are winning with a bid of %d\n", responses[0].bid)
			} else {
				fmt.Printf("Your bid of %d could not competet with %d from %v\n", bid, responses[0].bid, responses[0].bidder)
			}

		case "status":
			// Create a slice of status responses
			responses := make([]statusResponse, len(clientConnections))

			// Get the status from all nodes
			for _, client := range clientConnections {
				rsp, err := client.Status(ctx, &Node.StatusRequest{Bidder: *clientName})
				if err != nil {
					log.Printf("Failed to get status: %v", err)
				}
				responses = append(responses, statusResponse{bidder: rsp.Bidder, bid: rsp.Bid, status: rsp.Status, winner: rsp.Winner})
			}

			// Check if all responses are the same and print the status
			for _, rsp := range responses {
				// Check if all responses are the same

				if rsp.bidder != responses[0].bidder || rsp.bid != responses[0].bid || rsp.status != responses[0].status {
					log.Printf("Error: Not all responses are the same")
				}
			}

			if responses[0].winner {
				fmt.Printf("You are winning with a bid of %d\n", responses[0].bid)
			} else {
				fmt.Printf("You are not winning against a bid of %d\n", responses[0].bid)
			}
		case "join":
			if joined {
				fmt.Printf("You have already joined the auction\n")
			} else {
				joined = true
				fmt.Printf("You have joined the auction\n")
				go func() {
					// Create wait group
					wg := sync.WaitGroup{}

					// Create a slice of result responses
					// responses := make([]resultResponse, len(clientConnections))

					// Create mutex
					// arbiter := sync.Mutex{}

					var winner string
					var bid int64

					for _, client := range clientConnections {
						wg.Add(1)
						go func(client Node.AuctionClient) {
							rsp, err := client.Result(ctx, &Node.ResultRequest{Bidder: *clientName})
							if err != nil {
								log.Printf("Failed to get result: %v", err)
							}
							// Print the result of the auction
							fmt.Printf("The winner is %s with a bid of %d\n", rsp.Bidder, rsp.Bid)

							winner = rsp.Bidder
							bid = rsp.Bid

							// Add the response to the slice
							// response := resultResponse{bidder: rsp.Bidder, bid: rsp.Bid, status: rsp.Status}
							// arbiter.Lock()
							// responses = append(responses, response)
							// arbiter.Unlock()
							wg.Done()
						}(client)

					}
					wg.Wait()

					// Check if all responses are the same and print the status

					// for _, rsp := range responses {
					// 	// Print all responses
					// 	fmt.Printf("Bidder: %s, Bid: %d, Status: %s\n", rsp.bidder, rsp.bid, rsp.status)

					// 	//Check if data is the same
					// 	if rsp.bidder != responses[0].bidder || rsp.bid != responses[0].bid || rsp.status != responses[0].status {
					// 		log.Printf("Error: Not all responses are the same")
					// 	}
					// }

					fmt.Printf("The auction is over\n")
					fmt.Printf("The winner is %s with a bid of %d\n", winner, bid)

				}()
			}
		case "exit":
			os.Exit(0)
		default:
			fmt.Printf("Invalid command\n Commands:\n bid - Bid on the auction\n status - Get the status of the auction\n Join - Join the auction\n")
		}

	}

}
