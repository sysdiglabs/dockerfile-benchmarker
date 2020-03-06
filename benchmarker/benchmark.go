package benchmarker

import (
	"fmt"
	"os"
	"strings"

	"github.com/sysdiglabs/dockerfile-benchmarker/utils"

	"github.com/sysdiglabs/dockerfile-benchmarker/pkg/benchmark"
	"github.com/sysdiglabs/dockerfile-benchmarker/pkg/dockerfile"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

type DockerBenchmarker struct {
	dfiles             map[string]*dockerfile.Dockerfile
	violationReport    *benchmark.ViolationReport
	trustedBaseImages  map[string]bool
	disallowedPackages map[string]bool
	secretPatterns     map[string]bool
	debugMode          bool
}

// NewDockerBenchmarker returns a bm object
func NewDockerBenchmarker() *DockerBenchmarker {
	return &DockerBenchmarker{
		dfiles:             map[string]*dockerfile.Dockerfile{},
		violationReport:    benchmark.NewBenchmarkViolationReport(),
		trustedBaseImages:  nil,
		disallowedPackages: nil,
		secretPatterns:     nil,
		debugMode:          false,
	}
}

func (bm *DockerBenchmarker) SetTrustedBaseImages(images []string) {
	if len(images) == 0 {
		return
	}

	bm.trustedBaseImages = map[string]bool{}

	for _, image := range images {
		bm.trustedBaseImages[image] = true
	}
}

func (bm *DockerBenchmarker) SetDisallowedPackages(packages []string) {
	if len(packages) == 0 {
		return
	}

	bm.disallowedPackages = map[string]bool{}

	for _, pkg := range packages {
		bm.disallowedPackages[pkg] = true
	}
}

func (bm *DockerBenchmarker) SetSecretPattern(patterns []string) {
	if len(patterns) == 0 {
		return
	}

	bm.secretPatterns = map[string]bool{}

	for _, pattern := range patterns {
		bm.secretPatterns[pattern] = true
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

	if bm.debugMode {
		fmt.Println(result.AST.Dump())
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

	// CIS 4.2 Use trusted base images for containers
	bm.CheckTrustedBaseImages()

	// CIS 4.3 Do not install unnecessary packages in the container
	bm.CheckDisallowedPackages()

	// CIS 4.6 add HEALTHCHECK instruction to the container image
	bm.CheckHealthCheck()

	// CIS 4.7 Do not use update instructions alone in the Dockerfile
	bm.CheckRunUpdateOnly()

	// CIS 4.9 Use COPY instead of ADD in Dockerfile
	bm.CheckAdd()

	// CIS 4.10 Do not store secrets in Dockerfiles
	bm.CheckSecretsInsideImage()
}

// CIS 4.1 Create a user for the container
func (bm *DockerBenchmarker) CheckNonRootUser() {
	dfiles := []string{}

	for file, df := range bm.dfiles {
		nonRootUserCreated := false
		for _, di := range df.Instructions {
			if di.Instruction == dockerfile.USER {
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
}

// CIS 4.2 Use trusted base images for containers
func (bm *DockerBenchmarker) CheckTrustedBaseImages() {
	// no trusted base images are provided
	if bm.trustedBaseImages == nil {
		return
	}

	violationMap := map[string]bool{}

	for file, df := range bm.dfiles {
		baseImages := df.GetBaseImages()

		for _, image := range baseImages {
			if !bm.IsTrustedBaseImage(image) {
				violation := createViolation(file, image)
				violationMap[violation] = true
			}
		}
	}

	bm.violationReport.AddViolation(benchmark.CIS_4_2, utils.MapToArray(violationMap))
}

// CIS 4.3 Do not install unnecessary packages in the container
func (bm *DockerBenchmarker) CheckDisallowedPackages() {
	// no disallowed packages are provided
	if bm.disallowedPackages == nil {
		return
	}

	violationMap := map[string]bool{}

	for file, df := range bm.dfiles {
		for disallowedPkg := range bm.disallowedPackages {
			// apt
			idxs := df.LookupInstructionAndContent(dockerfile.Run, `apt\s+install\s+[^;|&]+`+disallowedPkg)
			if len(idxs) > 0 {
				violation := createViolation(file, disallowedPkg)
				violationMap[violation] = true
			}

			// apt-get
			idxs = df.LookupInstructionAndContent(dockerfile.Run, `apt-get\s+install\s+[^;|&]+`+disallowedPkg)
			if len(idxs) > 0 {
				violation := createViolation(file, disallowedPkg)
				violationMap[violation] = true
			}

			// apk
			idxs = df.LookupInstructionAndContent(dockerfile.Run, `apk\s+add\s+[^;|&]+`+disallowedPkg)
			if len(idxs) > 0 {
				violation := createViolation(file, disallowedPkg)
				violationMap[violation] = true
			}
		}
	}

	bm.violationReport.AddViolation(benchmark.CIS_4_3, utils.MapToArray(violationMap))
}

// CIS 4.6 add HEALTHCHECK instruction to the container image
func (bm *DockerBenchmarker) CheckHealthCheck() {
	dfiles := []string{}

	for file, df := range bm.dfiles {
		found := df.LookupInstruction(dockerfile.HEALTHCHECK)

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

func (bm *DockerBenchmarker) IsTrustedBaseImage(image string) bool {
	// always return true if there is no trusted base image provided
	if bm.trustedBaseImages == nil {
		return true
	}

	if _, exists := bm.trustedBaseImages[image]; exists {
		return true
	}

	return false
}

// CIS 4.9 Use COPY instead of ADD in Dockerfile
func (bm *DockerBenchmarker) CheckAdd() {
	dfiles := []string{}

	for file, df := range bm.dfiles {
		found := df.LookupInstruction(dockerfile.ADD)

		if found {
			dfiles = append(dfiles, file)
		}
	}

	bm.violationReport.AddViolation(benchmark.CIS_4_9, dfiles)
}

// CIS 4.10 Do not store secrets in Dockerfiles (check label and env instructions only)
func (bm *DockerBenchmarker) CheckSecretsInsideImage() {
	if bm.secretPatterns == nil {
		return
	}

	violationMap := map[string]bool{}

	for file, df := range bm.dfiles {
		for secretPattern := range bm.secretPatterns {
			// ENV
			idxs := df.LookupInstructionAndContent(dockerfile.ENV, secretPattern)
			if len(idxs) > 0 {
				violation := createViolation(file, fmt.Sprintf("ENV contains '%s'", secretPattern))
				violationMap[violation] = true
			}

			// LABEL
			idxs = df.LookupInstructionAndContent(dockerfile.LABEL, secretPattern)
			if len(idxs) > 0 {
				violation := createViolation(file, fmt.Sprintf("LABEL contains '%s'", secretPattern))
				violationMap[violation] = true
			}
		}
	}

	bm.violationReport.AddViolation(benchmark.CIS_4_10, utils.MapToArray(violationMap))
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

func createViolation(dockerfile, detail string) string {
	return fmt.Sprintf("%s: %s", dockerfile, detail)
}
