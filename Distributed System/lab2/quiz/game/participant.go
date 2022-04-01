package game

import (
	"context"
	"dat520/lab2/quiz"
	pb "dat520/lab2/quiz/proto"
	"dat520/lab2/quiz/utils"
	"fmt"
	"net"
	"os"

	"github.com/relab/gorums"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/protobuf/types/known/emptypb"
)

// Participant structure contains the details of the participant
// name contains the name of the participant
// address is the gorums server address started by the participant,
// quizmaster sends the questions to this server.
// quizDone is used to communicate the completion of the quiz.
// nodeId is received as response to the register and passed as
// input to the MarkReady RPC.
type participant struct {
	name             string
	address          string
	quizDone         chan bool
	nodeId           uint32
	inputChannel     chan int32
	clientConnection *grpc.ClientConn
	client           pb.ParticipantServiceClient
}

// NewParticipant initializes the participant
func NewParticipant(name string) *participant {
	quizDone := make(chan bool)
	inputChannel := make(chan int32)
	return &participant{name: name, quizDone: quizDone, inputChannel: inputChannel}
}

// RegisterParticipant is the first operation of the participant,
// This function perform the following operations.
// 1) Start a gorums server to receive the questions.
// 2) Create a connection to the quizmaster on the serverPort and
//    create a client with this connection.
// 3) Send this server address to the quizmaster with Register RPC.
// 4) Store the node id returned in the register RPC response.
func (p *participant) RegisterParticipant(serverPort int) {

	// Step 1 : Start a gorum server and listen on any available port.
	gorumSrv := gorums.NewServer()
	lis, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		fmt.Printf("Failed to listen on port: %v. Cannot start participant", err)
		os.Exit(1)
	}
	p.address = lis.Addr().String()
	go func() {
		pb.RegisterQuizMasterServer(gorumSrv, p)
		err = gorumSrv.Serve(lis)
		if err != nil {
			fmt.Printf("Error listening gorum server on participant")
		}
	}()

	var opts []grpc.DialOption
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	conn, err := grpc.Dial("0.0.0.0:"+fmt.Sprintf("%d", serverPort), opts...)
	if err != nil {
		fmt.Printf("Cannot start participant.Error: %v\n", err)
		os.Exit(1)
	}
	p.clientConnection = conn
	p.client = pb.NewParticipantServiceClient(conn)

	ctx := context.TODO()
	nodeId, err := p.client.Register(ctx, &pb.Participant{Name: p.name, Score: 0, Address: p.address})
	if err != nil {
		fmt.Printf("Cannot start participant.Error: %v\n", err)
		os.Exit(1)
	}
	p.nodeId = nodeId.NodeId

}

// StartQuiz starts the quiz for the participant and waits until it is completed.
// 1. Call the "MarkReady" RPC with the node id received in Register RPC.
// 2. Wait until the completion of the all the questions (can be done through a channel).
func (p *participant) StartQuiz(serverPort int) {
	p.client.MarkReady(context.Background(), &pb.RegisterResponse{
		NodeId: p.nodeId,
	})
	<-p.quizDone
}

// Answer function receives the question from quizmaster.
// 1. Display the questions and options to the user
// 2. Waits for ParticipantReadTimeout seconds for user input.
// 3. If the user input is received within the timeout, then the response is sent
// as ParticipantAnswer. If timeout happens then InvalidAnswer is sent as response.
// You can use the ReadOptionFromUser to read the participant's response.
func (p *participant) Answer(ctx gorums.ServerCtx, in *pb.Question) (*pb.ParticipantAnswer, error) {
	fmt.Printf("Question no. %d . %s\n", in.GetId(), in.QuestionText)

	for i, ans := range in.AnswerText {
		fmt.Printf("%d) %s\n", i+1, ans)
	}
	ans := p.ReadOptionFromUser()

	fmt.Printf("Anser selected: %d\n", ans)
	if in.IsLastQuestion {
		p.quizDone <- true
	}

	return &pb.ParticipantAnswer{
		QuestionId: in.GetId(),
		Answer:     ans,
		Participant: &pb.Participant{
			Name:    p.name,
			Address: p.address,
			Score:   0,
		},
	}, nil
}

// GetResults fetches the results from the quizmaster.
// 1. Call the "GetResults" RPC to fetch the results of all the participants.
// 2. Waits only ParticipantGetResultTimeout seconds for the RPC completion.
// 3. After the RPC completion, display the results on the STDOUT
func (p *participant) GetResults(serverPort int) {
	ctx, cancel := context.WithTimeout(context.Background(), quiz.ParticipantGetResultTimeout)
	defer cancel()
	result, err := p.client.GetResults(ctx, &emptypb.Empty{})
	if err != nil {
		fmt.Printf("Error getting result : %s", err.Error())
	}

	participants := result.GetParticipants()
	for _, user := range participants {
		fmt.Printf("Name : %s , Address: %s ,Score : %d\n", user.Name, user.Address, user.Score)
	}
}

// ReadOptionFromUser reads the option from an user from the STDIN and pass it to the
// inputChannel. It reads the input within ParticipantReadTimeout seconds
// and if a reply is done by a participant, it sends the reply to the inputChannel.
// In case of timeout, the function should send InvalidAnswer in the inputChannel.
func (p *participant) ReadOptionFromUser() int32 {
	val := utils.ReadAnswerWithTimeout(quiz.ParticipantReadTimeout, p.inputChannel)
	return val
}

// CloseConnection closes the client connection at the end of quiz
func (p *participant) CloseConnection() {
	p.clientConnection.Close()
}
