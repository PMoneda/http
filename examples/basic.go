package examples

import (
	"github.com/PMoneda/http"
)

func postToPruu() string {
	response, _ := http.Post("https://pruu.herokuapp.com/dump/test", "my body message", http.Header{Key: "X-Custom-Header", Value: "Custom-Value"})
	return response
}

func getFromPruu() string {
	response, _ := http.Get("https://pruu.herokuapp.com/dump/test")
	return response
}
