package system

import (
	"testing"
)

func (s *SystemTests) TestGetPlatforms(t *testing.T) {
	err := s.Application.Controllers.System.GetDrupalPlatforms()
	if err != nil {
		t.Error(err)
	}
}
