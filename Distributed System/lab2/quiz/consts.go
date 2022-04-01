package quiz

import "time"

// participant wait for the user input
const ParticipantReadTimeout = 10 * time.Second

// participant wait for GetResult RPC
const ParticipantGetResultTimeout = 10 * time.Second

// Default port for quizmaster GRPC server
const QuizmasterDefaultPort = 50001

// Score added for correct answer
const CorrectAnswerScore = 100

// The amount of time quizmaster waits for responses from participants
const QuizmasterQuestionTimeout = 20 * time.Second

// Used in creating the manager for configuration in quizmaster before each question
const QuizmasterParticipantConnect = 1 * time.Second

// Quizmaster sleeps this time while checking for active, ready participants
const QuizmasterSleepTimer = 5 * time.Second

// Quizmaster waits for the GRPC server to comeup
const QuizmasterServerWaitTimer = 20 * time.Second

// Default value for invalid answer
const InvalidAnswer = -1
