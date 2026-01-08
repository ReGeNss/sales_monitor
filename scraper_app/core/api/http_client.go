package api

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type HTTPClient interface {
	Get(url string) (*http.Response, error)
	Post(url string, body map[string]interface{}) (*http.Response, error)
}

func NewHTTPClient() HTTPClient {
	return &httpClientImpl{
		client: &http.Client{},
	}
}

type httpClientImpl struct {
	client *http.Client
}

func (c *httpClientImpl) Get(url string) (*http.Response, error) {
	return c.client.Get(url)
}

func (c *httpClientImpl) Post(url string, body map[string]interface{}) (*http.Response, error) {
	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	return c.client.Post(url, "application/json", bytes.NewReader(jsonBody))
}