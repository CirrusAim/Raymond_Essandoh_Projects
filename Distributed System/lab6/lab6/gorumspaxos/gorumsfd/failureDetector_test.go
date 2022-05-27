package gorumsfd

import (
	pb "dat520/lab6/gorumspaxos/fdproto"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/relab/gorums"
)

type mockLD struct {
	suspectCalled int
	suspectNodes  []int
	restoreCalled int
	restoreNodes  []int
}

func (mld *mockLD) Suspect(id int) {
	mld.suspectCalled++
	mld.suspectNodes = append(mld.suspectNodes, id)
}

func (mld *mockLD) Restore(id int) {
	mld.restoreCalled++
	mld.restoreNodes = append(mld.restoreNodes, id)
}

func (mld *mockLD) Subscribe() <-chan int {
	ch := make(chan int)
	return ch
}

func (mld *mockLD) Leader() int {
	return 0
}

type MockConfiguration struct {
	pb.Configuration
}

func (q *MockConfiguration) Ping(ctx gorums.ServerCtx, in *pb.HeartBeat) (resp *pb.HeartBeat, err error) {
	return nil, nil
}

var FDPFDTests = []FDPerformFailureDetectionTest{
	{
		replies: []*pb.HeartBeat{
			{
				Id: 1,
			},
			{
				Id: 2,
			},
		},
		suspectedNodes: []int{3},
		restoreNodes:   []int{},
		wantOutput:     false,
		delay:          1,
		description:    "insufficient replies with some suspect nodes",
	},
	{
		replies: []*pb.HeartBeat{
			{
				Id: 1,
			},
			{
				Id: 3,
			},
		},
		suspectedNodes: []int{2},
		restoreNodes:   []int{3},
		wantOutput:     false,
		delay:          2,
		description:    "insufficient replies with some suspect nodes and restore nodes",
	},
	{
		replies: []*pb.HeartBeat{
			{
				Id: 1,
			},
			{
				Id: 2,
			},
			{
				Id: 3,
			},
		},
		suspectedNodes: []int{},
		restoreNodes:   []int{2},
		wantOutput:     true,
		delay:          3,
		description:    "sufficient replies with no suspect nodes and some restore nodes",
	},
	{
		replies: []*pb.HeartBeat{
			{
				Id: 1,
			},
			{
				Id: 2,
			},
			{
				Id: 3,
			},
		},
		suspectedNodes: []int{},
		restoreNodes:   []int{},
		wantOutput:     true,
		delay:          3,
		description:    "sufficient replies with no suspect nodes and no restore nodes",
	}, {
		replies: []*pb.HeartBeat{
			{
				Id: 1,
			},
			{
				Id: 3,
			},
		},
		suspectedNodes: []int{2},
		restoreNodes:   []int{},
		wantOutput:     false,
		delay:          3,
		description:    "insufficient replies with some suspect nodes from restore nodes",
	},
}

func TestFDPerformFailureDetection(t *testing.T) {
	mockSR := mockLD{}
	e := NewEvtFailureDetector(1, &mockSR, map[string]uint32{"1": 1, "2": 2, "3": 3},
		1*time.Second, 1*time.Second)
	for _, testCase := range FDPFDTests {
		mockSR.suspectNodes = make([]int, 0)
		mockSR.restoreNodes = make([]int, 0)
		replies := make(map[uint32]*pb.HeartBeat)
		for _, reply := range testCase.replies {
			replies[uint32(reply.Id)] = reply
		}
		resp, gotOutput := newEvtFailureDetectorQSpec(e).PingQF(&pb.HeartBeat{}, replies)
		e.SendStatusOfNodes()
		if testCase.wantOutput != gotOutput {
			t.Errorf("\nPingQF\ntest: %s\ngot %v\n", testCase.description, resp)
		}
		if diff := cmp.Diff(testCase.restoreNodes, mockSR.restoreNodes); diff != "" {
			t.Errorf("\nPerformFailureDetection: restoreNodes mismatch\ntest: %s\n mismatch (-want +got):\n%s",
				testCase.description, diff)
		}
		if diff := cmp.Diff(testCase.suspectedNodes, mockSR.suspectNodes); diff != "" {
			t.Errorf("\nPerformFailureDetection: suspectedNodes mismatch\ntest: %s\n mismatch (-want +got):\n%s",
				testCase.description, diff)
		}
		if int(time.Second)*testCase.delay != int(e.delay) {
			t.Errorf("\nPerformFailureDetection: mismatch delay\ntest: %s\n mismatch (-want  %d +got %d):\n",
				testCase.description, int(time.Second)*testCase.delay, int(e.delay))
		}
		mockSR = mockLD{}
	}
}
