package main

import (
	"go/ast"
	"strings"
)

// SturctureExtractorTempalte extract request params into struct
const SturctureExtractorTempalte = `{{define "SturctureExtractorTempalte"}}
{{range $structureName, $name := .}}
func {{$structureName}}Extractor(r *http.Request) ({{$structureName}}, ApiError) {
	var requestData map[string][]string
	var paramName string
	var ok bool
	var err ApiError
	var result = {{$structureName}}{}
	if(r.Method == http.MethodGet) {
		requestData = r.URL.Query()
	} else {
		r.ParseForm()
		requestData = r.Form
	}

	{{range $field, $validators := $name.FieldList}}
		paramName = strings.ToLower("{{$field}}")		
		{{if eq $validators.FieldType "string"}}			
			{{if ne (index $validators.Validators "paramname") "" }}
				paramName = "{{ index $validators.Validators "paramname" }}"
			{{end}}		
			{{if ne (index $validators.Validators "required") "" }}
				ok, err = RequiredStringValidator(requestData, paramName)
				if !ok {
					return result, err
				}
			{{end}}					
			{{if ne (index $validators.Validators "default") "" }}
				{{$field}}, ok := SetDefaultStringValue(requestData, paramName, "{{index $validators.Validators "default"}}")				
				if !ok {
					{{$field}} = requestData[paramName][0]
				}
			{{else}}
				{{$field}} := requestData[paramName][0]	
			{{end}}	

			

			{{if ne (index $validators.Validators "min") "" }}
				ok, err = MinStringValidator({{$field}}, {{index $validators.Validators "min"}}, paramName)
				if !ok {
					return result, err
				}
			{{end}}
			{{if ne (index $validators.Validators "max") "" }}
				ok, err = MaxStringValidator({{$field}}, {{index $validators.Validators "max"}}, paramName)
				if !ok {
					return result, err
				}
			{{end}}
			{{if ne (index $validators.Validators "enum") "" }}
				ok, err = EnumValidator({{$field}}, "{{index $validators.Validators "enum"}}", paramName)
				if !ok {
					return result, err
				}
			{{end}}
		{{else}}
			{{if ne (index $validators.Validators "paramname") "" }}
				paramName = "{{ index $validators.Validators "paramname" }}"
			{{end}}		
			{{if ne (index $validators.Validators "required") "" }}
				ok, err = RequiredIntValidator(requestData, paramName)
				if !ok {
					return result, err
				}
			{{end}}
			{{$field}}, err := GetParamInt(requestData, paramName)
			if err != (ApiError{}) {
				return result, err
			}		
			{{if ne (index $validators.Validators "default") "" }}
				{{$field}} = SetDefaultIntValue(requestData, paramName, "{{index $validators.Validators "default"}}")				
			{{end}}	

			{{if ne (index $validators.Validators "min") "" }}
				ok, err = MinIntValidator({{$field}}, {{index $validators.Validators "min"}}, paramName)
				if !ok {
					return result, err
				}
			{{end}}
			{{if ne (index $validators.Validators "max") "" }}
				ok, err = MaxIntValidator({{$field}}, {{index $validators.Validators "max"}}, paramName)
				if !ok {
					return result, err
				}
			{{end}}
			{{if ne (index $validators.Validators "enum") "" }}
				ok, err = EnumValidator({{$field}}, "{{index $validators.Validators "enum"}}", paramName)
				if !ok {
					return result, err
				}
			{{end}}
		{{end}}
		result.{{$field}} = {{$field}}
	{{end}}

	return result, err
}
{{end}}
{{end}}`

// StructGenerator generate structs
type StructGenerator struct {
	GenerationLabel string
	Declararion     *ast.GenDecl
	structList      map[string]*StructField
	structType      string
}

// StructField types with fields of structs
type StructField struct {
	FieldList map[string]*StructValidator
}

// StructValidator types and validators
type StructValidator struct {
	FieldType  string
	Validators map[string]string
}

// NewStructGenerator set and return struct for generation structs
func (g *Generator) NewStructGenerator() *StructGenerator {
	g.structGenereator = StructGenerator{
		GenerationLabel: g.StructGenerationLabel,
		structList:      make(map[string]*StructField),
	}

	return &g.structGenereator
}

// FillStruct define and set struct fields and validators
func (sg *StructGenerator) FillStruct() {
	for _, spec := range sg.Declararion.Specs {
		declType, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		structType, ok := declType.Type.(*ast.StructType)
		if !ok {
			continue
		}

		structField := &StructField{
			FieldList: make(map[string]*StructValidator, 0),
		}

		for _, field := range structType.Fields.List {
			if field.Tag != nil && strings.HasPrefix(field.Tag.Value, sg.GenerationLabel) {
				structValidator := &StructValidator{}
				tags := strings.ReplaceAll(field.Tag.Value, sg.GenerationLabel, "")
				fieldType := ""
				ident, ok := field.Type.(*ast.Ident)
				if ok {
					fieldType = ident.String()
				}
				tags = strings.ReplaceAll(tags, "`", "")
				tags = strings.ReplaceAll(tags, "\"", "")
				splitedTags := strings.Split(tags, ",")
				validators := make(map[string]string)
				for _, tag := range splitedTags {
					keyValue := strings.Split(tag, "=")
					if len(keyValue) > 1 {
						validators[keyValue[0]] = keyValue[1]
					} else {
						validators[keyValue[0]] = "1"
					}
				}
				structValidator.FieldType = fieldType
				structValidator.Validators = validators
				structField.FieldList[field.Names[0].Name] = structValidator
			}
		}

		if len(structField.FieldList) > 0 {
			sg.structList[declType.Name.Name] = structField
		}
	}
}
