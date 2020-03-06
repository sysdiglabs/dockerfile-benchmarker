package benchmarker

import (
	"testing"
)

var (
	bm = NewDockerBenchmarker()
)

func init() {
	bm.SetTrustedBaseImages([]string{"golang:1.12-alpine", "alpine"})
	bm.SetDisallowedPackages([]string{"netcat"})
	bm.SetSecretPattern([]string{"key", "secret"})
	bm.debugMode = true
}

func TestPassedDockerfile(t *testing.T) {
	err := bm.ParseDockerfile("../test/Dockerfile_pass")

	if err != nil {
		t.Errorf(err.Error())
		return
	}

	bm.RunBenchmark()

	report := bm.GetViolationReport()

	for _, v := range report.Violations {
		if len(v.Violations) > 0 {
			t.Errorf("rule '%s' failed: %s", v.Rule, v.Violations)
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
		if len(v.Violations) == 0 {
			t.Errorf("rule '%s' failed: %s", v.Rule, v.Violations)
		}
	}
}
