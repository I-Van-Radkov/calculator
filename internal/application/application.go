package application

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	calculation "github.com/I-Van-Radkov/calculator/pkg"
)

type CalculateRequest struct {
	Expression string `json:"expression"`
}

type CalculateResponse struct {
	Result string `json:"result"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type Application struct {
}

func New() *Application {
	return &Application{}
}

func (a *Application) Run() error {
	http.HandleFunc("/api/v1/calculate", LoggingMiddleware(CalculatorHandler))

	log.Println("Сервер запущен на адресе :8080")
	return http.ListenAndServe(":8080", nil)
}

func CalculatorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "{\"error\": \"Internal server error\"}", http.StatusInternalServerError)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "{\"error\": \"Internal server error\"}", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	var request CalculateRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, "{\"error\": \"Internal server error\"}", http.StatusInternalServerError)
		return
	}

	result, err := calculation.Calc(request.Expression)
	if err != nil {
		var errorResponse ErrorResponse
		if err.Error() == "expression is not valid" {
			errorResponse = ErrorResponse{
				Error: "Expression is not valid",
			}
			w.WriteHeader(http.StatusUnprocessableEntity)
		} else {
			errorResponse = ErrorResponse{
				Error: "Internal server error",
			}
			w.WriteHeader(http.StatusInternalServerError)
		}
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	response := CalculateResponse{
		Result: fmt.Sprintf("%v", result),
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)
}

func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		next.ServeHTTP(w, r)

		elapsed := time.Since(start)
		log.Printf("Обработка запроса заняла %s", elapsed)
	}
}
