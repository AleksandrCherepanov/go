package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
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

type Middleware struct {
	process func(w http.ResponseWriter, r *http.Request)
}

func (m Middleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	m.process(w, r)
}

func (h *MyApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/user/profile":
		HTTPMethodCheckMiddleware := Middleware{}
		HTTPMethodCheckMiddleware.process = h.wrapperProfile
		h.wrapperProfile(w, r)
	case "/user/create":
		HTTPMethodCheckMiddleware := Middleware{}
		HTTPMethodCheckMiddleware.process = h.wrapperCreate
		h.wrapperCreate(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		error := BuildHttpError("unknown method")
		io.WriteString(w, error)
	}
}

func (h *OtherApi) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/user/create":
		HTTPMethodCheckMiddleware := Middleware{}
		HTTPMethodCheckMiddleware.process = h.wrapperCreate
		h.wrapperCreate(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
		error := BuildHttpError("unknown method")
		io.WriteString(w, error)
	}
}

func (h *MyApi) wrapperProfile(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	fmt.Println("QUERY")
	fmt.Println(r.URL.Query())

	ok, apiError := RequiredValidator(r.URL.Query(), "login")
	if !ok {
		w.WriteHeader(apiError.HTTPStatus)
		io.WriteString(w, BuildHttpError(apiError.Error()))
		return
	}

	params := ProfileParams{}
	res, err := h.Profile(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, BuildHttpError("bad request"))
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, BuildResponse(res))

}

func (h *MyApi) wrapperCreate(w http.ResponseWriter, r *http.Request) {

	// func CheckPOSTMethodMiddleware(next http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		if "POST" != r.Method {
	// 			w.WriteHeader(http.StatusNotAcceptable)
	// 			io.WriteString(w, BuildHttpError("bad method"))
	// 			return
	// 		}

	// 		next.ServeHTTP(w, r)
	// 	})
	// }
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	ok, apiError := RequiredValidator(r.Form, "login")
	if !ok {
		w.WriteHeader(apiError.HTTPStatus)
		io.WriteString(w, BuildHttpError(apiError.Error()))
		return
	}

	params := CreateParams{}
	res, err := h.Create(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, BuildHttpError("bad request"))
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, BuildResponse(res))

}

func (h *OtherApi) wrapperCreate(w http.ResponseWriter, r *http.Request) {

	// func CheckPOSTMethodMiddleware(next http.Handler) http.Handler {
	// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		if "POST" != r.Method {
	// 			w.WriteHeader(http.StatusNotAcceptable)
	// 			io.WriteString(w, BuildHttpError("bad method"))
	// 			return
	// 		}

	// 		next.ServeHTTP(w, r)
	// 	})
	// }
	err := r.ParseForm()
	if err != nil {
		panic(err)
	}

	fmt.Println(r.Form)

	params := OtherCreateParams{}
	res, err := h.Create(r.Context(), params)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		io.WriteString(w, BuildHttpError("bad request"))
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, BuildResponse(res))

}

func RequiredValidator(request map[string][]string, paramName string) (bool, ApiError) {
	param, ok := request[paramName]
	fmt.Println(request)
	if !ok {
		apiError := ApiError{}
		apiError.HTTPStatus = http.StatusBadRequest
		apiError.Err = errors.New(paramName + " must me not empty")
		return false, apiError
	}

	for _, value := range param {
		if value == "" {
			apiError := ApiError{}
			apiError.HTTPStatus = http.StatusBadRequest
			apiError.Err = errors.New(paramName + " must me not empty")
			return false, apiError
		}
	}

	return true, ApiError{}
}
