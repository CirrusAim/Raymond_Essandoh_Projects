package quiz

import pb "dat520/lab2/quiz/proto"

var Questions = []*pb.Question{
	{
		Id:           1,
		QuestionText: "Which one of these characters is not friends with Harry Potter?",
		AnswerText: []string{
			"Ron Weasley",
			"Neville Longbottom",
			"Draco Malfoy",
			"Hermione Granger",
		},
		CorrectAnswer: 2,
	},
	{
		Id:            2,
		QuestionText:  "What is the color of Donald Duck's bowtie?",
		AnswerText:    []string{"Red", "Yellow", "Blue", "White"},
		CorrectAnswer: 0,
	},
	{
		Id:           3,
		QuestionText: "Who was the only U.S. President to resign?",
		AnswerText: []string{
			"Herbert Hoover",
			"Richard Nixon",
			"George W. Bush",
			"Donald Trump",
		},
		CorrectAnswer: 1,
	},
	{
		Id:           4,
		QuestionText: "In Pirates of the Caribbean, what was the name of Captain Jack Sparrow's ship?",
		AnswerText: []string{
			"The Marauder",
			"The Slytherin",
			"The Black Python",
			"The Black Pearl",
		},
		CorrectAnswer: 3,
	},
	{
		Id:           5,
		QuestionText: "Which one of these characters aren't a part of the Friends group?",
		AnswerText: []string{
			"Rachel",
			"Joey",
			"Gunther",
			"Monica",
		},
		CorrectAnswer: 2,
	},
	{
		Id:            6,
		QuestionText:  "Who is not a computer scientist?",
		AnswerText:    []string{"Edsger Dijkstra", "Leslie Lamport", "Barbara Liskov", "Richard Feynman"},
		CorrectAnswer: 3,
	},
}
