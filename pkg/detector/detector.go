package detector

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"sync"

	"sentinel/pkg/analyzer"
	"sentinel/pkg/ml"
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
	// MACHINE LEARNING CONFIG
	ModelPath string
	UseML     bool
	MLWeight  float64
	MLOnly    bool
}

type Detector struct {
	config     Config
	analyzer   *analyzer.Analyzer
	mlDetector *ml.MLDetector
}

type FileResult struct {
	Path     string
	Score    float64
	Signals  []models.Signal
	Language string
	Error    error
	// MACHINE LEARNING SPECIFIC FIELDS
	MLScore        float64
	HeuristicScore float64
	UsedML         bool
	Confidence     float64
}

type ScanResults struct {
	Files         []FileResult
	TotalFiles    int
	DetectedFiles int
	AverageScore  float64
	UsedML        bool
}

func New(config Config) *Detector {
	detector := &Detector{
		config:   config,
		analyzer: analyzer.New(),
	}
	if config.UseML && config.ModelPath != "" {
		mlConfig := ml.DetectionConfig{
			ModelPath:     config.ModelPath,
			MLWeight:      config.MLWeight,
			FallbackMode:  !config.MLOnly,
			RequireML:     config.MLOnly,
			MinConfidence: config.Threshold,
			NumFeatures:   25,
		}
		mlDetector, err := ml.NewMLDetector(mlConfig)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: ML initialization failed: %v\n", err)
			fmt.Fprintf(os.Stderr, "Falling back to heuristics-only mode\n")
		} else {
			detector.mlDetector = mlDetector
			if config.Verbose {
				fmt.Fprintf(os.Stderr, "ML model loaded: %s\n", config.ModelPath)
			}
		}
	}
	return detector
}

func (d *Detector) ScanFiles(files []string) (*ScanResults, error) {
	results := &ScanResults{
		Files:      make([]FileResult, 0, len(files)),
		TotalFiles: len(files),
		UsedML:     d.mlDetector != nil,
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
		result.Error = fmt.Errorf("invalid file: %w", err)
		return result
	}

	content, err := os.ReadFile(path)
	if err != nil {
		result.Error = err
		return result
	}

	codeString := string(content)

	// Run analysis (get heuristic signals)
	signals := d.analyzer.Analyze(codeString, result.Language)
	result.Signals = signals

	// Calculate heuristic score (using maximum signal score)
	heuristicScore := calculateMaxScore(signals)
	result.HeuristicScore = heuristicScore

	// Use ML detection if available
	if d.mlDetector != nil {
		mlResult, err := d.mlDetector.Detect(signals, result.Language, codeString)
		if err != nil {
			// ML failed, use heuristics only
			result.Score = heuristicScore
			result.MLScore = 0.0
			result.UsedML = false
		} else {
			// ML succeeded
			result.Score = mlResult.FinalScore
			result.MLScore = mlResult.MLScore
			result.HeuristicScore = mlResult.HeuristicScore
			result.UsedML = mlResult.UsedML
			result.Confidence = mlResult.Confidence
		}
	} else {
		result.Score = heuristicScore
		result.MLScore = 0.0
		result.UsedML = false
	}

	return result
}

func calculateMaxScore(signals []models.Signal) float64 {
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

/*
 * LOG TRAINING DATA
 * Saves scan results to CSV for ML training
 */
func (d *Detector) LogTrainingData(results *ScanResults, isAI bool) error {
	file, err := os.OpenFile("sentinel_training.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	info, _ := file.Stat()
	numFeatures := 25
	if info.Size() == 0 {
		header := []string{"language_id"}
		for i := 1; i <= numFeatures; i++ {
			header = append(header, fmt.Sprintf("f%d", i))
		}
		header = append(header, "label")
		writer.Write(header)
	}

	labelStr := "0"
	if isAI {
		labelStr = "1"
	}

	for _, f := range results.Files {
		if f.Error != nil {
			continue
		}

		langID := languageIDMap[f.Language]
		if langID == "" {
			langID = "0"
		}

		features := make([]float64, numFeatures)
		for i, s := range f.Signals {
			if i < numFeatures {
				features[i] = s.Score
			}
		}

		row := []string{langID}
		for _, val := range features {
			row = append(row, strconv.FormatFloat(val, 'f', 4, 64))
		}
		row = append(row, labelStr)

		writer.Write(row)
	}
	return nil
}
