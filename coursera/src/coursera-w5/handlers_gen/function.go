package main

import (
	"encoding/json"
	"go/ast"
	"strings"
)

// FuncGenerator generate functions
type FuncGenerator struct {
	GenerationLabel string
	Declararion     *ast.FuncDecl
	Wrappers        map[string][]Wrapper
}

// Wrapper contains information about function for building wrapper
type Wrapper struct {
	Method    string
	URL       URL
	ParamType string
}

// URL consists information about http request
type URL struct {
	Path   string `json:"url"`
	Auth   bool   `json:"auth"`
	Method string `json:"method"`
}

// NewFuncGenerator set and return struct for generation wrappers
func (g *Generator) NewFuncGenerator() *FuncGenerator {
	g.funcGenerator = FuncGenerator{
		GenerationLabel: g.FuncGenerationLabel,
		Wrappers:        make(map[string][]Wrapper),
	}

	return &g.funcGenerator
}

func (fg *FuncGenerator) getURLInfo() URL {
	doc := fg.Declararion.Doc
	url := URL{}

	if doc != nil {
		var urlJSON string
		for _, comment := range doc.List {
			urlJSON = strings.ReplaceAll(comment.Text, fg.GenerationLabel, "")
		}

		json.Unmarshal([]byte(urlJSON), &url)
	}

	return url
}

func (fg FuncGenerator) defineType(expr interface{}) string {
	ident, ok := expr.(*ast.Ident)
	if ok {
		return ident.String()
	}
	return ""
}

func (fg *FuncGenerator) addWrapper(expr ast.Expr, url URL, paramType string) {
	wrappers := fg.Wrappers[fg.defineType(expr)]
	wp := Wrapper{}
	wp.Method = fg.Declararion.Name.Name
	wp.URL = url
	wp.ParamType = paramType
	wrappers = append(wrappers, wp)
	fg.Wrappers[fg.defineType(expr)] = wrappers
}

// IsGenerationNeed if true will generate
func (fg FuncGenerator) IsGenerationNeed() bool {
	doc := fg.Declararion.Doc
	if doc == nil {
		return false
	}

	needCodegen := false
	for _, comment := range doc.List {
		needCodegen = needCodegen || strings.HasPrefix(comment.Text, fg.GenerationLabel)
	}

	if !needCodegen {
		return false
	}

	return true
}

// FillWrappers define and set recv type of function
func (fg *FuncGenerator) FillWrappers() {
	if !fg.IsGenerationNeed() {
		return
	}

	recv := fg.Declararion.Recv
	if recv == nil {
		return
	}

	url := fg.getURLInfo()

	var paramType string
	params := fg.Declararion.Type.Params
	if params != nil {
		paramType = fg.defineType(params.List[1].Type)
	}

	for _, r := range recv.List {
		expr, ok := r.Type.(*ast.StarExpr)
		if ok {
			fg.addWrapper(expr.X, url, paramType)
			continue
		}

		fg.addWrapper(r.Type, url, paramType)
	}
}
