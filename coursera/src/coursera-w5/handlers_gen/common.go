package main

// MainTemplate template for imports
const MainTemplate = `package {{.PackageName}}

import ({{range .Imports}}
	"{{.}}"{{end}}
)

func BuildHttpError(message string) string {
	apiError := make(map[string]interface{}, 0)
	apiError["error"] = message
	response, err := json.Marshal(apiError)
	if err != nil {
		panic(err)
	}

	return string(response)
}

func BuildResponse(res interface{}) string {
	response := make(map[string]interface{}, 0)
	response["error"] = ""
	response["response"] = res

	body, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	return string(body)
}

type Handler struct {
	process func(w http.ResponseWriter, r *http.Request)
}

func (h Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.process(w, r)
}

{{template "HTTPMethodCheckMiddlewareTemplate"}}
{{template "AuthCheckMiddlewareTemplate"}}
{{template "GetParamIntTemplate"}}
{{template "SetDefaultValueIntTemplate"}}
{{template "SetDefaultValueStringTemplate"}}
{{template "RequiredIntValidatorTemplate"}}
{{template "RequiredStringValidatorTemplate"}}
{{template "EnumValidatorTemplate"}}
{{template "MinIntValidatorTemplate"}}
{{template "MaxIntValidatorTemplate"}}
{{template "MinStringValidatorTemplate"}}
{{template "MaxStringValidatorTemplate"}}
`
