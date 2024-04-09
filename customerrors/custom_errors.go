package customerrors

type InvalidExpressionError struct {
	Expression string
}

func (m InvalidExpressionError) Error() string {
	return "invalid expression: " + m.Expression
}

func NewInvalidExpressionError(expression string) InvalidExpressionError {
	return InvalidExpressionError{Expression: expression}
}

type InvalidTokenError struct {
	Operator string
}

func NewInvalidTokenError(operator string) InvalidTokenError {
	return InvalidTokenError{Operator: operator}
}

func (m InvalidTokenError) Error() string {
	return "invalid token: " + m.Operator
}
