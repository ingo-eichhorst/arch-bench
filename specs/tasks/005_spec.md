# Task #5 - Metric: CLI commands

Implement the CLI commands for the benchmark service

Run specific test suite:
go run main.go run <benchmark-name> --test-suite <test-suite-name>

List available benchmarks:
go run main.go list benchmarks

List available test suites
go run main.go list test-suites <benchmark-name>

List available provider and models
go run main.go list providers

## Objective
  - All commands should be executable from the CLI
	- The results should be returned in a CLI friendly easy to read format

## Details
  - Use adapters/cli/commands.go file for implementing the commands
