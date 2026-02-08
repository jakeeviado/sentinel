package detector

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"sentinel/pkg/analyzer"
	"sentinel/pkg/models"
)

var languageIDMap = map[string]string{
	"unknown":    "0",
	"python":     "1",
	"java":       "2",
	"javascript": "3",
	"typescript": "4",
	"go":         "5",
	"rust":       "6",
	"cpp":        "7",
	"c":          "8",
	"ruby":       "9",
	"php":        "10",
	"csharp":     "11",
	"kotlin":     "12",
	"swift":      "13",
}

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

	info, err := os.Stat(path)
	if err != nil || info.Size() == 0 || info.IsDir() {
		result.Error = os.ErrInvalid
		return result
	}

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

	maxScore := 0.0
	for _, signal := range signals {
		if signal.Score > maxScore {
			maxScore = signal.Score
		}
	}
	return maxScore
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

func (d *Detector) LogTrainingData(results *ScanResults, isAI bool) error {
	file, err := os.OpenFile("sentinel_training.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	labelStr := "0"
	if isAI {
		labelStr = "1"
	}

	for _, f := range results.Files {
		if f.Error != nil || len(f.Signals) == 0 {
			continue
		}

		langID := languageIDMap[f.Language]
		if langID == "" {
			langID = "0"
		}

		row := []string{langID}
		for _, s := range f.Signals {
			row = append(row, strconv.FormatFloat(s.Score, 'f', 4, 64))
		}
		row = append(row, labelStr)

		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}
