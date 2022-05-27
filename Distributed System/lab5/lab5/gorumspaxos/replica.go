package gorumspaxos

import (
	"container/list"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"sync"

	"dat520/lab3/leaderdetector"
	fd "dat520/lab5/gorumspaxos/gorumsfd"
	pb "dat520/lab5/gorumspaxos/proto"

	"github.com/relab/gorums"
)

// constants used in implementation
const (
	// NoRoundId indicates invalid value for round
	NoRoundId int32 = -1
	// Ignore indicates invalid slotID
	Ignore uint32 = 0
	// waitTimeForPhaseOne, is the time which the phase1 quorum functions wait for reply
	waitTimeForPhaseOne time.Duration = 5 * time.Second
	// waitTimeForPhaseTwo, is the time which the phase2 quorum functions wait for reply
	waitTimeForPhaseTwo time.Duration = 5 * time.Second
	// failureDetectordDelay is the default value of the delay parameter in failure detector
	failureDetectordDelay time.Duration = 5 * time.Second
	// failureDetectordDelta is the default value of the delta parameter in failure detector
	failureDetectordDelta time.Duration = 1 * time.Second
)

// PaxosReplica is the structure composing the Proposer and Acceptor.
// failureDetector implementation of the failure detector
// DO NOT remove the existing fields in the structure
type PaxosReplica struct {
	pb.MultiPaxos
	sync.Mutex
	*Acceptor
	*Proposer
	failureDetector fd.FailureDetector
	id              int        // id is the id of the node
	localAddr       string     // localAddr is the local address of the replica
	responseList    *list.List // responseChan is the channel used by replica to deliver the response to proposer
	timings         map[string]float64
}

// NewPaxosReplica returns a new Paxos replica with a configuration as provided
// by the input addrs. This replica will run on the given port.
func NewPaxosReplica(args NewPaxosReplicaArgs) *PaxosReplica {
	acceptor := NewAcceptor()
	quorumSize := (len(args.NodeMap)-1)/2 + 1
	qspec := NewPaxosQSpec(quorumSize)
	nodeIds := make([]int, 0)
	for _, id := range args.NodeMap {
		nodeIds = append(nodeIds, int(id))
	}
	ld := leaderdetector.NewMonLeaderDetector(nodeIds)
	proposerArgs := NewProposerArgs{
		id:               args.Id,
		aduSlotID:        Ignore,
		leaderDetector:   ld,
		qspec:            qspec,
		nodeMap:          args.NodeMap,
		phaseOneWaitTime: waitTimeForPhaseOne,
		phaseTwoWaitTime: waitTimeForPhaseTwo,
	}
	proposer := NewProposer(proposerArgs)
	failureDetector := fd.NewEvtFailureDetector(args.Id, ld, args.NodeMap,
		failureDetectordDelay, failureDetectordDelta)
	responseList := list.New()
	return &PaxosReplica{
		localAddr:       args.LocalAddr,
		Acceptor:        acceptor,
		Proposer:        proposer,
		failureDetector: failureDetector,
		id:              args.Id,
		responseList:    responseList,
		timings:         make(map[string]float64),
	}
}

// ServerStart starts the replica
// 1. Invokes the start function of the proposer
// 2. Create a new gorums server
// 3. Register MultiPaxos server
// 4. Start failure detector
// 5. Call Serve on gorums server
func (replica *PaxosReplica) ServerStart(lis net.Listener) {
	replica.Start()
	gorumSrv := gorums.NewServer()
	pb.RegisterMultiPaxosServer(gorumSrv, replica)
	replica.failureDetector.StartFailureDetector(gorumSrv)
	err := gorumSrv.Serve(lis)
	if err != nil {
		log.Printf("Error serving connections : %s", err)
		gorumSrv.GracefulStop()
	}

}

// Prepare handles the prepare quorum calls from the proposer by passing the received messages to its acceptor.
// It receives prepare massages and pass them to handlePrepare method of acceptor.
// It returns promise messages back to the proposer by its acceptor.
func (r *PaxosReplica) Prepare(ctx gorums.ServerCtx, prepMsg *pb.PrepareMsg) (*pb.PromiseMsg, error) {
	//log.Printf("Node id %d \t Acceptor: Prepare(%v) received", r.id, prepMsg)
	prm := r.handlePrepare(prepMsg)
	return prm, nil
}

// Accept handles the accept quorum calls from the proposer by passing the received messages to its acceptor.
// It receives Accept massages and pass them to handleAccept method of acceptor.
// It returns learn massages back to the proposer by its acceptor
func (r *PaxosReplica) Accept(ctx gorums.ServerCtx, accMsg *pb.AcceptMsg) (*pb.LearnMsg, error) {
	//log.Printf("Node id %d \t Acceptor: Accept(%v) received", r.id, accMsg)
	lrn := r.handleAccept(accMsg)
	return lrn, nil
}

// Commit is invoked when the proposer calls the commit RPC on the configuration.
// It receives a learn massage from proposer, this means the request is decided and
// the replica can commit the request. It returns an empty massage back.
func (r *PaxosReplica) Commit(ctx gorums.ServerCtx, lrnMsg *pb.LearnMsg) (*pb.Empty, error) {
	//log.Printf("Node id %d \t Learner:cCommit(%v) received", r.id, lrnMsg)
	r.IncrementAllDecidedUpTo()
	r.Lock()
	r.responseList.PushBack(lrnMsg.Val)
	r.Unlock()
	return &pb.Empty{}, nil
}

// ClientHandle method is invoked when the client calls the ClientHandle RPC on all the replicas.
// This method may be called multiple times from different clients.
// If the client request is committed on the replica then it should send the response.
// P.S. Since the method is called by multiple clients, do remember to return the matching reply to
// the client. Example If Client C send the replica the request M1 then return to the client when M1 is decided.
// While waiting for M1 to get committed, M2 may be proposed and committed by the replicas.
// getResponse method helps you to match the request to the response.
func (r *PaxosReplica) ClientHandle(ctx gorums.ServerCtx, req *pb.Value) (rsp *pb.Response, err error) {
	//log.Printf("Node id %d\t Replica: ClientHandle(%v) received", r.id, req)
	start := time.Now()
	r.AddRequestToQ(req)
	resp, err := r.getResponse(req)
	if err == nil {
		end := time.Since(start).Seconds()
		r.Lock()
		key := req.ClientID + "_" + strconv.Itoa(int(req.ClientSeq)) + "_" + req.ClientCommand
		r.timings[key] = float64(end)
		r.Unlock()
	}
	return resp, err

}

// getResponse: is called after adding the client request to the queue.
// This function waits on the responseList for the response
// matching the request. This involves creating a new go routine and making it
// repeatedly check the responseList for the matched response. If a response
// is not present within the phaseTwoWaitTime, clean up the created goroutine
// and return error.
// P.S. Channels can also be used to achieve this functionality, if required you have complete
// freedom to change the definition of this method.
func (r *PaxosReplica) getResponse(request *pb.Value) (rsp *pb.Response, err error) {

	// This function checks if pb.Value is present in responseList. This checks continuously for value in the
	// list and finally timeout after waitTimeForPhaseTwo duration. if the value is found return the Reponse,else
	//return nil and false

	for {
		select {
		case <-time.After(waitTimeForPhaseTwo):
			return nil, errors.New("no response found")
		default:
			r.Lock()
			for msg := r.responseList.Front(); msg != nil; msg = msg.Next() {
				val, _ := msg.Value.(*pb.Value)

				if val == nil && request == nil {
					return &pb.Response{}, nil
				}
				// Comparing value with the incoming request.
				if val != nil && request != nil &&
					val.ClientSeq == request.ClientSeq &&
					val.ClientID == request.ClientID &&
					val.ClientCommand == request.ClientCommand &&
					val.IsNoop == request.IsNoop {
					r.Unlock()
					return &pb.Response{
						ClientID:      val.ClientID,
						ClientSeq:     val.ClientSeq,
						ClientCommand: val.ClientCommand,
					}, nil
				}

			}
			r.Unlock()
		}
	}
}

func (r *PaxosReplica) GetThroughPut() {
	//data := make([]float64, 0)
	var tp float64
	for _, v := range r.timings {
		tp += v
	}

	tp = float64(len(r.timings)) / tp
	//mean, _ := st.Mean(data)
	//median, _ := st.Median(data)
	//sd, _ := st.StandardDeviation(data)
	//max, _ := st.Max(data)
	//min, _ := st.Min(data)

	logfile, err := os.Create(fmt.Sprintf("/Users/hysensw/Videos/lab5/lab5/lab5/gorumspaxos/cmd/paxosclient/server_%s.json", strconv.Itoa(r.id)))
	if err != nil {
		panic(err)
	}
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(logfile)
	mp := make(map[string][]float64)
	mp["data"] = []float64{tp}
	//mp["mean"] = mean
	//mp["median"] = median
	//mp["sd"] = sd
	//mp["max"] = max
	//mp["min"] = min
	bt, err := json.Marshal(mp)

	log.Printf("%v", string(bt))
	log.SetOutput(os.Stdout)
}
