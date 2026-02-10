package analyzer

import (
	"testing"
)

func TestASTAnalyzerPython(t *testing.T) {
	analyzer := NewAST()

	// AI-generated Python code with obvious patterns
	code := `
def process_data(data):
    """Process the input data and return the result."""
    # Initialize the result variable
    result = []

    # TODO: Process each item in the data
    for item in data:
        # Create a temporary variable to store the processed item
        temp = item * 2

        # Append the temporary result to the result list
        result.append(temp)

    # Return the final result
    return result

def calculate_value(value):
    """Calculate and return the computed value."""
    # Initialize variable
    temp = value

    # Perform calculation
    temp = temp + 10

    # Return result
    return temp
`

	signals := analyzer.AnalyzeWithAST(code, "python")

	if len(signals) == 0 {
		t.Fatal("Expected AST signals but got none")
	}

	// Check that we got AST-specific signals
	hasASTSignal := false
	for _, sig := range signals {
		if sig.Name == "ast_comment_density" ||
			sig.Name == "ast_generic_naming" ||
			sig.Name == "function_uniformity" {
			hasASTSignal = true
			break
		}
	}

	if !hasASTSignal {
		t.Error("Expected at least one AST-specific signal")
	}

	// Should detect high comment density
	foundCommentSignal := false
	for _, sig := range signals {
		if sig.Name == "ast_comment_density" && sig.Score > 0.5 {
			foundCommentSignal = true
		}
	}
	if !foundCommentSignal {
		t.Error("Expected to detect high comment density")
	}

	// Should detect generic naming
	foundNamingSignal := false
	for _, sig := range signals {
		if sig.Name == "ast_generic_naming" && sig.Score > 0.5 {
			foundNamingSignal = true
		}
	}
	if !foundNamingSignal {
		t.Error("Expected to detect generic variable naming")
	}
}

func TestASTAnalyzerJavaScript(t *testing.T) {
	analyzer := NewAST()

	code := `
function getData() {
    const data = fetch();
    const result = process(data);
    return result;
}

function getItems() {
    const data = fetch();
    const result = process(data);
    return result;
}
`

	signals := analyzer.AnalyzeWithAST(code, "javascript")

	if len(signals) == 0 {
		t.Fatal("Expected AST signals for JavaScript")
	}

	// Should detect structural repetition
	foundRepetition := false
	for _, sig := range signals {
		if sig.Name == "structural_repetition" && sig.Score > 0 {
			foundRepetition = true
		}
	}
	if !foundRepetition {
		t.Error("Expected to detect structural repetition")
	}
}

func TestASTAnalyzerUnsupportedLanguage(t *testing.T) {
	analyzer := NewAST()

	code := `
public class Test {
    public static void main(String[] args) {
        System.out.println("Hello");
    }
}
`

	signals := analyzer.AnalyzeWithAST(code, "java")

	// Should return empty signals for unsupported language
	if len(signals) != 0 {
		t.Error("Expected no AST signals for unsupported language")
	}
}

func TestFunctionUniformity(t *testing.T) {
	analyzer := NewAST()

	// Code with very uniform functions (AI pattern)
	uniformCode := `
def func1():
    x = 1
    y = 2
    z = x + y
    return z

def func2():
    a = 3
    b = 4
    c = a + b
    return c

def func3():
    p = 5
    q = 6
    r = p + q
    return r

def func4():
    m = 7
    n = 8
    o = m + n
    return o

def func5():
    i = 9
    j = 10
    k = i + j
    return k
`

	signals := analyzer.AnalyzeWithAST(uniformCode, "python")

	foundUniformity := false
	for _, sig := range signals {
		if sig.Name == "function_uniformity" && sig.Score > 0.5 {
			foundUniformity = true
		}
	}

	if !foundUniformity {
		t.Error("Expected to detect function uniformity")
	}
}

func TestTrivialFunctions(t *testing.T) {
	analyzer := NewAST()

	// Function with no control flow (trivial)
	trivialCode := `
def process(items):
    """Process all items."""
    result = []
    for item in items:
        result.append(item)
    return result
`

	signals := analyzer.AnalyzeWithAST(trivialCode, "python")

	foundTrivial := false
	for _, sig := range signals {
		if sig.Name == "trivial_functions" && sig.Score > 0 {
			foundTrivial = true
		}
	}

	if !foundTrivial {
		t.Error("Expected to detect trivial function")
	}
}
