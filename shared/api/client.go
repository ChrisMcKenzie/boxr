package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Secret-Ironman/boxr/shared/types"
)

type apiClient struct {
	// url of the boxr api
	Url string
	// port of the boxr api
	Port int
	// whether or not to use ssl
	Ssl bool
	// http client
	client *http.Client
}

func NewApiClient(url string, port int, ssl bool) *apiClient {
	c := new(apiClient)
	c.Url = url
	c.Port = port
	c.Ssl = ssl
	c.client = &http.Client{}
	return c
}

func (c *apiClient) CreatePallet(pallet *types.Pallet) (resp *http.Response, err error) {
	return c.CallApi("POST", "pallet", pallet)
}

func (c *apiClient) CallApi(method string, path string, payload interface{}) (resp *http.Response, err error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, c.ParseUrl(path), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")
	if err != nil {
		return nil, err
	}
	return c.client.Do(req)
}

func (c *apiClient) ParseUrl(uri string) (url string) {
	return fmt.Sprintf("http://%s:%d/api/%s", c.Url, c.Port, uri)
}
