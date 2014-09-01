package process_state

import (
	"testing"
)

type ProcessState struct {
}

func Init() *ProcessState {
	return &ProcessState{}
}

func (p *ProcessState) RunTests(t *testing.T) {
	p.TestProcessStartsAndEnds(t)
}
