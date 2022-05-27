package main

import (
	"context"
	paxos "dat520/lab5/gorumspaxos"
	pb "dat520/lab5/gorumspaxos/proto"
	"encoding/json"
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

func main() {
	var (
		saddrs        = flag.String("addrs", "", "server addresses separated by ','")
		clientRequest = flag.String("clientRequest", "", "client requests seperated by ','")
		clientId      = flag.String("clientId", "", "Client Id, different for each client")
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

	clientRequests := strings.Split(*clientRequest, ",")
	if len(clientRequests) == 0 {
		log.Fatalln("no client requests are provided")
	}

	// start a initial proposer
	ClientStart(addrs, clientRequests, clientId)
}

// ClientStart creates the configuration with the list of replicas addresses, which are read from the
// commandline. From the list of clientRequests, send each request to the configuration and
// wait for the reply. Upon receiving the reply send the next request.
func ClientStart(addrs []string, clientRequests []string, clientId *string) {
	log.Printf("Connecting to %d Paxos replicas: %v", len(addrs), addrs)
	config, mgr := createConfiguration(addrs)
	defer mgr.Close()
	data := make([]float64, 0)
	for index, request := range clientRequests {
		start := time.Now()
		req := pb.Value{ClientID: *clientId, ClientSeq: uint32(index), ClientCommand: request}
		//resp := doSendRequest(config, &req)
		doSendRequest(config, &req)

		data = append(data, float64(time.Since(start).Seconds()))

		//log.Printf("response: %v\t for the client request: %v", resp, &req)
	}

	//mean, _ := st.Mean(data)
	//median, _ := st.Median(data)
	//sd, _ := st.StandardDeviation(data)
	//max, _ := st.Max(data)
	//min, _ := st.Min(data)
	logfile, err := os.Create(fmt.Sprintf("/Users/hysensw/Videos/lab5/lab5/lab5/gorumspaxos/cmd/paxosclient/client_%s.json", *clientId))
	if err != nil {
		panic(err)
	}
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(logfile)
	mp := make(map[string][]float64)
	mp["data"] = data
	//mp["median"] = median
	//mp["sd"] = sd
	//mp["max"] = max
	//mp["min"] = min
	bt, err := json.Marshal(mp)

	log.Printf("%v", string(bt))
	log.SetOutput(os.Stdout)

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
