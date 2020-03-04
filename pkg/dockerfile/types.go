package dockerfile

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/moby/buildkit/frontend/dockerfile/parser"
)

const (
	Healthcheck = "healthcheck"
	User        = "user"
	Root        = "root"
	Add         = "add"
	Run         = "run"
)

// DockerInstruction example:
// - Instruction: run
// - Flags: --from=builder
// - Content: apt-get update
type DockerInstruction struct {
	Instruction string
	Flags       []string
	Content     []string
}

func InitDockerInstruction() *DockerInstruction {
	return &DockerInstruction{
		Instruction: "",
		Flags:       []string{},
		Content:     []string{},
	}
}

type Dockerfile struct {
	File         string
	Instructions []*DockerInstruction
}

func NewDockerfile(file string) *Dockerfile {
	return &Dockerfile{
		File:         file,
		Instructions: []*DockerInstruction{},
	}
}

func (df *Dockerfile) AddNode(node *parser.Node) {
	di := InitDockerInstruction()
	di.Instruction = node.Value
	di.Flags = append(di.Flags, node.Flags...)
	for n := node.Next; n != nil; n = n.Next {
		di.Content = append(di.Content, n.Value)
	}

	df.Instructions = append(df.Instructions, di)
}

func (df *Dockerfile) String() string {
	ret := ""
	ret = fmt.Sprintf("Dockerfile: %s\n", df.File)

	for _, di := range df.Instructions {
		ret += di.Instruction
		if len(di.Flags) > 0 {
			ret += " "
			ret += strings.Join(di.Flags, " ")
		}

		if len(di.Content) > 0 {
			ret += " "
			ret += strings.Join(di.Content, " ")
		}
		ret += "\n"
	}

	return ret
}

func (df *Dockerfile) LookupInstruction(inst string) bool {
	i := strings.ToLower(inst)

	for _, di := range df.Instructions {
		if di.Instruction == i {
			return true
		}
	}

	return false
}

func (df *Dockerfile) LookupInstructionAndContent(inst, cont string) []int {
	indexList := []int{}
	i := strings.ToLower(inst)
	c := strings.ToLower(cont)

	re, err := regexp.Compile(c)

	if err != nil {
		return []int{}
	}

	for idx, di := range df.Instructions {
		if di.Instruction == i {
			for _, content := range di.Content {
				if re.MatchString(content) {
					indexList = append(indexList, idx)
				}
			}
		}
	}

	return indexList
}
