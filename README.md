# ⌀ SENTINEL

SENTINEL is a multi-language static analysis CLI tool designed to assist code reviewers by identifying risky, inconsistent, or unconventional code patterns, especially in AI-assisted development environments.

It integrates seamlessly into CI/CD pipelines (GitHub Actions, GitLab CI, Jenkins) to enforce code quality, maintainability, and reliability at scale.

Built for modern engineering teams embracing AI-assisted development workflows, Sentinel provides an additional layer of insight to support consistent and high-quality codebases.

**Features**

- **CI/CD Ready**
- **Fast & Lightweight**
- **Multi-Language** - Python, Java, JavaScript, TypeScript, Go, Rust, C/C++, Ruby, PHP, C#, Kotlin, Swift (early support - actively improving)
- **Hybrid Detection:** Sentinel utilizes a dual-layer verification pipeline, merging **Heuristics-Driven Analysis** with **Machine Learning (ONNX)**.

**Disclaimer**

Sentinel identifies unconventional or high-risk code patterns using heuristic and machine learning techniques. These patterns may appear in both human-written and AI-assisted code.

Results are probabilistic and may include false positives and false negatives. Sentinel is designed to support code quality and review processes, not to determine authorship or replace human judgment.


---

# Quickstart

**Sample Command:**

```bash
./sentinel scan --path ./training-sets/ai --verbose --collect --label ai
./sentinel scan --path ./training-sets/human --verbose --collect --label human
```

**Sample Output:**

```
================================================================================
                       ⌀ SENTINEL - Code Detection Report
================================================================================

Detection Mode:       Heuristics Only
Total Files Scanned:  6
Files Detected:       1
Average Score:        0.62
Detection Threshold:  0.70

[!] DETECTED FILES (above threshold):
--------------------------------------------------------------------------------
training-sets/ai/java/Calculator.java
   Score: 0.80 | Language: java
   Signals:
     • comment_density (0.80): Excessive comment density detected (>30%)
       Comment lines: 47 / 107 (43.9%)
     • boilerplate_patterns (0.30): Some boilerplate patterns present
       Boilerplate matches: 2

[?] SUSPICIOUS FILES (below threshold but noteworthy):
--------------------------------------------------------------------------------
training-sets/ai/python/ai_generated_calculator.py
   Score: 0.60 | Language: python
   Signals:
     • boilerplate_patterns (0.60): Multiple boilerplate patterns detected
       Boilerplate matches: 3
     • formatting_consistency (0.50): Very consistent indentation (possibly AI-generated)
       Unique indentation levels: 3
     • code_complexity (0.30): Low cyclomatic complexity
       Control flow density: 0.06
training-sets/ai/python/ai_generated_data_processor.py
   Score: 0.60 | Language: python
   Signals:
     • generic_naming (0.60): High use of generic variable names
       Generic name occurrences: 42
     • boilerplate_patterns (0.60): Multiple boilerplate patterns detected
       Boilerplate matches: 3
     • code_complexity (0.30): Low cyclomatic complexity
       Control flow density: 0.07
training-sets/ai/python/ai_generated_string_utils.py
   Score: 0.60 | Language: python
   Signals:
     • boilerplate_patterns (0.60): Multiple boilerplate patterns detected
       Boilerplate matches: 3
     • code_complexity (0.30): Low cyclomatic complexity
       Control flow density: 0.10
training-sets/ai/java/DataProcessor.java
   Score: 0.60 | Language: java
   Signals:
     • generic_naming (0.60): High use of generic variable names
       Generic name occurrences: 68
     • boilerplate_patterns (0.60): Multiple boilerplate patterns detected
       Boilerplate matches: 3
     • comment_density (0.50): High comment density (>20%)
       Comment lines: 58 / 205 (28.3%)
     • code_complexity (0.30): Low cyclomatic complexity
       Control flow density: 0.07
training-sets/ai/java/StringUtils.java
   Score: 0.50 | Language: java
   Signals:
     • comment_density (0.50): High comment density (>20%)
       Comment lines: 45 / 177 (25.4%)
     • boilerplate_patterns (0.30): Some boilerplate patterns present
       Boilerplate matches: 2

================================================================================
FAILED: 1 file(s) detected as likely AI-generated
================================================================================
Scanning path: ./training-sets/human
Threshold: 0.70
Found 2 files to scan
Successfully logged results to sentinel_training.csv
================================================================================
                       ⌀ SENTINEL - Code Detection Report
================================================================================

Detection Mode:       Heuristics Only
Total Files Scanned:  2
Files Detected:       0
Average Score:        0.25
Detection Threshold:  0.70

[?] SUSPICIOUS FILES (below threshold but noteworthy):
--------------------------------------------------------------------------------
training-sets/human/python/human_written_calculator.py
   Score: 0.50 | Language: python
   Signals:
     • formatting_consistency (0.50): Very consistent indentation (possibly AI-generated)
       Unique indentation levels: 3
     • code_complexity (0.30): Low cyclomatic complexity
       Control flow density: 0.08

================================================================================
PASSED: No AI-generated code detected above threshold
================================================================================
```

**Other Commands:**

```bash
./sentinel

# Heuristics Only
# Scan current directory
./sentinel scan --path .

# Scan with specific treshold
./sentinel scan --path ./src --threshold 0.8

# Scan specific language/s
./sentinel scan --path . --languages python,javascript

# Output as JSON
./sentinel scan --path . --json

# Build fails if a possible AI generated code detected
./sentinel scan --path . --fail-on-detection --threshold 0.75

# Collect training data for Machine Learning Model (directory should be ready inside the project)
./sentinel scan --path ./examples/ai --collect --label ai
./sentinel scan --path ./examples/human --collect --label human

# Hybrid Mode (Default when --model is provided)
# Blends ML probability with heuristic patterns using a weighted average.
./sentinel scan --path . --model ./model/model.onnx --ml-weight 0.7 --verbose

# Heuristics Only
# Ignores the ML model entirely and looks for raw code patterns.
./sentinel scan --path . --no-ml --verbose

# ML Only
# Forces the detector to rely primarily on the model. 
./sentinel scan --path . --ml-only --model ./model/model.onnx --verbose
```

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

---

# Installation

**Build from the Source:**

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

# CI/CD Integration

**GitHub Actions**

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

**GitLab CI**

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

**Jenkins**

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

**Option for Git Diff Scanning (for CI/CD)**

```bash
# Scan only files changed in PR
./sentinel scan --git-diff origin/main --threshold 0.75
```

---

# License

MIT License - see [LICENSE](LICENSE) file for details.
