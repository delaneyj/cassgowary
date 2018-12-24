package cassgowary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVariableLessThanEqualTo(t *testing.T) {
	x := NewVariable("x")
	solver := NewSolver()
	err := solver.AddConstraint(x.LessThanOrEqualToFloat(100))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.True(t, x.Value <= 100)
	err = solver.AddConstraint(x.EqualsFloat(90))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.InDelta(t, x.Value, 90, Epsilon)
}

func TestVariableLessThanEqualToUnsatisfiable(t *testing.T) {
	x := NewVariable("x")
	solver := NewSolver()
	err := solver.AddConstraint(x.LessThanOrEqualToFloat(100))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.True(t, x.Value <= 100)
	err = solver.AddConstraint(x.EqualsFloat(110))
	assert.Error(t, err)
	solver.UpdateVariables()
	assert.InDelta(t, x.Value, 110, Epsilon)
}

func TestVariableGreaterThanEqualTo(t *testing.T) {
	x := NewVariable("x")
	solver := NewSolver()
	err := solver.AddConstraint(x.GreaterThanOrEqualToFloat(100))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.True(t, x.Value >= 100)
	err = solver.AddConstraint(x.EqualsFloat(110))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.InDelta(t, x.Value, 110, Epsilon)
}

func TestVariableGreaterThanEqualToUnsatisfiable(t *testing.T) {
	x := NewVariable("x")
	solver := NewSolver()
	err := solver.AddConstraint(x.GreaterThanOrEqualToFloat(100))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.True(t, x.Value >= 100)
	err = solver.AddConstraint(x.EqualsFloat(90))
	assert.Error(t, err)
	solver.UpdateVariables()
}
