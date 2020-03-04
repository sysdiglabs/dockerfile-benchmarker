# dockerfile-benchmarker
CIS Docker Benchmark for dockerfiles

## Use cases
Run CIS Docker Benchmark rules for dockerfiles. The following CIS rules are applicable:
1. CIS 4.1 Create a user for the container
2. CIS 4.6 Add HEALTHCHECK instruction to the container image
3. CIS 4.7 Do not use update instructions alone in the Dockerfile
4. CIS 4.9 Use COPY instead of ADD in Dockerfile

## Build
`make build`

## Usage
`dockerfile-benchmarker -d <directory of dockerfiles> -p <pattern>`
```
$ ./dockerfile-benchmarker -h
dockerfile-benchmarker runs CIS Docker Benchmark for dockerfiles. Rule applicable are 4.1, 4.6. 4.7 and 4.9.

Usage:
  dockerfile-benchmarker [flags]

Flags:
  -d, --directory string   directory to lookup for dockerfile (default: ./)
  -h, --help               help for dockerfile-benchmarker
      --level string       Log level (default "info")
  -p, --pattern string     dockerfile name pattern (default: dockerfile)
  ```
