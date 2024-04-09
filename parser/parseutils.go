package parser

import (
	"shuntingyard/customerrors"
	"strconv"
	"strings"
	"text/scanner"
)

var operators = []string{"+", "-", "*", "/", "%", "^", "(", ")"}

func isNumeric(token string) bool {
	_, err := strconv.ParseFloat(token, 64)
	return err == nil
}

func isValidToken(token string) bool {
	for _, op := range operators {
		if token == op {
			return true
		}
	}
	return isNumeric(token)
}

func Tokenize(input string) ([]string, error) {
	var s scanner.Scanner
	s.Init(strings.NewReader(input))

	var tok rune
	var result = make([]string, 0)
	for tok != scanner.EOF {
		tok = s.Scan()
		value := strings.TrimSpace(s.TokenText())
		if len(value) > 0 {
			result = append(result, s.TokenText())
		}
	}

	for _, token := range result {
		if !isValidToken(token) {
			return nil, customerrors.NewInvalidTokenError(token)
		}
	}
	return result, nil
}
