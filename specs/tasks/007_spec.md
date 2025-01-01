# Task #6 - Image Support

Support test cases that have an image as part of their context

## Objectives
- One test suite in the demo folder gets an image added to it. Adjust the json config of the test suite @benchmarks/demo/test_suite_001/test_001_001 to include an config.json. This config.json should reference all the available files to be used as context and output examples.
- The benchmark_config needs to be adjusted to reflekt this new style of test case handling
- this should be consistent and all the other test cases should get an config as well
- The llm ports and implementation should accept images as well