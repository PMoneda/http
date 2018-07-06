package examples

import (
	"github.com/PMoneda/http"
)

func postToPruu() *http.HTTPResponse {
	response, _ := http.Post("https://pruu.herokuapp.com/dump/test", "my body message", http.Header{Key: "X-Custom-Header", Value: "Custom-Value"})
	return response
}

func getFromPruu() *http.HTTPResponse {
	response, _ := http.Get("https://pruu.herokuapp.com/dump/test")
	return response
}
