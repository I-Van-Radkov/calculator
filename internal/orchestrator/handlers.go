package orchestrator

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/google/uuid"
)

type CalculateRequest struct {
	Expression string `json:"expression"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func (o *Orchestrator) getTaskHandler(w http.ResponseWriter, _ *http.Request) {
	select {
	case task := <-o.tasks:
		log.Printf("New Task: ID: %s; Arg1: %v; Arg2: %v; Operation: %s; OperationTime: %d", task.ID, task.Arg1, task.Arg2, task.Operation, task.OperationTime)
		json.NewEncoder(w).Encode(task)
	default:
		http.Error(w, "No task available", http.StatusNotFound)
	}
}

func (o *Orchestrator) postTaskHandler(w http.ResponseWriter, r *http.Request) {
	var result Result
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid data", http.StatusUnprocessableEntity)
		return
	}

	o.results <- result
	w.WriteHeader(http.StatusOK)
}

func (o *Orchestrator) calculateHandler(w http.ResponseWriter, r *http.Request) {
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

	rpn, err := infixToRPN(request.Expression) // Перевод в RPN
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
		log.Printf("[ERR] infixToRPN: %v", err)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	id := uuid.New().String() // Генерация ID
	o.expressions[id] = Expression{ID: id, Status: "pending"}

	taskList, err := o.createTasksFromRPN(rpn) // Разбиение на таски
	if err != nil {
		errorResponse := ErrorResponse{
			Error: "Expression is not valid",
		}
		log.Printf("[ERR] createTasksFromRPN: %v", err)
		w.WriteHeader(http.StatusUnprocessableEntity)
		json.NewEncoder(w).Encode(errorResponse)
		return
	}

	go o.processTasksSequentially(id, taskList)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{"id": id})
	log.Printf("Успешный результат /calculate: id : %s", id)
}

func (o *Orchestrator) expressionsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "{\"error\": \"Internal server error\"}", http.StatusInternalServerError)
		return
	}

	o.mu.Lock()
	defer o.mu.Unlock()

	expressionList := make([]Expression, 0, len(o.expressions))
	for _, expr := range o.expressions {
		expressionList = append(expressionList, expr)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]Expression{"expressions": expressionList})
}

func (o *Orchestrator) expressionIDHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "{\"error\": \"Internal server error\"}", http.StatusInternalServerError)
		return
	}

	id := r.URL.Path[len("/api/v1/expressions/"):]

	o.mu.Lock()
	defer o.mu.Unlock()

	expr, found := o.expressions[id]
	if !found {
		http.Error(w, "{\"error\": \"Expression not found\"}", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]Expression{"expression": expr})
}
