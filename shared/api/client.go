package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (c *apiClient) CreatePallet(pallet *types.Pallet) (resp *Response, err error) {
	return c.SendData("POST", "pallets", pallet)
}

func (c *apiClient) GetAllPallets() (resp *Response, err error) {
	return c.GetData("GET", "pallets")
}

func (c *apiClient) SendData(method string, path string, payload interface{}) (resp *Response, err error) {
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(method, c.ParseUrl(path), bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	return c.callApi(req)
}

func (c *apiClient) GetData(method string, path string) (resp *Response, err error) {
	req, _ := http.NewRequest(method, c.ParseUrl(path), nil)
	req.Header.Set("Content-Type", "application/json")

	return c.callApi(req)
}

func (c *apiClient) callApi(req *http.Request) (resp *Response, err error) {
	response, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()
	body, _ := ioutil.ReadAll(response.Body)
	e := json.Unmarshal(body, &resp)

	if e != nil {
		return nil, e
	}
	return resp, nil
}

func (c *apiClient) ParseUrl(uri string) (url string) {
	return fmt.Sprintf("http://%s:%d/api/%s", c.Url, c.Port, uri)
}
