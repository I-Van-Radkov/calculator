package orchestrator

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/google/uuid"
)

// createTasksFromRPN создает слайс последовательных подзадач из RPN
func (o *Orchestrator) createTasksFromRPN(postfix []string) ([]Task, error) {
	var taskList []Task
	stack := []string{}

	if len(postfix) == 1 { // если одно число в выражении
		arg, _ := o.isFloatToken(postfix[0])
		taskList = append(taskList, Task{
			ID:            uuid.New().String(),
			Arg1:          arg,
			Arg2:          0,
			Operation:     "+",
			OperationTime: 100,
		})

		return taskList, nil
	}

	for _, token := range postfix {
		if isNumber(token) {
			stack = append(stack, token)
		} else {

			arg2 := stack[len(stack)-1]
			arg1 := stack[len(stack)-2]
			stack = stack[:len(stack)-2]

			operationTime := o.config.getOperationTime(token)
			taskID := uuid.New().String()

			idForArg1 := ""
			idForArg2 := ""
			var (
				arg1Float float64
				arg2Float float64
				isFloat   bool
			)

			arg1Float, isFloat = o.isFloatToken(arg1)
			if !isFloat {
				idForArg1 = arg1 // если в стэке токен - не число, то это id таски
			}

			arg2Float, isFloat = o.isFloatToken(arg2)
			if !isFloat {
				idForArg2 = arg2 // если в стэке токен - не число, то это id таски
			}

			task := Task{
				ID:            taskID,
				Arg1:          arg1Float,
				Arg2:          arg2Float,
				Operation:     token,
				OperationTime: operationTime,
				idForArg1:     idForArg1,
				idForArg2:     idForArg2,
			}

			taskList = append(taskList, task)

			stack = append(stack, taskID)
		}
	}

	return taskList, nil
}

// parseFloatOrFetchResult выясняет, токен из стэка RPN является числом или ID для таски
func (o *Orchestrator) isFloatToken(token string) (float64, bool) {
	if val, err := strconv.ParseFloat(token, 64); err == nil {
		return val, true
	}
	// Если это не число, то это taskID
	return 0, false
}

// processTasksSequentially последовательно вычисляет задачи через агентов
func (o *Orchestrator) processTasksSequentially(id string, taskList []Task) {
	for i := 0; i < len(taskList); i++ {
		task := taskList[i]

		switch {
		case task.idForArg1 != "" && task.idForArg2 == "": // если аргумент1 - это результат таски с ID
			task.Arg1 = o.taskResults[task.idForArg1]
			task.idForArg1 = ""
		case task.idForArg1 == "" && task.idForArg2 != "": // если аргумент2 - это результат таски с ID
			task.Arg2 = o.taskResults[task.idForArg2]
			task.idForArg2 = ""
		case task.idForArg1 != "" && task.idForArg2 != "": // если оба аргумента - это результаты тасок
			task.Arg1 = o.taskResults[task.idForArg1]
			task.idForArg1 = ""
			task.Arg2 = o.taskResults[task.idForArg2]
			task.idForArg2 = ""
		}

		o.tasks <- task

		// Ожидание результата

		result := <-o.results
		if result.ID == task.ID {
			o.mu.Lock()

			o.taskResults[task.ID] = result.Result

			o.mu.Unlock()
		}
	}

	// Обновление статуса выражения
	finalTaskID := taskList[len(taskList)-1].ID
	if o.taskResults[finalTaskID] != 0 {
		o.mu.Lock()

		exprId := o.expressions[id]
		exprId.Status = "completed"
		exprId.Result = o.taskResults[finalTaskID]
		o.expressions[id] = exprId

		log.Printf("Получен результат для выражения с ID %s: %v", id, exprId.Result)

		o.mu.Unlock()
	}
}

// getOperationTime возвращает время выполнения операции
func (c *config) getOperationTime(operation string) int {
	switch operation {
	case "+":
		return c.timeAddition
	case "-":
		return c.timeSubtraction
	case "*":
		return c.timeMultiplications
	case "/":
		return c.timeDivisions
	default:
		return 100 //дефолтное
	}
}

// getTokens разделяет строку на токены для преобразование записи в RPN
func getTokens(expression string) ([]string, error) {
	tokens := []string{}

	strRez := ""
	isLastDigit := false
	lastBracket := ""

	for _, el := range expression {
		if el >= '0' && el <= '9' || el == '.' {
			if isLastDigit {
				strRez += string(el)
			} else {
				isLastDigit = true
				if strRez != "" {
					tokens = append(tokens, strRez)
				}
				strRez = string(el)
			}
			lastBracket = ""
		} else if el == '-' || el == '+' || el == '*' || el == '/' {
			if !isLastDigit && lastBracket == "(" && el == '-' { // Для отрицательных чисел
				isLastDigit = true
				lastBracket = ""
				if strRez != "" {
					tokens = append(tokens, strRez)
				}
				strRez = string(el)
			} else if lastBracket == "(" && el != '-' {
				fmt.Println("q")
				return nil, fmt.Errorf("expression is not valid")
			} else {
				if !isLastDigit && lastBracket == "" && el != '-' {
					fmt.Println(isLastDigit, string(el), "1")
					return nil, fmt.Errorf("expression is not valid")
				}

				lastBracket = ""
				isLastDigit = false
				if strRez != "" {
					tokens = append(tokens, strRez)
				}
				strRez = string(el)
			}
		} else if el == '(' || el == ')' {
			if !isLastDigit && lastBracket == "" && el == ')' {
				return nil, fmt.Errorf("expression is not valid")
			}

			if el == '(' {
				if lastBracket == "(" || (!isLastDigit && lastBracket == "") {
					isLastDigit = false
					lastBracket = "("
					if strRez != "" {
						tokens = append(tokens, strRez)
					}
					strRez = string(el)
				} else {
					return nil, fmt.Errorf("expression is not valid")
				}
			} else {
				if lastBracket == "(" && !isLastDigit {
					return nil, fmt.Errorf("expression is not valid")
				}

				isLastDigit = false
				lastBracket = ")"
				if strRez != "" {
					tokens = append(tokens, strRez)
				}
				strRez = string(el)
			}
		} else {
			return nil, fmt.Errorf("expression is not valid")
		}
	}

	if strRez != "" {
		tokens = append(tokens, strRez)
	}

	if !isLastDigit && lastBracket == "" { // если выражение заканчивается не числом
		return nil, fmt.Errorf("expression is not valid")
	}

	return tokens, nil
}

// infixToRPN переводит строку в запись RPN
func infixToRPN(expression string) ([]string, error) {
	var rpn []string
	var stack []string
	operatorPrecedence := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2}

	expression = strings.ReplaceAll(expression, ",", ".")
	expression = strings.ReplaceAll(expression, " ", "")

	// Разделение выражения на токены
	tokens, err := getTokens(expression)
	if err != nil {
		return nil, fmt.Errorf(err.Error())
	}

	//Проверка на деление на 0
	if checkDivByZero(tokens) {
		return nil, fmt.Errorf("division by zero")
	}

	for _, token := range tokens {
		// Если токен - число, добавляем в RPN
		if isNumber(token) {
			rpn = append(rpn, token)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			// Извлекаем операторы из стека до левой скобки
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				operator := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				rpn = append(rpn, operator)
			}
			if len(stack) == 0 {
				return nil, fmt.Errorf("expression is not valid")
			}
			// Удаляем левую скобку
			stack = stack[:len(stack)-1]
		} else { // Оператор
			// Сравниваем приоритеты с операторами в стеке
			for len(stack) > 0 && getPrecedence(token, operatorPrecedence) <= getPrecedence(stack[len(stack)-1], operatorPrecedence) {
				operator := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				rpn = append(rpn, operator)
			}
			stack = append(stack, token)
		}
	}

	// Переносим оставшиеся операторы из стека в RPN
	for len(stack) > 0 {
		operator := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		rpn = append(rpn, operator)
	}

	return rpn, nil
}

// Проверка деления на 0
func checkDivByZero(tokens []string) bool {
	for i := 0; i < len(tokens)-1; i++ {
		if tokens[i] == "/" && tokens[i+1] == "0" {
			return true
		}
	}

	return false
}

// Проверка, является ли строка числом
func isNumber(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

// Получение приоритета оператора
func getPrecedence(operator string, operatorPrecedence map[string]int) int {
	precedence, ok := operatorPrecedence[operator]
	if !ok {
		return 0
	}
	return precedence
}
