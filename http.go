package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

//Header is a http header protocol abstraction
type Header struct {
	Key   string
	Value string
}

//Put do a PUT request
func Put(url string, body interface{}, headers ...Header) (string, error) {
	return doRequest("PUT", url, body, headers...)
}

//Post do a POST request
func Post(url string, body interface{}, headers ...Header) (string, error) {
	return doRequest("POST", url, body, headers...)
}

//Get do a GET request
func Get(url string) (string, error) {
	if currentContext != nil {
		return doRequestMock("GET", url, nil)
	}
	if resp, err := http.Get(url); err != nil {
		return "", err
	} else if response, err := ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	} else {
		return string(response), nil
	}
}

//GetJSON do a GET request and unmarshal response to JSON
func GetJSON(url string, obj interface{}) error {
	response, err := Get(url)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(response), obj)

}

func doRequest(method, url string, body interface{}, headers ...Header) (string, error) {
	if currentContext != nil {
		return doRequestMock(method, url, body)
	}
	return httpRequest(method, url, body, headers...)
}

func httpRequest(method, url string, body interface{}, headers ...Header) (string, error) {
	client := http.DefaultClient
	reqBody := ""
	switch v := body.(type) {
	case string:
		reqBody = v
	default:
		j, err := json.Marshal(body)
		if err != nil {
			return "", err
		}
		reqBody = string(j)
	}
	req, err := http.NewRequest(method, url, strings.NewReader(reqBody))
	if err != nil {
		return "", err
	}
	if headers == nil {
		req.Header["Content-Type"] = []string{"application/json"}
	} else {
		for _, header := range headers {
			req.Header[header.Key] = []string{header.Value}
		}
	}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	if response, err := ioutil.ReadAll(resp.Body); err != nil {
		return "", err
	} else if resp.StatusCode >= 300 {
		return "", fmt.Errorf("Status %s: response: %s", resp.Status, string(response))
	} else {
		return string(response), nil
	}
}
