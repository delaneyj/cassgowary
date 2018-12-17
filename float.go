package cassgowary

import (
	"math"
)

const (
	FloatMin     = Float(-math.MaxFloat64)
	FloatMax     = Float(math.MaxFloat64)
	FloatEpsilon = 7/3.0 - 4/3.0 - 1
)

type Float float64

func (f Float) Raw() float64 {
	return float64(f)
}

func (f Float) Equals(other Float) bool {
	if math.Abs((f - other).Raw()) < FloatEpsilon {
		return true
	}
	return false
}

func (f Float) Strength() Strength {
	return Strength(f)
}

func (f Float) NearZero() bool {
	return f.Equals(0)
}

func (f Float) EqualsExpression(e *Expression) *Constraint {
	return e.EqualsFloat(f)
}

func (f Float) EqualsTerm(t Term) *Constraint {
	return t.EqualsFloat(f)
}

func (f Float) EqualsVariable(v *Variable) *Constraint {
	return v.EqualsFloat(f)
}

func (f Float) LessThanOrEqualToExpression(e *Expression) *Constraint {
	ne := NewExpression(f)
	c := ne.LessThanOrEqualTo(e)
	return c
}

func (f Float) LessThanOrEqualToTerm(t *Term) *Constraint {
	e := NewExpressionFrom(t)
	c := f.LessThanOrEqualToExpression(e)
	return c
}

func (f Float) LessThanOrEqualToVariable(v *Variable) *Constraint {
	t := NewTermFrom(v)
	c := f.LessThanOrEqualToTerm(t)
	return c
}

func (f Float) GreaterThanOrEqualToTerm(t *Term) *Constraint {
	e := NewExpression(f)
	c := e.GreaterThanOrEqualToTerm(t)
	return c
}

func (f Float) GreaterThanOrEqualToVariable(v *Variable) *Constraint {
	t := NewTermFrom(v)
	c := t.GreaterThanOrEqualToFloat(f)
	return c
}

func (f Float) ModifyStrength(c *Constraint) *Constraint {
	s := f.Strength()
	c2 := c.NewModifyStrength(s)
	return c2
}
