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

#### reporter/
- **reporter.go**: Formats and displays results
  - Text output with color coding
  - JSON output for CI/CD integration
  - Detailed signal reporting
