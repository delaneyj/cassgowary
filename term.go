package cassgowary

import (
	"fmt"
	"log"
)

type Term struct {
	Variable    *Variable
	Coefficient Float
}
type Terms []*Term

func NewTerm(v *Variable, coefficient Float) *Term {
	return &Term{
		Variable:    v,
		Coefficient: coefficient,
	}
}

func NewTermFrom(variable *Variable) *Term {
	return NewTerm(variable, 1.0)
}

func (t *Term) Value() Float {
	return t.Coefficient * t.Value()
}

func (t *Term) String() string {
	return fmt.Sprintf(
		`variable: (%s) coefficient:%f`,
		t.Variable, t.Coefficient,
	)
}

func (t *Term) Multiply(coefficient Float) *Term {
	return NewTerm(t.Variable, t.Coefficient*coefficient)
}

func (t *Term) Divide(denominator Float) *Term {
	return t.Multiply(1 / denominator)
}

func (t *Term) Negate() *Term {
	return t.Multiply(-1)
}

func (t *Term) AddExpression(e *Expression) *Expression {
	return e.AddTerm(t)
}

func (t *Term) Add(other *Term) *Expression {
	terms := Terms{t, other}
	return NewExpressionFrom(terms...)
}

func (t *Term) AddVariable(v *Variable) *Expression {
	return t.Add(NewTermFrom(v))
}

func (t *Term) AddFloat(constant Float) *Expression {
	return NewExpression(constant, t)
}

func (t *Term) SubtractExpression(e *Expression) *Expression {
	return e.Negate().AddTerm(t)
}

func (t *Term) Subtract(other *Term) *Expression {
	return t.Add(other.Negate())
}

func (t *Term) subtract(v *Variable) *Expression {
	return t.Add(v.Negate())
}

func (t *Term) SubtractFloat(constant Float) *Expression {
	return t.AddFloat(-constant)
}

func (t *Term) EqualsExpression(e *Expression) *Constraint {
	return e.EqualsTerm(t)
}

func (t *Term) Equals(other *Term) *Constraint {
	e := NewExpressionFrom(t)
	c := e.EqualsTerm(other)
	return c
}

func (t *Term) EqualsVariable(v *Variable) *Constraint {
	return NewExpressionFrom(t).EqualsVariable(v)
}

func (t *Term) EqualsFloat(constant Float) *Constraint {
	log.Print(t)
	e := NewExpressionFrom(t)
	c := e.EqualsFloat(constant)
	return c
}

func (t *Term) LessThanOrEqualToExpression(e *Expression) *Constraint {
	te := NewExpressionFrom(t)
	c := te.LessThanOrEqualTo(e)
	return c
}

func (t *Term) LessThanOrEqualTo(other *Term) *Constraint {
	e := NewExpressionFrom(t)
	c := e.LessThanOrEqualToTerm(other)
	return c
}

func (t *Term) LessThanOrEqualToVariable(v *Variable) *Constraint {
	e := NewExpressionFrom(t)
	c := e.LessThanOrEqualToVariable(v)
	return c
}

func (t *Term) LessThanOrEqualToFloat(constant Float) *Constraint {
	e := NewExpressionFrom(t)
	c := e.LessThanOrEqualToFloat(constant)
	return c
}

func (t *Term) GreaterThanOrEqualToExpression(e *Expression) *Constraint {
	te := NewExpressionFrom(t)
	c := te.GreaterThanOrEqualTo(e)
	return c
}

func (t *Term) GreaterThanOrEqualTo(other *Term) *Constraint {
	e := NewExpressionFrom(t)
	c := e.GreaterThanOrEqualToTerm(other)
	return c
}

func (t *Term) GreaterThanOrEqualToVariable(v *Variable) *Constraint {
	e := NewExpressionFrom(t)
	c := e.GreaterThanOrEqualToVariable(v)
	return c
}

func (t *Term) GreaterThanOrEqualToFloat(constant Float) *Constraint {
	e := NewExpressionFrom(t)
	c := e.GreaterThanOrEqualToFloat(constant)
	return c
}
