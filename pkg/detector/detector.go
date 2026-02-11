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

type DetectorConfiguration struct {
	Threshold    float64
	Languages    []string
	ExcludePaths []string
	IsVerbose    bool
	ModelPath    string
	UseML        bool
	MLWeight     float64
	IsMLOnly     bool
}

type Detector struct {
	detectorConfig DetectorConfiguration
	analyzer       *analyzer.Analyzer
	mlDetector     *ml.MLDetector
}

type FileResult struct {
	Path           string
	Score          float64
	Signals        []models.Signal
	Language       string
	Error          error
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

func New(config DetectorConfiguration) *Detector {
	detector := &Detector{
		detectorConfig: config,
		analyzer:       analyzer.New(),
	}

	if config.UseML && config.ModelPath != "" {
		mlConfig := ml.DetectionConfig{
			ModelPath:     config.ModelPath,
			MLWeight:      config.MLWeight,
			FallbackMode:  !config.IsMLOnly,
			RequireML:     config.IsMLOnly,
			MinConfidence: config.Threshold,
			NumFeatures:   ml.DefaultNumFeatures,
		}

		machineLearningDetector, err := ml.NewMLDetector(mlConfig)

		if err != nil {
			fmt.Fprintf(os.Stderr, "Warning: ML initialization failed: %v\n", err)
			fmt.Fprintf(os.Stderr, "Falling back to heuristics-only mode\n")
		} else {
			detector.mlDetector = machineLearningDetector

			if config.IsVerbose {
				fmt.Fprintf(os.Stderr, "ML model loaded: %s\n", config.ModelPath)
			}
		}
	}
	return detector
}

func (d *Detector) ScanFiles(files []string) (*ScanResults, error) {
	report := &ScanResults{
		Files:      make([]FileResult, 0, len(files)),
		TotalFiles: len(files),
		UsedML:     d.mlDetector != nil,
	}

	resultsChannel := make(chan FileResult, len(files))
	var waitGroup sync.WaitGroup
	semaphore := make(chan struct{}, 10)

	for _, file := range files {
		waitGroup.Add(1)
		go func(path string) {
			defer waitGroup.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()
			result := d.scanFile(path)
			resultsChannel <- result
		}(file)
	}

	go func() {
		waitGroup.Wait()
		close(resultsChannel)
	}()

	totalScore := 0.0
	for result := range resultsChannel {
		report.Files = append(report.Files, result)
		totalScore += result.Score

		if result.Score >= d.detectorConfig.Threshold {
			report.DetectedFiles++
		}
	}

	if len(report.Files) > 0 {
		report.AverageScore = totalScore / float64(len(report.Files))
	}

	return report, nil
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

	fileContent, err := os.ReadFile(path)
	if err != nil {
		result.Error = err
		return result
	}

	codeString := string(fileContent)

	heuristicSignals := d.analyzer.Analyze(codeString, result.Language)
	result.Signals = heuristicSignals
	hScore := calculateHeuristicMaxScore(heuristicSignals)
	result.HeuristicScore = hScore

	if d.mlDetector != nil {
		mlResult, err := d.mlDetector.Detect(heuristicSignals, result.Language, codeString)

		if err != nil {
			if d.detectorConfig.IsMLOnly {
				result.Error = fmt.Errorf("ML inference failed: %w", err)
				return result
			}
			result.Score = hScore
			result.UsedML = false
		} else {
			result.Score = mlResult.FinalScore
			result.MLScore = mlResult.MLScore
			result.UsedML = mlResult.UsedML
			result.Confidence = mlResult.Confidence
		}
	} else {
		result.Score = hScore
		result.UsedML = false
	}

	return result
}

func calculateHeuristicMaxScore(signals []models.Signal) float64 {
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

// LogTrainingData saves scan results to a CSV file for Machine Learning training purposes.
// It appends to "sentinel_training.csv" and creates the file with a header if it doesn't exist.
func (d *Detector) LogTrainingData(results *ScanResults, isAI bool) error {
	trainingDataFile, err := os.OpenFile("sentinel_training.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}

	defer func() {
		if err := trainingDataFile.Close(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to close training data file: %v\n", err)
		}
	}()

	writer := csv.NewWriter(trainingDataFile)
	defer writer.Flush()

	fe := ml.NewFeatureExtractor()
	info, err := trainingDataFile.Stat()
	if err != nil {
		return fmt.Errorf("failed to stat training file: %w", err)
	}

	if info.Size() == 0 {
		dummyVector := fe.Extract([]models.Signal{}, "unknown", ml.CodeMetrics{})
		header := append(dummyVector.Names, "label")
		if err := writer.Write(header); err != nil {
			return err
		}
	}

	labelStr := "0"
	if isAI {
		labelStr = "1"
	}

	for _, f := range results.Files {
		if f.Error != nil {
			continue
		}
		content, err := os.ReadFile(f.Path)
		if err != nil {
			continue
		}
		metrics := ml.ExtractCodeMetrics(string(content))
		vector := fe.Extract(f.Signals, f.Language, metrics)
		row := make([]string, 0, len(vector.Features)+1)
		for _, val := range vector.Features {
			row = append(row, strconv.FormatFloat(float64(val), 'f', 4, 32))
		}
		row = append(row, labelStr)

		if err := writer.Write(row); err != nil {
			return err
		}
	}
	return nil
}
