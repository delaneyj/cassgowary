package cassgowary

import (
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSimpleNew(t *testing.T) {
	solver := NewSolver()
	x := NewVariable("x")
	e := x.AddFloat(2)
	c := e.EqualsFloat(20.0)
	solver.AddConstraint(c)
	solver.UpdateVariables()
	assert.InDelta(t, 18.0, x.Value.Raw(), FloatEpsilon)
}

func TestSimple0(t *testing.T) {
	solver := NewSolver()
	x := NewVariable("x")
	y := NewVariable("y")
	solver.AddConstraint(x.EqualsFloat(20))
	solver.AddConstraint(x.AddFloat(2).Equals(y.AddFloat(10)))
	solver.UpdateVariables()
	assert.InDelta(t, 12, y.Value.Raw(), FloatEpsilon)
	assert.InDelta(t, 20, x.Value.Raw(), FloatEpsilon)
}

func TestSimple1(t *testing.T) {
	x := NewVariable("x")
	y := NewVariable("y")
	solver := NewSolver()
	err := solver.AddConstraint(x.Equals(y))
	assert.NoError(t, err)
	solver.UpdateVariables()
	assert.InDelta(t, x.Value.Raw(), y.Value.Raw(), FloatEpsilon)
}

func TestCasso1(t *testing.T) {
	x := NewVariable("x")
	y := NewVariable("y")
	solver := NewSolver()

	err := solver.AddConstraint(x.LessThanOrEqualTo(y))
	assert.NoError(t, err)
	err = solver.AddConstraint(y.EqualsExpression(x.AddFloat(3.0)))
	assert.NoError(t, err)
	err = solver.AddConstraint(x.EqualsFloat(10.0).NewModifyStrength(Weak))
	assert.NoError(t, err)
	err = solver.AddConstraint(y.EqualsFloat(10.0).NewModifyStrength(Weak))
	assert.NoError(t, err)

	solver.UpdateVariables()

	if math.Abs(x.Value.Raw()-10) < FloatEpsilon {
		assert.InDelta(t, 10, x.Value.Raw(), FloatEpsilon)
		assert.InDelta(t, 13, y.Value.Raw(), FloatEpsilon)
	} else {
		assert.InDelta(t, 7, x.Value.Raw(), FloatEpsilon)
		assert.InDelta(t, 10, y.Value.Raw(), FloatEpsilon)
	}
}

func TestAddDelete1(t *testing.T) {
	x := NewVariable("x")
	solver := NewSolver()

	err := solver.AddConstraint(x.LessThanOrEqualToFloat(100).NewModifyStrength(Weak))
	assert.NoError(t, err)

	solver.UpdateVariables()
	assert.InDelta(t, 100, x.Value.Raw(), FloatEpsilon)

	c10 := x.LessThanOrEqualToFloat(10)
	c20 := x.LessThanOrEqualToFloat(20)

	err = solver.AddConstraint(c10)
	assert.NoError(t, err)
	err = solver.AddConstraint(c20)
	assert.NoError(t, err)

	solver.UpdateVariables()

	assert.InDelta(t, 10, x.Value.Raw(), FloatEpsilon)

	solver.RemoveConstraint(c10)

	solver.UpdateVariables()

	assert.InDelta(t, 20, x.Value.Raw(), FloatEpsilon)

	solver.RemoveConstraint(c20)
	solver.UpdateVariables()

	assert.InDelta(t, 100, x.Value.Raw(), FloatEpsilon)

	c10again := x.LessThanOrEqualToFloat(10)

	err = solver.AddConstraint(c10again)
	assert.NoError(t, err)
	err = solver.AddConstraint(c10)
	assert.NoError(t, err)
	solver.UpdateVariables()

	assert.InDelta(t, 10, x.Value.Raw(), FloatEpsilon)

	solver.RemoveConstraint(c10)
	solver.UpdateVariables()
	assert.InDelta(t, 10, x.Value.Raw(), FloatEpsilon)

	solver.RemoveConstraint(c10again)
	solver.UpdateVariables()
	assert.InDelta(t, 100, x.Value.Raw(), FloatEpsilon)
}

func TestAddDelete2(t *testing.T) {
	x := NewVariable("x")
	y := NewVariable("y")
	solver := NewSolver()

	err := solver.AddConstraint(x.EqualsFloat(100).NewModifyStrength(Weak))
	assert.NoError(t, err)
	err = solver.AddConstraint(y.EqualsFloat(120).NewModifyStrength(Strong))
	assert.NoError(t, err)

	c10 := x.LessThanOrEqualToFloat(10)
	c20 := x.LessThanOrEqualToFloat(20)

	err = solver.AddConstraint(c10)
	assert.NoError(t, err)
	err = solver.AddConstraint(c20)
	assert.NoError(t, err)
	solver.UpdateVariables()

	assert.InDelta(t, 10, x.Value.Raw(), FloatEpsilon)
	assert.InDelta(t, 120, y.Value.Raw(), FloatEpsilon)

	err = solver.RemoveConstraint(c10)
	assert.NoError(t, err)
	solver.UpdateVariables()

	assert.InDelta(t, 20, x.Value.Raw(), FloatEpsilon)
	assert.InDelta(t, 120, y.Value.Raw(), FloatEpsilon)

	cxy := x.Multiply(2).EqualsVariable(y)
	err = solver.AddConstraint(cxy)
	assert.NoError(t, err)
	solver.UpdateVariables()

	assert.InDelta(t, 20, x.Value.Raw(), FloatEpsilon)
	assert.InDelta(t, 40, y.Value.Raw(), FloatEpsilon)

	err = solver.RemoveConstraint(c20)
	assert.NoError(t, err)
	solver.UpdateVariables()

	assert.InDelta(t, 60, x.Value.Raw(), FloatEpsilon)
	assert.InDelta(t, 120, y.Value.Raw(), FloatEpsilon)

	err = solver.RemoveConstraint(cxy)
	assert.NoError(t, err)
	solver.UpdateVariables()

	assert.InDelta(t, 100, x.Value.Raw(), FloatEpsilon)
	assert.InDelta(t, 120, y.Value.Raw(), FloatEpsilon)
}

func TestInconsistent1(t *testing.T) {
	x := NewVariable("x")
	solver := NewSolver()

	err := solver.AddConstraint(x.EqualsFloat(10))
	assert.NoError(t, err)
	err = solver.AddConstraint(x.EqualsFloat(5))
	assert.NoError(t, err)

	solver.UpdateVariables()
}

func TestInconsistent2(t *testing.T) {
	x := NewVariable("x")
	solver := NewSolver()

	err := solver.AddConstraint(x.GreaterThanOrEqualToFloat(10))
	assert.NoError(t, err)
	err = solver.AddConstraint(x.LessThanOrEqualToFloat(5))
	assert.Error(t, err)
}

func TestInconsistent3(t *testing.T) {
	w := NewVariable("w")
	x := NewVariable("x")
	y := NewVariable("y")
	z := NewVariable("z")
	solver := NewSolver()

	err := solver.AddConstraint(w.GreaterThanOrEqualToFloat(10))
	assert.NoError(t, err)
	err = solver.AddConstraint(x.GreaterThanOrEqualTo(w))
	assert.NoError(t, err)
	err = solver.AddConstraint(y.GreaterThanOrEqualTo(x))
	assert.NoError(t, err)
	err = solver.AddConstraint(z.GreaterThanOrEqualTo(y))
	assert.NoError(t, err)
	err = solver.AddConstraint(z.GreaterThanOrEqualToFloat(8))
	assert.NoError(t, err)
	err = solver.AddConstraint(z.LessThanOrEqualToFloat(4))
	assert.Error(t, err)
}
