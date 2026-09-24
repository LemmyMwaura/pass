package reader

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"syscall"

	"golang.org/x/term"
)

// InputReader reads prompts from stdin.
type InputReader struct {
	reader *bufio.Reader
}

func NewInputReader() *InputReader {
	return &InputReader{
		reader: bufio.NewReader(os.Stdin),
	}
}

// ReadUserInput reads a line of visible text from the user.
func (r *InputReader) ReadUserInput(prompt string) (string, error) {
	fmt.Print(prompt)

	text, err := r.reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(text), nil
}

// ReadPassword reads a line without echoing it to the terminal.
func (r *InputReader) ReadPassword(prompt string) (string, error) {
	fmt.Print(prompt)

	bytes, err := term.ReadPassword(int(syscall.Stdin))
	fmt.Println()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(bytes)), nil
}
