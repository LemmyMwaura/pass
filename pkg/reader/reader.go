package reader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// reader from stdIn
type InputReader struct {
	reader *bufio.Reader
}

func NewInputReader() *InputReader {
	r := bufio.NewReader(os.Stdin)

	return &InputReader{
		reader: r,
	}
}

// ReadUserInput reads input from the user with the provided prompt and returns it.
func (r *InputReader) ReadUserInput(prompt string) (string, error) {
	fmt.Print(prompt)

	text, err := r.reader.ReadString('\n')

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(text), nil
}
