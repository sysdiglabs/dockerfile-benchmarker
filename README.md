# dockerfile-benchmarker
CIS Docker Benchmark for dockerfiles

## Use cases
Run CIS Docker Benchmark rules for dockerfiles. The following CIS rules are applicable:
1. CIS 4.1 Create a user for the container
2. CIS 4.2 Use trusted base images for containers (user provide trusted base image list)
3. CIS 4.3 Do not install unnecessary packages in the container (user provide the disallowed package list)
4. CIS 4.6 Add HEALTHCHECK instruction to the container image
5. CIS 4.7 Do not use update instructions alone in the Dockerfile
6. CIS 4.9 Use COPY instead of ADD in Dockerfile
7. CIS 4.10 Do not store secrets in Dockerfiles (user provide the secret pattern, only checks contents in `ENV` and `LABEL` instructions)

## Build
`make build`

## Usage
```
$ ./dockerfile-benchmarker -h
dockerfile-benchmarker runs CIS Docker Benchmark for dockerfiles. Rule applicable are 4.1, 4.2, 4.3, 4.6. 4.7, 4.9 and 4.10.

Usage:
  dockerfile-benchmarker [flags]

Flags:
  -d, --directory string             directory to lookup for dockerfile (default "./")
  -p, --disallowed-packages string   list of disallowed packages separated by comma
  -f, --dockerfile-pattern string    dockerfile name pattern (default "dockerfile")
  -h, --help                         help for dockerfile-benchmarker
      --level string                 Log level (default "info")
  -s, --secret-patterns string       list of secret patterns separated by comma
  -b, --trusted-base-images string   list of trusted base images separated by comma
  ```

## Example output
```
$ ./dockerfile-benchmarker -p "netcat" -s "secret, key" -b "alpine,golang:1.12-alpine" | jq .
INFO[2020-03-05T15:52:50-08:00] Trusted base images: [alpine golang:1.12-alpine] 
INFO[2020-03-05T15:52:50-08:00] Disallowed packages: [netcat]                
INFO[2020-03-05T15:52:50-08:00] Secret patterns: [secret key]                
{
  "benchmark_violation_report": [
    {
      "cis_rule": "CIS 4.1 Create a user for the container",
      "files": [
        "test/Dockerfile_fail"
      ]
    },
    {
      "cis_rule": "CIS 4.2 Use trusted base images for containers",
      "files": [
        "container/Dockerfile",
        "test/Dockerfile_fail"
      ]
    },
    {
      "cis_rule": "CIS 4.3 Do not install unnecessary packages in the container",
      "files": [
        "test/Dockerfile_fail"
      ]
    },
    {
      "cis_rule": "CIS 4.6 Add HEALTHCHECK instruction to the container image",
      "files": [
        "test/Dockerfile_fail"
      ]
    },
    {
      "cis_rule": "CIS 4.7 Do not use update instructions alone in the Dockerfile",
      "files": [
        "test/Dockerfile_fail"
      ]
    },
    {
      "cis_rule": "CIS 4.9 Use COPY instead of ADD in Dockerfile",
      "files": [
        "test/Dockerfile_fail"
      ]
    },
    {
      "cis_rule": "CIS 4.10 Do not store secrets in Dockerfiles",
      "files": [
        "test/Dockerfile_fail"
      ]
    }
  ]
}
```
