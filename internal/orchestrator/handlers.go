package orchestrator

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type CalculateRequest struct {
	Expression string `json:"expression"`
}

func (o *Orchestrator) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "{\"error\": \"Internal server error\"}", http.StatusInternalServerError)
		return
	}

	
}

func calculateHandler(w http.ResponseWriter, r *http.Request) {
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

	rpnData, err := infixToRPN(request.Expression)
	if err != nil {
		// TODO: допилить ошибку
	}

	
}

func expressionsHandler(w http.ResponseWriter, r *http.Request) {

}

func expressionsIDHandler(w http.ResponseWriter, r *http.Request) {

}
