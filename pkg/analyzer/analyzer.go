package analyzer

import (
	"fmt"
	"regexp"
	"strings"

	"sentinel/pkg/models"
)

type Analyzer struct {
	// CONFIGURATION FOR ANALYSIS
}

func New() *Analyzer {
	return &Analyzer{}
}

func (a *Analyzer) Analyze(code string, language string) []models.Signal {
	signals := make([]models.Signal, 0)

	// RUN HEURISTIC CHCKS
	signals = append(signals, a.checkCommentDensity(code))
	signals = append(signals, a.checkGenericNaming(code))
	signals = append(signals, a.checkRepetitivePatterns(code))
	signals = append(signals, a.checkCodeComplexity(code))
	signals = append(signals, a.checkFormattingConsistency(code))
	signals = append(signals, a.checkBoilerplatePatterns(code, language))

	return signals
}

/*
 * Check for excessive comments
 * (AI often over-comments)
 */
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

/*
 * Check for generic variable names
 * (common in AI code)
 */
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

/*
 * Checks repetitive code patterns
 */
func (a *Analyzer) checkRepetitivePatterns(code string) models.Signal {
	lines := strings.Split(code, "\n")
	if len(lines) < 10 {
		return models.Signal{Name: "repetitive_patterns", Score: 0.0}
	}

	// CHECK FOR REPEATED LINES
	lineCount := make(map[string]int)
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		// SHORT LINES ARE IGNORED
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

/*
 * Checks the code complexity
 * (AI code tends to be simpler but shittier)
 */
func (a *Analyzer) checkCodeComplexity(code string) models.Signal {
	// COUNTS CONTROL FLOW STATEMENTS
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

	// LOW COMPLEXITY MIGHT INDICATE AN AI-GEN CODE
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

/*
 * Check formatting consistency
 * (AI is very consistent)
 */
func (a *Analyzer) checkFormattingConsistency(code string) models.Signal {
	lines := strings.Split(code, "\n")
	if len(lines) < 5 {
		return models.Signal{Name: "formatting_consistency", Score: 0.0}
	}

	// CHECK INDENTAION CONSISTENCY
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

	// PERFECT CONSISTENCY MIGHT INDICATE AI-GEN CODE
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

/*
 * Checks for the most common AI boilerplate patterns (more improvements soon)
 */
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
