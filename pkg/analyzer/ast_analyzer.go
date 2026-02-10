package analyzer

import (
	"fmt"
	"math"
	"strings"

	"sentinel/pkg/models"
	"sentinel/pkg/parser"

	sitter "github.com/smacker/go-tree-sitter"
)

// ASTAnalyzer performs AST-based code analysis using tree-sitter
type ASTAnalyzer struct {
	parser *parser.TreeSitterParser
}

// NewAST creates a new AST-based analyzer
func NewAST() *ASTAnalyzer {
	return &ASTAnalyzer{
		parser: parser.New(),
	}
}

/*
 * AnalyzeWithAST performs AST-based analysis on code
 */
func (a *ASTAnalyzer) AnalyzeWithAST(code string, language string) []models.Signal {
	return nil
	signals := make([]models.Signal, 0)

	if !a.parser.IsSupported(language) {
		return nil
	}

	parsed, err := a.parser.Parse(code, language)
	if err != nil {
		return nil
	}
	defer parsed.Close()

	// Run AST-based checks
	signals = append(signals, a.checkASTCommentDensity(parsed))
	signals = append(signals, a.checkASTGenericNaming(parsed))
	signals = append(signals, a.checkFunctionUniformity(parsed))
	signals = append(signals, a.checkTrivialFunctions(parsed))
	signals = append(signals, a.checkStructuralRepetition(parsed))
	signals = append(signals, a.checkOverDocumentation(parsed))

	return signals
}

// checkASTCommentDensity uses AST to accurately count comments
func (a *ASTAnalyzer) checkASTCommentDensity(pc *parser.ParsedCode) models.Signal {
	comments := pc.GetComments()

	// Count total lines of code (non-comment, non-blank)
	lines := strings.Split(string(pc.Code), "\n")
	codeLines := 0
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) > 0 {
			codeLines++
		}
	}

	commentLines := 0
	for _, comment := range comments {
		commentText := pc.GetNodeText(comment)
		commentLines += strings.Count(commentText, "\n") + 1
	}

	if codeLines == 0 {
		return models.Signal{Name: "ast_comment_density", Score: 0.0}
	}

	density := float64(commentLines) / float64(codeLines)

	var score float64
	var description string

	if density > 0.4 {
		score = 0.9
		description = fmt.Sprintf("Excessive comment density: %.1f%% (AI often over-comments)", density*100)
	} else if density > 0.25 {
		score = 0.6
		description = fmt.Sprintf("High comment density: %.1f%%", density*100)
	} else {
		score = 0.0
		description = "Normal comment density"
	}

	return models.Signal{
		Name:        "ast_comment_density",
		Score:       score,
		Description: description,
		Evidence:    fmt.Sprintf("%d comment lines / %d code lines", commentLines, codeLines),
	}
}

// checkASTGenericNaming analyzes actual variable/parameter names from AST
func (a *ASTAnalyzer) checkASTGenericNaming(pc *parser.ParsedCode) models.Signal {
	genericNames := map[string]bool{
		"temp": true, "tmp": true, "data": true, "result": true,
		"item": true, "value": true, "obj": true, "elem": true,
		"variable": true, "parameter": true, "arr": true, "list": true,
	}

	variables, err := pc.GetVariables()
	if err != nil {
		return models.Signal{Name: "ast_generic_naming", Score: 0.0}
	}

	genericCount := 0
	uniqueVars := make(map[string]bool)

	for _, varNode := range variables {
		varName := strings.ToLower(pc.GetNodeText(varNode))

		// Skip duplicates and very short names
		if len(varName) < 2 || uniqueVars[varName] {
			continue
		}
		uniqueVars[varName] = true

		if genericNames[varName] {
			genericCount++
		}
	}

	totalVars := len(uniqueVars)
	if totalVars == 0 {
		return models.Signal{Name: "ast_generic_naming", Score: 0.0}
	}

	ratio := float64(genericCount) / float64(totalVars)

	var score float64
	var description string

	if ratio > 0.5 {
		score = 0.95
		description = fmt.Sprintf("Very high generic naming: %.0f%% of variables", ratio*100)
	} else if ratio > 0.3 {
		score = 0.7
		description = fmt.Sprintf("High generic naming: %.0f%% of variables", ratio*100)
	} else if ratio > 0.15 {
		score = 0.4
		description = fmt.Sprintf("Moderate generic naming: %.0f%% of variables", ratio*100)
	} else {
		score = 0.0
		description = "Good variable naming"
	}

	return models.Signal{
		Name:        "ast_generic_naming",
		Score:       score,
		Description: description,
		Evidence:    fmt.Sprintf("%d generic names out of %d unique variables", genericCount, totalVars),
	}
}

// checkFunctionUniformity detects suspiciously uniform function sizes (AI pattern)
func (a *ASTAnalyzer) checkFunctionUniformity(pc *parser.ParsedCode) models.Signal {
	functions, err := pc.GetFunctions()
	if err != nil || len(functions) < 3 {
		return models.Signal{Name: "function_uniformity", Score: 0.0}
	}

	// Calculate line counts for each function
	sizes := make([]int, 0, len(functions))
	for _, fn := range functions {
		fnText := pc.GetNodeText(fn)
		lineCount := strings.Count(fnText, "\n") + 1
		sizes = append(sizes, lineCount)
	}

	// Calculate mean and standard deviation
	mean := 0.0
	for _, size := range sizes {
		mean += float64(size)
	}
	mean /= float64(len(sizes))

	variance := 0.0
	for _, size := range sizes {
		diff := float64(size) - mean
		variance += diff * diff
	}
	variance /= float64(len(sizes))
	stdDev := math.Sqrt(variance)

	// Coefficient of variation (lower = more uniform)
	cv := stdDev / mean

	var score float64
	var description string

	if cv < 0.15 && len(functions) >= 5 {
		score = 0.85
		description = fmt.Sprintf("Suspiciously uniform function sizes (CV: %.2f)", cv)
	} else if cv < 0.25 && len(functions) >= 5 {
		score = 0.5
		description = fmt.Sprintf("Very consistent function sizes (CV: %.2f)", cv)
	} else {
		score = 0.0
		description = "Normal function size variation"
	}

	return models.Signal{
		Name:        "function_uniformity",
		Score:       score,
		Description: description,
		Evidence:    fmt.Sprintf("%d functions, mean: %.1f lines, std dev: %.1f", len(functions), mean, stdDev),
	}
}

// checkTrivialFunctions detects functions that are overly simple (AI pattern)
func (a *ASTAnalyzer) checkTrivialFunctions(pc *parser.ParsedCode) models.Signal {
	functions, err := pc.GetFunctions()
	if err != nil || len(functions) == 0 {
		return models.Signal{Name: "trivial_functions", Score: 0.0}
	}

	trivialCount := 0

	for _, fn := range functions {
		// Count control flow statements in function
		controlFlowCount := 0
		pc.WalkAST(fn, func(node *sitter.Node) bool {
			nodeType := node.Type()
			controlFlowTypes := map[string]bool{
				"if_statement":     true,
				"for_statement":    true,
				"while_statement":  true,
				"switch_statement": true,
			}
			if controlFlowTypes[nodeType] {
				controlFlowCount++
			}
			return true
		})

		// Get function size
		fnText := pc.GetNodeText(fn)
		lineCount := strings.Count(fnText, "\n") + 1

		// Trivial if: large (>5 lines) but no control flow
		if lineCount > 5 && controlFlowCount == 0 {
			trivialCount++
		}
	}

	ratio := float64(trivialCount) / float64(len(functions))

	var score float64
	var description string

	if ratio > 0.5 {
		score = 0.8
		description = fmt.Sprintf("Many trivial functions: %.0f%% have no control flow", ratio*100)
	} else if ratio > 0.3 {
		score = 0.5
		description = fmt.Sprintf("Some trivial functions: %.0f%% have no control flow", ratio*100)
	} else {
		score = 0.0
		description = "Normal function complexity"
	}

	return models.Signal{
		Name:        "trivial_functions",
		Score:       score,
		Description: description,
		Evidence:    fmt.Sprintf("%d/%d functions are trivial", trivialCount, len(functions)),
	}
}

// checkStructuralRepetition detects repeated AST patterns
func (a *ASTAnalyzer) checkStructuralRepetition(pc *parser.ParsedCode) models.Signal {
	functions, err := pc.GetFunctions()
	if err != nil || len(functions) < 3 {
		return models.Signal{Name: "structural_repetition", Score: 0.0}
	}

	// Hash function structures
	structureHashes := make(map[string]int)

	for _, fn := range functions {
		hash := hashASTStructure(fn)
		structureHashes[hash]++
	}

	// Count duplicates
	duplicates := 0
	for _, count := range structureHashes {
		if count > 1 {
			duplicates += count - 1
		}
	}

	ratio := float64(duplicates) / float64(len(functions))

	var score float64
	var description string

	if ratio > 0.4 {
		score = 0.8
		description = fmt.Sprintf("High structural repetition: %.0f%% duplicate patterns", ratio*100)
	} else if ratio > 0.2 {
		score = 0.5
		description = fmt.Sprintf("Moderate structural repetition: %.0f%% duplicate patterns", ratio*100)
	} else {
		score = 0.0
		description = "Low structural repetition"
	}

	return models.Signal{
		Name:        "structural_repetition",
		Score:       score,
		Description: description,
		Evidence:    fmt.Sprintf("%d/%d functions have duplicate structures", duplicates, len(functions)),
	}
}

// checkOverDocumentation detects when comments are longer than code (AI pattern)
func (a *ASTAnalyzer) checkOverDocumentation(pc *parser.ParsedCode) models.Signal {
	functions, err := pc.GetFunctions()
	if err != nil || len(functions) == 0 {
		return models.Signal{Name: "over_documentation", Score: 0.0}
	}

	overDocumentedCount := 0

	for _, fn := range functions {
		// Find comments immediately before/in function
		commentSize := 0
		pc.WalkAST(fn, func(node *sitter.Node) bool {
			nodeType := node.Type()
			if nodeType == "comment" || nodeType == "line_comment" || nodeType == "block_comment" {
				commentSize += len(pc.GetNodeText(node))
			}
			return true
		})

		fnText := pc.GetNodeText(fn)
		codeSize := len(fnText) - commentSize

		// Over-documented if comments are >70% of function size
		if codeSize > 0 && float64(commentSize)/float64(codeSize) > 0.7 {
			overDocumentedCount++
		}
	}

	ratio := float64(overDocumentedCount) / float64(len(functions))

	var score float64
	var description string

	if ratio > 0.5 {
		score = 0.85
		description = fmt.Sprintf("Excessive documentation: %.0f%% of functions over-commented", ratio*100)
	} else if ratio > 0.3 {
		score = 0.6
		description = fmt.Sprintf("Heavy documentation: %.0f%% of functions over-commented", ratio*100)
	} else {
		score = 0.0
		description = "Reasonable documentation level"
	}

	return models.Signal{
		Name:        "over_documentation",
		Score:       score,
		Description: description,
		Evidence:    fmt.Sprintf("%d/%d functions have comments longer than code", overDocumentedCount, len(functions)),
	}
}

// hashASTStructure creates a simple hash of AST structure
func hashASTStructure(node *sitter.Node) string {
	if node == nil {
		return ""
	}

	// Create a simple structural signature
	signature := node.Type()

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		signature += "|" + child.Type()
	}

	return signature
}
