package process_state

import (
	"github.com/tuomasvapaavuori/site_installer/app/models"
	"testing"
)

func (p *ProcessState) TestProcessStartsAndEnds(t *testing.T) {
	outChan := make(chan string)

	process := models.NewProcess("Test Process", outChan)

	go func() {
		finished := p.ListenOutChannel(outChan, t)
		if !finished {
			t.Error("Error finishing process.")
		}
	}()

	process.Start()
	process.Finish()
}

func (p *ProcessState) ListenOutChannel(channel chan string, t *testing.T) bool {
	for {
		msg := <-channel
		switch msg {
		case "Process: Test Process started.":
			t.Log(msg)
			break
		case "Process: Test Process aborted.":
			t.Log(msg)
			break
		case "Process: Test Process finished.":
			t.Log(msg)
			return true
		}
	}
}
