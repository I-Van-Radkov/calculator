package agent

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

const (
	orchestratorURL = "http://localhost:8080"
)

func NewAgent() *Agent {
	cp, err := strconv.Atoi(os.Getenv("COMPUTING_POWER"))
	if err != nil || cp < 1 {
		cp = 1
	}

	return &Agent{
		computingPower: cp,
	}
}

func (a *Agent) Run() {
	for i := 0; i < a.computingPower; i++ {
		log.Printf("Starting worker %d", i)
		go a.worker(i)
	}

	select {}
}

func performOperation(arg1, arg2 float64, operation string) float64 {
	switch operation {
	case "+":
		return arg1 + arg2
	case "-":
		return arg1 - arg2
	case "*":
		return arg1 * arg2
	case "/":
		return arg1 / arg2
	default:
		return 0
	}
}

func (a *Agent) worker(workerID int) {
	for {
		resp, err := http.Get(orchestratorURL + "/internal/task")
		if err != nil {
			resp.Body.Close()

			log.Printf("Worker %d: error getting task: %v", workerID, err)
			time.Sleep(time.Second)
			continue
		}

		if resp.StatusCode == http.StatusOK {
			var task Task

			err := json.NewDecoder(resp.Body).Decode(&task)
			resp.Body.Close()
			if err != nil {
				log.Printf("Worker %d: error decoding task: %v", workerID, err)
				time.Sleep(1 * time.Second)
				continue
			}

			/*fmt.Println("task:", task)
			log.Printf("worker %d получил выражение: %v %s %v", workerID, task.Arg1, task.Operation, task.Arg2)*/
			resultCalculation := performOperation(task.Arg1, task.Arg2, task.Operation)

			resultData := Result{
				ID:     task.ID,
				Result: resultCalculation,
			}
			resultDataBytes, _ := json.Marshal(resultData)

			http.Post(orchestratorURL+"/internal/task", "application/json", bytes.NewBuffer(resultDataBytes))
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}
