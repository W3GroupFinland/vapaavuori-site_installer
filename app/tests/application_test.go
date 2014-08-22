package tests

import (
	"github.com/tuomasvapaavuori/site_installer/app/tests/hostmaster"
	"github.com/tuomasvapaavuori/site_installer/app/tests/hosts_domains"
	"testing"
)

func TestRunApplicationTests(t *testing.T) {
	hosts_domains.RunTests(t)
	hostmaster.RunTests(t)
}
