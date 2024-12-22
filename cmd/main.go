package main

import (
	"net/http"
)

func main() {
	http.HandleFunc("/api/v1/calculate", CalculatorHandler)
	http.ListenAndServe(":8080", nil)
}
