package cassgowary

import (
	"errors"
	"fmt"
)

var (
	DuplicateConstraintErr     = constraintError("unsatisfiable constraint")
	DuplicateEditVariableErr   = errors.New("duplicate edit variable")
	InternalSolverErr          = errors.New("internal solver error")
	NonLinearExpressionErr     = errors.New("non-linear expression")
	RequiredFailureErr         = errors.New("required failure")
	UnknownConstraintErr       = constraintError("unknown constraint")
	UnknownEditVariableErr     = errors.New("unknown edit variable")
	UnsatisfiableConstraintErr = constraintError("unsatisfiable constraint")
)

func constraintError(prefix string) func(c *Constraint) error {
	return func(c *Constraint) error {
		return fmt.Errorf("%s %s", prefix, c)
	}
}
