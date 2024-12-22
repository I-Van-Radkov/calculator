package main

import (
	"log"
	"net/http"

	application "github.com/I-Van-Radkov/calculator/internal/application"
)

func main() {
	http.HandleFunc("/api/v1/calculate", application.LoggingMiddleware(application.CalculatorHandler))

	log.Println("Сервер запущен на адресе :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("Ошибка запуска сервера: %v", err)
	}
}
