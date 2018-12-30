package model

import "strings"

// Problem is a question to ask and the correct answer to the question.
type Problem struct {
	// Question asked
	Question string
	// Answer to the question
	Answer string
}

// CheckAnswer checks whether the given answer is correct for this Problem.
func (q *Problem) CheckAnswer(answer string) bool {
	return strings.TrimSpace(q.Answer) == strings.TrimSpace(answer)
}
