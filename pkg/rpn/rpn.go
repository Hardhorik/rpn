package rpn

import (
	"errors"
	"fmt"
	"strconv"
	"unicode"
)

var priorities = map[string]int{"(": 0, "+": 1, "-": 1, "*": 2, "/": 2, "~": 3}

func toPostfix(infix string) string {

	postfix := ""
	var stack []string

	if len(infix) == 0 {
		return "error"
	}

	last_symb := string(infix[len(infix)-1])
	ops_counter := 0

	for i := 0; i < len(infix); i++ {
		c := string(infix[i])

		if _, err := strconv.Atoi(last_symb); err != nil {
			if last_symb != ")" {
				return "error"
			}
		}

		if _, isOperator := priorities[c]; isOperator {
			ops_counter++

			if ops_counter > 1 {
				return "error"
			}
		} else {
			ops_counter = 0
		}

		if unicode.IsDigit(rune(infix[i])) {
			for ; i < len(infix) && (unicode.IsDigit(rune(infix[i])) || infix[i] == '.' || infix[i] == ','); i++ {
				postfix += string(infix[i])
			}
			postfix += " "
			i--
		} else if c == "(" {
			stack = append(stack, c)
		} else if c == ")" {
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top == "(" {
					break
				}
				postfix += top
				//postfix += top + " "
			}
		} else if _, isOperator := priorities[c]; isOperator {
			if c == "-" {
				if i == 0 || infix[i-1] == '(' || infix[i-1] == ',' || infix[i-1] == '.' {
					c = "~"
				}
			}

			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if top == "(" || priorities[top] < priorities[c] {
					break
				}
				postfix += top
				//postfix += top + " "
				stack = stack[:len(stack)-1]
			}
			stack = append(stack, c)
		}
	}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		postfix += top
		//postfix += top + " "
	}

	return postfix
}

func execute(operation string, firstNum float64, secondNum float64) float64 {
	output := 0.00
	switch operation {
	case "+":
		output = firstNum + secondNum
	case "-":
		output = firstNum - secondNum
	case "*":
		output = firstNum * secondNum
	case "/":
		if secondNum != 0 {
			output = firstNum / secondNum
		} else {
			output = 0
		}
	default:
		output = 0.00
	}
	return output
}

func Calc(expression string) (float64, error) {
	var locals []float64

	postfix := toPostfix(expression)

	switch postfix {
	case "error":
		return 0, errors.New("ERROR")
	}

	for i := 0; i < len(postfix); i++ {
		c := string(postfix[i])

		if unicode.IsDigit(rune(postfix[i])) {
			num := ""

			for ; i < len(postfix) && (unicode.IsDigit(rune(postfix[i])) || postfix[i] == '.' || postfix[i] == ','); i++ {
				num += string(postfix[i])
			}
			i--

			val, err := strconv.ParseFloat(num, 64)
			if err != nil {
				return 0, err
			}
			locals = append(locals, val)
		} else if _, isOperator := priorities[c]; isOperator {
			var first, second float64

			if len(locals) > 0 {
				second = locals[len(locals)-1]
				locals = locals[:len(locals)-1]
			} else {
				second = 0
			}

			if len(locals) > 0 {
				first = locals[len(locals)-1]
				locals = locals[:len(locals)-1]
			} else {
				first = 0
			}

			result := execute(c, first, second)
			locals = append(locals, result)
		}
	}

	if len(locals) == 0 {
		return 0, fmt.Errorf("результат не найден")
	}
	return locals[len(locals)-1], nil
}
