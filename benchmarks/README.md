# Benchmarks Directory

This directory contains all available benchmarks, each with its own set of test suites and test cases.

Structure:
```
benchmarks/
├── benchmark1/
│   ├── config.json
│   ├── testsuite1/
│   │   ├── config.json
│   │   ├── testcase1/
│   │   │   ├── input.txt
│   │   │   └── expected_output.txt
│   │   └── testcase2/
│   │       ├── input.txt
│   │       └── expected_output.txt
│   └── testsuite2/
│       ├── config.json
│       └── ...
├── benchmark2/
│   ├── config.json
│   └── ...
└── ...
```

## Adding a new benchmark
1. Create a new directory under `benchmarks/` with your benchmark name.
2. Add a `config.json` file in the benchmark directory to define global benchmark settings.
3. Create subdirectories for each test suite within the benchmark.
4. In each test suite directory, add a `config.json` file for suite-specific settings.
5. Create subdirectories for each test case within the test suite directories.
6. Add necessary input and expected output files for each test case.

## Configuration files
- `benchmark/config.json`: Define benchmark-wide settings, metrics, and models to be used.
- `testsuite/config.json`: Define suite-specific settings, overrides, or additional configurations.

Please refer to the example benchmarks and test suites for more detailed structure and configuration options.