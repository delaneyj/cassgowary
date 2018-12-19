package cassgowary

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const ConstantVariableTestEpsilon = 1.0e-8

func TestFloatLessThanEqualTo(t *testing.T) {
	solver := NewSolver()
	x := NewVariable("x")
	f := Float(100)
	c := f.LessThanOrEqualToVariable(x)
	solver.AddConstraint(c)
	solver.UpdateVariables()
	assert.True(t, 100 <= x.Value)
	c2 := x.EqualsFloat(110)
	solver.AddConstraint(c2)
	solver.UpdateVariables()
	assert.InDelta(t, 110, x.Value.Raw(), ConstantVariableTestEpsilon)
}

func TestLessThanEqualToUnsatisfiable(t *testing.T) {
	solver := NewSolver()
	x := NewVariable("x")
	f := Float(100)
	c := f.LessThanOrEqualToVariable(x)
	err := solver.AddConstraint(c)
	assert.NoError(t, err)
	solver.UpdateVariables()

	assert.True(t, 100 <= x.Value)
	c2 := x.EqualsFloat(10)
	err = solver.AddConstraint(c2)
	assert.Error(t, err)
}

func TestGreaterThanEqualTo(t *testing.T) {
	solver := NewSolver()
	x := NewVariable("x")
	solver.AddConstraint(Float(100).GreaterThanOrEqualToVariable(x))
	solver.UpdateVariables()
	assert.True(t, 100 >= x.Value)
	solver.AddConstraint(x.EqualsFloat(90))
	solver.UpdateVariables()
	assert.InDelta(t, 90, x.Value.Raw(), ConstantVariableTestEpsilon)
}

func TestGreaterThanEqualToUnsatisfiable(t *testing.T) {
	solver := NewSolver()
	x := NewVariable("x")
	solver.AddConstraint(Float(100).GreaterThanOrEqualToVariable(x))
	solver.UpdateVariables()
	assert.True(t, 100 >= x.Value)
	err := solver.AddConstraint(x.EqualsFloat(110))
	assert.Error(t, err)
}
