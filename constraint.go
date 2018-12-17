package cassgowary

import (
	"fmt"
)

type Constraint struct {
	expression *Expression
	Strength   Strength
	Op         RelationalOperator
}

func NewConstraintRequired(expr *Expression, op RelationalOperator) *Constraint {
	return NewConstraint(expr, op, Required)
}

func NewConstraint(expr *Expression, op RelationalOperator, strength Strength) *Constraint {
	return &Constraint{
		expression: expr.Reduce(),
		Op:         op,
		Strength:   ClipStrength(strength),
	}
}

func NewConstraintFrom(other *Constraint, s Strength) *Constraint {
	return NewConstraint(other.expression, other.Op, s)
}

func (c *Constraint) String() string {
	return fmt.Sprintf(
		"expression: (%v) strength:%f operator:%v",
		c.expression, c.Strength, c.Op,
	)
}

func (c *Constraint) NewModifyStrength(strength Strength) *Constraint {
	return NewConstraintFrom(c, strength)
}
