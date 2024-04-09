package main

import (
	"fmt"
	"log"
	"shuntingyard/parser"
	"shuntingyard/tree"
)

func Menu() {
	log.Println("Enter expression")
	var expression string
	_, err := fmt.Scanln(&expression)
	if err != nil {
		log.Println("Invalid expression")
		return
	}

	tokens, err := parser.Tokenize(expression)
	if err != nil {
		log.Println(err.Error())
		return
	}
	root, err := tree.InfixToAST(tokens)
	if err != nil {
		log.Println(err.Error())
		return
	}
	ExpressionMenu(root)
}

func ExpressionMenu(root tree.Node) {
	log.Println("[1] Evaluate an expression")
	log.Println("[2] Convert an expression to infix, prefix, and postfix notations")
	log.Println("[3] Exit")

	var choice int
	print("Enter choice: ")
	_, err := fmt.Scanf("%d", &choice)
	if err != nil {
		log.Println("Invalid choice")
		return
	}

	switch choice {
	case 1:
		result, err := root.Evaluate()
		if err != nil {
			log.Println(err.Error())
			return
		}
		log.Println("Result:", result)
	case 2:
		log.Println("Infix:", root.Infix())
		log.Println("Prefix:", root.Prefix())
		log.Println("Postfix:", root.Postfix())
	case 3:
		return
	}
}

func main() {
	defer func() {
		if err := recover(); err != nil {
			log.Println("panic occurred:", err)
		}
	}()

	log.SetFlags(0)
	for {
		Menu()
	}
}
