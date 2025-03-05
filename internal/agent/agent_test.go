package agent_test

import (
	"testing"

	"github.com/I-Van-Radkov/calculator/internal/agent"
)

func TestCalc(t *testing.T) {
	testCases := []struct {
		name           string
		arg1           float64
		arg2           float64
		operation      string
		expectedResult float64
	}{
		{
			name:           "addition",
			arg1:           23,
			arg2:           9,
			operation:      "+",
			expectedResult: 32,
		},
		{
			name:           "addition with negative number",
			arg1:           -2,
			arg2:           3,
			operation:      "+",
			expectedResult: 1,
		},
		{
			name:           "multiplication",
			arg1:           2,
			arg2:           10,
			operation:      "*",
			expectedResult: 20,
		},
		{
			name:           "division",
			arg1:           8,
			arg2:           4,
			operation:      "/",
			expectedResult: 2,
		},
		{
			name:           "addition with float",
			arg1:           9.0,
			arg2:           0.7,
			operation:      "+",
			expectedResult: 9.7,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			val := agent.PerformOperation(testCase.arg1, testCase.arg2, testCase.operation)

			if val != testCase.expectedResult {
				t.Fatalf("%f should be equal %f", val, testCase.expectedResult)
			}
		})
	}
}
