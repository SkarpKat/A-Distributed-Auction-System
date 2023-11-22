# A-Distributed-Auction-System

## Implementation
The server is based on **leader-less(active) replication**. The downside for using this is the chance of latency from one of the nodes. This however will not be present in this implementaion (*could be implemented with a delay on the server side to simulate, can use a boolean flag to activate a slow server*)

## Client
The client is able to bid on the auction that the server is hosting.

## Server
The "*server*" is build up of several nodes that the client multicasts to to keep them all up yo date on what's happening in the auction.