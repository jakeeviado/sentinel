# ⌀ Sentinel Quick Start Guide

Get up and running with Sentinel in 5 minutes!

## Installation

### Option 1: Download Pre-built Binary (Recommended)

```bash
# Linux
curl -L https://github.com/jakeeviado/sentinel/releases/latest/download/sentinel-linux-amd64 -o sentinel
chmod +x sentinel
sudo mv sentinel /usr/local/bin/

# macOS
curl -L https://github.com/jakeeviado/sentinel/releases/latest/download/sentinel-darwin-amd64 -o sentinel
chmod +x sentinel
sudo mv sentinel /usr/local/bin/

# Windows (PowerShell)
Invoke-WebRequest -Uri "https://github.com/jakeeviado/sentinel/releases/latest/download/sentinel-windows-amd64.exe" -OutFile "sentinel.exe"
```

### Option 2: Build from Source

```bash
git clone https://github.com/jakeeviado/sentinel.git
cd sentinel
make build
sudo cp bin/sentinel /usr/local/bin/
```

## First Scan

Scan your current project:

```bash
cd /path/to/your/project
sentinel scan --path .
```

You'll see output like:

```
================================================================================
                        ⌀ SENTINEL - Code Detection Report
================================================================================

Total Files Scanned:  42
Files Detected:       3
Average Score:        0.45
Detection Threshold:  0.70
```

## Common Use Cases

### 1. Scan Specific Directory

```bash
sentinel scan --path ./src
```

### 2. Scan Only Certain Languages

```bash
sentinel scan --path . --languages python,javascript
```

### 3. Adjust Detection Threshold

```bash
# More strict (fewer false positives)
sentinel scan --path . --threshold 0.85

# Less strict (catch more potential AI code)
sentinel scan --path . --threshold 0.60
```

### 4. Get Detailed Analysis

```bash
sentinel scan --path . --verbose
```

This shows why each file was flagged:

```
   src/helper.py
   Score: 0.82 | Language: python
   Signals:
     • comment_density (0.80): Excessive comment density detected (>30%)
       Comment lines: 45 / 120 (37.5%)
     • generic_naming (0.90): Very high use of generic variable names
       Generic name occurrences: 23
```

### 5. Output as JSON

Perfect for integrating with other tools:

```bash
sentinel scan --path . --json > results.json
```

## CI/CD Integration

### GitHub Actions

Add to `.github/workflows/sentinel.yml`:

```yaml
name: AI Code Check
on: [pull_request]

jobs:
  sentinel:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      
      - name: Download Sentinel
        run: |
          curl -L https://github.com/jakeeviado/sentinel/releases/latest/download/sentinel-linux-amd64 -o sentinel
          chmod +x sentinel
      
      - name: Scan Code
        run: ./sentinel scan --path . --threshold 0.75 --fail-on-detection
```

### GitLab CI

Add to `.gitlab-ci.yml`:

```yaml
sentinel:
  stage: test
  script:
    - curl -L https://github.com/jakeeviado/sentinel/releases/latest/download/sentinel-linux-amd64 -o sentinel
    - chmod +x sentinel
    - ./sentinel scan --path . --threshold 0.75 --fail-on-detection
```

## Configuration File

Create `.sentinel.yaml` in your project root:

```yaml
threshold: 0.75
languages:
  - python
  - javascript
  - go
exclude:
  - "vendor/*"
  - "node_modules/*"
  - "*.test.js"
verbose: true
fail_on_detection: true
```

Now just run:

```bash
sentinel scan --path .
```

## Understanding Scores

- **0.0 - 0.4**: Likely human-written
- **0.4 - 0.6**: Uncertain, manual review recommended
- **0.6 - 0.8**: Suspicious, likely AI-assisted
- **0.8 - 1.0**: Very likely AI-generated

## Detection Signals Explained

Sentinel looks for these patterns:

1. **Comment Density** - AI code often has excessive comments
2. **Generic Naming** - Overuse of names like `temp`, `data`, `result`
3. **Repetitive Patterns** - Unusual code repetition
4. **Low Complexity** - Overly simplified logic
5. **Perfect Formatting** - Too-consistent formatting
6. **Boilerplate** - Common AI-generated patterns

## Tips for Best Results

### DO:
- Set appropriate thresholds for your codebase
- Use `--verbose` to understand why code was flagged
- Exclude generated code and dependencies
- Review results manually - Sentinel is a detection aid

### DON'T:
- Rely solely on automated detection
- Use extremely low thresholds (<0.5)
- Forget to exclude test files and mocks
- Ignore context - some patterns are legitimate

## Common Issues

### "Too many false positives"

Increase threshold:
```bash
sentinel scan --threshold 0.85
```

### "Not detecting obvious AI code"

Lower threshold and check specific files:
```bash
sentinel scan --threshold 0.60 --verbose
```

### "Scanning takes too long"

Exclude unnecessary directories:
```bash
sentinel scan --exclude "vendor/*,node_modules/*,dist/*"
```

Hell Yeah!
