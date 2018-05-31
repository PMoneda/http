[![GoDoc](https://godoc.org/github.com/PMoneda/http?status.svg)](https://godoc.org/github.com/PMoneda/http)
[![Go Report Card](https://goreportcard.com/badge/github.com/PMoneda/http)](https://goreportcard.com/report/github.com/PMoneda/http)
# http
Http Wrapper over default http lib with high level abstraction and mock capabilities


Install
1) go get
```shell
$ go get github.com/PMoneda/http
```
2) dep
```shell
dep ensure -add github.com/PMoneda/http
```

Examples

GET Request

```golang
package main

import (
  "github.com/PMoneda/http"
  "fmt"
)

func main() {
  response, _ := http.Get("http://my-awesome-api.com");
  fmt.Println(response)
}


```

POST Request with headers

```golang
package main

import (
  "github.com/PMoneda/http"
  "fmt"
)

func main() {
  body := struct {
  ID string
  }{ID:"my-id"}

  response, _ := http.Post("http://my-awesome-api.com", body, http.Header{Key: "Content-Type", Value: "application/x-www-form-urlencoded"});
  fmt.Println(response)
}


```
Simple POST Request
```golang

package main

import (
  "github.com/PMoneda/http"
  "fmt"
)

func main() {
  body := struct {
  ID string
  }{ID:"my-id"}

  response, _ := http.Post("http://my-awesome-api.com", body);
  fmt.Println(response)
}

```

Testing with http Mock
```golang

func TestShouldRequestToAwesomeAPI(t *testing.T) {
  mock := http.ReponseMock{
    Method: "POST",
  }
  http.With(t, func(ctx *http.MockContext) {
    ctx.RegisterMock(&mock)
    response, err := http.Post("http://some-url","body")
    if err != nil {
       t.Fail()
    }
  })
}

```
