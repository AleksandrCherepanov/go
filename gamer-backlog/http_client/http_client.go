package http_client

import (
	"fmt"
	"io"
	"net/http"
)

func Request(method, url string, body io.Reader) ([]byte, *error) {
	httpClient := http.DefaultClient

	httpRequest, httpRequestCreationError := http.NewRequest(method, url, body)
	if httpRequestCreationError != nil {
		return make([]byte, 0), &httpRequestCreationError
	}

	httpResponse, httpRequestExecutionError := httpClient.Do(httpRequest)
	if httpRequestExecutionError != nil {
		fmt.Println(httpRequestExecutionError.Error())
		return make([]byte, 0), &httpRequestCreationError
	}

	defer httpResponse.Body.Close()
	httpResonseBody, httpResponseBodyReadError := io.ReadAll(httpResponse.Body)
	if httpResponseBodyReadError != nil {
		return make([]byte, 0), &httpResponseBodyReadError
	}

	return httpResonseBody, nil
}
