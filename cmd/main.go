package main

import (
	"log"

	application "github.com/I-Van-Radkov/calculator/internal/application"
)

func main() {
	app := application.New()

	err := app.Run()
	if err != nil {
		log.Fatal("Ошибка запуска сервера: %v", err)
	}
}
