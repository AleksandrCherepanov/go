package tinkoff

import (
	"errors"
	"io"
	"net/http"
)

const errorPrefix = "Portfolio: "

type API struct {
	token     string
	accountId string
	url       string
}

func New(token string, accountId string, url string) API {
	return API{
		token:     token,
		accountId: accountId,
		url:       url,
	}
}

func (api API) Portfolio() (string, error) {
	url := api.url + "/portfolio?brokerAccountId=" + api.accountId

	httpRequest, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return "", errors.New(errorPrefix + err.Error())
	}

	httpRequest.Header.Set("Authorization", "Bearer "+api.token)

	client := http.Client{}
	response, err := client.Do(httpRequest)
	if err != nil {
		return "", errors.New(errorPrefix + err.Error())
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", errors.New(errorPrefix + err.Error())
	}

	return string(body), nil
}
