package application

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

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

func CalculatorHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "{\"error\": \"Failed to read request body\"}", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	var request CalculateRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		http.Error(w, "{\"error\": \"Invalid request body\"}", http.StatusBadRequest)
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
