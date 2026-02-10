package ml

import (
	"math"
	"sentinel/pkg/models"
)

type FeatureVector struct {
	Features []float32
	Names    []string
}

type FeatureExtractor struct {
	languageIDMap map[string]int
}

func NewFeatureExtractor() *FeatureExtractor {
	return &FeatureExtractor{
		languageIDMap: map[string]int{
			"unknown":    0,
			"python":     1,
			"java":       2,
			"javascript": 3,
			"typescript": 4,
			"go":         5,
			"rust":       6,
			"cpp":        7,
			"c":          8,
			"ruby":       9,
			"php":        10,
			"csharp":     11,
			"kotlin":     12,
			"swift":      13,
		},
	}
}

func (fe *FeatureExtractor) Extract(signals []models.Signal, language string, codeMetrics CodeMetrics) *FeatureVector {
	features := make([]float32, 0, 50)
	names := make([]string, 0, 50)
	langID := fe.languageIDMap[language]
	features = append(features, float32(langID))
	names = append(names, "language_id")
	signalMap := make(map[string]float64)

	for _, sig := range signals {
		signalMap[sig.Name] = sig.Score
	}

	signalNames := []string{
		"comment_density",
		"generic_naming",
		"repetitive_patterns",
		"code_complexity",
		"formatting_consistency",
		"boilerplate_patterns",
		"ast_comment_density",
		"ast_generic_naming",
		"function_uniformity",
		"trivial_functions",
		"structural_repetition",
		"over_documentation",
	}

	for _, name := range signalNames {
		score := signalMap[name]
		features = append(features, float32(score))
		names = append(names, name)
	}

	features = append(features, float32(codeMetrics.LineCount))
	names = append(names, "line_count")

	features = append(features, float32(codeMetrics.CharCount))
	names = append(names, "char_count")

	features = append(features, float32(codeMetrics.AvgLineLength))
	names = append(names, "avg_line_length")

	features = append(features, float32(codeMetrics.EmptyLineRatio))
	names = append(names, "empty_line_ratio")

	signalScores := make([]float64, 0, len(signals))
	for _, sig := range signals {
		signalScores = append(signalScores, sig.Score)
	}

	stats := calculateStats(signalScores)
	features = append(features, float32(stats.Mean))
	names = append(names, "signal_mean")

	features = append(features, float32(stats.StdDev))
	names = append(names, "signal_stddev")

	features = append(features, float32(stats.Max))
	names = append(names, "signal_max")

	features = append(features, float32(stats.Min))
	names = append(names, "signal_min")

	features = append(features, float32(stats.Median))
	names = append(names, "signal_median")

	activeSignals := 0
	for _, sig := range signals {
		if sig.Score > 0.0 {
			activeSignals++
		}
	}
	features = append(features, float32(activeSignals))
	names = append(names, "active_signal_count")

	return &FeatureVector{
		Features: features,
		Names:    names,
	}
}

type CodeMetrics struct {
	LineCount      int
	CharCount      int
	AvgLineLength  float64
	EmptyLineRatio float64
}

type Stats struct {
	Mean   float64
	StdDev float64
	Max    float64
	Min    float64
	Median float64
}

func calculateStats(scores []float64) Stats {
	if len(scores) == 0 {
		return Stats{}
	}

	sum := 0.0
	for _, s := range scores {
		sum += s
	}
	mean := sum / float64(len(scores))

	variance := 0.0
	for _, s := range scores {
		diff := s - mean
		variance += diff * diff
	}
	variance /= float64(len(scores))
	stdDev := math.Sqrt(variance)

	max := scores[0]
	min := scores[0]
	for _, s := range scores {
		if s > max {
			max = s
		}
		if s < min {
			min = s
		}
	}

	sorted := make([]float64, len(scores))
	copy(sorted, scores)
	for i := 0; i < len(sorted); i++ {
		for j := i + 1; j < len(sorted); j++ {
			if sorted[i] > sorted[j] {
				sorted[i], sorted[j] = sorted[j], sorted[i]
			}
		}
	}
	median := sorted[len(sorted)/2]

	return Stats{
		Mean:   mean,
		StdDev: stdDev,
		Max:    max,
		Min:    min,
		Median: median,
	}
}

func ExtractCodeMetrics(code string) CodeMetrics {
	lines := splitLines(code)

	totalChars := len(code)
	lineCount := len(lines)

	if lineCount == 0 {
		return CodeMetrics{}
	}

	avgLineLength := float64(totalChars) / float64(lineCount)

	emptyLines := 0
	for _, line := range lines {
		if len(trimSpace(line)) == 0 {
			emptyLines++
		}
	}
	emptyLineRatio := float64(emptyLines) / float64(lineCount)

	return CodeMetrics{
		LineCount:      lineCount,
		CharCount:      totalChars,
		AvgLineLength:  avgLineLength,
		EmptyLineRatio: emptyLineRatio,
	}
}

func splitLines(s string) []string {
	lines := make([]string, 0)
	start := 0
	for i := 0; i < len(s); i++ {
		if s[i] == '\n' {
			lines = append(lines, s[start:i])
			start = i + 1
		}
	}
	if start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}

func trimSpace(s string) string {
	start := 0
	end := len(s)

	for start < end && isSpace(s[start]) {
		start++
	}
	for end > start && isSpace(s[end-1]) {
		end--
	}

	return s[start:end]
}

func isSpace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\n' || c == '\r'
}
