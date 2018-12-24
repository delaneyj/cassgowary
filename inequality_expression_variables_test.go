package cassgowary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpressionLessThanEqualTo(t *testing.T) {
	x := NewVariable("x")
	solver := NewSolver()
	err := solver.AddConstraint(NewExpression(100).LessThanOrEqualToVariable(x))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.True(t, 100 <= x.Value)
	err = solver.AddConstraint(x.EqualsFloat(110))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.InDelta(t, x.Value, 110, Epsilon)
}

func TestExpressionLessThanEqualToUnsatisfiable(t *testing.T) {
	x := NewVariable("x")
	solver := NewSolver()
	err := solver.AddConstraint(NewExpression(100).LessThanOrEqualToVariable(x))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.True(t, x.Value <= 100)
	err = solver.AddConstraint(x.EqualsFloat(10))
	assert.Error(t, err)
	solver.UpdateVariables()
	assert.InDelta(t, 10, x.Value, Epsilon)
}

func TestExpressionGreaterThanEqualTo(t *testing.T) {
	x := NewVariable("x")
	solver := NewSolver()
	err := solver.AddConstraint(NewExpression(100).GreaterThanOrEqualToVariable(x))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.True(t, 100 >= x.Value)
	err = solver.AddConstraint(x.EqualsFloat(90))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.InDelta(t, x.Value, 90, Epsilon)
}

func TestExpressionGreaterThanEqualToUnsatisfiable(t *testing.T) {
	x := NewVariable("x")
	solver := NewSolver()
	err := solver.AddConstraint(NewExpression(100).GreaterThanOrEqualToVariable(x))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.True(t, 100 >= x.Value)
	err = solver.AddConstraint(x.EqualsFloat(110))
	assert.Error(t, err)
	solver.UpdateVariables()
	assert.InDelta(t, 110, x.Value, Epsilon)
}
