package cassgowary

import "math"

const Epsilon = 1.0e-12

func FloatEquals(f, other float64) bool {
	if math.Abs(f-other) < Epsilon {
		return true
	}
	return false
}

func FloatNearZero(f float64) bool {
	return FloatEquals(f, 0)
}

func FloatEqualsExpression(f float64, e *Expression) *Constraint {
	return e.EqualsFloat(f)
}

func FloatEqualsTerm(f float64, t Term) *Constraint {
	return t.EqualsFloat(f)
}

func FloatEqualsVariable(f float64, v *Variable) *Constraint {
	return v.EqualsFloat(f)
}

func FloatLessThanOrEqualToExpression(f float64, e *Expression) *Constraint {
	ne := NewExpression(f)
	c := ne.LessThanOrEqualTo(e)
	return c
}

func FloatLessThanOrEqualToTerm(f float64, t *Term) *Constraint {
	e := NewExpressionFrom(t)
	c := FloatLessThanOrEqualToExpression(f, e)
	return c
}

func FloatLessThanOrEqualToVariable(f float64, v *Variable) *Constraint {
	t := NewTermFrom(v)
	c := FloatLessThanOrEqualToTerm(f, t)
	return c
}

func FloatGreaterThanOrEqualToTerm(f float64, t *Term) *Constraint {
	e := NewExpression(f)
	c := e.GreaterThanOrEqualToTerm(t)
	return c
}

func FloatGreaterThanOrEqualToVariable(f float64, v *Variable) *Constraint {
	t := NewTermFrom(v)
	c := FloatGreaterThanOrEqualToTerm(f, t)
	return c
}

func FloatModifyStrength(f float64, c *Constraint) *Constraint {
	s := Strength(f)
	c2 := c.NewModifyStrength(s)
	return c2
}
