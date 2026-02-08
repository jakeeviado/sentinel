# ⌀ Sentinel Project Structure

## Key Components

### Entry Point
- **main.go**: Initializes and executes the CLI

### Commands (cmd/)
- **root.go**: Defines the root `sentinel` command, handles configuration loading
- **scan.go**: Implements the `scan` subcommand with file collection and orchestration

### Core Logic (pkg/)

#### detector/
- **detector.go**: Main detection engine
  - Concurrent file scanning
  - Score calculation
  - Result aggregation

#### analyzer/
- **analyzer.go**: Implements detection signals
  - Comment density analysis
  - Generic naming detection
  - Pattern repetition checking
  - Code complexity measurement
  - Formatting consistency
  - Boilerplate pattern matching

#### parser/
  - Language parsers (future: Tree-sitter) to be implemented

#### reporter/
- **reporter.go**: Formats and displays results
  - Text output with color coding
  - JSON output for CI/CD integration
  - Detailed signal reporting

## For Future Development Workflows

### Adding a New Detection Signal

1. **Define Signal** in `pkg/analyzer/analyzer.go`:
```go
func (a *Analyzer) checkNewSignal(code string) detector.Signal {
    // Implementation
    return detector.Signal{
        Name:        "new_signal",
        Score:       calculatedScore,
        Description: "What this detects",
        Evidence:    "Supporting data",
    }
}
```

2. **Add to Analysis** in `Analyze()` method:
```go
signals = append(signals, a.checkNewSignal(code))
```

3. **Write Tests** in `pkg/analyzer/analyzer_test.go`

### Adding Language Support

1. Add file extension to `detectLanguage()` in `pkg/detector/detector.go`
2. Add language-specific patterns to `analyzer.go`
3. (Future) Add Tree-sitter parser for that language

## Configuration

Configuration sources (in priority order):
1. Command-line flags
2. Environment variables
3. `.sentinel.yaml` in current directory
4. `~/.sentinel.yaml` in home directory
5. Default values

## License

MIT License - See LICENSE file
