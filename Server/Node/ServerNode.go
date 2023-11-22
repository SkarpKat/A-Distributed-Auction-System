package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"time"

	Node "github.com/SkarpKat/A-Distributed-Auction-System/Server/proto"
	"google.golang.org/grpc"
)

var (
	nodePort = flag.String("port", "8080", "The port of the node")
	isSlow   = flag.Bool("slow", false, "Whether the node is slow or not")
	start    = false
	duration = flag.Int("duration", 0, "The duration of the auction")
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

	if in.Bid > s.currentbid {
		s.currentbid = in.Bid
		s.currentbidder = in.Bidder
		return &Node.BidResponse{Bidder: s.currentbidder, Bid: s.currentbid, Status: s.status}, nil
	} else {
		return &Node.BidResponse{Bidder: s.currentbidder, Bid: s.currentbid, Status: s.status}, nil
	}
}

func (s *AuctionServer) Status(ctx context.Context, in *Node.StatusRequest) (*Node.StatusResponse, error) {
	if in.Bidder == s.currentbidder {
		return &Node.StatusResponse{Bidder: s.currentbidder, Bid: s.currentbid, Status: s.status, Winner: true}, nil
	} else {
		return &Node.StatusResponse{Bidder: s.currentbidder, Bid: s.currentbid, Status: s.status, Winner: false}, nil
	}
}

func (s *AuctionServer) Result(ctx context.Context, in *Node.ResultRequest) (*Node.ResultResponse, error) {
	start = true

	if start {
		timer := *duration
		time.Sleep(time.Duration(timer) * time.Second)
	}

	start = false

	return &Node.ResultResponse{Bidder: s.currentbidder, Bid: s.currentbid, Status: s.status}, nil
}

func main() {
	// Define the address and port for the server node
	address := "localhost"

	flag.Parse()

	// Create a TCP listener on the specified address and port
	listener, err := net.Listen("tcp", address+":"+*nodePort)
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}

	fmt.Println("Server node started on", address+":"+*nodePort)

	grpcServer := grpc.NewServer()

	AuctionServer := &AuctionServer{port: *nodePort, isSlow: *isSlow, currentbid: 0, currentbidder: "", status: "open"}

	Node.RegisterAuctionServer(grpcServer, AuctionServer)

	log.Printf("Server node listening on port %v", *nodePort)

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}

}
