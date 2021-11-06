package main

import (
	"os"
	"text/template"
)

// Generator main struct for generation
type Generator struct {
	PackageName           string
	Imports               []string
	FuncGenerationLabel   string
	StructGenerationLabel string
	funcGenerator         FuncGenerator
	structGenereator      StructGenerator
}

// GenerateCommon generate common parts
func (g Generator) GenerateCommon(out *os.File) {
	tmpl, err := template.New("MainTemplate").Parse(MainTemplate)

	if err != nil {
		panic(err)
	}

	tmpl.New("HTTPMethodCheckMiddlewareTemplate").Parse(HTTPMethodCheckMiddlewareTemplate)
	tmpl.New("AuthCheckMiddlewareTemplate").Parse(AuthCheckMiddlewareTemplate)
	tmpl.New("GetParamIntTemplate").Parse(GetParamIntTemplate)
	tmpl.New("SetDefaultValueIntTemplate").Parse(SetDefaultValueIntTemplate)
	tmpl.New("SetDefaultValueStringTemplate").Parse(SetDefaultValueStringTemplate)
	tmpl.New("RequiredIntValidatorTemplate").Parse(RequiredIntValidatorTemplate)
	tmpl.New("RequiredStringValidatorTemplate").Parse(RequiredStringValidatorTemplate)
	tmpl.New("EnumValidatorTemplate").Parse(EnumValidatorTemplate)
	tmpl.New("MinIntValidatorTemplate").Parse(MinIntValidatorTemplate)
	tmpl.New("MaxIntValidatorTemplate").Parse(MaxIntValidatorTemplate)
	tmpl.New("MinStringValidatorTemplate").Parse(MinStringValidatorTemplate)
	tmpl.New("MaxStringValidatorTemplate").Parse(MaxStringValidatorTemplate)

	tmpl.ExecuteTemplate(out, "MainTemplate", g)
}
