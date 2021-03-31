module github.com/sysdiglabs/dockerfile-benchmarker

go 1.12

require (
	github.com/moby/buildkit v0.6.4
	github.com/sirupsen/logrus v1.4.2
	github.com/spf13/cobra v0.0.6
	golang.org/x/tools v0.1.0 // indirect
	k8s.io/apimachinery v0.17.3
)

replace (
	github.com/containerd/containerd v1.3.0-0.20190507210959-7c1e88399ec0 => github.com/containerd/containerd v1.3.0-beta.2.0.20190823190603-4a2f61c4f2b4
	github.com/docker/docker => github.com/moby/moby v0.7.3-0.20190826074503-38ab9da00309
	golang.org/x/crypto v0.0.0-20190129210102-0709b304e793 => golang.org/x/crypto v0.0.0-20180904163835-0709b304e793
)
