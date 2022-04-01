package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"strconv"
	"time"

	pb "dat520/lab2/grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var (
	help = flag.Bool(
		"help",
		false,
		"Show usage help",
	)
	endpoint = flag.String(
		"endpoint",
		"localhost:12111",
		"Endpoint on which server runs or to which client connects",
	)
)

// Usage prints usage info
func Usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "\nOptions:\n")
	flag.PrintDefaults()
}

//  return pb.InsertRequest object
func getInsertData(key, val string) *pb.InsertRequest {
	return &pb.InsertRequest{
		Key:   key,
		Value: val,
	}
}

func getLookUpData(key string) *pb.LookupRequest {
	return &pb.LookupRequest{
		Key: key,
	}
}

func genKVPairDataForTesting(count int) map[string]string {
	key := "Key"
	val := "Value"
	mp := map[string]string{}
	for i := 1; i <= count; i++ {
		iKey := key + strconv.Itoa(i)
		iVal := val + strconv.Itoa(i)
		mp[iKey] = iVal
	}
	return mp
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	flag.Usage = Usage
	flag.Parse()
	if *help {
		flag.Usage()
		return
	}

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial(*endpoint, opts...)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()

	// this is the generated interface for clients
	client := pb.NewKeyValueServiceClient(conn)

	count := 10

	// generate the data for testing
	mp := genKVPairDataForTesting(count)

	for k, v := range mp {
		fmt.Printf("Inserting  %s - %s pair\n", k, v)
		// call RPC method insert of server
		resp, _ := client.Insert(ctx, getInsertData(k, v))
		if !resp.Success {
			fmt.Printf("Failed to insert %s - %s pair\n", k, v)
		}
	}

	for k := range mp {
		resp, _ := client.Lookup(ctx, getLookUpData(k))
		if len(resp.GetValue()) >= 0 {
			fmt.Printf("Value found for Key  %s is : %s\n", k, resp.GetValue())
		}
	}

	allKeys, _ := client.Keys(ctx, &pb.KeysRequest{})

	fmt.Printf("All keys Present are: %+v\n", allKeys.Keys)

}
