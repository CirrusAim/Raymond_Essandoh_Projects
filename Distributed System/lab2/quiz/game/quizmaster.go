package game

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"sync"

	"dat520/lab2/quiz"
	pb "dat520/lab2/quiz/proto"
	"dat520/lab2/quiz/utils"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Quizmaster structure holds the data to conduct quiz
// participants map stores the nodeID and the participant
// serverPort stores the port on which the server is running
// correctAnswers map stores the question id and the correct answer for that question
// configuration is the nodes to send the question and can change for each question
type quizmaster struct {
	sync.Mutex
	pb.UnimplementedParticipantServiceServer
	mgr               *pb.Manager
	participants      map[uint32]*pb.Participant
	serverPort        int
	configuration     *pb.Configuration
	readyParticipants map[uint32]bool
	answerReceived    map[uint32][]int32
}

// NewQuizMaster initializes the quizmaster structure
func NewQuizMaster(port int) *quizmaster {
	participants := make(map[uint32]*pb.Participant)

	// init gorums manager
	mgr := pb.NewManager(
		gorums.WithDialTimeout(quiz.QuizmasterParticipantConnect),
		gorums.WithGrpcDialOptions(
			grpc.WithInsecure(), // disable TLS
			grpc.WithBlock(),    // block until connections are made
		),
	)
	configuration, _ := mgr.NewConfiguration()
	// TODO(student) Add any other necessary information to the quizmaster
	return &quizmaster{
		mgr:               mgr,
		serverPort:        port,
		participants:      participants,
		readyParticipants: make(map[uint32]bool),
		configuration:     configuration,
		answerReceived:    make(map[uint32][]int32),
	}
}

// Register handles the Register RPC called from the participants.
// 1. Check the contents of the pb.Participant sent from the participants.
// 2. If the details are incorrect, send error
// 3. If all the details are correct, then store the participant
// in the quizmaster structure.
func (q *quizmaster) Register(_ context.Context, p *pb.Participant) (*pb.RegisterResponse, error) {
	if p.Name == "" && p.Score != 0 {
		return &pb.RegisterResponse{}, errors.New("incorrect participant details")
	}
	var newNodeConfig *pb.Configuration
	var err error
	var nodeId uint32
	q.Lock()
	defer q.Unlock()
	mgr := pb.NewManager(
		gorums.WithDialTimeout(quiz.QuizmasterParticipantConnect),
		gorums.WithGrpcDialOptions(
			grpc.WithInsecure(), // disable TLS
			grpc.WithBlock(),    // block until connections are made
		),
	)

	// If it is the first participant, configure with only that node, else
	// append the new participant's address to already present addresses to make
	//  new configuration
	if q.configuration == nil {
		newNodeConfig, err = q.mgr.NewConfiguration(
			gorums.WithNodeList([]string{p.Address}),
			q,
		)
	} else {
		newNodeConfig, err = q.mgr.NewConfiguration(
			q.configuration.WithNewNodes(gorums.WithNodeList([]string{p.Address})),
			q,
		)
	}

	if err != nil {
		fmt.Printf("Unable to register . Error : %s\n", err)
		return &pb.RegisterResponse{}, err
	}
	q.mgr = mgr
	q.configuration = newNodeConfig

	// find the nodeID corresponding to the address of the participant and store in
	// participant map.
	// Participant is registered but is not ready.
	allNodes := q.configuration.Nodes()
	for _, node := range allNodes {
		if node.Address() == p.Address {
			nodeId = node.ID()
			q.participants[nodeId] = p
			q.readyParticipants[nodeId] = false
			q.answerReceived[nodeId] = []int32{}
			break
		}
	}

	return &pb.RegisterResponse{
		NodeId: nodeId,
	}, nil
}

// MarkReady handles the "MarkReady" RPC from the participants,
// This marks the participants are ready to participate in the
// Quiz.
// Caller Participant is included in the configuration of the next question.
func (q *quizmaster) MarkReady(ct context.Context,
	resp *pb.RegisterResponse) (*emptypb.Empty, error) {
	// Mark the participant ready
	q.Lock()
	q.readyParticipants[resp.NodeId] = true
	q.Unlock()
	return &emptypb.Empty{}, nil
}

// GetResults handles the "GetResults" RPC from the participants.
// pb.Result should contain the score of all the registered participants.
// If no participants are available return error.
func (q *quizmaster) GetResults(ct context.Context, e *emptypb.Empty) (res *pb.Result, err error) {
	// It will show the result of all participants even if they got disconnected before the
	// quiz completion or came after quiz started
	var par []*pb.Participant
	for _, v := range q.participants {
		par = append(par, &pb.Participant{
			Name:    v.Name,
			Address: v.Address,
			Score:   v.Score,
		})
	}
	return &pb.Result{Participants: par}, nil
}

// AnswerQF is the quorum function for the "Answer" RPC. It handles the
// responses for the question sent from the quizmaster.
// 1. It checks the answer for the question, if the answer is correct then
// updates the score of the participant.
// 2. Once all the responses are received, then return true,
// otherwise return false to further process the replies.
// P.S. This function is called multiple times with "replies" containing both
// old and new responses from the participants. While updating the scores,
// one should take care of not processing the reply twice.
func (q *quizmaster) AnswerQF(in *pb.Question, replies map[uint32]*pb.ParticipantAnswer) (*pb.ParticipantAnswer, bool) {
	// answerReceived
	alreadyAnswered := func(nodeId uint32, qId int32) bool {
		v := q.answerReceived[nodeId]
		for _, id := range v {
			if id == qId {
				return true
			}
		}
		return false
	}

	liveParticipants := q.getLiveParticipants()
	q.Lock()
	defer q.Unlock()

	// find correct answer of the question.
	var correctAns int32
	for _, q := range quiz.Questions {
		if q.Id == in.Id {
			correctAns = q.CorrectAnswer + 1
		}
	}

	// calculate result for every participant's answer.
	for id, rep := range replies {
		if alreadyAnswered(id, in.Id) {
			continue
		}
		ans := rep.Answer
		if rep.QuestionId == in.Id && ans == correctAns {
			q.participants[id].Score += quiz.CorrectAnswerScore
			q.answerReceived[id] = append(q.answerReceived[id], in.Id)
		}
	}
	if len(liveParticipants) == len(replies) {
		return &pb.ParticipantAnswer{}, true
	} else {
		return &pb.ParticipantAnswer{}, false
	}

}

// sendQuestions sends the questions from the questions variable from
// the questions.go.
// 1. Before sending the question, a new gorums client should be formed
// with the configuration of ready and live participants.
// A dead participant should not be included in the configuration.
// 2. It should call the "Answer" RPC on the participants
// with the question and wait for QuizmasterQuestionTimeout seconds.
// 3. After receiving the responses for the question,
// it should take user input from the user to "send next question" or "show results".
// This function returns after sending all questions.
func (q *quizmaster) sendQuestions() {
	lastIdx := len(quiz.Questions)
	isLastQ := false
	for i, v := range quiz.Questions {
		q.runGorumsClientForEachQuestion()
		if i == lastIdx-1 {
			isLastQ = true
		}
		question := &pb.Question{
			Id:             v.Id,
			QuestionText:   v.QuestionText,
			AnswerText:     v.AnswerText,
			CorrectAnswer:  -1,
			IsLastQuestion: isLastQ,
		}
		// wait some time before proceeding to next question
		ctx, cancel := context.WithTimeout(context.Background(), quiz.QuizmasterQuestionTimeout)
		defer cancel()
		q.configuration.Answer(ctx, question)

	}

}

// Checks whether the participant is still live by dialing
func isLive(address string) bool {
	_, err := net.DialTimeout("tcp", address, quiz.QuizmasterParticipantConnect)
	return err == nil
}

// getLiveParticipants returns the list of addresses of ready and live participants
// This function should check if a registered participant is ready a alive.
// You can use the isLive function to perform the latter.
// It should remove participants that are not alive.
func (q *quizmaster) getLiveParticipants() (addresses []string) {
	// TODO(student) Implement this function
	q.Lock()
	defer q.Unlock()
	for k, v := range q.participants {
		isReady, ok := q.readyParticipants[k]
		if isLive(v.Address) && ok && isReady {
			addresses = append(addresses, v.Address)
		}
	}
	return addresses
}

// runGorumsClientForEachQuestion creates a gorums manager
// and create a configuration with ready and live participants (you can use the getLiveParticipants)
// store this configuration object in quizmaster to use in Vote RPC.
func (q *quizmaster) runGorumsClientForEachQuestion() {
	// TODO(student) Implement this function
	addresses := q.getLiveParticipants()
	mgr := pb.NewManager(
		gorums.WithDialTimeout(quiz.QuizmasterParticipantConnect),
		gorums.WithGrpcDialOptions(
			grpc.WithInsecure(), // disable TLS
			grpc.WithBlock(),    // block until connections are made
		),
	)
	allNodesConfig, err := mgr.NewConfiguration(
		gorums.WithNodeList(addresses),
		q,
	)

	if err != nil {
		fmt.Printf("Error : %s\n", err)
	}
	q.configuration = allNodesConfig
}

// RunGrpcServer starts the GRPC server to receive the RPC from participants.
// 1. Starts a GRPC server with the user input port
// 2. register the quizmaster as the handler for the RPC of the participants.
func (q *quizmaster) RunGrpcServer() {
	// TODO(student) Implement this function
	listener, err := net.Listen("tcp", "localhost:"+fmt.Sprintf("%d", q.serverPort))
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}

	grpcServer := grpc.NewServer()
	pb.RegisterParticipantServiceServer(grpcServer, q)

	err = grpcServer.Serve(listener)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}

// Print the results of the participants
func (q *quizmaster) showResults() {
	q.Lock()
	defer q.Unlock()
	for _, participant := range q.participants {
		log.Printf("Participant %s score is %d\n\n", participant.Name, participant.Score)
	}
}

// startQuiz: Take the user input to start the quiz
// Wait until the participants are ready and live.
// If at least one live participant is ready then start
// the quiz by sending the questions.
// After sending all the questions, show options of
// "ShowResults" and "Exit".
// If "ShowResults" is select print the results and exit
// if the "Exit" option is selected.
func (q *quizmaster) StartQuiz() {
	fmt.Println("Press Enter to start the Quiz")
	if !utils.ReadEnter() {
		log.Fatal("unable to read input")
		return
	}

	// Wait until there is atleast one live participant
	for {
		if len(q.getLiveParticipants()) > 0 {
			break
		}
	}

	q.sendQuestions()
	for {
		expectedCommands := []string{"ShowResults", "Exit"}
		option := utils.ReadCommand(expectedCommands)
		if option == 0 {
			q.showResults()
		} else {
			q.reset()
			return
		}
	}
}

func (q *quizmaster) reset() {
	q.participants = make(map[uint32]*pb.Participant)
}
