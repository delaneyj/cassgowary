package tests

import (
	"testing"

	. "github.com/delaneyj/cassgowary"
	"github.com/stretchr/testify/assert"
)

func TestVariableVariableLessThanEqualTo(t *testing.T) {
	solver := NewSolver()

	x := NewVariable("x")
	y := NewVariable("y")

	err := solver.AddConstraint(x.EqualsFloat(100))
	assert.NoError(t, err)
	err = solver.AddConstraint(x.LessThanOrEqualTo(y))
	assert.NoError(t, err)

	solver.UpdateVariables()
	assert.True(t, x.Value <= 100)
	err = solver.AddConstraint(x.EqualsFloat(90))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.InDelta(t, x.Value.Raw(), 90, FloatEpsilon)
}

func TestVariableVariableLessThanEqualToUnsatisfiable(t *testing.T) {
	solver := NewSolver()

	x := NewVariable("x")
	y := NewVariable("y")

	solver.AddConstraint(y.EqualsFloat(100))
	solver.AddConstraint(x.LessThanOrEqualTo(y))

	solver.UpdateVariables()
	assert.True(t, x.Value <= 100)
	solver.AddConstraint(x.EqualsFloat(110))
	solver.UpdateVariables()
}

func TestVariableVariableGreaterThanEqualTo(t *testing.T) {
	solver := NewSolver()

	x := NewVariable("x")
	y := NewVariable("y")

	solver.AddConstraint(y.EqualsFloat(100))
	solver.AddConstraint(x.GreaterThanOrEqualTo(y))

	solver.UpdateVariables()
	assert.True(t, x.Value >= 100)
	solver.AddConstraint(x.EqualsFloat(110))
	solver.UpdateVariables()
	assert.InDelta(t, x.Value.Raw(), 110, FloatEpsilon)
}

func TestVariableVariableGreaterThanEqualToUnsatisfiable(t *testing.T) {
	solver := NewSolver()

	x := NewVariable("x")
	y := NewVariable("y")

	solver.AddConstraint(y.EqualsFloat(100))

	solver.AddConstraint(x.GreaterThanOrEqualTo(y))
	solver.UpdateVariables()
	assert.True(t, x.Value >= 100)
	solver.AddConstraint(x.EqualsFloat(90))
	solver.UpdateVariables()
}
