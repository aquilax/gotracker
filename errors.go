package main

import (
	"fmt"
)

type TrackerError struct {
	message string
}

func (te TrackerError) Error() string {
	return fmt.Sprintf("d14:failure reason%d:%se", len(te.message), te.message)
}