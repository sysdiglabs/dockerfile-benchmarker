package benchmarker

import (
	"testing"
)

var (
	bm = NewDockerBenchmarker()
)

func TestPassedDockerfile(t *testing.T) {
	err := bm.ParseDockerfile("../test/Dockerfile_pass")

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	bm.RunBenchmark()

	report := bm.GetViolationReport()

	for _, v := range report.Violations {
		if len(v.Files) > 0 {
			t.Errorf("rule '%s' failed: %s", v.Rule, v.Files)
		}
	}
}

func TestFailedDockerfile(t *testing.T) {
	err := bm.ParseDockerfile("../test/Dockerfile_fail")

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	bm.RunBenchmark()

	report := bm.GetViolationReport()

	for _, v := range report.Violations {
		if len(v.Files) == 0 {
			t.Errorf("rule '%s' failed: %s", v.Rule, v.Files)
		}
	}
}
