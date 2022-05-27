package gorumsfd

import (
	"context"
	"time"

	pb "dat520/lab5/gorumspaxos/fdproto"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	ManagerWaitTimeOut = 5 * time.Second
)

// EvtFailureDetector represents a Eventually Perfect Failure Detector as
// described at page 53 in:
// Christian Cachin, Rachid Guerraoui, and Luís Rodrigues: "Introduction to
// Reliable and Secure Distributed Programming" Springer, 2nd edition, 2011.

// DO NOT remove the existing fields in the structure
type EvtFailureDetector struct {
	id            int               // the id of this node
	alive         map[uint32]bool   // map of node ids considered alive
	suspected     map[uint32]bool   // map of node ids  considered suspected
	sr            SuspectRestorer   // Provided SuspectRestorer implementation
	delay         time.Duration     // the current delay for the timeout procedure
	delta         time.Duration     // the delta value to be used when increasing delay
	stop          chan struct{}     // channel for signaling a stop request to the main run loop
	nodeMap       map[string]uint32 // list of addresses of the nodes. {address localhost:port => hash(localhost:port}
	manager       *pb.Manager       // Manager to create the configuration
	configuration *pb.Configuration // Configuration to call the quorum calls
}

// NewEvtFailureDetector returns a new Eventual Failure Detector. It takes the
// following arguments:
// id: The id of the node running this instance of the failure detector.
// running this instance of the failure detector).
// sr: A leader detector implementing the SuspectRestorer interface.
// addrs: A list of address of the replicas.
// delay: The timeout delay after which the failure detection is performed.
// delta: The value to be used when increasing delay.
func NewEvtFailureDetector(id int, sr SuspectRestorer, nodeMap map[string]uint32,
	delay time.Duration, delta time.Duration) *EvtFailureDetector {
	suspected := make(map[uint32]bool)
	alive := make(map[uint32]bool)
	for _, id := range nodeMap {
		alive[id] = false
	}
	mgr := pb.NewManager(gorums.WithDialTimeout(ManagerWaitTimeOut),
		gorums.WithGrpcDialOptions(
			grpc.WithBlock(),                                         // block until connections are made
			grpc.WithTransportCredentials(insecure.NewCredentials()), // disable TLS
		),
	)
	return &EvtFailureDetector{
		id:        id,
		alive:     alive,
		suspected: suspected,
		sr:        sr,
		nodeMap:   nodeMap,
		delay:     delay,
		delta:     delta,
		manager:   mgr,
		stop:      make(chan struct{}),
	}
}

// StartFailureDetector starts main run loop in a separate goroutine.
// This function should perform the following functionalities
// 1. Register FailureDetectorServer implementation
// 2. The started Go Routine, after the e.delay, it should call PerformFailureDetection
// 3. Started Go Routine, should also wait to receive the signal to stop
func (e *EvtFailureDetector) StartFailureDetector(srv *gorums.Server) error {
	// TODO(student) complete StartFailureDetector
	pb.RegisterFailureDetectorServer(srv, newEvtFailureDetectorServer(e))

	go func() {
		for {
			select {
			case <-e.stop:
				return
			case <-time.After(e.delay):
				e.PerformFailureDetection()
			}
		}
	}()
	return nil
}

// PerformFailureDetection is the method used to perform ping
// operation for all nodes and report all suspected nodes.
// 1. Create configuration with all the nodes if not previously done.
// 2. call Ping rpc on the configuration
// 3. call SendStatusOfNodes to send suspect and restore notifications.
func (e *EvtFailureDetector) PerformFailureDetection() error {
	// TODO(student) complete PerformFailureDetection
	if e.configuration == nil {
		var err error
		e.configuration, err = e.manager.NewConfiguration(gorums.WithNodeMap(e.nodeMap))
		if err != nil {
			return err
		}
	}

	hb, err := newEvtFailureDetectorServer(e).Ping(gorums.ServerCtx{
		Context: context.TODO(),
	}, &pb.HeartBeat{})

	if err != nil {
		return err
	}

	// Use configuration object to call RPC methods
	e.configuration.Ping(context.Background(), hb)

	// Send status of nodes to Leader detector. Leader detector we're using from Lab3
	// It uses same algorithm as used in lab3 failure detector
	e.SendStatusOfNodes()
	return nil
}

// Stop stops e's main run loop.
func (e *EvtFailureDetector) Stop() {
	e.stop <- struct{}{}
}

// SendStatusOfNodes: reports the status of the nodes to the SuspectRestorer. This
// method is called after Ping RPC which marks all the live nodes.
// If a node which is previously suspected and now is alive then increase the e.delay by e.delta
// All non reachable nodes are reported and all previously reported,
// now live nodes are restored by calling the Suspect and Restore functions of the SuspectRestorer.
func (e *EvtFailureDetector) SendStatusOfNodes() {
	// TODO(student) complete PerformFailureDetection
	if e.isDelayIncreaseRequired() {
		e.delay += e.delta
	}
	for _, node := range e.nodeMap {
		vAlive, _ := e.alive[node]
		vSuspected, _ := e.suspected[node]
		if !vAlive && !vSuspected {
			e.suspected[node] = true
			e.sr.Suspect(int(node))
		} else if vAlive && vSuspected {
			delete(e.suspected, node)
			e.sr.Restore(int(node))
		}
	}
	//Set alive to empty
	for node := range e.alive {
		delete(e.alive, node)
	}
}

// evtFailureDetectorServer implements the FailureDetector RPC
type evtFailureDetectorServer struct {
	e *EvtFailureDetector
}

// newEvtFailureDetectorServer returns the evtFailureDetectorServer
func newEvtFailureDetectorServer(e *EvtFailureDetector) evtFailureDetectorServer {
	return evtFailureDetectorServer{e}
}

// Ping handles the Ping RPC from the other replicas. Reply contains the id of the node
func (srv evtFailureDetectorServer) Ping(ctx gorums.ServerCtx, in *pb.HeartBeat) (resp *pb.HeartBeat, err error) {
	resp = &pb.HeartBeat{Id: int32(srv.e.id)}
	return resp, err
}

// evtFailureDetectorQSpec implements the QuorumSpec for the RPC
type evtFailureDetectorQSpec struct {
	e *EvtFailureDetector
}

// newEvtFailureDetectorQSpec returns the evtFailureDetectorQSpec
func newEvtFailureDetectorQSpec(e *EvtFailureDetector) evtFailureDetectorQSpec {
	return evtFailureDetectorQSpec{e}
}

// PingQF is the quorum function to handle the replies to Ping RPC call. Nodes replied to the call
// are marked live.
// Quoram will Poll on PingQF. replies will contain the previous results as well.
// 1,2,3,4,5 nodes there. first poll result -> 2,3 . so replies will contain 2,3
// next poll fetches 4.  So the replies will contain 2,3,4
// that how the condition in 196 line will result to true eventually. and at that point quoram will no longer poll.
// This will not poll indefinately. there is timeout.
func (q evtFailureDetectorQSpec) PingQF(in *pb.HeartBeat, replies map[uint32]*pb.HeartBeat) (*pb.HeartBeat, bool) {
	// TODO(student) complete PerformFailureDetection
	// It mark
	for k := range replies {
		q.e.alive[k] = true
	}
	if len(replies) == len(q.e.nodeMap) {
		return nil, true
	}
	return nil, false
}

// returns true when a node was suspected before but now marked as alive by PingQF
func (e *EvtFailureDetector) isDelayIncreaseRequired() bool {
	if len(e.alive) > len(e.suspected) {
		for k := range e.alive {
			if v, ok := e.suspected[k]; ok && v {
				return true
			}
		}
	} else {
		for k := range e.suspected {
			if v, ok := e.alive[k]; ok && v {
				return true
			}
		}
	}
	return false
}
