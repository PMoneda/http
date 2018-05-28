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

//Put make a PUT request
func Put(url string, body interface{}, headers ...Header) (string, error) {
	return doRequest("PUT", url, body, headers...)
}

//Post make a POST request
func Post(url string, body interface{}, headers ...Header) (string, error) {
	return doRequest("POST", url, body, headers...)
}

//Get make a GET request
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

//GetJSON make a GET request and unmarshal response to JSON
func GetJSON(url string, obj interface{}) error {
	if response, err := Get(url); err != nil {
		return err
	} else {
		return json.Unmarshal([]byte(response), obj)
	}
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
		if j, err := json.Marshal(body); err != nil {
			return "", err
		} else {
			reqBody = string(j)
		}
	}
	if req, err := http.NewRequest(method, url, strings.NewReader(reqBody)); err != nil {
		return "", err
	} else {
		if headers == nil {
			req.Header["Content-Type"] = []string{"application/json"}
		} else {
			for _, header := range headers {
				req.Header[header.Key] = []string{header.Value}
			}
		}

		if resp, err := client.Do(req); err != nil {
			return "", err
		} else {
			if response, err := ioutil.ReadAll(resp.Body); err != nil {
				return "", err
			} else if resp.StatusCode >= 300 {
				return "", fmt.Errorf("Status %s: response: %s", resp.Status, string(response))
			} else {
				return string(response), nil
			}
		}
	}
}
