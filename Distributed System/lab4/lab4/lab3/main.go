package main

import (
	"bufio"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"lab3/failuredetector"
	"lab3/leaderdetector"
	"lab3/network"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var serversDetail = map[int]string{
	1: "localhost:46789",
	2: "localhost:46790",
	3: "localhost:46791",
	4: "localhost:46792",
}

//var nodeIds = []int{4, 1}

var nodeIds = []int{4, 3, 2, 1}

var myId int
var myAddr string
var delta = 10 * time.Second

var (
	help = flag.Bool(
		"help",
		false,
		"Show usage help",
	)
	id = flag.Int(
		"i",
		-1,
		"ID of the Node",
	)
)

func main() {

	// register interrupts
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	signal.Notify(sigint, syscall.SIGTERM)

	flag.Usage = usage
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0)
	}

	if *id < 1 && *id > len(serversDetail) {
		log.Fatal("Invalid id.")
		return
	}

	// assign id and ip for this instance
	myId = *id
	myAddr = serversDetail[myId]

	// create a context for clean exit
	ctx, stop := context.WithCancel(context.Background())

	// wait for user input before proceeding
	fmt.Println("Press Enter to start the node.")
	readEnter()

	// create a channel for outgoing heartbeats
	hbSend := make(chan failuredetector.Heartbeat, len(nodeIds))

	// create a new leader instance
	leaderDet := leaderdetector.NewMonLeaderDetector(nodeIds)
	log.Printf("Leader id without connecting to other nodes is : %v", leaderDet.Leader())

	// subscribe to the leader change event
	subscriber := leaderDet.Subscribe()

	// create a new failure detector
	fd := failuredetector.NewEvtFailureDetector(myId, nodeIds, leaderDet, delta, hbSend)

	// start a UDP server to serve incoming heartbeats and start serving the requests
	server := network.NewUDPServer(serversDetail[myId], fd, hbSend)
	go server.ServeUDP(ctx)

	// create client connection to each of the nodes including the current nodes as well to send heartbeats
	client := network.NewUDPClient(serversDetail)

	// start failure detector
	fd.Start()

	go notifyLeaderChanges(ctx, subscriber)
	go sendHeartbeats(ctx, client, hbSend)

	<-sigint
	stop()
	client.Close()
}

// poll for leader change. If changed, print it
func notifyLeaderChanges(ctx context.Context, subs <-chan int) {
	for {
		select {
		case <-ctx.Done():
			return
		case leader := <-subs:
			log.Println("New leader id : ", leader)
		}
	}
}

// send the heartbeats to the other nodes
func sendHeartbeats(ctx context.Context, client *network.UDPClient, hbSend chan failuredetector.Heartbeat) {
	// heartbeat encoded to json
	type EncodedHeartbeat struct {
		From    int  `json:"from"`
		To      int  `json:"to"`
		Request bool `json:"request"`
	}

	for {
		select {
		case <-ctx.Done():
			return
		case hb := <-hbSend:
			//log.Printf("Sending : %+v", hb)
			newHeartbeat := EncodedHeartbeat{
				From:    hb.From,
				To:      hb.To,
				Request: hb.Request,
			}
			payload, err := json.Marshal(newHeartbeat)
			if err != nil {
				log.Printf("Unable to Marshall json. Skipping heartbeat. Err : %v", err)
				continue
			}
			client.SendPayload(hb.To, payload)
		}
	}
}

func readEnter() {
	reader := bufio.NewReader(os.Stdin)
	reader.ReadLine()
}

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}
