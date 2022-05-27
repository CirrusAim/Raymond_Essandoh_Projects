package main

import (
	paxos "dat520/lab6/gorumspaxos"
	"flag"
	"hash/fnv"
	"log"
	"net"
	"os"
	"strings"
)

func main() {
	/* get port of each server from input */
	var (
		localAddr = flag.String("laddr", "localhost:8080", "localaddress to listen on")
		saddrs    = flag.String("addrs", "", "all other remaining replica addresses separated by ','")
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
	l, err := net.Listen("tcp", *localAddr)
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()

	log.Printf("Waiting for requests at %s", l.Addr().String())
	nodeMap := make(map[string]uint32)
	addrs = append(addrs, *localAddr)
	for _, addr := range addrs {
		id := calculateHash(addr)
		nodeMap[addr] = uint32(id)
	}

	args := paxos.NewPaxosReplicaArgs{
		LocalAddr: *localAddr,
		Id:        calculateHash(*localAddr),
		NodeMap:   nodeMap,
	}
	replica := paxos.NewPaxosReplica(args)
	replica.ServerStart(l)

}

// calculateHash calculates an integer hash for the address of the node
func calculateHash(address string) int {
	h := fnv.New32a()
	h.Write([]byte(address))
	return int(h.Sum32())
}
