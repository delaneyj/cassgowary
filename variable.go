package cassgowary

import (
	"fmt"
)

type Variable struct {
	Name  string
	Value float64
}

func NewVariable(name string) *Variable {
	return NewVariableWithValue(name, 0)
}

func NewVariableWithValue(name string, value float64) *Variable {
	return &Variable{
		Name:  name,
		Value: value,
	}
}

func (v *Variable) String() string {
	return fmt.Sprintf("%s:%f", v.Name, v.Value)
}

// Variable multiply, divide, and unary invert
func (v *Variable) Multiply(coefficient float64) *Term {
	return NewTerm(v, coefficient)
}

func (v *Variable) Divide(denominator float64) *Term {
	return v.Multiply(1 / denominator)
}

func (v *Variable) Negate() *Term {
	return v.Multiply(-1)
}

func (v *Variable) AddExpression(e *Expression) *Expression {
	return e.AddVariable(v)
}

func (v *Variable) AddTerm(t *Term) *Expression {
	return t.AddVariable(v)
}

func (v *Variable) Add(other *Variable) *Expression {
	return NewTermFrom(v).AddVariable(other)
}

func (v *Variable) AddFloat(constant float64) *Expression {
	return NewTermFrom(v).AddFloat(constant)
}

func (v *Variable) SubtractExpression(e *Expression) *Expression {
	return v.AddExpression(e.Negate())
}

func (v *Variable) SubtractTerm(t Term) *Expression {
	return v.AddTerm(t.Negate())
}

func (v *Variable) Subtract(other *Variable) *Expression {
	return v.AddTerm(other.Negate())
}

func (v *Variable) SubtractFloat(constant float64) *Expression {
	return v.AddFloat(-constant)
}

// Variable relations
func (v *Variable) EqualsExpression(e *Expression) *Constraint {
	return e.EqualsVariable(v)
}

func (v *Variable) EqualsTerm(t Term) *Constraint {
	return t.EqualsVariable(v)
}

func (v *Variable) Equals(other *Variable) *Constraint {
	t := NewTermFrom(v)
	c := t.EqualsVariable(other)
	return c
}

func (v *Variable) EqualsFloat(constant float64) *Constraint {
	t := NewTermFrom(v)
	c := t.EqualsFloat(constant)
	return c
}

func (v *Variable) LessThanOrEqualToExpression(e *Expression) *Constraint {
	t := NewTermFrom(v)
	c := t.LessThanOrEqualToExpression(e)
	return c
}

func (v *Variable) LessThanOrEqualToTerm(t *Term) *Constraint {
	t2 := NewTermFrom(v)
	c := t2.LessThanOrEqualTo(t)
	return c
}

func (v *Variable) LessThanOrEqualTo(other *Variable) *Constraint {
	t := NewTermFrom(v)
	c := t.LessThanOrEqualToVariable(other)
	return c
}

func (v *Variable) LessThanOrEqualToFloat(constant float64) *Constraint {
	t := NewTermFrom(v)
	c := t.LessThanOrEqualToFloat(constant)
	return c
}

func (v *Variable) GreaterThanOrEqualToExpression(e *Expression) *Constraint {
	t := NewTermFrom(v)
	c := t.GreaterThanOrEqualToExpression(e)
	return c
}

func (v *Variable) GreaterThanOrEqualToTerm(t Term) *Constraint {
	return t.GreaterThanOrEqualToVariable(v)
}

func (v *Variable) GreaterThanOrEqualTo(other *Variable) *Constraint {
	t := NewTermFrom(v)
	c := t.GreaterThanOrEqualToVariable(other)
	return c
}

func (v *Variable) GreaterThanOrEqualToFloat(constant float64) *Constraint {
	t := NewTermFrom(v)
	c := t.GreaterThanOrEqualToFloat(constant)
	return c
}
