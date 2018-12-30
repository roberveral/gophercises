package runner

import (
	"bufio"
	"fmt"
	"io"
)

// Runner is an interface which allows to define a quiz runner. A quiz runner
// manages tasks like asking the user a question, reading the answer and showing
// the results.
//
// There can be different implementations: using the stdin/stdout, network calls, etc.
type Runner interface {
	// Ask asks a question to the user and retrieves its answer.
	Ask(number int, question string) string
	// ShowResults shows the results of the quiz to the user.
	ShowResults(correctAnswers, totalAnswers int)
	// NotifyTimeout notifies the user that the time to complete the quiz has expired.
	NotifyTimeout(correctAnswers, totalAnswers int)
}

// IoRunner is a Runner which implements the interface by reading from a given io.Reader
// and writting to a given io.Writer. This can be used to run the quiz in the console by
// passing os.Stdin and os.Stdout as reader and writer respectively.
type IoRunner struct {
	reader *bufio.Reader
	writer io.Writer
}

// NewIoRunner creates a new Runner which implements the interface by reading from a given io.Reader
// and writting to a given io.Writer.
func NewIoRunner(reader io.Reader, writer io.Writer) Runner {
	return &IoRunner{bufio.NewReader(reader), writer}
}

// Ask asks a question to the user and retrieves its answer.
func (r *IoRunner) Ask(number int, question string) string {
	fmt.Fprintf(r.writer, "Problem #%v: %s = ", number, question)
	answer, _ := r.reader.ReadString('\n')
	return answer
}

// ShowResults shows the results of the quiz to the user.
func (r *IoRunner) ShowResults(correctAnswers, totalAnswers int) {
	fmt.Fprintf(r.writer, "Scored %v out of %v\n", correctAnswers, totalAnswers)
}

// NotifyTimeout notifies the user that the time to complete the quiz has expired.
func (r *IoRunner) NotifyTimeout(correctAnswers, totalAnswers int) {
	fmt.Fprintln(r.writer, "\nOooh! Time is past!")
	r.ShowResults(correctAnswers, totalAnswers)
}
