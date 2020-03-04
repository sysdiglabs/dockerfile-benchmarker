package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sysdiglabs/dockerfile-benchmarker/benchmarker"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/util/json"
)

func main() {
	bm := benchmarker.NewDockerBenchmarker()

	log.SetFormatter(&log.TextFormatter{
		FullTimestamp: true,
	})

	var logLevel string
	var dir string
	var pattern string

	var rootCmd = &cobra.Command{
		Use:   "dockerfile-benchmarker",
		Short: "dockerfile-benchmarker runs CIS Docker Benchmark for dockerfiles",
		Long:  "dockerfile-benchmarker runs CIS Docker Benchmark for dockerfiles. Rule applicable are 4.1, 4.6. 4.7 and 4.9.",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			lvl, err := log.ParseLevel(logLevel)
			if err != nil {
				log.Fatal(err)
			}

			log.SetLevel(lvl)
		},
		Run: func(cmd *cobra.Command, args []string) {
			dfs := getDockerfiles(dir, pattern)
			checkDockerfiles(bm, dfs)
		},
	}

	rootCmd.PersistentFlags().StringVar(&logLevel, "level", "info", "Log level")
	rootCmd.Flags().StringVarP(&dir, "directory", "d", "./", "directory to lookup for dockerfile (default: ./)")
	rootCmd.Flags().StringVarP(&pattern, "pattern", "p", "dockerfile", "dockerfile name pattern (default: dockerfile)")
	rootCmd.Execute()
}

func checkDockerfiles(bm *benchmarker.DockerBenchmarker, dfs []string) {
	for _, df := range dfs {
		err := bm.ParseDockerfile(df)

		if err != nil {
			log.Error(err)
			continue
		}
	}

	// run benchmark
	bm.RunBenchmark()

	jsonOutput, err := json.Marshal(bm.GetViolationReport())

	if err != nil {
		log.Error(err)
		return
	}

	fmt.Println(string(jsonOutput))
}

func getDockerfiles(dir, pattern string) []string {
	dfs := []string{}

	err := filepath.Walk(dir,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.Contains(strings.ToLower(filepath.Base(path)), pattern) {
				dfs = append(dfs, path)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}

	return dfs
}
