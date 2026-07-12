## Sentinel

Multi-language static analysis CLI tool that identifies risky, inconsistent, or unconventional code patterns. Designed to assist code reviewers and maintain quality in AI-assisted development environments.

**Features**
* **Multi-Language Support**: Python, Java, JS/TS, Go, Rust, C/C++, Ruby, PHP, C#, Kotlin, Swift.
* **Hybrid Detection**: Combines heuristic analysis with ONNX-based ML models.
* **CI/CD Ready**: Native integration with GitHub Actions, GitLab CI, and Jenkins.
* **Workflow Integration**: Supports Git Diff scanning and automated build gating (`--fail-on-detection`).
* **Extensible**: Configurable thresholds via CLI flags or `.sentinel.yaml`.
* **Training Loop**: Built-in tools to collect and label scan results for model retraining.
---
**Risk Signals** - Evaluates code against the following metrics (scored 0.0–1.0):

| Signal                   | Description                                      |
| :----------------------- | :----------------------------------------------- |
| `comment_density`        | Comment-to-code ratio.                           |
| `generic_naming`         | Use of placeholder names (e.g., `temp`, `data`). |
| `repetitive_patterns`    | File-level code duplication.                     |
| `code_complexity`        | Control flow density.                            |
| `formatting_consistency` | Indentation/style uniformity.                    |
| `comment_redundancy`     | Comments that mirror code logic.                 |
| `emoji_sentiment`        | Presence of non-standard emojis.                 |
| `identifier_order`       | Suspiciously perfect alphabetical ordering.      |
| `defensive_ratio`        | Ratio of error handling to functional logic.     |

---

**Sample CLI Output**
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
FAILED: 1 file(s) require attention (risk threshold exceeded)
Note: Scores represent heuristic and/or ML-based risk estimates. Review is recommended for flagged files.
================================================================================
```

**Available Commands:**

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

---

#### CI/CD Integration with Github Actions, GitLab CI, and Jenkins

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

### Disclaimer
Sentinel is an aid for code quality and review; it is not intended to determine authorship or replace human judgment. Results are probabilistic and may include false positives.

### License
MIT License - see [LICENSE](LICENSE) file for details.
