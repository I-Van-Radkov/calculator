package orchestrator

import (
	"log"
	"net/http"
	"sync"
)

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		config:      *newConfig(),
		expressions: make(map[string]Expression),
		tasks:       make(chan Task, 100),
		taskResults: make(map[string]float64),
		results:     make(chan Result, 100),
		mu:          sync.Mutex{},
	}
}

func (o *Orchestrator) Run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/calculate", o.CalculateHandler)
	mux.HandleFunc("/api/v1/expressions", o.ExpressionsHandler)
	mux.HandleFunc("/api/v1/expressions/", o.ExpressionIDHandler)

	mux.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			o.GetTaskHandler(w, r)
		} else if r.Method == http.MethodPost {
			o.PostTaskHandler(w, r)
		} else {
			http.Error(w, "{\"error\": \"Internal server error\"}", http.StatusInternalServerError)
		}
	})

	log.Println("Оркестратор запущен на адресе :8080")
	return http.ListenAndServe(":8080", mux)
}
