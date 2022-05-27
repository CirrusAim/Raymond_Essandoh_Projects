package gorumspaxos

import (
	pb "dat520/lab5/gorumspaxos/proto"
	"sort"
)

// PaxosQSpec is a quorum specification object for Paxos.
// It only holds the quorum size.
// DO NOT remove the existing fields in the structure
type PaxosQSpec struct {
	qSize int
}

// NewPaxosQSpec returns a quorum specification object for Paxos
// for the given quorum size.
func NewPaxosQSpec(quorumSize int) pb.QuorumSpec {
	return &PaxosQSpec{
		qSize: quorumSize,
	}
}

// PrepareQF is the quorum function to process the replies of the Prepare RPC call.
// Proposer handle PromiseMsgs returned by the Acceptors, and any promiseslots added
// to the PromiseMsg should be in the increasing order.
func (qs PaxosQSpec) PrepareQF(prepare *pb.PrepareMsg,
	replies map[uint32]*pb.PromiseMsg) (*pb.PromiseMsg, bool) {
	if len(replies) >= qs.qSize {

		// check after quoram, all promise msgs have same Round value.

		// map to hold the slot id vs PromiseSlot. In case multiple slots have same
		// ids, keep the PromiseSlot with higher Vrnd value
		var mp = make(map[uint32]*pb.PromiseSlot)

		for _, v := range replies {
			if v.Rnd == nil {
				continue
			}
			if v.Rnd.Id != prepare.Crnd.Id {
				return nil, false
			}

			if v == nil {
				continue
			}

			// check if there are multiple slots, pick the one which has higher vrnd value
			for _, s := range v.Slots {
				if s.Slot.Id <= prepare.Slot.Id {
					continue
				}
				v, ok := mp[s.Slot.Id]
				if ok {
					if s.Vrnd == nil && v.Vrnd == nil {
						continue
					} else if s.Vrnd == nil && v.Vrnd != nil {
						mp[s.Slot.Id] = s
					}
					if s.Vrnd.Id > v.Vrnd.Id {
						mp[s.Slot.Id] = s
					}
				} else {
					mp[s.Slot.Id] = s
				}
			}

		}

		// create missing entries and fill them with Noop value
		keys := make([]int, 0, len(mp))
		for k := range mp {
			keys = append(keys, int(k))
		}
		sort.Ints(keys)
		final := make([]uint32, 0)
		idx := 0
		if len(keys) > 0 {
			for i := keys[0]; i <= keys[len(keys)-1]; i++ {
				if keys[idx] == i {
					final = append(final, uint32(keys[idx]))
					idx++
				} else {
					final = append(final, uint32(i))
				}
			}
		}

		Slots := make([]*pb.PromiseSlot, 0)

		for _, k := range final {
			v, ok := mp[k]
			if !ok {
				Slots = append(Slots, &pb.PromiseSlot{
					Slot:  &pb.Slot{Id: k},
					Vrnd:  &pb.Round{Id: prepare.Crnd.Id},
					Value: &pb.Value{IsNoop: true},
				})

			} else {
				Slots = append(Slots, v)
			}
		}

		promise := pb.PromiseMsg{
			Rnd:   prepare.Crnd,
			Slots: Slots,
		}

		return &promise, true

	}
	return nil, false
}

// AcceptQF is the quorum function for the Accept quorum RPC call
// This is where the Proposer handle LearnMsgs to determine if a
// value has been decided by the Acceptors.
// The quorum function returns true if a value has been decided,
// and the corresponding LearnMsg holds the round number and value
// that was decided. If false is returned, no value was decided.
func (qs PaxosQSpec) AcceptQF(accMsg *pb.AcceptMsg, replies map[uint32]*pb.LearnMsg) (*pb.LearnMsg, bool) {
	if len(replies) < qs.qSize {
		return nil, false
	}

	// Process replies after sufficient replies are received. if the values are matching with the accMsg, return the Learn Msg
	// The return value of this is used in performAccept in proposer.go
	for _, v := range replies {
		if v.Rnd.Id != accMsg.Rnd.Id || v.Slot.Id != accMsg.Slot.Id {
			return nil, false
		}

		if (v.Val == nil) != (accMsg.Val == nil) {
			return nil, false
		}

		if v.Val == nil && accMsg.Val == nil {
			continue
		}

		if v.Val.ClientID != accMsg.Val.ClientID || v.Val.ClientSeq != accMsg.Val.ClientSeq || v.Val.ClientCommand != accMsg.Val.ClientCommand {
			return nil, false
		}

	}

	return &pb.LearnMsg{
		Rnd:  accMsg.Rnd,
		Slot: accMsg.Slot,
		Val:  accMsg.Val,
	}, true
}

// CommitQF is the quorum function for the Commit quorum RPC call.
// This function just waits for a quorum of empty replies,
// indicating that at least a quorum of Learners have committed
// the value decided by the Acceptors.
func (qs PaxosQSpec) CommitQF(_ *pb.LearnMsg, replies map[uint32]*pb.Empty) (*pb.Empty, bool) {
	// Check if quoram is formed. If true, return true with Empty value
	if len(replies) >= qs.qSize {
		return &pb.Empty{}, true
	}
	return nil, false
}

// ClientHandleQF is the quorum function  for the ClientHandle quorum RPC call.
// This functions waits for replies from the majority of replicas. Received replies
// should be validated before returning the response
func (qs PaxosQSpec) ClientHandleQF(in *pb.Value, replies map[uint32]*pb.Response) (*pb.Response, bool) {

	// Process reponses from the replica instances and return value to client
	if len(replies) < qs.qSize {
		return nil, false
	}

	validData := 0

	for _, v := range replies {
		if in.ClientID == v.ClientID && in.ClientSeq == v.ClientSeq && in.ClientCommand == v.ClientCommand {
			validData++
		}

	}
	if validData >= qs.qSize {
		return &pb.Response{
			ClientID:      in.ClientID,
			ClientSeq:     in.ClientSeq,
			ClientCommand: in.ClientCommand,
		}, true
	}

	return nil, false

}
