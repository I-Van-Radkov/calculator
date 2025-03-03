package orchestrator

import (
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
	timeAdd := getOperationTimeFromEnv("+")
	timeSub := getOperationTimeFromEnv("-")
	timeMult := getOperationTimeFromEnv("*")
	timeDev := getOperationTimeFromEnv("/")

	return &config{
		timeAddition:        timeAdd,
		timeSubtraction:     timeSub,
		timeMultiplications: timeMult,
		timeDivisions:       timeDev,
	}
}

func getOperationTimeFromEnv(operation string) int {
	var envVarName string
	switch operation {
	case "+":
		envVarName = "TIME_ADDITION_MS"
	case "-":
		envVarName = "TIME_SUBTRACTION_MS"
	case "*":
		envVarName = "TIME_MULTIPLICATIONS_MS"
	case "/":
		envVarName = "TIME_DIVISIONS_MS"
	default:
		return 0
	}

	timeStr := os.Getenv(envVarName)
	timeMs, err := strconv.Atoi(timeStr)
	if err != nil {
		return 100
	}

	return timeMs
}
