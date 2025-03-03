package agent

type Agent struct {
	computingPower int
}

type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
}

// Result представляет результат выполнения задачи
type Result struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}
