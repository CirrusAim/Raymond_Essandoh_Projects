package main

import (
	"dat520/lab2/quiz"
	"dat520/lab2/quiz/game"
	"flag"
	"fmt"
	"log"
	"time"
)

// Quizmaster takes the port number as a command line argument,
// this port number is used to run the server for the participants
// to connect.
// runGrpcServer: starts the server for the participants to connect
// QuizmasterServerWaitTimer seconds of wait is done to start the GRPC server
// StartQuiz: starts the quiz for the registered, ready and live participants.
func main() {
	port := flag.Int("port", quiz.QuizmasterDefaultPort, "Quizmaster's port listening for participants.")
	flag.Parse()

	if err := checkPort(*port); err != nil {
		log.Fatal(err)
	}
	q := game.NewQuizMaster(*port)
	go q.RunGrpcServer()
	time.Sleep(quiz.QuizmasterServerWaitTimer)
	log.Println("Waiting for the participants to register")
	q.StartQuiz()
}

func checkPort(port int) error {
	if port < 1 || port > 65536 {
		return fmt.Errorf("invalid port number %d", port)
	}
	return nil
}
