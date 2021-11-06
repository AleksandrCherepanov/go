package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"text/template"
)

// HTTPServeTemplate template
const HTTPServeTemplate = `
{{range $key, $value := .Wrappers}}
func (h * {{$key}}) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path { {{range $value}}
	case "{{.URL.Path}}":
		handler := Handler{}
		handler.process = h.wrapper{{.Method}}
		
		{{ if and (ne .URL.Method "") (eq .URL.Auth true) }}
			HTTPMethodCheckMiddleware(AuthCheckMiddleware(handler), "{{.URL.Method}}").ServeHTTP(w, r)
		{{else if ne .URL.Method ""}}
			HTTPMethodCheckMiddleware(handler, "{{.URL.Method}}").ServeHTTP(w, r)
		{{else if eq .URL.Auth true}}
			AuthCheckMiddleware(handler), .URL.Method).ServeHTTP(w, r)
		{{else}}
			handler.ServeHTTP(w,r)
		{{end}}

		{{end}}
	default:
		w.WriteHeader(http.StatusNotFound)
		error := BuildHttpError("unknown method")
		io.WriteString(w, error)
	}
}
{{end}}
`

// WrapperTemplate template
const WrapperTemplate = `{{range $key, $value := .Wrappers}}{{range $value}}
func (h *{{$key}}) wrapper{{.Method}}(w http.ResponseWriter, r *http.Request) {
	{{template "CallAPIMethodTemplate" .}}
}
{{end}}{{end}}
`

// CallAPIMethodTemplate template
const CallAPIMethodTemplate = `{{define "CallAPIMethodTemplate"}}
	defer(func() {
		if r := recover(); r != nil {
			w.WriteHeader(http.StatusInternalServerError)
			io.WriteString(w, BuildHttpError(fmt.Sprintf("%v", r)))
		}
	})()	


	params, ae := {{.ParamType}}Extractor(r)
	if (ae != (ApiError{})) {
		w.WriteHeader(ae.HTTPStatus)
		io.WriteString(w, BuildHttpError(ae.Error()))
		return
	}

	res, err := h.{{.Method}}(r.Context(), params)
	if err != nil {
		message := err.Error()
		code := http.StatusInternalServerError
		ae, ok := err.(ApiError)
		if (ok) {
			code = ae.HTTPStatus
			message = ae.Error()
		}

		w.WriteHeader(code)
		io.WriteString(w, BuildHttpError(message))
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, BuildResponse(res))
{{end}}`

// GenerateHTTPServe genereates functions for serve http requests
func (g Generator) GenerateHTTPServe(out *os.File) {
	tmplServe, err := template.New("httpServe").Parse(HTTPServeTemplate)

	if err != nil {
		panic(err)
	}

	tmplServe.Execute(out, g.funcGenerator)
}

// GenerateWrappers generate wrappers
func (g Generator) GenerateWrappers(out *os.File) {
	wrapperTmpl, err := template.New("WrapperTemplate").Parse(WrapperTemplate)
	// wrapperTmpl.New("HTTPMethodCheckTemplate").Parse(HTTPMethodCheckTemplate)
	wrapperTmpl.New("CallAPIMethodTemplate").Parse(CallAPIMethodTemplate)

	if err != nil {
		panic(err)
	}

	wrapperTmpl.ExecuteTemplate(out, "WrapperTemplate", g.funcGenerator)
}

func main() {
	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, os.Args[1], nil, parser.ParseComments)
	if err != nil {
		panic(err)
	}

	out, err := os.Create(os.Args[2])
	if err != nil {
		panic(err)
	}

	generator := Generator{
		PackageName: node.Name.Name,
		Imports: []string{
			"net/http",
			"encoding/json",
			"io",
			"strconv",
			"errors",
			"strings",
			"fmt",
		},
		FuncGenerationLabel:   "// apigen:api ",
		StructGenerationLabel: "`apivalidator:",
	}

	generator.GenerateCommon(out)

	fg := generator.NewFuncGenerator()
	sg := generator.NewStructGenerator()

	for _, decl := range node.Decls {
		genDecl, ok := decl.(*ast.GenDecl)
		if ok {
			sg.Declararion = genDecl
			sg.FillStruct()
			continue
		}
	}

	for _, decl := range node.Decls {
		funcDecl, ok := decl.(*ast.FuncDecl)
		if ok {
			fg.Declararion = funcDecl
			fg.FillWrappers()
			continue
		}
	}

	tmpl, err := template.New("SturctureExtractorTempalte").Parse(SturctureExtractorTempalte)
	if err != nil {
		panic(err)
	}

	tmpl.Execute(out, generator.structGenereator.structList)

	generator.GenerateHTTPServe(out)
	generator.GenerateWrappers(out)
}
