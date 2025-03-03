package orchestrator

import "sync"

type Expression struct {
	ID     string  `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}

type Task struct {
	ID            string  `json:"id"`
	Arg1          float64 `json:"arg1"`
	Arg2          float64 `json:"arg2"`
	Operation     string  `json:"operation"`
	OperationTime int     `json:"operation_time"`
	// idForArg1 и idForArg2 необходимы для создания подзадач
	idForArg1 string // если не пустой, то Arg1 - результат выражения с ID, которое хранится в idForArg1
	idForArg2 string // если не пустой, то Arg2 - результат выражения с ID, которое хранится в idForArg2
}

type Result struct {
	ID     string  `json:"id"`
	Result float64 `json:"result"`
}

type Orchestrator struct {
	config      config
	expressions map[string]Expression // мапа с выражениями
	tasks       chan Task             // канал для передачи подзадач для агента
	taskResults map[string]float64    // хранение результатов определенных подзадач с доступом по taskID
	results     chan Result           // канал для получения результатов вычисления агентами
	mu          sync.Mutex
}
