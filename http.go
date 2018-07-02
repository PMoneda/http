package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"sync"
	"testing"
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
	if currentContext.mocks != nil {
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
	if currentContext.mocks != nil {
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

//ReponseMock is mock configure struct
type ReponseMock struct {
	Method        string
	URL           string
	ReponseBody   string
	requestBody   string
	ResponseError error
	executed      int
}

//CalledTimes return how many times mock was called
func (resp *ReponseMock) CalledTimes() int {
	return resp.executed
}

//RequestBody returns request body that mock received
func (resp *ReponseMock) RequestBody() string {
	return resp.requestBody
}

var currentContext MockContext
var mutex sync.Mutex

//MockContext keep current state of mocks
type MockContext struct {
	mocks map[string]*ReponseMock
	test  *testing.T
}

//RegisterMock register a new mock response to current context
func (c *MockContext) RegisterMock(mock *ReponseMock) {
	key := fmt.Sprintf("%s:%s", mock.Method, mock.URL)
	c.mocks[key] = mock
}

//Fail fail test
func (c *MockContext) Fail() {
	c.test.Fail()
}

//With creates a new context to register mocks
func With(t *testing.T, callback func(*MockContext)) {
	mutex.Lock()
	defer mutex.Unlock()
	currentContext.test = t
	currentContext.mocks = make(map[string]*ReponseMock)
	callback(&currentContext)
	currentContext.mocks = nil
}

func getMock(method, url string) *ReponseMock {
	key := fmt.Sprintf("%s:%s", method, url)
	for k, v := range currentContext.mocks {
		if k == ":" {
			return v
		} else if v.URL == "*" && v.Method == method {
			return v
		} else if k == key {
			return v
		}
	}
	return nil
}

func doRequestMock(method, url string, body interface{}) (string, error) {
	mock := getMock(method, url)
	if mock == nil {
		return "", fmt.Errorf("mock for %s %s not defined exception", method, url)
	}
	mock.executed++
	jsonBody, _ := json.Marshal(body)
	mock.requestBody = string(jsonBody)
	return mock.ReponseBody, mock.ResponseError
}
