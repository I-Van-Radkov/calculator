package main

import (
	"log"

	orch "github.com/I-Van-Radkov/calculator/internal/orchestrator"
)

func main() {
	orchestrator := orch.NewOrchestrator()
	if err := orchestrator.Run(); err != nil {
		log.Fatalf("Ошибка запуска оркестратора: %v", err)
	}
}
