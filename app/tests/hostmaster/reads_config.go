package hostmaster

import (
	"testing"
)

func (a *ApplicationTests) TestReadsConfig(t *testing.T) {
	// If application mysql user is empty test fails.
	if a.Application.Base.Config.Mysql.User == "" {
		t.Error("app.Base.Config.Mysql.User == \"\", expected not empty value.")
	}
}
