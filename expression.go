package cassgowary

import (
	"fmt"
	"log"
	"strings"

	"github.com/emirpasic/gods/maps/linkedhashmap"
)

type Expression struct {
	Terms    Terms
	Constant Float
}

func NewExpression(constant Float, terms ...*Term) *Expression {
	log.Print(constant, terms)
	return &Expression{
		Constant: constant,
		Terms:    terms,
	}
}

func NewExpressionFrom(terms ...*Term) *Expression {
	log.Printf("%+v", terms)
	return NewExpression(0, terms...)
}

func (e *Expression) Value() float64 {
	result := e.Constant
	for _, t := range e.Terms {
		result += t.Value()
	}
	return float64(result)
}

func (e *Expression) IsConstant() bool {
	return len(e.Terms) == 0
}

func (e *Expression) Reduce() *Expression {
	vars := linkedhashmap.New() //*Variable,Float

	for _, t := range e.Terms {
		value := t.Coefficient
		if tv, exists := vars.Get(t.Variable); exists {
			value += tv.(Float)
		}
		vars.Put(t.Variable, value)
	}

	reducedTerms := make(Terms, 0, vars.Size())
	vars.Each(func(k, v interface{}) {
		variable := k.(*Variable)
		value := v.(Float)
		t := NewTerm(variable, value)
		reducedTerms = append(reducedTerms, t)
	})

	return &Expression{
		Terms:    reducedTerms,
		Constant: e.Constant,
	}
}

func (e *Expression) String() string {
	var (
		sb         strings.Builder
		IsConstant = e.IsConstant()
	)
	sb.WriteString(fmt.Sprintf(
		"IsConstant:%t constant:%f",
		IsConstant, e.Constant,
	))

	if !e.IsConstant() {
		sb.WriteString(" Terms: [")
		for _, t := range e.Terms {
			sb.WriteString("(")
			sb.WriteString(t.String())
			sb.WriteString(")")
		}
		sb.WriteString("] ")
	}
	return sb.String()
}

// Expression multiply, divide, and unary invert
func (e *Expression) Multiply(coefficient Float) *Expression {
	terms := make(Terms, len(e.Terms))
	for i, t := range e.Terms {
		terms[i] = t.Multiply(coefficient)
	}
	return NewExpression(e.Constant*coefficient, terms...)
}

func (e *Expression) MultiplyExpression(other Expression) (*Expression, error) {
	switch {
	case e.IsConstant():
		return other.Multiply(e.Constant), nil
	case other.IsConstant():
		return other.Multiply(other.Constant), nil
	default:
		return nil, NonLinearExpressionErr
	}
}

func (e *Expression) DivideFloat(denominator Float) *Expression {
	return e.Multiply(1 / denominator)
}

func (e *Expression) Divide(other Expression) (*Expression, error) {
	if other.IsConstant() {
		return e.DivideFloat(other.Constant), nil
	}
	return nil, NonLinearExpressionErr
}

func (e *Expression) Negate() *Expression {
	return e.Multiply(-1)
}

func (e *Expression) Add(other *Expression) *Expression {
	constant, tc := e.Constant+other.Constant, len(e.Terms)
	terms := make(Terms, tc+len(other.Terms))
	copy(terms, e.Terms)
	copy(terms[tc:], other.Terms)
	return NewExpression(constant, terms...)
}

func (e *Expression) AddTerm(t *Term) *Expression {
	tc := len(e.Terms)
	terms := make(Terms, tc+1)
	copy(terms, e.Terms)
	terms[tc] = t
	return NewExpression(e.Constant, terms...)
}

func (e *Expression) AddVariable(v *Variable) *Expression {
	return e.AddTerm(NewTermFrom(v))
}

func (e *Expression) AddFloat(constant Float) *Expression {
	return NewExpression(e.Constant+constant, e.Terms...)
}

func (e *Expression) Subtract(other *Expression) *Expression {
	negated := other.Negate()
	return e.Add(negated)
}

func (e *Expression) SubtractTerm(t *Term) *Expression {
	negated := t.Negate()
	return e.AddTerm(negated)
}

func (e *Expression) SubtractVariable(v *Variable) *Expression {
	negated := v.Negate()
	return e.AddTerm(negated)
}

func (e *Expression) SubtractFloat(constant Float) *Expression {
	return e.AddFloat(-constant)
}

// Expression relations
func (e *Expression) Equals(other *Expression) *Constraint {
	e2 := e.Subtract(other)
	log.Printf("e:%+v", e)
	log.Printf("other:%+v", other)
	log.Printf("e2:%+v", e2)
	return NewConstraintRequired(e2, OP_EQ)
}

func (e *Expression) EqualsTerm(t *Term) *Constraint {
	return e.Equals(NewExpressionFrom(t))
}

func (e *Expression) EqualsVariable(v *Variable) *Constraint {
	return e.EqualsTerm(NewTermFrom(v))
}

func (e *Expression) EqualsFloat(constant Float) *Constraint {
	e2 := NewExpression(constant)
	c := e.Equals(e2)
	return c
}

func (e *Expression) LessThanOrEqualTo(other *Expression) *Constraint {
	e2 := e.Subtract(other)
	log.Print(e2)
	return NewConstraintRequired(e2, OP_LE)
}

func (e *Expression) LessThanOrEqualToTerm(t *Term) *Constraint {
	return e.LessThanOrEqualTo(NewExpressionFrom(t))
}

func (e *Expression) LessThanOrEqualToVariable(v *Variable) *Constraint {
	return e.LessThanOrEqualToTerm(NewTermFrom(v))
}

func (e *Expression) LessThanOrEqualToFloat(constant Float) *Constraint {
	return e.LessThanOrEqualTo(NewExpression(constant))
}

func (e *Expression) GreaterThanOrEqualTo(other *Expression) *Constraint {
	return NewConstraintRequired(e.Subtract(other), OP_GE)
}

func (e *Expression) GreaterThanOrEqualToTerm(t *Term) *Constraint {
	return e.GreaterThanOrEqualTo(NewExpressionFrom(t))
}

func (e *Expression) GreaterThanOrEqualToVariable(v *Variable) *Constraint {
	return e.GreaterThanOrEqualToTerm(NewTermFrom(v))
}

func (e *Expression) GreaterThanOrEqualToFloat(constant Float) *Constraint {
	return e.GreaterThanOrEqualTo(NewExpression(constant))
}
