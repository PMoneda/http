package http

import (
	"encoding/json"
	"fmt"
	"testing"
)

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

var currentContext *MockContext

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
	ctx := new(MockContext)
	ctx.test = t
	ctx.mocks = make(map[string]*ReponseMock)
	currentContext = ctx
	callback(ctx)
	currentContext = nil
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
