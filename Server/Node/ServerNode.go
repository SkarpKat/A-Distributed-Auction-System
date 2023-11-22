package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	Node "github.com/SkarpKat/A-Distributed-Auction-System/Server/proto"
	"google.golang.org/grpc"
)

var (
	nodePort = flag.String("port", "8080", "The port of the node")
	isSlow   = flag.Bool("slow", false, "Whether the node is slow or not")
	start    = false
	duration = flag.Int("duration", 60, "The duration of the auction in seconds")
)

type AuctionServer struct {
	Node.UnimplementedAuctionServer

	port          string // The port of the node
	isSlow        bool   // Whether the node is slow or not
	currentbid    int64  // The current bid of the auction
	currentbidder string // The current bidder of the auction
	status        string // The status of the auction
}

func (s *AuctionServer) Bid(ctx context.Context, in *Node.BidRequest) (*Node.BidResponse, error) {
	// Print the recieved bid and bidder
	fmt.Printf("Bid recieved: %d from %s\n", in.Bid, in.Bidder)

	if in.Bid > s.currentbid && start {
		s.currentbid = in.Bid
		s.currentbidder = in.Bidder
		return &Node.BidResponse{Bidder: s.currentbidder, Bid: s.currentbid, Status: s.status}, nil
	} else {
		return &Node.BidResponse{Bidder: s.currentbidder, Bid: s.currentbid, Status: s.status}, nil
	}
}

func (s *AuctionServer) Status(ctx context.Context, in *Node.StatusRequest) (*Node.StatusResponse, error) {
	// Print the requester
	fmt.Printf("Status requested by %s\n", in.Bidder)

	if in.Bidder == s.currentbidder {
		return &Node.StatusResponse{Bidder: s.currentbidder, Bid: s.currentbid, Status: s.status, Winner: true}, nil
	} else {
		return &Node.StatusResponse{Bidder: s.currentbidder, Bid: s.currentbid, Status: s.status, Winner: false}, nil
	}
}

func (s *AuctionServer) Result(ctx context.Context, in *Node.ResultRequest) (*Node.ResultResponse, error) {
	// Print the requester who joined the auction
	fmt.Printf("%s has joined the auction\n", in.Bidder)
	start = true

	if start {
		timer := *duration
		time.Sleep(time.Duration(timer) * time.Second)
	}

	start = false

	// Print the result of the auction
	fmt.Printf("The winner is %s with a bid of %d\n", s.currentbidder, s.currentbid)

	return &Node.ResultResponse{Bidder: s.currentbidder, Bid: s.currentbid, Status: s.status}, nil
}

func main() {

	flag.Parse()

	// Create a TCP listener on the specified address and port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", *nodePort))
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}

	fmt.Println("Server node started on localhost:" + *nodePort)

	grpcServer := grpc.NewServer()

	AuctionServer := &AuctionServer{port: *nodePort, isSlow: *isSlow, currentbid: 0, currentbidder: "", status: "open"}

	Node.RegisterAuctionServer(grpcServer, AuctionServer)

	log.Printf("Server node listening on port %v", *nodePort)

	go func() {
		//Make commands to check values of currentbid, currentbidder, status
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			command := scanner.Text()
			switch command {
			case "currentbid":
				fmt.Printf("Current bid is: %d\n", AuctionServer.currentbid)
			case "currentbidder":
				fmt.Printf("Current bidder is: %s\n", AuctionServer.currentbidder)
			case "status":
				fmt.Printf("Status is: %s and start bool is %v\n", AuctionServer.status, start)
			case "start":
				fmt.Printf("Starting auction\n")
				AuctionServer.status = "open"
				start = true
			case "stop":
				fmt.Printf("Stopping auction\n")
				AuctionServer.status = "closed"
				start = false
			case "shutdown":
				os.Exit(0)
			default:
				fmt.Printf("Command not recognized\n Commands are: currentbid, currentbidder, status, shutdown\n")
			}
		}
	}()

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
