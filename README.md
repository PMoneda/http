# http
Http Helper Lib for golang with mock support


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

GET

```golang

response, err := http.Get("http://my-awesome-api.com");

```

POST with headers

```golang

body := struct {
  ID string
}{ID:"my-id"}

http.Post("http://my-awesome-api.com", body, map[string]string{"Content-Type":"application/json});

```
Simple POST
```golang

body := struct {
  ID string
}{ID:"my-id"}

response, err := http.Post("http://my-awesome-api.com", body);

```

With Mock
```golang
  mock := http.ReponseMock{
    Method: "POST",
  }
  http.With(func(ctx *http.MockContext) {
    ctx.RegisterMock(&mock)
    response, err := http.Post("http://some-url","body")
  })
```
