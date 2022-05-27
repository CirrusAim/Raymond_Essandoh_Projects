package gorumspaxos

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"
	"testing"
	"time"

	pb "dat520/lab6/gorumspaxos/proto"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func TestProposer(t *testing.T) {
	numberOfServers := 5
	numberOfClients := 3
	numberOfRequest := 5
	serverAddresses := make([]string, 0)
	nodeMap := make(map[string]uint32)
	lisMap := make(map[string]net.Listener)
	replicas := make([]*PaxosReplica, 0)
	for i := 0; i < numberOfServers; i++ {
		lis, err := net.Listen("tcp", "127.0.0.1:0")
		if err != nil {
			log.Fatal(err)
		}
		addr := lis.Addr().String()
		serverAddresses = append(serverAddresses, addr)
		nodeMap[addr] = uint32(i)
		lisMap[addr] = lis
	}

	for addr, id := range nodeMap {
		args := NewPaxosReplicaArgs{
			LocalAddr: addr,
			Id:        int(id),
			NodeMap:   nodeMap,
		}
		replica := NewPaxosReplica(args)
		replicas = append(replicas, replica)
		go replica.ServerStart(lisMap[addr])
	}
	time.Sleep(5 * time.Second)
	var wg sync.WaitGroup
	wg.Add(numberOfClients)
	for i := 0; i < numberOfClients; i++ {
		go func(id int) {
			mgr := pb.NewManager(gorums.WithDialTimeout(5*time.Second),
				gorums.WithGrpcDialOptions(
					grpc.WithBlock(), // block until connections are made
					grpc.WithTransportCredentials(insecure.NewCredentials()), // disable TLS
				),
			)
			quorumSize := (numberOfServers-1)/2 + 1
			qfspec := NewPaxosQSpec(quorumSize)
			config, err := mgr.NewConfiguration(qfspec, gorums.WithNodeMap(nodeMap))
			if err != nil {
				log.Fatalf("Error in forming the configuration: %v\n", err)
			}
			for k := 0; k < numberOfRequest; k++ {
				waitTimeForRequest := 5 * time.Second
				ctx, cancel := context.WithTimeout(context.Background(), waitTimeForRequest)
				defer cancel()
				req := pb.Value{
					ClientID:      fmt.Sprint(id),
					ClientSeq:     uint32(k),
					AccountNumber: 101,
					Tx:            &pb.Transaction{},
				}
				resp, err := config.ClientHandle(ctx, &req)
				if err != nil {
					t.Errorf("Got Error while waiting for the reply%v\n", err)
				}
				if resp.TxResult.AccountNumber != req.AccountNumber {
					t.Errorf("Response client command is incorrect -want %v +got %v\n", req.AccountNumber, resp.TxResult.AccountNumber)
				} else if resp.ClientID != req.ClientID {
					t.Errorf("Response client ID is incorrect -want %s +got %s\n", req.ClientID, resp.ClientID)
				} else if resp.ClientSeq != req.ClientSeq {
					t.Errorf("Response client ID is incorrect -want %d +got %d\n", req.ClientSeq, resp.ClientSeq)
				}
			}
			mgr.Close()
			wg.Done()
		}(i)
	}
	wg.Wait()
	for _, replica := range replicas {
		replica.Stop()
	}
}
