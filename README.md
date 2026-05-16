# ⌀ SENTINEL

SENTINEL is a multi-language static analysis CLI tool designed to assist code reviewers by identifying risky, inconsistent, or unconventional code patterns, especially in AI-assisted development environments.

It integrates seamlessly into CI/CD pipelines (GitHub Actions, GitLab CI, Jenkins) to enforce code quality, maintainability, and reliability at scale.

Built for modern engineering teams embracing AI-assisted development workflows, Sentinel provides an additional layer of insight to support consistent and high-quality codebases.

**Features**

- **Multi-Language** — Python, Java, JavaScript, TypeScript, Go, Rust, C/C++, Ruby, PHP, C#, Kotlin, Swift
- **Hybrid Detection** — Combines heuristic signal analysis with an ONNX ML model; supports heuristics-only, ML-only, or blended scoring
- **Git Diff Scanning** — Scan only files changed in a PR or branch (`--git-diff`)
- **CI/CD Integration** — Works out of the box with GitHub Actions, GitLab CI, and Jenkins; supports `--fail-on-detection` for automated gates
- **JSON Output** — Machine-readable results for pipeline consumption
- **Configurable Thresholds** — Tune sensitivity per project via flags or `.sentinel.yaml`
- **Training Data Collection** — Built-in tooling to collect labeled scan results for model retraining
- **Fast & Lightweight** — Single binary, no runtime dependencies

--- 

**Disclaimer**

Sentinel identifies unconventional or high-risk code patterns using heuristic and machine learning techniques. These patterns may appear in both human-written and AI-assisted code.

Results are probabilistic and may include false positives and false negatives. Sentinel is designed to support code quality and review processes — **not to determine authorship or replace human judgment.**.


---

# Heuristics

Sentinel evaluates code against the following risk signals:

| Signal | Description |
|---|---|
| `comment_density` | Calculates comment-to-code ratio. Unusually high density can indicate over-documented or templated code. |
| `generic_naming` | Flags heavy use of placeholder names like `temp`, `data`, or `obj`. High frequency raises maintainability and clarity risk. |
| `repetitive_patterns` | Looks for duplicated lines across the file. High repetition often signals copy-paste patterns or low code quality. |
| `code_complexity` | Measures control flow density. Unusually low complexity can indicate shallow or auto-generated logic. |
| `formatting_consistency` | Evaluates indentation uniformity. Suspiciously rigid consistency across a large file is worth flagging. |
| `comment_redundancy` | Flags inline comments that closely mirror the surrounding code. High overlap may indicate over-explained or low-signal documentation. |
| `emoji_sentiment` | Scans for clusters of informal emojis in comments or log messages, which are uncommon in production codebases. |
| `identifier_order` | Checks whether large blocks of identifiers are in perfect alphabetical order. Rarely occurs naturally; worth a closer look. |
| `defensive_ratio` | Evaluates the ratio of defensive checks to functional logic. A heavily skewed ratio can indicate templated error handling. |

> All signals produce a score from `0.0` to `1.0`. Scores are combined into a final risk score compared against the configured `--threshold`.

---

# Usage Example

**Sample Command:**

```bash
./sentinel scan --path ./examples --verbose
```

**Sample Output:**

```
Scanning path: ./examples
Threshold: 0.70
Found 5 files to scan

================================================================================
                       ⌀ SENTINEL - Code Analysis Report
================================================================================

Detection Mode:       Heuristics Only
Total Files Scanned:  5
Files Detected:       2
Average Score:        0.54
Detection Threshold:  0.70

[!] HIGH-RISK FILES (above threshold):
--------------------------------------------------------------------------------
examples\flagged\python\flagged_1.py
   Score: 0.90 | Language: python
   Signals:
     • generic_naming (0.90): Very high use of generic variable names
       Generic name occurrences: 47
     • formatting_consistency (0.50): Unusually rigid indentation uniformity across a large file
       Unique indentation levels: 3
     • code_complexity (0.30): Low cyclomatic complexity
       Control flow density: 0.07
examples\flagged\python\flagged_2.py
   Score: 0.80 | Language: python
   Signals:
     • emoji_sentiment (0.80): High density of informal emojis in production code
       Emojis found: 🚀 (1), ✨ (1), ✅ (1)
     • comment_density (0.50): High comment density (>20%)
       Comment lines: 5 / 22 (22.7%)
     • formatting_consistency (0.50): Unusually rigid indentation uniformity across a large file
       Unique indentation levels: 3

[?] REVIEW RECOMMENDED (below threshold but noteworthy):
--------------------------------------------------------------------------------
examples\flagged\python\flagged_3.py
   Score: 0.60 | Language: python
   Signals:
     • identifier_order (0.60): Large identifier blocks in perfect alphabetical order warrant a closer look
       Perfectly sorted blocks (6+ items): 1
     • formatting_consistency (0.50): Unusually rigid indentation uniformity across a large file
       Unique indentation levels: 2
     • emoji_sentiment (0.30): Informal emojis detected in comments or log messages
       Emojis found: ✨ (1), 🤖 (1)

================================================================================
FAILED: 2 file(s) require attention (risk threshold exceeded)
Note: Scores represent heuristic and/or ML-based risk estimates. Review is recommended for flagged files.
================================================================================
```

**Other Commands:**

```bash
./sentinel

# Scan current directory
./sentinel scan --path .

# Scan with specific threshold
./sentinel scan --path ./src --threshold 0.8

# Scan specific language/s
./sentinel scan --path . --languages python,javascript

# Output as JSON
./sentinel scan --path . --json

# Fail build if high-risk patterns are detected above threshold
./sentinel scan --path . --fail-on-detection --threshold 0.75

# Collect training data for Machine Learning Model (directory should be ready inside the project)
./sentinel scan --path ./examples/flagged --collect --label flagged
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
- `--fail-on-detection` - Exit with error code if a risky code is detected
- `--exclude` - Paths to exclude from scanning
- `--verbose` - Verbose output with detailed signals
- `--json` - Output results in JSON format
- `--collect` - Save scan results to local training CSV
- `--label` - Label for training data: 'flagged' or 'human'
- `--model` - Path to ONNX model file
- `--no-ml` - Disable ML inference (heuristics only)
- `--ml-only` - Use ML only (fail if model not available)
- `--ml-weight` - Weight given to ML score (0.0-1.0)
- `--git-diff` - Scan only files changed against the specified branch


# Build & Installation

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
```

### Installation (optional)

**Linux/macOS:**
```bash
sudo mv sentinel /usr/local/bin/
```

**Windows (PowerShell as Administrator):**
```powershell
# Create a tools directory and add to PATH (one-time setup)
New-Item -ItemType Directory -Force -Path "C:\tools"
Move-Item sentinel.exe "C:\tools\sentinel.exe"
[Environment]::SetEnvironmentVariable("Path", $env:Path + ";C:\tools", "Machine")
```

**To uninstall (Windows):**
```powershell
Remove-Item "C:\tools\sentinel.exe"
```
```bash
# Linux/macOS
sudo rm /usr/local/bin/sentinel
```

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

# Project Configuration

Sentinel supports a `.sentinel.yaml` config file for persistent settings, so you don't have to pass flags every time.

**Setup:**
```bash
cp .sentinel.example.yaml .sentinel.yaml
```

Sentinel automatically loads `.sentinel.yaml` from the current directory or `$HOME`.

**Example `.sentinel.yaml`:**

```yaml
# Detection threshold (0.0 - 1.0)
# Files scoring above this value are flagged for review.
threshold: 0.70

# Languages to scan (leave empty to scan all supported languages)
languages:
  - python
  - java
  - javascript
  - typescript
  - go

# Paths to exclude from scanning (glob patterns supported)
exclude:
  - "vendor/*"
  - "node_modules/*"
  - ".git/*"
  - "*.min.js"
  - "*.test.js"
  - "*.spec.py"
  - "*_test.go"
  - "dist/*"
  - "build/*"

# Verbose output (show detailed signal breakdown per file)
verbose: false

# JSON output format (useful for CI/CD pipelines)
json: false

# Fail build if any file exceeds the threshold
fail_on_detection: true
```

> CLI flags always take precedence over config file values.

---

# License

MIT License - see [LICENSE](LICENSE) file for details.
