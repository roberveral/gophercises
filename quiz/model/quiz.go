package model

import (
	"bufio"
	"encoding/csv"
	"math/rand"
	"os"
	"time"

	"github.com/roberveral/gophercises/quiz/runner"

	"github.com/pkg/errors"
)

// Quiz is the representation of a set of questions which have to be answered.
type Quiz struct {
	Problems []Problem
}

// LoadFromCSV loads a Quiz from a CSV file which contains records with two fields, the question and the answer.
func LoadFromCSV(path string) (*Quiz, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, errors.Wrap(err, "Unable to open CSV file")
	}
	defer file.Close()

	reader := csv.NewReader(bufio.NewReader(file))
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.Wrap(err, "Malformed CSV file")
	}

	problems := make([]Problem, len(records))

	for i, record := range records {
		problems[i] = Problem{record[0], record[1]}
	}

	return &Quiz{problems}, nil
}

// Execute executes the quiz asking the user for the answers.
// It will show the questions in the given writer and will retrieve the answers from the given reader.
// To complete the quiz the user has to answer the questions before the given timer goes off, so
// all the unanswered questions are considered incorrect
func (q *Quiz) Execute(quizRunner runner.Runner, timer *time.Timer) {
	total := len(q.Problems)
	correct := 0
	answerChannel := make(chan string)

	for i, problem := range q.Problems {
		// Answer has to be processed in another goroutine so it can be dropped when the timer goes off
		go func() { answerChannel <- quizRunner.Ask(i, problem.Question) }()

		// Let's see what happens first, either time runs out or the user places an answer in time
		select {
		case <-timer.C:
			quizRunner.NotifyTimeout(correct, total)
			return
		case answer := <-answerChannel:
			if problem.CheckAnswer(answer) {
				correct++
			}
		}
	}

	quizRunner.ShowResults(correct, total)
}

// Shuffle reorders the quiz problems randomly
func (q *Quiz) Shuffle() {
	rand.Shuffle(len(q.Problems), func(i, j int) {
		q.Problems[i], q.Problems[j] = q.Problems[j], q.Problems[i]
	})
}
