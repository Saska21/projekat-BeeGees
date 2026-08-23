# go-distributed-algorithms

This repository contains Go exercises and projects for learning distributed algorithms.

## Repository Layout

- `week01_02_foundations/`
	- Go fundamentals labs and extra exercises.
- `week03_04_basic_abstractions/`
	- Single-threaded simulator for process/link/failure/crypto abstractions.
	- Link stack: `Fair-Loss -> Stubborn -> Perfect -> Authenticated`.
- `week03_04_parallel/`
	- Parallel (goroutines/channels) version of the same abstractions.

## Setup Go in VS Code

1. Install Go:
	- https://go.dev/dl/
2. Install VS Code:
	- https://code.visualstudio.com/
3. Open VS Code and install the Go extension:
	- https://marketplace.visualstudio.com/items?itemName=golang.Go
4. Verify Go installation in terminal:
	- `go version`
5. Open this project folder in VS Code.
6. Let the Go extension install recommended tools (click Install All if prompted).
7. Run a file from terminal:
	- `go run path/to/main.go`

## Learn Go Online (Editor + Fundamentals)

- A Tour of Go (best for fundamentals with interactive editor):
  - https://go.dev/tour/welcome/1
- Go Playground (online editor to run short Go programs):
  - https://go.dev/play/

## Optional Helpful Docs

- Official Go documentation:
  - https://go.dev/doc/
- Effective Go:
  - https://go.dev/doc/effective_go

## Run Week 3-4 Simulators

From the repository root:

1. Basic abstractions (single-threaded):
	 - `cd week03_04_basic_abstractions && go run .`
2. Parallel version:
	 - `cd week03_04_parallel && go run .`

## Quick Validation

- Build all packages in a module:
	- `go build ./...`
- Run tests in the parallel module:
	- `cd week03_04_parallel && go test ./...`
