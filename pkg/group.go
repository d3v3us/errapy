package pkg

import (
	"bytes"
	"fmt"
	"io"
)

type Group []error
type formattedGroup []error

type stateWriter struct {
	fmt.State
	io.Writer
}

func Grouped(errs ...error) error {
	var group Group
	group.Add(errs...)
	return group.Err()
}

func (group *Group) Add(errs ...error) {
	for _, err := range errs {
		if err != nil {
			*group = append(*group, err)
		}
	}
}

func (group Group) Err() error {
	return formattedGroup(group)
}

func (sw stateWriter) Write(p []byte) (n int, err error) {
	return sw.Writer.Write(p)
}

func (group formattedGroup) Format(f fmt.State, c rune) {
	const defaultDelim = "; "
	const newlineDelim = "\n--- "

	delim := defaultDelim
	if f.Flag(int('+')) {
		io.WriteString(f, "group:\n--- ")
		delim = newlineDelim
	}

	var buffer bytes.Buffer
	sw := stateWriter{State: f, Writer: &buffer}

	for i, err := range group {
		if i != 0 {
			buffer.WriteString(delim)
		}
		if formatter, ok := err.(fmt.Formatter); ok {
			formatter.Format(sw, c)
		} else {
			fmt.Fprintf(sw, "%v", err)
		}
	}

	io.WriteString(f, buffer.String())
}

func (group formattedGroup) Error() string { return fmt.Sprintf("%v", group) }
