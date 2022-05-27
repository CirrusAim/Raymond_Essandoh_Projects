package failuredetector

import (
	"time"
)

// EvtFailureDetector represents a Eventually Perfect Failure Detector as
// described at page 53 in:
// Christian Cachin, Rachid Guerraoui, and Luís Rodrigues: "Introduction to
// Reliable and Secure Distributed Programming" Springer, 2nd edition, 2011.
type EvtFailureDetector struct {
	id        int          // the id of this node
	nodeIDs   []int        // node ids for every node in cluster
	alive     map[int]bool // map of node ids considered alive
	suspected map[int]bool // map of node ids  considered suspected

	sr SuspectRestorer // Provided SuspectRestorer implementation

	delay         time.Duration // the current delay for the timeout procedure
	delta         time.Duration // the delta value to be used when increasing delay
	timeoutSignal *time.Ticker  // the timeout procedure ticker

	hbSend chan<- Heartbeat // channel for sending outgoing heartbeat messages
	hbIn   chan Heartbeat   // channel for receiving incoming heartbeat messages
	stop   chan struct{}    // channel for signaling a stop request to the main run loop

	testingHook func() // DO NOT REMOVE THIS LINE. A no-op when not testing.
}

// NewEvtFailureDetector returns a new Eventual Failure Detector. It takes the
// following arguments:
//
// id: The id of the node running this instance of the failure detector.
//
// nodeIDs: A list of ids for every node in the cluster (including the node
// running this instance of the failure detector).
//
// ld: A leader detector implementing the SuspectRestorer interface.
//
// delta: The initial value for the timeout interval. Also the value to be used
// when increasing delay.
//
// hbSend: A send only channel used to send heartbeats to other nodes.
func NewEvtFailureDetector(id int, nodeIDs []int, sr SuspectRestorer, delta time.Duration, hbSend chan<- Heartbeat) *EvtFailureDetector {
	suspected := make(map[int]bool)
	alive := make(map[int]bool)

	for x := range nodeIDs {
		alive[x] = true
	}

	return &EvtFailureDetector{
		id:        id,
		nodeIDs:   nodeIDs,
		alive:     alive,
		suspected: suspected,

		sr: sr,

		delay: delta,
		delta: delta,

		hbSend: hbSend,
		hbIn:   make(chan Heartbeat, 8),
		stop:   make(chan struct{}),

		testingHook: func() {}, // DO NOT REMOVE THIS LINE. A no-op when not testing.
	}
}

// Start starts e's main run loop as a separate goroutine. The main run loop
// handles incoming heartbeat requests and responses. The loop also trigger e's
// timeout procedure at an interval corresponding to e's internal delay
// duration variable.
func (e *EvtFailureDetector) Start() {
	e.timeoutSignal = time.NewTicker(e.delay)
	go func(timeout <-chan time.Time) {
		for {
			e.testingHook() // DO NOT REMOVE THIS LINE. A no-op when not testing.
			select {
			case hb := <-e.hbIn:
				//  true => request, false => reply
				if hb.Request {
					hbResp := Heartbeat{
						From:    hb.To,
						To:      hb.From,
						Request: false,
					}
					e.hbSend <- hbResp
				} else {
					e.alive[hb.From] = true
				}

			case <-timeout:
				e.timeout()
			case <-e.stop:
				return
			}
		}
	}(e.timeoutSignal.C)
}

// DeliverHeartbeat delivers heartbeat hb to failure detector e.
func (e *EvtFailureDetector) DeliverHeartbeat(hb Heartbeat) {
	e.hbIn <- hb
}

// Stop stops e's main run loop.
func (e *EvtFailureDetector) Stop() {
	e.stop <- struct{}{}
}

// Internal: timeout runs e's timeout procedure.
func (e *EvtFailureDetector) timeout() {
	if e.isDelayIncreaseRequired() {
		e.delay += e.delta
	}
	for _, node := range e.nodeIDs {
		_, okAlive := e.alive[node]
		_, okSuspected := e.suspected[node]
		if !okAlive && !okSuspected {
			e.suspected[node] = true
			e.sr.Suspect(node)
		} else if okAlive && okSuspected {
			//e.suspected[node] = false
			delete(e.suspected, node)
			e.sr.Restore(node)
		}
		e.hbSend <- Heartbeat{
			From:    e.id,
			To:      node,
			Request: true,
		}
	}
	//Set alive to empty
	for node := range e.alive {
		delete(e.alive, node)
	}
}

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
