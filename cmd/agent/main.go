package main

import (
	"log"

	"github.com/I-Van-Radkov/calculator/internal/agent"
)

func main() {
	worker := agent.NewAgent()
	log.Println("Starting Agent")
	worker.Run()
}
