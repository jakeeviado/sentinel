package ml

import (
	"fmt"
	"math"
	"sentinel/pkg/models"
)

const (
	// DefaultNumFeatures is the total count of features in the vector.
	// 1 (Lang) + 10 (Heuristics) + 4 (Structural) + 6 (Stats) = 21
	DefaultNumFeatures = 21

	// DefaultMLWeight represents the influence of the ML score in the final calculation.
	// A value of 0.7 means the final score is 70% ML and 30% heuristic signals.
	DefaultMLWeight = 0.7

	// DefaultMinConfidence is the threshold used to determine if a file is
	// flagged as AI-generated when ML results are available.
	DefaultMinConfidence = 0.5
)

type MLDetector struct {
	onnxDetector     *ONNXDetector
	featureExtractor *FeatureExtractor
	config           DetectionConfig
}

type DetectionConfig struct {
	ModelPath     string
	MLWeight      float64
	FallbackMode  bool
	RequireML     bool
	MinConfidence float64
	NumFeatures   int
}

func NewMLDetector(config DetectionConfig) (*MLDetector, error) {
	if config.MLWeight == 0.0 {
		config.MLWeight = DefaultMLWeight
	}
	if config.MinConfidence == 0.0 {
		config.MinConfidence = DefaultMinConfidence
	}
	if config.NumFeatures == 0 {
		config.NumFeatures = DefaultNumFeatures
	}

	featureExtractor := NewFeatureExtractor()

	var onnxDetector *ONNXDetector
	var err error

	if config.ModelPath != "" {
		modelConfig := ModelConfig{
			ModelPath:      config.ModelPath,
			InputName:      "input",
			OutputName:     "output",
			NumFeatures:    config.NumFeatures,
			OptimizeMemory: true,
		}

		onnxDetector, err = NewONNXDetector(modelConfig)
		if err != nil {
			if config.RequireML {
				return nil, fmt.Errorf("ML model required but failed to load: %w", err)
			}
			onnxDetector = &ONNXDetector{initialized: false}
		}
	} else {
		onnxDetector = &ONNXDetector{initialized: false}
	}

	return &MLDetector{
		onnxDetector:     onnxDetector,
		featureExtractor: featureExtractor,
		config:           config,
	}, nil
}

type DetectionResult struct {
	FinalScore     float64
	HeuristicScore float64
	MLScore        float64
	UsedML         bool
	Confidence     float64
	IsAIGenerated  bool
}

func (d *MLDetector) Detect(signals []models.Signal, language string, code string) (*DetectionResult, error) {
	result := &DetectionResult{
		UsedML: false,
	}

	maxHeuristic := 0.0
	for _, sig := range signals {
		if sig.Score > maxHeuristic {
			maxHeuristic = sig.Score
		}
	}
	result.HeuristicScore = maxHeuristic

	if d.onnxDetector != nil && d.onnxDetector.IsInitialized() {
		metrics := ExtractCodeMetrics(code)
		vector := d.featureExtractor.Extract(signals, language, metrics)

		features := padFeatures(vector.Features, d.config.NumFeatures)

		mlProb, err := d.onnxDetector.Predict(features)
		if err != nil {
			if !d.config.FallbackMode {
				return nil, fmt.Errorf("ML inference failed: %w", err)
			}
			result.MLScore = 0.0
			result.FinalScore = maxHeuristic
		} else {
			result.MLScore = float64(mlProb)
			result.UsedML = true
			mlWeight := d.config.MLWeight
			heuristicWeight := 1.0 - mlWeight
			result.FinalScore = (mlWeight * result.MLScore) + (heuristicWeight * maxHeuristic)
		}
	} else {
		result.FinalScore = maxHeuristic
	}

	if result.UsedML {
		agreement := 1.0 - math.Abs(result.MLScore-result.HeuristicScore)
		clarity := math.Abs(result.FinalScore-0.5) * 2.0
		result.Confidence = (agreement * 0.6) + (clarity * 0.4)
	} else {
		result.Confidence = 0.70
	}
	result.IsAIGenerated = result.FinalScore >= d.config.MinConfidence
	return result, nil
}

func (d *MLDetector) GetModelInfo() map[string]interface{} {
	info := map[string]interface{}{
		"ml_enabled":     d.onnxDetector.IsInitialized(),
		"ml_weight":      d.config.MLWeight,
		"min_confidence": d.config.MinConfidence,
		"fallback_mode":  d.config.FallbackMode,
	}
	if d.onnxDetector.IsInitialized() {
		modelInfo := d.onnxDetector.GetModelInfo()
		for k, v := range modelInfo {
			info[k] = v
		}
	}
	return info
}

func (d *MLDetector) Close() error {
	if d.onnxDetector != nil {
		return d.onnxDetector.Close()
	}
	return nil
}

func abs(x float64) float64 {
	if x < 0 {
		return -x
	}
	return x
}

func padFeatures(features []float32, targetSize int) []float32 {
	if len(features) == targetSize {
		return features
	}
	padded := make([]float32, targetSize)
	copyLen := len(features)
	if copyLen > targetSize {
		copyLen = targetSize
	}
	copy(padded, features[:copyLen])
	for i := copyLen; i < targetSize; i++ {
		padded[i] = 0.0
	}
	return padded
}
