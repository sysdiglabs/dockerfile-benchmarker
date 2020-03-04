package benchmarker

import (
	"os"
	"strings"

	"github.com/sysdiglabs/dockerfile-benchmarker/pkg/benchmark"
	"github.com/sysdiglabs/dockerfile-benchmarker/pkg/dockerfile"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type DockerBenchmarker struct {
	dfiles          map[string]*dockerfile.Dockerfile
	violationReport *benchmark.ViolationReport
}

// NewDockerBenchmarker returns a bm object
func NewDockerBenchmarker() *DockerBenchmarker {
	return &DockerBenchmarker{
		dfiles:          map[string]*dockerfile.Dockerfile{},
		violationReport: benchmark.NewBenchmarkViolationReport(),
	}
}

func (bm *DockerBenchmarker) ParseDockerfile(file string) error {
	df, err := os.Open(file)

	if err != nil {
		return err
	}
	defer df.Close()

	result, err := parser.Parse(df)

	if err != nil {
		return err
	}

	if _, exists := bm.dfiles[file]; !exists {
		bm.dfiles[file] = dockerfile.NewDockerfile(file)
	}

	for _, node := range result.AST.Children {
		bm.dfiles[file].AddNode(node)
	}

	return nil
}

// GetViolationReport returns the benchmark violation report
func (bm *DockerBenchmarker) GetViolationReport() benchmark.ViolationReport {
	return *bm.violationReport
}

// RunBenchmark runs benchmark check
func (bm *DockerBenchmarker) RunBenchmark() {

	// CIS 4.1 Create a user for the container
	bm.CheckNonRootUser()

	// CIS 4.6 Add HEALTHCHECK instruction to the container image
	bm.CheckHealthCheck()

	// CIS 4.7 Do not use update instructions alone in the Dockerfile
	bm.CheckRunUpdateOnly()

	// CIS 4.9 Use COPY instead of ADD in Dockerfile
	bm.CheckAdd()
}

// CIS 4.1 Create a user for the container
func (bm *DockerBenchmarker) CheckNonRootUser() []string {
	dfiles := []string{}

	for file, df := range bm.dfiles {
		nonRootUserCreated := false
		for _, di := range df.Instructions {
			if di.Instruction == dockerfile.User {
				if len(di.Content) > 0 {
					content := strings.ToLower(di.Content[0])
					if strings.Contains(content, ":") {
						list := strings.Split(content, ":")
						if list[0] != dockerfile.Root && list[1] != dockerfile.Root && list[0] != "0" && list[1] != "0" {
							nonRootUserCreated = true
							break
						}
					} else {
						if content != dockerfile.Root && content != "0" {
							nonRootUserCreated = true
							break
						}
					}
				}
			}
		}

		if !nonRootUserCreated {
			dfiles = append(dfiles, file)
		}
	}

	bm.violationReport.AddViolation(benchmark.CIS_4_1, dfiles)

	return dfiles
}

// CIS 4.6 Add HEALTHCHECK instruction to the container image
func (bm *DockerBenchmarker) CheckHealthCheck() {
	dfiles := []string{}

	for file, df := range bm.dfiles {
		found := df.LookupInstruction(dockerfile.Healthcheck)

		if !found {
			dfiles = append(dfiles, file)
		}
	}

	bm.violationReport.AddViolation(benchmark.CIS_4_6, dfiles)
}

// CIS 4.7 Do not use update instructions alone in the Dockerfile
func (bm *DockerBenchmarker) CheckRunUpdateOnly() {
	dfiles := []string{}

	for file, df := range bm.dfiles {
		// apt
		updateIdxs := df.LookupInstructionAndContent(dockerfile.Run, `apt\s+update`)
		installIdx := df.LookupInstructionAndContent(dockerfile.Run, `apt\s+install`)

		updateOnly, _ := diffArray(updateIdxs, installIdx)
		if len(updateOnly) > 0 {
			dfiles = append(dfiles, file)
		}

		// apt-get
		updateIdxs = df.LookupInstructionAndContent(dockerfile.Run, `apt-get\s+update`)
		installIdx = df.LookupInstructionAndContent(dockerfile.Run, `apt-get\s+install`)

		updateOnly, _ = diffArray(updateIdxs, installIdx)
		if len(updateOnly) > 0 {
			dfiles = append(dfiles, file)
		}

		// apk
		updateIdxs = df.LookupInstructionAndContent(dockerfile.Run, `apk\s+update`)
		installIdx = df.LookupInstructionAndContent(dockerfile.Run, `apk\s+add`)

		updateOnly, _ = diffArray(updateIdxs, installIdx)
		if len(updateOnly) > 0 {
			dfiles = append(dfiles, file)
		}

	}

	bm.violationReport.AddViolation(benchmark.CIS_4_7, dfiles)
}

// CIS 4.9 Use COPY instead of ADD in Dockerfile
func (bm *DockerBenchmarker) CheckAdd() {
	dfiles := []string{}

	for file, df := range bm.dfiles {
		found := df.LookupInstruction(dockerfile.Add)

		if found {
			dfiles = append(dfiles, file)
		}
	}

	bm.violationReport.AddViolation(benchmark.CIS_4_9, dfiles)
}

func diffArray(arr1, arr2 []int) (arr3, arr4 []int) {
	m1 := map[int]bool{}
	m2 := map[int]bool{}

	for _, i := range arr1 {
		m1[i] = true
	}

	for _, i := range arr2 {
		m2[i] = true
	}

	for k := range m1 {
		if _, exists := m2[k]; !exists {
			arr3 = append(arr3, k)
		}
	}

	for k := range m2 {
		if _, exists := m1[k]; !exists {
			arr4 = append(arr4, k)
		}
	}

	return
}
