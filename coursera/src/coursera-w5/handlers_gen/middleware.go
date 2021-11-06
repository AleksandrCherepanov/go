package main

// HTTPMethodCheckMiddlewareTemplate template for checking http method type
const HTTPMethodCheckMiddlewareTemplate = `{{define "HTTPMethodCheckMiddlewareTemplate"}}
func HTTPMethodCheckMiddleware(next http.Handler, method string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusNotAcceptable)
			io.WriteString(w, BuildHttpError("bad method"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
{{end}}`

// AuthCheckMiddlewareTemplate check authorization of request
const AuthCheckMiddlewareTemplate = `{{define "AuthCheckMiddlewareTemplate"}}
func AuthCheckMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("X-Auth");

		if auth == "" || auth != "100500" {
			w.WriteHeader(http.StatusForbidden)
			io.WriteString(w, BuildHttpError("unauthorized"))
			return
		}

		next.ServeHTTP(w, r)
	})
}
{{end}}`

// GetParamIntTemplate get param as int
const GetParamIntTemplate = `{{define "GetParamIntTemplate"}}
func GetParamInt(params map[string][]string, paramName string) (int, ApiError) {
	value := params[paramName][0]
	ae := ApiError{}
	result, err := strconv.Atoi(value)
	if err != nil {		
		ae.HTTPStatus = http.StatusBadRequest
		ae.Err = errors.New(paramName + " must be int")
		return -1, ae
	}

	return result, ae
}
{{end}}`

// SetDefaultValueIntTemplate set default value int
const SetDefaultValueIntTemplate = `{{define "SetDefaultValueIntTemplate"}}
func SetDefaultIntValue(params map[string][]string, paramName string, defaultValue int) int {
	value, err := strconv.Atoi(params[paramName][0])
	if err != nil {
		panic(err)
	}

	if value == 0 {
		return defaultValue
	}

	return value
}
{{end}}`

// SetDefaultValueStringTemplate set default value string
const SetDefaultValueStringTemplate = `{{define "SetDefaultValueStringTemplate"}}
func SetDefaultStringValue(params map[string][]string, paramName string, defaultValue string) (string, bool) {
	value, ok := params[paramName]
	endValue := ""
	if (ok) {
		endValue = value[0]
	}

	if endValue == "" {
		return defaultValue, true
	}

	return endValue, false
}
{{end}}`
