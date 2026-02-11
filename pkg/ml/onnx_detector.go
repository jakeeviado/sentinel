package ml

import (
	"fmt"
	"os"

	onnxruntime "github.com/yalue/onnxruntime_go"
)

type ONNXDetector struct {
	session     *onnxruntime.DynamicSession[float32, float32]
	inputName   string
	outputName  string
	inputShape  []int64
	initialized bool
}

type ModelConfig struct {
	ModelPath      string
	InputName      string
	OutputName     string
	NumFeatures    int
	UseGPU         bool
	OptimizeMemory bool
}

func NewONNXDetector(config ModelConfig) (*ONNXDetector, error) {
	if _, err := os.Stat(config.ModelPath); os.IsNotExist(err) {
		return &ONNXDetector{initialized: false}, nil
	}

	if err := onnxruntime.InitializeEnvironment(); err != nil {
		return nil, fmt.Errorf("failed to initialize ONNX runtime: %w", err)
	}

	if config.InputName == "" {
		config.InputName = "input"
	}
	if config.OutputName == "" {
		config.OutputName = "output"
	}
	if config.NumFeatures == 0 {
		config.NumFeatures = DefaultNumFeatures
	}

	options, err := onnxruntime.NewSessionOptions()
	if err != nil {
		return nil, fmt.Errorf("failed to create session options: %w", err)
	}
	defer func() {
		if err := options.Destroy(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to destroy options: %v\n", err)
		}
	}()

	if config.OptimizeMemory {
		if err := options.SetMemPattern(false); err != nil {
			return nil, fmt.Errorf("failed to set memory pattern: %w", err)
		}
	}

	session, err := onnxruntime.NewDynamicSession[float32, float32](
		config.ModelPath,
		[]string{config.InputName},
		[]string{config.OutputName},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to load ONNX model: %w", err)
	}

	detector := &ONNXDetector{
		session:     session,
		inputName:   config.InputName,
		outputName:  config.OutputName,
		inputShape:  []int64{1, int64(config.NumFeatures)},
		initialized: true,
	}

	return detector, nil
}

func (d *ONNXDetector) Predict(features []float32) (float32, error) {
	if !d.initialized {
		return 0.0, fmt.Errorf("ONNX detector not initialized")
	}

	inputTensor, err := onnxruntime.NewTensor(d.inputShape, features)
	if err != nil {
		return 0.0, err
	}

	defer func() {
		if err := inputTensor.Destroy(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to destroy input tensor: %v\n", err)
		}
	}()

	outputTensors := []*onnxruntime.Tensor[float32]{nil, nil}

	err = d.session.Run([]*onnxruntime.Tensor[float32]{inputTensor}, outputTensors)
	if err != nil {
		return 0.0, fmt.Errorf("inference failed: %w", err)
	}

	defer func() {
		for i, t := range outputTensors {
			if t != nil {
				if err := t.Destroy(); err != nil {
					fmt.Fprintf(os.Stderr, "Warning: failed to destroy output tensor %d: %v\n", i, err)
				}
			}
		}
	}()

	if outputTensors[1] == nil {
		return 0.0, fmt.Errorf("model produced no probability tensor at index 1")
	}

	probData := outputTensors[1].GetData()
	if len(probData) < 2 {
		return 0.0, fmt.Errorf("invalid probability data length")
	}

	return probData[1], nil
}

func (d *ONNXDetector) IsInitialized() bool {
	return d.initialized
}

func (d *ONNXDetector) Close() error {
	if d.session != nil {
		if err := d.session.Destroy(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to destroy session: %v\n", err)
		}
	}
	return onnxruntime.DestroyEnvironment()
}

func (d *ONNXDetector) PredictBatch(featureBatch [][]float32) ([]float32, error) {
	if !d.initialized || len(featureBatch) == 0 {
		return nil, fmt.Errorf("ONNX detector not initialized or empty batch")
	}

	numSamples := int64(len(featureBatch))
	numFeatures := d.inputShape[1]

	flatFeatures := make([]float32, 0, numSamples*numFeatures)
	for _, f := range featureBatch {
		flatFeatures = append(flatFeatures, f...)
	}

	batchShape := []int64{numSamples, numFeatures}
	inputTensor, err := onnxruntime.NewTensor(batchShape, flatFeatures)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err := inputTensor.Destroy(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: failed to destroy input tensor: %v\n", err)
		}
	}()

	outputTensors := make([]*onnxruntime.Tensor[float32], 1)
	err = d.session.Run([]*onnxruntime.Tensor[float32]{inputTensor}, outputTensors)
	if err != nil {
		return nil, err
	}
	defer func() {
		if outputTensors[0] != nil {
			if err := outputTensors[0].Destroy(); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: failed to destroy batch output tensor: %v\n", err)
			}
		}
	}()

	return outputTensors[0].GetData(), nil
}

func (d *ONNXDetector) GetModelInfo() map[string]interface{} {
	return map[string]interface{}{
		"initialized":  d.initialized,
		"input_name":   d.inputName,
		"output_name":  d.outputName,
		"input_shape":  d.inputShape,
		"num_features": d.inputShape[1],
	}
}
