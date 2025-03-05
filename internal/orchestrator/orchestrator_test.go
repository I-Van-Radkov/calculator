package orchestrator_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/I-Van-Radkov/calculator/internal/orchestrator"
)

func TestCalcHandler(t *testing.T) {
	tests := []struct {
		name       string
		method     string
		body       interface{}
		statusCode int
		isError    interface{}
	}{
		{
			name:       "valid",
			method:     "POST",
			body:       map[string]string{"expression": "3+5*2"},
			statusCode: http.StatusCreated,
			isError:    false,
		},
		{
			name:       "invalid",
			method:     "POST",
			body:       map[string]string{"expression": "4+8//]"},
			statusCode: http.StatusUnprocessableEntity,
			isError:    true,
		},
		{
			name:       "invalid",
			method:     "POST",
			body:       map[string]string{"expression": "3*(2+)"},
			statusCode: http.StatusUnprocessableEntity,
			isError:    true,
		},
		{
			name:       "division by zero",
			method:     "POST",
			body:       map[string]string{"expression": "3/0"},
			statusCode: http.StatusInternalServerError,
			isError:    true,
		},
		{
			name:       "letters",
			method:     "POST",
			body:       map[string]string{"expression": "gahsd"},
			statusCode: http.StatusUnprocessableEntity,
			isError:    true,
		},
		{
			name:       "method",
			method:     "GET",
			body:       nil,
			statusCode: http.StatusInternalServerError,
			isError:    true,
		},
		{
			name:       "bad request",
			method:     "POST",
			body:       "invalid json",
			statusCode: http.StatusInternalServerError,
			isError:    true,
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
			handler := http.HandlerFunc(orchestrator.NewOrchestrator().CalculateHandler)
			handler.ServeHTTP(rec, req)

			if rec.Code != test.statusCode {
				t.Errorf("Expected status code %d, got %d", test.statusCode, rec.Code)
			}
		})
	}
}
