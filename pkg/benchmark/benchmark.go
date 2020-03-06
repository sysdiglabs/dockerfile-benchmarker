package benchmark

const (
	CIS_4_1  = "CIS 4.1 Create a user for the container"
	CIS_4_2  = "CIS 4.2 Use trusted base images for containers"
	CIS_4_3  = "CIS 4.3 Do not install unnecessary packages in the container"
	CIS_4_6  = "CIS 4.6 Add HEALTHCHECK instruction to the container image"
	CIS_4_7  = "CIS 4.7 Do not use update instructions alone in the Dockerfile"
	CIS_4_9  = "CIS 4.9 Use COPY instead of ADD in Dockerfile"
	CIS_4_10 = "CIS 4.10 Do not store secrets in Dockerfiles"
)

type ViolationReport struct {
	Violations []RuleViolation `json:"cis_docker_benchmark_violation_report"`
}

type RuleViolation struct {
	Rule       string   `json:"rule"`
	Violations []string `json:"violations"`
}

func NewBenchmarkViolation(rule string, violations []string) RuleViolation {
	return RuleViolation{
		Rule:       rule,
		Violations: violations,
	}
}

func NewBenchmarkViolationReport() *ViolationReport {
	return &ViolationReport{
		Violations: []RuleViolation{},
	}
}

func (vr *ViolationReport) AddViolation(rule string, files []string) {
	violation := NewBenchmarkViolation(rule, files)

	vr.Violations = append(vr.Violations, violation)
}
