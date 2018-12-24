package cassgowary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFloatLessThanEqualTo(t *testing.T) {
	solver := NewSolver()
	x := NewVariable("x")
	c := FloatLessThanOrEqualToVariable(100, x)
	err := solver.AddConstraint(c)
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.True(t, 100 <= x.Value)
	err = solver.AddConstraint(x.EqualsFloat(110))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.InDelta(t, 110, x.Value, Epsilon)
}

func TestFloatLessThanEqualToUnsatisfiable(t *testing.T) {
	x := NewVariable("x")
	solver := NewSolver()
	err := solver.AddConstraint(FloatLessThanOrEqualToVariable(100, x))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.True(t, 100 <= x.Value)
	err = solver.AddConstraint(x.EqualsFloat(10))
	assert.Error(t, err)
	solver.UpdateVariables()
	assert.InDelta(t, 10, x.Value, Epsilon)
}

func TestFloatGreaterThanEqualTo(t *testing.T) {
	solver := NewSolver()
	x := NewVariable("x")
	err := solver.AddConstraint(FloatGreaterThanOrEqualToVariable(100, x))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.True(t, 100 >= x.Value)
	err = solver.AddConstraint(x.EqualsFloat(90))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.InDelta(t, 90, x.Value, Epsilon)
}

func TestFloatGreaterThanEqualToUnsatisfiable(t *testing.T) {
	x := NewVariable("x")
	solver := NewSolver()
	err := solver.AddConstraint(FloatGreaterThanOrEqualToVariable(100, x))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.True(t, 100 >= x.Value)
	err = solver.AddConstraint(x.EqualsFloat(110))
	assert.Error(t, err)
	solver.UpdateVariables()
	assert.InDelta(t, 110, x.Value, Epsilon)
}
