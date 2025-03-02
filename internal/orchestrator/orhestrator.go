package orchestrator

import (
	"log"
	"net/http"
	"os"
	"strconv"
)

type config struct {
	timeAddition        int
	timeSubtraction     int
	timeMultiplications int
	timeDivisions       int
}

func newConfig() *config {
	ta, err := strconv.Atoi(os.Getenv("TIME_ADDITION_MS"))
	if err != nil {

	}

	ts, err := strconv.Atoi(os.Getenv("TIME_SUBTRACTION_MS"))
	if err != nil {

	}

	tm, err := strconv.Atoi(os.Getenv("TIME_MULTIPLICATIONS_MS"))
	if err != nil {

	}

	td, err := strconv.Atoi(os.Getenv("TIME_DIVISIONS_MS"))
	if err != nil {

	}

	return &config{
		timeAddition:        ta,
		timeSubtraction:     ts,
		timeMultiplications: tm,
		timeDivisions:       td,
	}
}

type Orchestrator struct {
	config config
	mu sync.Mutex
}

func NewOrchestrator() *Orchestrator {
	return &Orchestrator{
		config: *newConfig(),
	}
}

func (o *Orchestrator) Run() error {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/v1/calculate", calculateHandler)
	mux.HandleFunc("/api/v1/expressions", expressionsHandler)
	mux.HandleFunc("/api/v1/expressions/", expressionsIDHandler)

	mux.HandleFunc("/internal/task", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			o.getTaskHandler(w, r)
		} else if r.Method == http.MethodPost {
			//o.postTaskHandler(w, r)
		} else {
			http.Error(w, `{"error":"Wrong Method"}`, http.StatusMethodNotAllowed)
		}
	})

	log.Println("Оркестратор запущен на адресе :8080")
	return http.ListenAndServe(":8080", mux)
}
