package ml

import (
	"fmt"
	"sentinel/pkg/models"
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
	// Set defaults
	if config.MLWeight == 0.0 {
		config.MLWeight = 0.7
	}
	if config.MinConfidence == 0.0 {
		config.MinConfidence = 0.5
	}
	if config.NumFeatures == 0 {
		config.NumFeatures = 25
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
	result := &DetectionResult{}

	heuristicScore := 0.0
	for _, sig := range signals {
		if sig.Score > heuristicScore {
			heuristicScore = sig.Score
		}
	}
	result.HeuristicScore = heuristicScore

	if d.onnxDetector.IsInitialized() {
		codeMetrics := ExtractCodeMetrics(code)
		featureVector := d.featureExtractor.Extract(signals, language, codeMetrics)
		features := padFeatures(featureVector.Features, d.config.NumFeatures)

		mlScore, err := d.onnxDetector.Predict(features)
		if err != nil {
			if !d.config.FallbackMode {
				return nil, fmt.Errorf("ML inference failed: %w", err)
			}
			result.MLScore = 0.0
			result.UsedML = false
			result.FinalScore = heuristicScore
		} else {
			result.MLScore = float64(mlScore)
			result.UsedML = true
			mlWeight := d.config.MLWeight
			heuristicWeight := 1.0 - mlWeight
			result.FinalScore = (mlWeight * result.MLScore) + (heuristicWeight * heuristicScore)
		}
	} else {
		result.MLScore = 0.0
		result.UsedML = false
		result.FinalScore = heuristicScore
	}

	if result.UsedML {
		scoreDiff := abs(result.MLScore - result.HeuristicScore)
		result.Confidence = 1.0 - scoreDiff
	} else {
		result.Confidence = 0.7
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
