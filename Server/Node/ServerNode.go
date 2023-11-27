package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
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
	log.Printf("Bid recieved: %d from %s\n", in.Bid, in.Bidder)

	if !start && s.status == "closed" {
		start = true
		s.status = "open"

		go func() {
			timer := *duration
			time.Sleep(time.Duration(timer) * time.Second)
			start = false
			s.status = "closed"
			log.Printf("Auction over\n")
		}()
	}

	if s.isSlow {
		delay := rand.Int63n(5000)
		log.Printf("Slow node sleeping for %d milliseconds\n", delay)
		time.Sleep(time.Duration(delay) * time.Millisecond)
		log.Printf("Slow node finished sleeping\n")
	}

	if in.Bid > s.currentbid && start && s.status == "open" {
		s.currentbid = in.Bid
		s.currentbidder = in.Bidder
		return &Node.BidResponse{Bidder: s.currentbidder, Bid: s.currentbid, Status: s.status}, nil
	} else {
		return &Node.BidResponse{Bidder: s.currentbidder, Bid: s.currentbid, Status: s.status}, nil
	}
}

func (s *AuctionServer) Result(ctx context.Context, in *Node.ResultRequest) (*Node.ResultResponse, error) {

	// Print the result of the auction
	log.Printf("Result request recieved from %s\n", in.Bidder)

	if s.isSlow {
		delay := rand.Int63n(5000)
		log.Printf("Slow node sleeping for %d milliseconds\n", delay)
		time.Sleep(time.Duration(delay) * time.Millisecond)
		log.Printf("Slow node finished sleeping\n")
	}

	return &Node.ResultResponse{Bidder: s.currentbidder, Bid: s.currentbid, Status: s.status}, nil
}

func main() {

	flag.Parse()

	logPath := fmt.Sprintf("Server/Logs/Node%s.log", *nodePort)

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}

	log.SetOutput(io.MultiWriter(file, os.Stdout))

	// Create a TCP listener on the specified address and port
	listener, err := net.Listen("tcp", fmt.Sprintf(":%v", *nodePort))
	if err != nil {
		log.Fatal("Error starting the server:", err)
	}

	fmt.Println("Server node started on localhost:" + *nodePort)

	grpcServer := grpc.NewServer()

	AuctionServer := &AuctionServer{port: *nodePort, isSlow: *isSlow, currentbid: 0, currentbidder: "", status: "closed"}

	Node.RegisterAuctionServer(grpcServer, AuctionServer)

	log.Printf("Server node listening on port %v", *nodePort)

	go func() {
		//Make commands to check values of currentbid, currentbidder, status
		scanner := bufio.NewScanner(os.Stdin)

		for scanner.Scan() {
			command := scanner.Text()
			switch command {
			case "currentbid":
				log.Printf("Current bid is: %d\n", AuctionServer.currentbid)
			case "currentbidder":
				log.Printf("Current bidder is: %s\n", AuctionServer.currentbidder)
			case "status":
				log.Printf("Status is: %s and start bool is %v\n", AuctionServer.status, start)
			case "restart":
				start = false
				AuctionServer.status = "closed"
				AuctionServer.currentbid = 0
				AuctionServer.currentbidder = ""
				log.Printf("")
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
