package orchestrator

import (
	"fmt"
	"github.com/google/uuid"
	"strconv"
	"unicode"

	"github.com/artemkkkkk/DistributedCalculator/internal/config_manager"
	"github.com/artemkkkkk/DistributedCalculator/internal/custom_errors"
	"github.com/artemkkkkk/DistributedCalculator/internal/structs"
)

func tokenize(expression string) ([]string, error) {
	var tokens []string
	var current string

	for _, char := range expression {
		if unicode.IsDigit(char) || char == '.' {
			current += string(char)
		} else if isOperator(string(char)) || isParenthesis(string(char)) {
			if current != "" {
				tokens = append(tokens, current)
				current = ""
			}
			tokens = append(tokens, string(char))
		} else if unicode.IsSpace(char) {
			continue
		} else {
			return nil, custom_errors.InvalidExpressionError
		}
	}

	if current != "" {
		tokens = append(tokens, current)
	}

	return tokens, nil
}

func isOperator(token string) bool {
	return token == "+" || token == "-" || token == "*" || token == "/"
}

func isParenthesis(token string) bool {
	return token == "(" || token == ")"
}

func infixToPostfix(tokens []string) ([]string, error) {
	var output []string
	var stack []string

	for _, token := range tokens {
		if isNumber(token) {
			output = append(output, token)
		} else if isOperator(token) {
			for len(stack) > 0 && precedence(stack[len(stack)-1]) >= precedence(token) {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, token)
		} else if token == "(" {
			stack = append(stack, token)
		} else if token == ")" {
			for len(stack) > 0 && stack[len(stack)-1] != "(" {
				output = append(output, stack[len(stack)-1])
				stack = stack[:len(stack)-1]
			}
			if len(stack) == 0 {
				return nil, custom_errors.InvalidExpressionError
			}
			stack = stack[:len(stack)-1]
		} else {
			return nil, custom_errors.InvalidExpressionError
		}
	}

	for len(stack) > 0 {
		if isParenthesis(stack[len(stack)-1]) {
			return nil, custom_errors.InvalidExpressionError
		}
		output = append(output, stack[len(stack)-1])
		stack = stack[:len(stack)-1]
	}

	return output, nil
}

func isNumber(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

func precedence(op string) int {
	switch op {
	case "+", "-":
		return 1
	case "*", "/":
		return 2
	default:
		return 0
	}
}

func getOperationTime(operation string) int {
	switch operation {
	case "+":
		return config_manager.AdditionMS
	case "-":
		return config_manager.SubtractionMS
	case "*":
		return config_manager.MultiplicationMS
	case "/":
		return config_manager.DivisionMS
	default:
		return 0
	}
}

func buildTasks(rpn []string, expId string) ([]*structs.Task, error) {
	var stack []interface{}
	var tasks []*structs.Task

	for _, token := range rpn {
		if isNumber(token) {
			stack = append(stack, token)
		} else if isOperator(token) {
			if len(stack) < 2 {
				return nil, custom_errors.InvalidExpressionError
			}

			arg2 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			arg1 := stack[len(stack)-1]
			stack = stack[:len(stack)-1]

			taskID := "task" + uuid.New().String()
			operationTime := getOperationTime(token)

			task := &structs.Task{
				Id:              taskID,
				ExpId:           expId,
				Arg1:            fmt.Sprintf("%v", arg1),
				Arg2:            fmt.Sprintf("%v", arg2),
				Operation:       token,
				OperationTimeMS: fmt.Sprintf("%d", operationTime),
			}

			tasks = append(tasks, task)
			stack = append(stack, taskID)
		} else {
			return nil, custom_errors.InvalidExpressionError
		}
	}

	if len(stack) != 1 {
		return nil, custom_errors.InvalidExpressionError
	}

	return tasks, nil
}

func CreateTasks(expression string, expId string) ([]*structs.Task, error) {
	tokens, err := tokenize(expression)
	if err != nil {
		return nil, err
	}

	rpn, err := infixToPostfix(tokens)
	if err != nil {
		return nil, err
	}

	tasks, err := buildTasks(rpn, expId)
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
