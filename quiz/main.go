package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/roberveral/gophercises/quiz/runner"

	"github.com/roberveral/gophercises/quiz/model"
)

func main() {
	csvPath := flag.String("csv", "problems.csv", "Path to the CSV file with the problems in the form 'question,answer'")
	timeout := flag.Duration("timeout", 30*time.Second, "Timeout for the user to complete the quiz")
	shuffle := flag.Bool("shuffle", false, "Whether to shuffle the quiz problems or not")
	flag.Parse()

	quiz, err := model.LoadFromCSV(*csvPath)
	if err != nil {
		fmt.Printf("An error occured: %v\n", err)
		os.Exit(1)
	}

	if *shuffle {
		quiz.Shuffle()
	}

	quiz.Execute(runner.NewIoRunner(os.Stdin, os.Stdout), time.NewTimer(*timeout))
}
