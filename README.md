# A-Distributed-Auction-System

## Implementation
The server is based on **leader-less(active) replication**. The downside for using this is the chance of latency from one of the nodes. This however will not be present in this implementaion (*could be implemented with a delay on the server side to simulate, can use a boolean flag to activate a slow server*)

## Client
The client is able to bid on the auction that the server is hosting and ask for the result.

## Server
The "*server*" is build up of several nodes that the client multicasts to to keep them all up yo date on what's happening in the auction.

## How To Use
Here we will go over how to use the program and what flags there are and a basic setup to get the idea on how it works.
### Client
Here we will go over the client and how it's used
#### Flags
The Client is simple and the flags are as followed:
- ```-name```: which is the name of the client.
The ports are defined the code in an array. You can expand it if you want more nodes.
#### Commands
The commands the client can use in the terminal is:
- ```bid```: this command make it possible to make a bid. After the command you will put the amount. You can not do ```bid 1000``` this will not work. If it's the first bid it will start the auction.
- ```result```: this will give the result of the auction. If the auction is still underway it will say so.
- ```exit```: this will kill the client.
### Server node
Here we will go over the node.
#### Flags
The flags that are on the ServerNode is as followed:
- ```-port```: Here you set the port that the node will be listining on.
- ```-slow```: This will be off by default it's here to simulate the down side of having to multicast since one of the downfall is that one node could be slow an you need to wait for it. It can be turned on by setting it to "true" and there will be a random delay between 0 - 5 seconds.
- ```-duration```: This is the duration of the auction. By default the auction will be 1 min but this can be changed. Be aware that it all nodes need to have the same duration.
#### Commands
The node has several commands that can be executed.
- ```currentbid```: this will print the current bid which is the highest.
- ```currentbidder```: this prints the current bidder which has the highest bid
- ```status```: this will print the status of the auction. This will be the status and the start bool value.
- ```restart```: this will restart the auction so the clients can run it again. **IMPORTANT:** all nodes should be restarted before doing a bid to keep nodes consistent.
- ```shutdown```: this will shutdown the node. You can also use CTRL + C.