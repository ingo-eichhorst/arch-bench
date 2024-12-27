# Software Architecture Benchmark CLI Tool Specification

## 1. Overview
A command-line tool for benchmarking Large Language Models (LLMs) and Vision Language Models (VLMs) on software architecture tasks. The tool executes predefined benchmark scenarios, collects results, and generates evaluation reports.

### 1.1 Key Features
- Execute benchmarks across multiple evaluation dimensions
- Support for both text-only and vision-based architectural tasks (Input can be PDF, text, diagram-images or a mixure of them)
- Extensible benchmark scenario definition
- Standardized result collection and reporting
- Model-agnostic design with pluggable model interfaces

## 2. Quality Attributes

### 2.1 Maintainability
- Hexagonal architecture to separate core logic from external dependencies
- Clear separation between benchmark definitions and execution logic
- Pluggable model interfaces for easy addition of new LLMs/VLMs

### 2.2 Extensibility
- Support for adding new benchmark scenarios without code changes
- Pluggable scoring mechanisms for different evaluation metrics
- Flexible report generation formats

### 2.3 Reliability
- Robust error handling for API failures and timeout scenarios
- Result persistence to prevent data loss
- Validation of benchmark definitions and configurations

### 2.4 Performance
- Parallel execution of independent benchmark tasks
- Efficient handling of large architectural diagrams
- Optimized model API usage

## 3. Technical Architecture

The Programming language used is Golang aka Go. The style of writing the app is hexagonal to seperate the core domain from the implementation of external dependencies like LLMs.

### 3.1 Core Components
1. **CLI Interface**
   - Command parsing and validation
   - Configuration management
   - Progress reporting

2. **Benchmark Engine**
   - Benchmark scenario loading (It's planned to have 10 TestSuites with 100 test cases each)
   - Task execution orchestration
   - In error cases one benchmark run can be retried from where it left of
   - Result collection

3. **Model Adapters**
   - LLM/VLM API integration
   - Response parsing and validation
   - Error handling

4. **Storage**
   - Benchmark status should be saved for one run so it can be picked up or we can rerun some parts of the scenario
   - Storage starts as local files
   - Report generation on CLI and in the stored files

## 4. Data Models

Benchmark has 1..n TestSuites
TestSuite has 1..n TestCases
TestSuite has 1..n Metrics (that are used to test every TestCase in the TestSuite)
TestCase has 1..1 DataSet (that consists of a set of inputs and expected outputs)
A TestCase generates 1..1 EvaluationResult that is later on used to generate the complete report
Provider has 1..n Model
A TestSuite has 1..n Models are going to be tested in parallel
A Benchmark has 1..1 Model configured for evaluating the results of the TestCases

## 5. Implementation Guidelines

### 5.1 Directory Structure for the start
```
/cmd
  /bench
    main.go
/internal
  /app
    /ports           # Primary ports
    /adapters        # Secondary adapters
    /core           # Core domain logic
  /pkg
    /models         # Data models
    /utils          # Shared utilities
/configs            # Configuration files
/scenarios         # Benchmark scenario definitions
```

### 5.2 Configuration
- YAML-based configuration for:
  - Model API credentials
  - Benchmark parameters
  - Scoring criteria
  - Output formats

### 5.3 Error Handling
- Detailed error types for different failure scenarios
- Graceful degradation for non-critical failures
- Comprehensive error logging

## 6. CLI Interface

```bash
# Run specific benchmark
benchmark run <benchmark-name>

# Run specific test suite
benchmark run <benchmark-name> --test-suite <test-suite-name>

# List available benchmarks
benchmark list benchmarks

# List available test suites
benchmark list test-suites <benchmark-name>

# List available provider and models
benchmark list providers
```

## 7. Future Considerations
- Support for batch benchmark execution
- Integration with CI/CD pipelines
- Web interface for result visualization
- Benchmark scenario sharing mechanism
- Support for custom scoring plugins