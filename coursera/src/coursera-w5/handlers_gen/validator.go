package main

// RequiredIntValidatorTemplate check if param int is required
const RequiredIntValidatorTemplate = `{{define "RequiredIntValidatorTemplate"}}
func RequiredIntValidator(params map[string][]string, paramName string) (bool, ApiError) {
	value, err := strconv.Atoi(params[paramName][0])
	if err != nil {
		panic(err)
	}
	error := ApiError{}
	if value == 0 {	
		error.HTTPStatus = http.StatusBadRequest
		error.Err = errors.New(paramName + " must me not empty")
		return false, error
	}

	return true, error
}
{{end}}`

// RequiredStringValidatorTemplate check if param string is required
const RequiredStringValidatorTemplate = `{{define "RequiredStringValidatorTemplate"}}
func RequiredStringValidator(params map[string][]string, paramName string) (bool, ApiError) {
	value, ok := params[paramName]
	endValue := ""
	if ok {
		endValue = value[0]
	}

	error := ApiError{}
	if endValue == "" {	
		error.HTTPStatus = http.StatusBadRequest
		error.Err = errors.New(paramName + " must me not empty")
		return false, error
	}

	return true, error
}
{{end}}`

// EnumValidatorTemplate check if param is in range
const EnumValidatorTemplate = `{{define "EnumValidatorTemplate"}}
func EnumValidator(value, defaultValue, paramName string) (bool, ApiError) {
	error := ApiError{}
	
	allowed := strings.Split(defaultValue, "|")

	found := false;
	for _, v := range allowed {
		if value == v {
			found = true
			break;
		}
	}

	if found {
		return true, error 
	}

	error.HTTPStatus = http.StatusBadRequest
	error.Err = errors.New(paramName + " must be one of [" + strings.Join(allowed, ", ") + "]")
	return false, error
}
{{end}}`

// MinIntValidatorTemplate check if value grate or equals
const MinIntValidatorTemplate = `{{define "MinIntValidatorTemplate"}}
func MinIntValidator(value, min int, paramName string) (bool, ApiError) {	
	error := ApiError{}

	if value >= min {
		return true, error
	}

	error.HTTPStatus = http.StatusBadRequest
	error.Err = errors.New(paramName + " must be >= " + fmt.Sprintf("%v", min))
	return false, error
}
{{end}}`

// MaxIntValidatorTemplate check if value less or equals
const MaxIntValidatorTemplate = `{{define "MaxIntValidatorTemplate"}}
func MaxIntValidator(value, max int, paramName string) (bool, ApiError) {
	error := ApiError{}

	if value <= max {
		return true, error
	}

	error.HTTPStatus = http.StatusBadRequest
	error.Err = errors.New(paramName + " must be <= " + fmt.Sprintf("%v", max))
	return false, error
}
{{end}}`

// MinStringValidatorTemplate check if value grate or equals
const MinStringValidatorTemplate = `{{define "MinStringValidatorTemplate"}}
func MinStringValidator(value string, min int, paramName string) (bool, ApiError) {
	error := ApiError{}

	if len(value) >= min {
		return true, error
	}

	error.HTTPStatus = http.StatusBadRequest
	error.Err = errors.New(paramName + " len must be >= " + fmt.Sprintf("%v", min))
	return false, error
}
{{end}}`

// MaxStringValidatorTemplate check if value less or equals
const MaxStringValidatorTemplate = `{{define "MaxStringValidatorTemplate"}}
func MaxStringValidator(value string, max int, paramName string) (bool, ApiError) {	
	error := ApiError{}

	if len(value) <= max {
		return true, error
	}

	error.HTTPStatus = http.StatusBadRequest
	error.Err = errors.New(paramName + " len must be <= " + fmt.Sprintf("%v", max))
	return false, error
}
{{end}}`
