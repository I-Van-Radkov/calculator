package orchestrator

import (
	"fmt"
	"strconv"
	"strings"
)

func getTokens(expression string) []string {
	tokens := []string{}

	strRez := ""
	isLastDigit := false

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
		} else {
			if !isLastDigit && el == '-' { // Для отрицательных чисел
				isLastDigit = true
				if strRez != "" {
					tokens = append(tokens, strRez)
				}
				strRez = string(el)
			} else {
				isLastDigit = false
				if strRez != "" {
					tokens = append(tokens, strRez)
				}
				strRez = string(el)
			}
		}
	}
	if strRez != "" {
		tokens = append(tokens, strRez)
	}
	return tokens
}

func infixToRPN(expression string) ([]string, error) {
	var rpn []string
	var stack []string
	operatorPrecedence := map[string]int{"+": 1, "-": 1, "*": 2, "/": 2}

	expression = strings.ReplaceAll(expression, ",", ".")
	expression = strings.ReplaceAll(expression, " ", "")

	// Разделение выражения на токены
	tokens := getTokens(expression)

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
