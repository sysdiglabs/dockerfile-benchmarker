package benchmark

const (
	CIS_4_1 = "CIS 4.1 Create a user for the container"
	CIS_4_6 = "CIS 4.6 Add HEALTHCHECK instruction to the container image"
	CIS_4_7 = "CIS 4.7 Do not use update instructions alone in the Dockerfile"
	CIS_4_9 = "CIS 4.9 Use COPY instead of ADD in Dockerfile"
)

type ViolationReport struct {
	Violations []Violation `json:"benchmark_violation_report"`
}

type Violation struct {
	Rule  string   `json:"cis_rule"`
	Files []string `json:"files"`
}

func NewBenchmarkViolation(rule string, files []string) Violation {
	return Violation{
		Rule:  rule,
		Files: files,
	}
}

func NewBenchmarkViolationReport() *ViolationReport {
	return &ViolationReport{
		Violations: []Violation{},
	}
}

func (vr *ViolationReport) AddViolation(rule string, files []string) {
	violation := NewBenchmarkViolation(rule, files)

	vr.Violations = append(vr.Violations, violation)
}