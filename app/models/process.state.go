package models

import (
	"errors"
	"fmt"
)

const (
	ProcessNotStarted = 0
	ProcessStarted    = 1
	ProcessFinished   = 2
	ProcessAborted    = 3
)

type Process struct {
	ChannelOut      chan string
	ProcessMessage  chan string
	StateChannel    chan SubProcessState
	State           int
	NumberProcessed int
	FailedProcesses int
	TotalSub        int
	SubProcesses    map[string]*SubProcess
	ProcessName     string
}

type SubProcess struct {
	ProcessMessage chan string
	StateChannel   chan SubProcessState
	State          int
	ProcessName    string
	Message        string
	Current        int
	Total          int
}

type SubProcessState struct {
	State       int
	ProcessName string
}

func NewProcess(processName string, channel chan string) *Process {
	return &Process{
		ChannelOut:      channel,
		State:           ProcessNotStarted,
		ProcessName:     processName,
		FailedProcesses: 0,
		TotalSub:        0,
		ProcessMessage:  make(chan string),
		SubProcesses:    make(map[string]*SubProcess),
		StateChannel:    make(chan SubProcessState),
	}
}

func (p *Process) Start() {
	p.State = ProcessStarted

	// Send process started message.
	p.ChannelOut <- fmt.Sprintf("Process: %v started.", p.ProcessName)

	go func() {
		for {
			if p.State == ProcessFinished || p.State == ProcessAborted {
				break
			}
			select {
			case message := <-p.ProcessMessage:
				// Send channel out
				p.ChannelOut <- message
			case subState := <-p.StateChannel:
				p.SubProcessState(subState)
			}
		}
	}()
}

func (p *Process) SubProcessState(subState SubProcessState) {
	switch subState.State {
	case ProcessFinished:
		p.NumberProcessed++
		p.ChannelOut <- fmt.Sprintf("Sub process: %v finished.", subState.ProcessName)
		break
	case ProcessAborted:
		p.ChannelOut <- fmt.Sprintf("Sub process: %v aborted.", subState.ProcessName)
		p.FailedProcesses++
		break
	case ProcessStarted:
		p.ChannelOut <- fmt.Sprintf("Sub process: %v started.", subState.ProcessName)
	}
}

func (p *Process) Abort() {
	p.State = ProcessAborted
	p.ChannelOut <- fmt.Sprintf("Process: %v aborted.", p.ProcessName)
}

func (p *Process) Finish() {
	p.State = ProcessFinished
	p.ChannelOut <- fmt.Sprintf("Process: %v finished.", p.ProcessName)
}

func (p *Process) Update() {

}

func (p *Process) AddSubProcess(processName string, channel chan string) (*Process, error) {
	if _, exists := p.SubProcesses[processName]; exists {
		return p, errors.New("Sub process exists already.")
	}

	p.SubProcesses[processName] = &SubProcess{
		ProcessMessage: channel,
		State:          ProcessNotStarted,
		ProcessName:    processName,
	}

	return p, nil
}

func (p *Process) GetSubProcess(processName string) (*SubProcess, error) {
	if sub, exists := p.SubProcesses[processName]; exists {
		return sub, nil
	}

	return &SubProcess{}, errors.New("Sub process doesn't exist.")
}
