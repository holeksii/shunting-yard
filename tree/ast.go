package tree

import (
	"fmt"
	"shuntingyard/customerrors"
	"strconv"
	"strings"
)

var precedence = map[string]int{
	"+": 1, "-": 1,
	"*": 2, "/": 2,
	"%": 3, "^": 3,
}

// Define Node interface as a pointer interface
type Node interface {
	GetValue() string
	isOperator() bool
	Evaluate() (float64, error)
	Infix() string
	Prefix() string
	Postfix() string
}

// OperandNode represents a numeric value in the expression
type OperandNode struct {
	Value string
}

// GetValue returns the value of the operand node
func (o *OperandNode) GetValue() string {
	return o.Value
}

// isOperator returns false for OperandNode
func (o *OperandNode) isOperator() bool {
	return false
}

// Evaluate returns the numeric value of the operand node
func (o *OperandNode) Evaluate() (float64, error) {
	// parse and return the float64 value of the operand node
	return strconv.ParseFloat(o.Value, 64)
}

// Infix returns the infix notation of the operand node
func (o *OperandNode) Infix() string {
	return o.Value
}

// Prefix returns the prefix notation of the operand node
func (o *OperandNode) Prefix() string {
	return o.Value
}

// Postfix returns the postfix notation of the operand node
func (o *OperandNode) Postfix() string {
	return o.Value
}

// OperatorNode represents an operator in the expression
type OperatorNode struct {
	Value     string
	LeftNode  *Node
	RightNode *Node
}

// GetValue returns the value of the operator node
func (o *OperatorNode) GetValue() string {
	return o.Value
}

// isOperator returns true for OperatorNode
func (o *OperatorNode) isOperator() bool {
	return true
}

// Evaluate returns the numeric value of the operator node
func (o *OperatorNode) Evaluate() (float64, error) {
	left, err := (*o.LeftNode).Evaluate()
	if err != nil {
		return 0, err
	}
	right, err := (*o.RightNode).Evaluate()
	if err != nil {
		return 0, err
	}
	switch o.Value {
	case "+":
		return left + right, nil
	case "-":
		return left - right, nil
	case "*":
		return left * right, nil
	case "/":
		if right == 0 {
			panic("zero division occured")
		}
		return left / right, nil
	case "%":
		return float64(int(left) % int(right)), nil
	case "^":
		return float64(int(left) ^ int(right)), nil
	default:
		return 0, fmt.Errorf("invalid operator %s", o.Value)
	}
}

// Infix returns the infix notation of the operator node
func (o *OperatorNode) Infix() string {
	return fmt.Sprintf("(%s %s %s)", (*o.LeftNode).Infix(), o.Value, (*o.RightNode).Infix())
}

// Prefix returns the prefix notation of the operator node
func (o *OperatorNode) Prefix() string {
	return fmt.Sprintf("%s %s %s", o.Value, (*o.LeftNode).Prefix(), (*o.RightNode).Prefix())
}

// Postfix returns the postfix notation of the operator node
func (o *OperatorNode) Postfix() string {
	return fmt.Sprintf("%s %s %s", (*o.LeftNode).Postfix(), (*o.RightNode).Postfix(), o.Value)
}

// Function to convert infix notation to AST
func InfixToAST(tokens []string) (Node, error) {
	var stack []Node
	var output []Node

	for _, token := range tokens {
		if token == "(" {
			stack = append(stack, &OperatorNode{Value: token})
		} else if token == ")" {
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				if top.GetValue() == "(" {
					break
				}
				output = append(output, top)
			}
		} else if opPrecedence, isOp := precedence[token]; isOp {
			for len(stack) > 0 {
				top := stack[len(stack)-1]
				if top.isOperator() && precedence[top.GetValue()] >= opPrecedence {
					stack = stack[:len(stack)-1]
					output = append(output, top)
				} else {
					break
				}
			}
			stack = append(stack, &OperatorNode{Value: token})
		} else {
			operandNode := &OperandNode{Value: token}
			output = append(output, operandNode)
		}
	}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		output = append(output, top)
	}

	// Convert the output stack to a single node (the root of the AST)
	astStack := []Node{}
	for _, token := range output {
		if token.isOperator() {
			if len(astStack) < 2 {
				return nil, customerrors.InvalidExpressionError{Expression: strings.Join(tokens, " ")}
			}
			right := astStack[len(astStack)-1]
			astStack = astStack[:len(astStack)-1]
			left := astStack[len(astStack)-1]
			astStack = astStack[:len(astStack)-1]

			operatorNode := &OperatorNode{Value: token.GetValue()}
			operatorNode.RightNode = &right
			operatorNode.LeftNode = &left

			astStack = append(astStack, operatorNode)
		} else {
			astStack = append(astStack, token)
		}
	}

	if len(astStack) != 1 {
		return nil, customerrors.InvalidExpressionError{Expression: strings.Join(tokens, " ")}
	}

	return astStack[0], nil
}



// print node reccursively
func printNode(node Node, level int) {
	
