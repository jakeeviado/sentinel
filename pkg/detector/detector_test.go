package detector

import (
	"sentinel/pkg/models"
	"testing"
)

func TestNew(t *testing.T) {
	config := Config{
		Threshold: 0.7,
		Languages: []string{"python", "go"},
		Verbose:   true,
	}

	detector := New(config)

	if detector == nil {
		t.Fatal("Expected detector to be initialized")
	}

	if detector.config.Threshold != 0.7 {
		t.Errorf("Expected threshold 0.7, got %f", detector.config.Threshold)
	}
}

func TestDetectLanguage(t *testing.T) {
	tests := []struct {
		path     string
		expected string
	}{
		{"test.py", "python"},
		{"main.go", "go"},
		{"app.js", "javascript"},
		{"Main.java", "java"},
		{"test.rs", "rust"},
		{"script.rb", "ruby"},
		{"Program.cs", "csharp"},
		{"unknown.xyz", "unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := detectLanguage(tt.path)
			if result != tt.expected {
				t.Errorf("detectLanguage(%s) = %s, want %s", tt.path, result, tt.expected)
			}
		})
	}
}

func TestCalculateScore(t *testing.T) {
	tests := []struct {
		name     string
		signals  []models.Signal
		expected float64
	}{
		{
			name:     "empty signals",
			signals:  []models.Signal{},
			expected: 0.0,
		},
		{
			name: "single signal",
			signals: []models.Signal{
				{Name: "test", Score: 0.8},
			},
			expected: 0.8,
		},
		{
			name: "multiple signals - picks maximum",
			signals: []models.Signal{
				{Name: "test1", Score: 0.6},
				{Name: "test2", Score: 0.9},
				{Name: "test3", Score: 0.4},
			},
			expected: 0.9,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calculateMaxScore(tt.signals)
			if result != tt.expected {
				t.Errorf("calculateMaxScore() = %f, want %f", result, tt.expected)
			}
		})
	}
}

func TestScanResultsHasDetections(t *testing.T) {
	results := &ScanResults{
		Files: []FileResult{
			{Path: "test1.py", Score: 0.5},
			{Path: "test2.py", Score: 0.8},
			{Path: "test3.py", Score: 0.6},
		},
	}

	if !results.HasDetections(0.7) {
		t.Error("Expected HasDetections(0.7) to be true")
	}

	if results.HasDetections(0.9) {
		t.Error("Expected HasDetections(0.9) to be false")
	}
}

func TestLanguageIDMapConsistency(t *testing.T) {
	testLangs := []string{"python", "go", "java", "rust", "unknown"}

	for _, lang := range testLangs {
		id, ok := languageIDMap[lang]
		if !ok {
			t.Errorf("Language %s found by detectLanguage but missing from languageIDMap", lang)
		}
		if id == "" {
			t.Errorf("Language %s has an empty ID in languageIDMap", lang)
		}
	}
}
