# arch-bench: Software Architecture Benchmarking Tool

This repository contains a command-line tool for benchmarking Large Language Models (LLMs) and Vision Language Models (VLMs) on software architecture tasks.  The tool is written in Go and utilizes a hexagonal architecture for maintainability and extensibility.

## Usage

1. **Installation:**

- Install golang from [here](https://go.dev/doc/install)

- Checkout this repository
```sh
  git clone https://github.com/ingo-eichhorst/arch-bench.git
```

- Install the dependencies:
```sh
  go mod install tidy
```

2. **Configuration:**

- Copy the `.env.example` file to `.env` file in the root directory and fill out the environment variables.

```sh
EVAL_API_KEY=xxxxx
EVAL_MODEL=gpt-4o-mini
EVAL_PROVIDER=openai

# Depending on the model your test suites are testing you need to set the API key for that provider.
OPENAI_API_KEY=xxxxx
```

3. **Running Benchmarks:**

Note: Later on this will be a CLI tool called arch-bench. For now we checkout the source code and build and run it manually.

```bash
cd cmd/benchmark

# Run the demo benchmark
go run main.go run demo

# You can copy any available benchmark from the `benchmarks` directory to crete your own.
cp benchmarks/demo benchmarks/my-own
go run main.go run my-own

# Run a specific test suite within a benchmark
go run main.go run <benchmark-name> --test-suite <test-suite-name>

# List available benchmarks
go run main.go list benchmarks
# List test suites in a benchmark
go run main.go list test-suites <benchmark-name>  
# List supported providers and models
go run main.go list providers
```

## Contributions 

TbD: Please contact me or simply create an issue.

## License

MIT
