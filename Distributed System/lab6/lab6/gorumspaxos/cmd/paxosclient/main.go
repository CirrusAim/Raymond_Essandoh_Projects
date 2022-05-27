package main

import (
	"context"
	paxos "dat520/lab6/gorumspaxos"
	pb "dat520/lab6/gorumspaxos/proto"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	gorums "github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	BALANCE    = 0
	DEPOSIT    = 1
	WITHDRAWAL = 2
)

func showUsage() {
	fmt.Println("Select Transaction type : ")
	fmt.Println("1. Get Balance.")
	fmt.Println("2. Deposit.")
	fmt.Println("3. Withdrawal.")
	fmt.Println("4. Quit.")

}

func main() {
	var (
		saddrs   = flag.String("addrs", "", "server addresses separated by ','")
		clientId = flag.String("clientId", "", "Client Id, different for each client")
	)

	flag.Usage = func() {
		log.Printf("Usage: %s [OPTIONS]\nOptions:", os.Args[0])
		flag.PrintDefaults()
	}
	flag.Parse()

	addrs := strings.Split(*saddrs, ",")
	if len(addrs) == 0 {
		log.Fatalln("no server addresses provided")
	}
	config, mgr := createConfiguration(addrs)
	defer mgr.Close()

	clientSeq := 1

	for {
		showUsage()
		var option int
		fmt.Scan(&option)
		switch option {
		case 1:
			{
				fmt.Print("Enter account number: ")
				var accNum int
				fmt.Scan(&accNum)
				req := pb.Value{
					ClientID:      *clientId,
					ClientSeq:     uint32(clientSeq),
					AccountNumber: uint32(accNum),
					Tx: &pb.Transaction{
						Op:     BALANCE,
						Amount: 0,
					}}
				resp := doSendRequest(config, &req)
				log.Printf("response: %+v\t for the client request: %+v", resp, &req)
				clientSeq++
			}

		case 2:
			fmt.Print("Enter account number: ")
			var accNum int
			fmt.Scan(&accNum)
			fmt.Print("\nEnter amount to deposit: ")
			var Amt int32
			fmt.Scan(&Amt)
			req := pb.Value{
				ClientID:      *clientId,
				ClientSeq:     uint32(clientSeq),
				AccountNumber: uint32(accNum),
				Tx: &pb.Transaction{
					Op:     DEPOSIT,
					Amount: Amt,
				}}
			resp := doSendRequest(config, &req)
			log.Printf("response: %+v\t for the client request: %+v", resp, &req)
			clientSeq++
		case 3:
			fmt.Print("Enter account number: ")
			var accNum int
			fmt.Scan(&accNum)
			fmt.Print("\nEnter amount to withdraw: ")
			var Amt int32
			fmt.Scan(&Amt)
			req := pb.Value{
				ClientID:      *clientId,
				ClientSeq:     uint32(clientSeq),
				AccountNumber: uint32(accNum),
				Tx: &pb.Transaction{
					Op:     WITHDRAWAL,
					Amount: Amt,
				}}
			resp := doSendRequest(config, &req)
			log.Printf("response: %v\t for the client request: %+v", resp, &req)
			clientSeq++
		case 4:
			os.Exit(0)
		default:
			fmt.Println("Invalid option")
		}
	}

	// resp := doSendRequest(config, &req)
	// log.Printf("response: %v\t for the client request: %v", resp, &req)

	// clientRequests := strings.Split(*clientRequest, ",")
	// if len(clientRequests) == 0 {
	// 	log.Fatalln("no client requests are provided")
	// }

	// // start a initial proposer
	// ClientStart(addrs, clientRequests, clientId)
}

// ClientStart creates the configuration with the list of replicas addresses, which are read from the
// commandline. From the list of clientRequests, send each request to the configuration and
// wait for the reply. Upon receiving the reply send the next request.
func ClientStart(addrs []string, clientRequests []string, clientId *string) {
	log.Printf("Connecting to %d Paxos replicas: %v", len(addrs), addrs)
	config, mgr := createConfiguration(addrs)
	defer mgr.Close()
	for index := range clientRequests {
		req := pb.Value{ClientID: *clientId, ClientSeq: uint32(index), Tx: &pb.Transaction{}}
		resp := doSendRequest(config, &req)
		log.Printf("response: %v\t for the client request: %v", resp, &req)
	}
	log.Printf("Successfully completed client id %s", *clientId)

}

// Internal: doSendRequest can send requests to paxos servers by quorum call and
// for the response from the quorum function.
func doSendRequest(config *pb.Configuration, value *pb.Value) *pb.Response {
	waitTimeForRequest := 5 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), waitTimeForRequest)
	defer cancel()
	resp, err := config.ClientHandle(ctx, value)
	if err != nil {
		log.Fatalf("ClientHandle quorum call error: %v", err)
	}
	if resp == nil {
		log.Println("response is nil")
	}
	return resp
}

// createConfiguration creates the gorums configuration with the list of addresses.
func createConfiguration(addrs []string) (configuration *pb.Configuration, manager *pb.Manager) {
	mgr := pb.NewManager(gorums.WithDialTimeout(5*time.Second),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(), // block until connections are made
			grpc.WithTransportCredentials(insecure.NewCredentials()), // disable TLS
		),
	)
	quorumSize := (len(addrs)-1)/2 + 1
	qfspec := paxos.NewPaxosQSpec(quorumSize)
	config, err := mgr.NewConfiguration(qfspec, gorums.WithNodeList(addrs))
	if err != nil {
		log.Fatalf("Error in forming the configuration: %v\n", err)
		return nil, nil
	}
	return config, mgr
}
