package custom_errors

import "errors"

var (
	ZeroDivisionError      = errors.New("ZeroDivisionError")
	InvalidExpressionError = errors.New("InvalidExpressionError")
	ExpressionNotFound     = errors.New("ExpressionNotFound")
)
