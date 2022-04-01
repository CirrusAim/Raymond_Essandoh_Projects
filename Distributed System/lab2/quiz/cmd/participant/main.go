package main

import (
	"dat520/lab2/quiz"
	"dat520/lab2/quiz/game"
	"dat520/lab2/quiz/utils"
	"flag"
	"fmt"
	"log"
)

// Participant takes two arguments, first one is the name of the participant and
// the second one is the port on which the quizmaster is listening for the participants.
// It initializes a new participant with the name provided in the command line
// argument and register with the quizmaster. It starts answering the questions
// after calling the startQuiz. After completion of the quiz, the participant
// should provide two options to the user, 1 to GetResults and the 2 to Exit.
// Do not change this function
func main() {
	var (
		participant = flag.String("name", "", "Name of participant for quiz.")
		qmPort      = flag.Int("port", quiz.QuizmasterDefaultPort, "Quizmaster's port listening for participants.")
	)
	flag.Parse()

	if err := checkPort(*qmPort); err != nil {
		log.Fatal(err)
	}
	if len(*participant) == 0 {
		log.Fatal("participant name must be provided.")
	}

	p := game.NewParticipant(*participant)
	fmt.Println("Enter to register with quiz master")
	if !utils.ReadEnter() {
		log.Fatal("unable to read input")
		return
	}
	p.RegisterParticipant(*qmPort)
	fmt.Println("Enter to start the quiz")
	if !utils.ReadEnter() {
		log.Fatal("unable to read input")
		return
	}
	p.StartQuiz(*qmPort)
	for {
		fmt.Println()
		fmt.Println("Select one of the operation")

		expectedCommands := []string{"ShowResults", "Exit"}
		/*for num, command := range []string{"GetResults", "Exit"} {
			fmt.Printf("%d) %s", num+1, command)
			fmt.Println()
		}*/
		option := utils.ReadCommand(expectedCommands)
		//option := p.ReadOptionFromUser()
		if option == 0 {
			p.GetResults(*qmPort)
		} else {
			break
		}
	}
	p.CloseConnection()
}

func checkPort(port int) error {
	if port < 1 || port > 65536 {
		return fmt.Errorf("invalid port number %d", port)
	}
	return nil
}
