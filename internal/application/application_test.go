package application_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/I-Van-Radkov/calculator/internal/application"
)

func TestCalcHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       interface{}
		statusCode int
		response   interface{}
	}{
		{
			name:       "valid",
			method:     "POST",
			body:       map[string]string{"expression": "3+5*2"},
			statusCode: http.StatusOK,
			response:   map[string]interface{}{"result": "13"},
		},
		{
			name:       "invalid",
			method:     "POST",
			body:       map[string]string{"expression": "4+8//]"},
			statusCode: http.StatusUnprocessableEntity,
			response:   map[string]interface{}{"error": "Expression is not valid"},
		},
		{
			name:       "empty",
			method:     "POST",
			body:       map[string]string{"expression": ""},
			statusCode: http.StatusUnprocessableEntity,
			response:   map[string]interface{}{"error": "Expression is not valid"},
		},
		{
			name:       "division by zero",
			method:     "POST",
			body:       map[string]string{"expression": "3/0"},
			statusCode: http.StatusUnprocessableEntity,
			response:   map[string]interface{}{"error": "Expression is not valid"},
		},
		{
			name:       "letters",
			method:     "POST",
			body:       map[string]string{"expression": "gahsd"},
			statusCode: http.StatusUnprocessableEntity,
			response:   map[string]interface{}{"error": "Expression is not valid"},
		},
		{
			name:       "method",
			method:     "GET",
			body:       nil,
			statusCode: http.StatusInternalServerError,
			response:   map[string]interface{}{"error": "Internal server error"},
		},
		{
			name:       "bad request",
			method:     "POST",
			body:       "invalid json",
			statusCode: http.StatusInternalServerError,
			response:   map[string]interface{}{"error": "Internal server error"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var reqBody []byte
			if test.body != nil {
				reqBody, _ = json.Marshal(test.body)
			}
			req := httptest.NewRequest(test.method, "/", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rec := httptest.NewRecorder()
			handler := http.HandlerFunc(application.CalculatorHandler)
			handler.ServeHTTP(rec, req)

			if rec.Code != test.statusCode {
				t.Errorf("Expected status code %d, got %d", test.statusCode, rec.Code)
			}

			var actualResponse map[string]interface{}
			json.Unmarshal(rec.Body.Bytes(), &actualResponse)
			expectedResponse := test.response.(map[string]interface{})
			for key, value := range expectedResponse {
				if actualResponse[key] != value {
					t.Errorf("Expected %s: %v, got %v", key, value, actualResponse[key])
				}
			}
		})
	}
}
