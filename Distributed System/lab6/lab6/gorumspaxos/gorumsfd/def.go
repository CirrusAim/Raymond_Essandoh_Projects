package gorumsfd

import (
	pb "dat520/lab6/gorumspaxos/fdproto"

	"github.com/relab/gorums"
)

// FailureDetector is the interface implemented by a failure detector.
type FailureDetector interface {
	StartFailureDetector(*gorums.Server) error
	Stop()
}

// Suspecter is the interface that wraps the Suspect method. Suspect indicates
// that the node with identifier id should be considered suspected.
type Suspecter interface {
	Suspect(id int)
}

// Restorer is the interface that wraps the Restore method. Restore indicates
// that the node with identifier id should be considered restored.
type Restorer interface {
	Restore(id int)
}

// SuspectRestorer is the interface that groups the Suspect and Restore
// methods.
type SuspectRestorer interface {
	Suspecter
	Restorer
}

// struct for testing purposes, DO NOT change
type FDPerformFailureDetectionTest struct {
	replies        []*pb.HeartBeat
	suspectedNodes []int
	restoreNodes   []int
	wantOutput     bool
	delay          int
	description    string
}
