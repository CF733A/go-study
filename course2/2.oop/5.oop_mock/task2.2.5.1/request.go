package main

import (
	"fmt"
	"io"
	"net/url"
)

func (e *Exmo) DoRequest(endpoint string, params url.Values) ([]byte, error) {
	reqURL := e.url + endpoint

	if params != nil {
		reqURL += "?" + params.Encode()
	}

	response, err := e.client.Get(reqURL)
	if err != nil {
		return nil, fmt.Errorf("ошибка отправки запроса %s", err)
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("ошибка чтения ответа %s", err)
	}

	return body, nil
}
