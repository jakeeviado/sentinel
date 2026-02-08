package detector

import (
	"os"
	"path/filepath"
	"sync"

	"sentinel/pkg/analyzer"
	"sentinel/pkg/models"
)

type Config struct {
	Threshold    float64
	Languages    []string
	ExcludePaths []string
	Verbose      bool
}

type Detector struct {
	config   Config
	analyzer *analyzer.Analyzer
}

type FileResult struct {
	Path     string
	Score    float64
	Signals  []models.Signal
	Language string
	Error    error
}

type ScanResults struct {
	Files         []FileResult
	TotalFiles    int
	DetectedFiles int
	AverageScore  float64
}

func New(config Config) *Detector {
	return &Detector{
		config:   config,
		analyzer: analyzer.New(),
	}
}

func (d *Detector) ScanFiles(files []string) (*ScanResults, error) {
	results := &ScanResults{
		Files:      make([]FileResult, 0, len(files)),
		TotalFiles: len(files),
	}

	resultsChan := make(chan FileResult, len(files))
	var wg sync.WaitGroup
	semaphore := make(chan struct{}, 10)

	for _, file := range files {
		wg.Add(1)
		go func(path string) {
			defer wg.Done()
			semaphore <- struct{}{}  
			defer func() { <-semaphore }()

			result := d.scanFile(path)
			resultsChan <- result
		}(file)
	}
	go func() {
		wg.Wait()
		close(resultsChan)
	}()

	totalScore := 0.0
	for result := range resultsChan {
		results.Files = append(results.Files, result)
		totalScore += result.Score

		if result.Score >= d.config.Threshold {
			results.DetectedFiles++
		}
	}

	if len(results.Files) > 0 {
		results.AverageScore = totalScore / float64(len(results.Files))
	}

	return results, nil
}

func (d *Detector) scanFile(path string) FileResult {
	result := FileResult{
		Path:    path,
		Signals: make([]models.Signal, 0),
	}

	result.Language = detectLanguage(path)

	content, err := os.ReadFile(path)
	if err != nil {
		result.Error = err
		return result
	}

	signals := d.analyzer.Analyze(string(content), result.Language)
	result.Signals = signals

	result.Score = d.calculateScore(signals)

	return result
}

func (d *Detector) calculateScore(signals []models.Signal) float64 {
	if len(signals) == 0 {
		return 0.0
	}

	totalWeight := 0.0
	weightedSum := 0.0

	for _, signal := range signals {
		if signal.Score > 0.0 {
			weight := 1.0
			weightedSum += signal.Score * weight
			totalWeight += weight
		}
	}

	if totalWeight == 0 {
		return 0.0
	}

	return weightedSum / totalWeight
}

func (r *ScanResults) HasDetections(threshold float64) bool {
	for _, file := range r.Files {
		if file.Score >= threshold {
			return true
		}
	}
	return false
}

func detectLanguage(path string) string {
	ext := filepath.Ext(path)
	languageMap := map[string]string{
		".py":    "python",
		".java":  "java",
		".js":    "javascript",
		".ts":    "typescript",
		".go":    "go",
		".rs":    "rust",
		".cpp":   "cpp",
		".c":     "c",
		".h":     "c",
		".rb":    "ruby",
		".php":   "php",
		".cs":    "csharp",
		".kt":    "kotlin",
		".swift": "swift",
	}

	if lang, ok := languageMap[ext]; ok {
		return lang
	}
	return "unknown"
}
