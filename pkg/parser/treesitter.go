package parser

import (
	"fmt"
	"sync"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
)

var (
	treeSitterMu sync.Mutex
)

type TreeSitterParser struct{}

type ParsedCode struct {
	AST      *sitter.Node
	Code     []byte
	Language string
	tree     *sitter.Tree // Store tree so we can close it later
}

func New() *TreeSitterParser {
	return &TreeSitterParser{}
}

func (p *TreeSitterParser) Parse(code string, language string) (*ParsedCode, error) {
	treeSitterMu.Lock()
	defer treeSitterMu.Unlock()

	var lang *sitter.Language

	switch language {
	case "python":
		lang = python.GetLanguage()
	case "javascript", "typescript":
		lang = javascript.GetLanguage()
	case "go":
		lang = golang.GetLanguage()
	default:
		return nil, fmt.Errorf("unsupported language: %s", language)
	}

	if lang == nil {
		return nil, fmt.Errorf("failed to get language for: %s", language)
	}

	parser := sitter.NewParser()
	defer parser.Close() // CRITICAL: Close parser when done

	parser.SetLanguage(lang)

	codeBytes := []byte(code)

	if len(codeBytes) == 0 {
		return nil, fmt.Errorf("empty code")
	}

	tree, err := parser.ParseCtx(nil, nil, codeBytes)
	if err != nil {
		return nil, fmt.Errorf("parse error: %w", err)
	}

	if tree == nil {
		return nil, fmt.Errorf("parser returned nil tree")
	}

	rootNode := tree.RootNode()
	if rootNode == nil {
		tree.Close() // Don't forget to close on error
		return nil, fmt.Errorf("tree returned nil root node")
	}

	return &ParsedCode{
		AST:      rootNode,
		Code:     codeBytes,
		Language: language,
		tree:     tree, // Store tree reference
	}, nil
}

// Close releases C memory - MUST be called when done with ParsedCode
func (pc *ParsedCode) Close() {
	if pc.tree != nil {
		pc.tree.Close()
		pc.tree = nil
	}
}

func (p *TreeSitterParser) IsSupported(language string) bool {
	supportedLanguages := map[string]bool{
		"python":     true,
		"javascript": true,
		"typescript": true,
		"go":         true,
	}
	return supportedLanguages[language]
}

func (pc *ParsedCode) Query(queryStr string) ([]*sitter.Node, error) {
	treeSitterMu.Lock()
	defer treeSitterMu.Unlock()

	var lang *sitter.Language
	switch pc.Language {
	case "python":
		lang = python.GetLanguage()
	case "javascript", "typescript":
		lang = javascript.GetLanguage()
	case "go":
		lang = golang.GetLanguage()
	default:
		return nil, fmt.Errorf("unsupported language for query: %s", pc.Language)
	}

	q, err := sitter.NewQuery([]byte(queryStr), lang)
	if err != nil {
		return nil, fmt.Errorf("invalid query: %w", err)
	}
	defer q.Close() // CRITICAL: Close query when done

	qc := sitter.NewQueryCursor()
	defer qc.Close() // CRITICAL: Close query cursor when done

	qc.Exec(q, pc.AST)

	nodes := []*sitter.Node{}
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		for _, c := range m.Captures {
			nodes = append(nodes, c.Node)
		}
	}

	return nodes, nil
}

func (pc *ParsedCode) GetNodeText(node *sitter.Node) string {
	if node == nil {
		return ""
	}
	return string(pc.Code[node.StartByte():node.EndByte()])
}

func (pc *ParsedCode) CountNodes(queryStr string) (int, error) {
	nodes, err := pc.Query(queryStr)
	if err != nil {
		return 0, err
	}
	return len(nodes), nil
}

func (pc *ParsedCode) WalkAST(node *sitter.Node, callback func(*sitter.Node) bool) {
	if node == nil {
		return
	}

	if !callback(node) {
		return
	}

	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		pc.WalkAST(child, callback)
	}
}

func (pc *ParsedCode) GetFunctions() ([]*sitter.Node, error) {
	var query string
	switch pc.Language {
	case "python":
		query = "(function_definition) @func"
	case "javascript", "typescript":
		query = "(function_declaration) @func"
	case "go":
		query = "(function_declaration) @func"
	default:
		return nil, fmt.Errorf("functions query not implemented for: %s", pc.Language)
	}

	return pc.Query(query)
}

func (pc *ParsedCode) GetVariables() ([]*sitter.Node, error) {
	var query string
	switch pc.Language {
	case "python":
		query = "(identifier) @var"
	case "javascript", "typescript":
		query = "(identifier) @var"
	case "go":
		query = "(identifier) @var"
	default:
		return nil, fmt.Errorf("variables query not implemented for: %s", pc.Language)
	}

	return pc.Query(query)
}

func (pc *ParsedCode) GetComments() []*sitter.Node {
	comments := []*sitter.Node{}

	pc.WalkAST(pc.AST, func(node *sitter.Node) bool {
		nodeType := node.Type()
		if nodeType == "comment" || nodeType == "line_comment" || nodeType == "block_comment" {
			comments = append(comments, node)
		}
		return true
	})

	return comments
}
