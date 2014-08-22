package system

import (
	"github.com/tuomasvapaavuori/site_installer/app"
	"testing"
)

type SystemTests struct {
	Application *app.Application
}

func Init(a *app.Application) *SystemTests {
	return &SystemTests{Application: a}
}

func (s *SystemTests) RunTests(t *testing.T) {
	s.TestGetPlatforms(t)
}
