# ⌀ Sentinel

**AI-Generated Code Detector for CI/CD Pipelines**

Sentinel is a multi-language static analysis tool that detects patterns indicative of AI-generated code. It integrates seamlessly into your build process and CI/CD pipelines to help maintain code quality and authenticity.

## Features

- **Fast & Lightweight** - Single binary, no dependencies
- **Multi-Language Support** - Python, Java, JavaScript, TypeScript, Go, Rust, C/C++, and more
- **Heuristic Detection** - Pattern-based analysis without ML overhead
- **Detailed Reporting** - JSON and human-readable output formats
- **CI/CD Ready** - GitHub Actions, GitLab CI, Jenkins integration
- **Configurable** - Adjustable thresholds and language filtering

## Installation

### Pre-built Binaries

Download the latest release for your platform:

```bash
# Linux
curl -L https://github.com/jakeeviado/sentinel/releases/latest/download/sentinel-linux-amd64 -o sentinel
chmod +x sentinel

# macOS
curl -L https://github.com/jakeeviado/sentinel/releases/latest/download/sentinel-darwin-amd64 -o sentinel
chmod +x sentinel

# Windows
curl -L https://github.com/jakeeviado/sentinel/releases/latest/download/sentinel-windows-amd64.exe -o sentinel.exe
```

### Build from Source

```bash
# Clone or extract the project
cd sentinel

# IMPORTANT: Generate go.sum with correct checksums
go mod tidy

# Download dependencies
go mod download

# Build
go build -o sentinel

# Install (optional)
sudo mv sentinel /usr/local/bin/
```

## Quick Start

```bash
# Scan current directory
sentinel scan --path .

# Scan with specific threshold
sentinel scan --path ./src --threshold 0.8

# Scan only Python and JavaScript files
sentinel scan --path . --languages python,javascript

# Output results as JSON
sentinel scan --path . --json

# Fail build if AI code detected
sentinel scan --path . --fail-on-detection --threshold 0.75
```

## Usage

### Basic Scanning

```bash
sentinel scan [flags]
```

**Flags:**
- `--path, -p` - Path to scan (default: current directory)
- `--languages, -l` - Comma-separated list of languages to scan
- `--threshold, -t` - Detection threshold 0.0-1.0 (default: 0.7)
- `--fail-on-detection` - Exit with error code if AI code detected
- `--exclude` - Paths to exclude from scanning
- `--verbose` - Verbose output with detailed signals
- `--json` - Output results in JSON format

### Git Diff Scanning

```bash
# Scan only files changed in PR
sentinel scan --git-diff origin/main --threshold 0.75
```

## CI/CD Integration

### GitHub Actions

```yaml
name: AI Code Detection
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

### GitLab CI

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

### Jenkins

```groovy
pipeline {
    agent any
    stages {
        stage('AI Code Detection') {
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

## Detection Methods

Sentinel uses multiple heuristic signals to detect AI-generated code:

### 1. Comment Density
An AI generated code often contains excessive intergalactic explanatory comments.

### 2. Generic Naming Patterns
Common use of generic variable names like `temp`, `data`, `result`.

### 3. Repetitive Patterns
Unusual amounts of repeated code structures.

### 4. Code Complexity
AI code tends to have lower cyclomatic complexity.

### 5. Formatting Consistency
Perfect formatting consistency can indicate AI generation.

### 6. Boilerplate Patterns
Common AI generated boilerplate and TODO comments (i know some of us use this TODO, and I am planning for re-adjustments).

## Configuration

Create a `.sentinel.yaml` file in your project root:

```yaml
threshold: 0.75
languages:
  - python
  - java
  - javascript
exclude:
  - "vendor/*"
  - "node_modules/*"
  - "*.test.js"
  - "*.spec.py"
verbose: true
```

## Output Examples

### Text Output

```
================================================================================
                       ⌀ SENTINEL - Code Detection Report
================================================================================

Total Files Scanned:  42
Files Detected:       3
Average Score:        0.45
Detection Threshold:  0.70

   DETECTED FILES (above threshold):
--------------------------------------------------------------------------------

   src/utils/helper.py
   Score: 0.85 | Language: python
   Signals:
     • comment_density (0.80): Excessive comment density detected (>30%)
       Comment lines: 45 / 120 (37.5%)
     • generic_naming (0.90): Very high use of generic variable names
       Generic name occurrences: 23
     • boilerplate_patterns (0.60): Multiple boilerplate patterns detected
       Boilerplate matches: 4

================================================================================
   FAILED: 3 file(s) detected as likely AI-generated
================================================================================
```

### JSON Output

```json
{
  "total_files": 42,
  "detected_files": 3,
  "average_score": 0.45,
  "threshold": 0.70,
  "files": [
    {
      "path": "src/utils/helper.py",
      "score": 0.85,
      "language": "python",
      "detected": true,
      "signals": [...]
    }
  ]
}
```

## Fully Supported Languages
- Python (.py) - Full pattern detection
- Java (.java) - Full pattern detection  
- JavaScript (.js) - Full pattern detection

## Partially Supported Languages
The following languages use generic heuristics only:
- TypeScript, Go, Rust, C/C++, Ruby, PHP, C#, Kotlin, Swift

## Future Plans

- [1] Adding more signals including (function/variable naming patterns, documentation style, error handling patterns)
- [2] Tree-sitter integration for AST-based analysis
- [3] Git diff scanning support
- [4] ML-based detection (onnxruntime-go)
- [5] Build tool plugins (Maven, Gradle, npm)
- [6] Custom rule definitions
- [7] Ignore files (`.sentinelignore`)

## License

MIT License - see [LICENSE](LICENSE) file for details.

## Disclaimer

Sentinel is a detection tool that identifies patterns commonly associated with AI-generated code. It should be used as one part of a comprehensive code review process, not as the sole arbiter of code authenticity.

---

This is personal project. Hell Yeah!
