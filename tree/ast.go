package tree

import (
	"fmt"
	"math"
	"shuntingyard/customerrors"
	"strconv"
	"strings"
)

var precedence = map[string]int{
	"+": 1, "-": 1,
	"*": 2, "/": 2,
	"%": 3, "^": 3,
}

type Node interface {
	GetValue() string
	isOperator() bool
	Evaluate() (float64, error)
	Infix() string
	Prefix() string
	Postfix() string
}

type OperandNode struct {
	Value string
}

func (o *OperandNode) GetValue() string {
	return o.Value
}

func (o *OperandNode) isOperator() bool {
	return false
}

func (o *OperandNode) Evaluate() (float64, error) {
	return strconv.ParseFloat(o.Value, 64)
}

func (o *OperandNode) Infix() string {
	return o.Value
}

func (o *OperandNode) Prefix() string {
	return o.Value
}

func (o *OperandNode) Postfix() string {
	return o.Value
}

type OperatorNode struct {
	Value     string
	LeftNode  *Node
	RightNode *Node
}

func (o *OperatorNode) GetValue() string {
	return o.Value
}

func (o *OperatorNode) isOperator() bool {
	return true
}

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
		return math.Pow(left, right), nil
	default:
		return 0, fmt.Errorf("invalid operator %s", o.Value)
	}
}

func (o *OperatorNode) Infix() string {
	return fmt.Sprintf("(%s %s %s)", (*o.LeftNode).Infix(), o.Value, (*o.RightNode).Infix())
}

func (o *OperatorNode) Prefix() string {
	return fmt.Sprintf("%s %s %s", o.Value, (*o.LeftNode).Prefix(), (*o.RightNode).Prefix())
}

func (o *OperatorNode) Postfix() string {
	return fmt.Sprintf("%s %s %s", (*o.LeftNode).Postfix(), (*o.RightNode).Postfix(), o.Value)
}

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

	astStack := []Node{}
	for _, token := range output {
		if token.isOperator() {
			if len(astStack) < 2 {
				print(1)
				return nil, customerrors.NewInvalidExpressionError(strings.Join(tokens, " "))
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
		return nil, customerrors.NewInvalidExpressionError(strings.Join(tokens, " "))
	}

	return astStack[0], nil
}

func printNode(node Node, level int) {
	if level > 0 {
		fmt.Print(strings.Repeat("│  ", level-1))
		fmt.Print("├─ ")
	}
	if node.isOperator() {
		opNode := node.(*OperatorNode)
		fmt.Println(opNode.GetValue())
		printNode(*opNode.LeftNode, level+1)
		printNode(*opNode.RightNode, level+1)
	} else {
		fmt.Println(node.GetValue())
	}
}

func PrintTree(root Node) {
	printNode(root, 0)
}
