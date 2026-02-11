// The package analyzer provides tools for 'heuristic analysis' of code to detect
// patterns commonly associated with AI-generated content.
//
// IMPROVING AND MAINTAINING THE HEURISTIC ANALYSIS
// When adding a New Detection Signal:
//
// sample:
//
// func (a *Analyzer) checkNewSignal(code string) detector.Signal {
//    // The implementation
//    return detector.Signal{
//        Name:        "new_signal",
//        Score:       calculatedScore,
//        Description: "What this detects",
//        Evidence:    "Supporting data",
//    }
// }
//
// Then add to Analyze method
//
// sample:
//
// signals = append(signals, a.checkNewSignal(code))

package analyzer

import (
	"fmt"
	"regexp"
	"strings"

	"sentinel/pkg/models"
)

type Analyzer struct{}

func New() *Analyzer {
	return &Analyzer{}
}

// Analyze runs the full suite of heuristic checks against the provided code.
// It returns a slice of 'Signals', each is representing a different detection metric.
func (a *Analyzer) Analyze(code string, language string) []models.Signal {
	signals := make([]models.Signal, 0)
	signals = append(signals, a.checkCommentDensity(code))
	signals = append(signals, a.checkGenericNaming(code))
	signals = append(signals, a.checkRepetitivePatterns(code))
	signals = append(signals, a.checkCodeComplexity(code))
	signals = append(signals, a.checkFormattingConsistency(code))
	signals = append(signals, a.checkBoilerplatePatterns(code, language))
	return signals
}

// This function calculates the ratio of comments to total lines.
// AI models tend to produce highly documented code, leading to higher density.
func (a *Analyzer) checkCommentDensity(code string) models.Signal {
	lines := strings.Split(code, "\n")
	if len(lines) == 0 {
		return models.Signal{Name: "comment_density", Score: 0.0}
	}

	commentLines := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "//") ||
			strings.HasPrefix(trimmed, "#") ||
			strings.HasPrefix(trimmed, "*") {
			commentLines++
		}
	}

	density := float64(commentLines) / float64(len(lines))

	var score float64
	var description string

	if density > 0.3 {
		score = 0.8
		description = "Excessive comment density detected (>30%)"
	} else if density > 0.2 {
		score = 0.5
		description = "High comment density (>20%)"
	} else {
		score = 0.0
		description = "Normal comment density"
	}

	return models.Signal{
		Name:        "comment_density",
		Score:       score,
		Description: description,
		Evidence:    formatEvidence("Comment lines: %d / %d (%.1f%%)", commentLines, len(lines), density*100),
	}
}

// This function identifies the frequency of placeholder variable names
// like 'temp', 'data', or 'obj', which are common in AI-generated snippets.
func (a *Analyzer) checkGenericNaming(code string) models.Signal {
	genericNames := []string{
		"temp", "tmp", "data", "result", "item", "value", "obj", "elem",
		"variable", "parameter", "foo", "bar", "test", "example",
	}

	matches := 0
	for _, name := range genericNames {
		pattern := `\b` + name + `\b`
		re := regexp.MustCompile(pattern)
		matches += len(re.FindAllString(strings.ToLower(code), -1))
	}

	lines := len(strings.Split(code, "\n"))
	if lines == 0 {
		return models.Signal{Name: "generic_naming", Score: 0.0}
	}

	density := float64(matches) / float64(lines)

	var score float64
	var description string

	if density > 0.5 {
		score = 0.9
		description = "Very high use of generic variable names"
	} else if density > 0.2 {
		score = 0.6
		description = "High use of generic variable names"
	} else {
		score = 0.0
		description = "Normal variable naming patterns"
	}

	return models.Signal{
		Name:        "generic_naming",
		Score:       score,
		Description: description,
		Evidence:    formatEvidence("Generic name occurrences: %d", matches),
	}
}

// This function looks for identical lines of code.
// High repetition can be a sign of limited variation in generated output.
func (a *Analyzer) checkRepetitivePatterns(code string) models.Signal {
	lines := strings.Split(code, "\n")
	if len(lines) < 10 {
		return models.Signal{Name: "repetitive_patterns", Score: 0.0}
	}

	lineCount := make(map[string]int)
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 10 {
			lineCount[trimmed]++
		}
	}

	repetitions := 0
	for _, count := range lineCount {
		if count > 2 {
			repetitions += count - 1
		}
	}

	repetitionRatio := float64(repetitions) / float64(len(lines))

	var score float64
	var description string

	if repetitionRatio > 0.3 {
		score = 0.7
		description = "High code repetition detected"
	} else if repetitionRatio > 0.15 {
		score = 0.4
		description = "Moderate code repetition"
	} else {
		score = 0.0
		description = "Normal code variation"
	}

	return models.Signal{
		Name:        "repetitive_patterns",
		Score:       score,
		Description: description,
		Evidence:    formatEvidence("Repetition ratio: %.1f%%", repetitionRatio*100),
	}
}

// This function measures control flow density.
// AI code often follows simpler, more linear paths compared to complex human logic.
func (a *Analyzer) checkCodeComplexity(code string) models.Signal {
	controlFlowKeywords := []string{
		"if", "else", "switch", "case", "for", "while", "do",
		"try", "catch", "finally", "throw",
	}

	lines := strings.Split(code, "\n")
	controlFlowCount := 0

	for _, line := range lines {
		lower := strings.ToLower(line)
		for _, keyword := range controlFlowKeywords {
			if strings.Contains(lower, keyword) {
				controlFlowCount++
				break
			}
		}
	}

	if len(lines) == 0 {
		return models.Signal{Name: "code_complexity", Score: 0.0}
	}

	complexity := float64(controlFlowCount) / float64(len(lines))

	var score float64
	var description string

	if complexity < 0.05 {
		score = 0.6
		description = "Unusually low cyclomatic complexity"
	} else if complexity < 0.1 {
		score = 0.3
		description = "Low cyclomatic complexity"
	} else {
		score = 0.0
		description = "Normal code complexity"
	}

	return models.Signal{
		Name:        "code_complexity",
		Score:       score,
		Description: description,
		Evidence:    formatEvidence("Control flow density: %.2f", complexity),
	}
}

// This fucntion evaluates how strictly indentation is followed.
// Perfect or near-perfect consistency is often a hallmark of programmatic generation.
func (a *Analyzer) checkFormattingConsistency(code string) models.Signal {
	lines := strings.Split(code, "\n")
	if len(lines) < 5 {
		return models.Signal{Name: "formatting_consistency", Score: 0.0}
	}

	indentations := make(map[int]int)
	for _, line := range lines {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}

		indent := 0
		for _, char := range line {
			if char == ' ' || char == '\t' {
				indent++
			} else {
				break
			}
		}
		indentations[indent]++
	}

	uniqueIndents := len(indentations)

	var score float64
	var description string

	if uniqueIndents <= 3 && len(lines) > 20 {
		score = 0.5
		description = "Very consistent indentation (possibly AI-generated)"
	} else {
		score = 0.0
		description = "Normal formatting variation"
	}

	return models.Signal{
		Name:        "formatting_consistency",
		Score:       score,
		Description: description,
		Evidence:    formatEvidence("Unique indentation levels: %d", uniqueIndents),
	}
}

// This function searches for common AI signatures and standard templates
// based on the specific programming language.
// (BAD IMPLEMENTATION FOR NOW)
func (a *Analyzer) checkBoilerplatePatterns(code string, language string) models.Signal {
	boilerplatePatterns := map[string][]string{
		"python": {
			"def main():",
			"if __name__ == \"__main__\":",
			"# TODO:",
			"# Example usage:",
			"# Initialize",
		},
		"java": {
			"public static void main",
			"// TODO:",
			"// Example:",
			"@Override",
		},
		"javascript": {
			"// TODO:",
			"// Example:",
			"export default",
			"const config =",
		},
	}

	patterns, ok := boilerplatePatterns[language]
	if !ok {
		return models.Signal{Name: "boilerplate_patterns", Score: 0.0}
	}

	matches := 0
	for _, pattern := range patterns {
		if strings.Contains(code, pattern) {
			matches++
		}
	}

	var score float64
	var description string

	if matches >= 3 {
		score = 0.6
		description = "Multiple boilerplate patterns detected"
	} else if matches >= 2 {
		score = 0.3
		description = "Some boilerplate patterns present"
	} else {
		score = 0.0
		description = "Minimal boilerplate"
	}

	return models.Signal{
		Name:        "boilerplate_patterns",
		Score:       score,
		Description: description,
		Evidence:    formatEvidence("Boilerplate matches: %d", matches),
	}
}

func formatEvidence(format string, args ...interface{}) string {
	return fmt.Sprintf(format, args...)
}
