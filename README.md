# ⌀ Sentinel

**The vigilant guard for code authenticity.**

Sentinel is a multi-language static analysis tool that detects patterns indicative of AI-generated code. It integrates seamlessly into your build process and CI/CD pipelines to help maintain code quality and authenticity.

> [!WARNING]
> **Experimental Build:** Sentinel is in active development as a hobby project.

## Core Features

- **CI/CD Ready** - GitHub Actions, GitLab CI, Jenkins integration
- **Fast & Lightweight** - Single binary, no dependencies. Built with Go for exceptional performance and minimal resource footprint
- **Multi-Language** - Python, Java, JavaScript, TypeScript, Go, Rust, C/C++, Ruby, PHP, C#, Kotlin, Swift (early support - actively improving)
- **Detailed Reporting** - JSON and human-readable output formats
- **Heuristic-Driven Analysis:** Sentinel currently relies on heuristic-based checking (predefined "hard-coded" logic). I am actively working and learning on integrating Machine Learning models to enhance detection.

## Sample Output

```bash
./sentinel scan --path ./examples/ai/python --verbose
```

```
Scanning path: ./examples/ai/python
Threshold: 0.70
Found 1 files to scan
================================================================================
                       ⌀ SENTINEL - Code Detection Report
================================================================================

Total Files Scanned:  1
Files Detected:       1
Average Score:        0.90
Detection Threshold:  0.70

[!] DETECTED FILES (above threshold):
--------------------------------------------------------------------------------
examples\ai\python\ai_generated.py
   Score: 0.90 | Language: python
   Signals:
     • generic_naming (0.90): Very high use of generic variable names
       Generic name occurrences: 47
     • boilerplate_patterns (0.60): Multiple boilerplate patterns detected
       Boilerplate matches: 5
     • formatting_consistency (0.50): Very consistent indentation (possibly AI-generated)
       Unique indentation levels: 3
     • code_complexity (0.30): Low cyclomatic complexity
       Control flow density: 0.07

================================================================================
FAILED: 1 file(s) detected as likely AI-generated
================================================================================
```

## Made with Go Language (Why Go?):

- Unlike Java (JVM) or Python, Go compiles directly to a single machine code binary. This allowed me to drop Sentinel into a CI environment and run it instantly without pre-installing any runtimes or libraries.
- Low memory usage compared to JVM or Node.js-based tools.
- Rich CLI ecosystem for building command-line interfaces

<img src="https://github.com/ashleymcnamara/gophers/blob/master/This_is_Fine_Gopher.png?raw=true" style="width: 100%;">

---

# Installation

## Build from the Source

```bash
# Clone or extract the project
cd sentinel

# IMPORTANT: Generate go.sum with correct checksums
go mod tidy

# Download dependencies
go mod download

# Build
go build -o sentinel.exe

# Install (optional)
sudo mv sentinel.exe /usr/local/bin/
```
---

# Quick Start

```bash
# SCAN CURRENT DIRECTORY
./sentinel scan --path .

# SCAN WITH SPECIFIC TRESHOLD
./sentinel scan --path ./src --threshold 0.8

# SCAN SPECIFIC LANGUAGES
./sentinel scan --path . --languages python,javascript

# JSON OUTPUT
./sentinel scan --path . --json

# FAIL BUILD IF AN AI GEN CODE DETECTED
./sentinel scan --path . --fail-on-detection --threshold 0.75

# COLLECT TRAINING DATA FOR ML MODELS (directory should be ready inside the project)
./sentinel scan --path ./examples/ai --collect --label ai
./sentinel scan --path ./examples/human --collect --label human

# Use ML + heuristics (MACHINE LEARNING MODEL IS REQUIRED!)
./sentinel scan --path . --hybrid --ml-weight 0.7 --verbose

# Heuristics only (no ML)
./sentinel scan --path . --no-ml --verbose

# ML only (MACHINE LEARNING MODEL IS REQUIRED!)
./sentinel scan --path . --ml-only --verbose
```
---

# Usage

**Flags:**
- `--path, -p` - Path to scan (default: current directory)
- `--languages, -l` - Comma-separated list of languages to scan
- `--threshold, -t` - Detection threshold 0.0-1.0 (default: 0.7)
- `--fail-on-detection` - Exit with error code if AI code detected
- `--exclude` - Paths to exclude from scanning
- `--verbose` - Verbose output with detailed signals
- `--json` - Output results in JSON format
- `--collect` - Save scan results to local training CSV
- `--label` - Label for training data: 'ai' or 'human'
- `--model` - Path to ONNX model file
- `--no-ml` - Disable ML inference (heuristics only)
- `--ml-only` - Use ML only (fail if model not available)
- `--ml-weight` - Weight given to ML score (0.0-1.0)

## Git Diff Scanning

```bash
# Scan only files changed in PR
./sentinel scan --git-diff origin/main --threshold 0.75
```

# CI/CD Integration

## GitHub Actions

```yaml
name: SENTINEL CODE SCANNER
on: [pull_request]

jobs:
  sentinel:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
        with:
          fetch-depth: 0
      
      - name: Download Sentinel
        run: |
          curl -L https://github.com/jakeeviado/sentinel/releases/latest/download/sentinel-linux-amd64 -o sentinel
          chmod +x sentinel
      
      - name: Scan Code
        run: ./sentinel scan --path . --threshold 0.75 --fail-on-detection
```

## GitLab CI

```yaml
sentinel:
  stage: test
  script:
    - curl -L https://github.com/jakeeviado/sentinel/releases/latest/download/sentinel-linux-amd64 -o sentinel
    - chmod +x sentinel
    - ./sentinel scan --path . --threshold 0.75 --fail-on-detection
  only:
    - merge_requests
```

## Jenkins

```groovy
pipeline {
    agent any
    stages {
        stage('SENTINEL CODE SCANNER') {
            steps {
                sh '''
                    curl -L https://github.com/jakeeviado/sentinel/releases/latest/download/sentinel-linux-amd64 -o sentinel
                    chmod +x sentinel
                    ./sentinel scan --path . --threshold 0.75 --fail-on-detection
                '''
            }
        }
    }
}
```

---

# License

MIT License - see [LICENSE](LICENSE) file for details.

# Disclaimer

Sentinel is a detection tool that identifies patterns commonly associated with AI-generated code. It should be used as one part of a comprehensive code review process, not as the sole arbiter of code authenticity.

---

Hell Yeah!
