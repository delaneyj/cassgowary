package cassgowary

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/emirpasic/gods/stacks/arraystack"
	"github.com/pkg/errors"
)

const constraintParserOperations = "-+/*^"

type ConstraintParser struct {
	pattern *regexp.Regexp
}

func NewConstraintParser() *ConstraintParser {
	p, err := regexp.Compile(`\s*(\S*)\s*(==|<=|>=)\s*([^!]*)(?:!(required|strong|medium|weak))?`)
	if err != nil {
		log.Fatal(err)
	}
	cp := &ConstraintParser{
		pattern: p,
	}
	return cp
}

type VariableResolver interface {
	ResolveVariable(name string) (*Variable, error)
	ResolveConstant(name string) (*Expression, error)
}

func (cp *ConstraintParser) ParseConstraint(rawConstraint string, variableResolver VariableResolver) (*Constraint, error) {
	matches := cp.pattern.FindStringSubmatch(rawConstraint)
	if len(matches) == 5 {
		variable, err := variableResolver.ResolveVariable(matches[1])
		if err != nil {
			return nil, errors.Wrap(err, "can' parse constraint")
		}
		operator, err := cp.parseOperator(matches[2])
		if err != nil {
			return nil, errors.Wrap(err, "can' parse constraint")
		}
		expression, err := cp.resolveExpression(matches[3], variableResolver)
		if err != nil {
			return nil, errors.Wrap(err, "can' parse constraint")
		}

		strength := cp.parseStrength(matches[4])

		e2 := variable.SubtractExpression(expression)
		return NewConstraint(e2, operator, strength), nil
	}

	return nil, fmt.Errorf("could not parse '%s'.", rawConstraint)
}

func (cp *ConstraintParser) parseOperator(rawOperator string) (RelationalOperator, error) {
	switch rawOperator {
	case "EQ", "==":
		return OP_EQ, nil
	case "GEQ", ">=":
		return OP_GE, nil
	case "LEQ", "<=":
		return OP_LE, nil
	default:
		return OP_EQ, fmt.Errorf("can't parse op string '%s'", rawOperator)
	}
}

func (cp *ConstraintParser) parseStrength(rawStrength string) Strength {
	switch rawStrength {
	case "!strong":
		return Strong
	case "!medium":
		return Medium
	case "!weak":
		return Weak
	default:
		return Required
	}
}

func (cp *ConstraintParser) resolveExpression(rawExpression string, variableResolver VariableResolver) (*Expression, error) {
	tokens := cp.tokenizeExpression(rawExpression)
	postFixExpression := cp.infixToPostfix(tokens)

	expressionStack := arraystack.New() // *Expression

	for _, e := range postFixExpression {
		switch e {
		case "+", "-", "*", "/":
			a, aFound := expressionStack.Pop()
			b, bFound := expressionStack.Pop()
			if !aFound || !bFound {
				return nil, errors.New("can't get expression from stack")
			}
			eA, eB := a.(*Expression), b.(*Expression)

			switch e {
			case "+":
				expressionStack.Push(eA.Add(eB))
			case "-":
				expressionStack.Push(eB.Subtract(eA))
			case "/":
				e, err := eB.Divide(eA)
				if err != nil {
					return nil, errors.Wrap(err, "can't divide expression")
				}
				expressionStack.Push(e)
			case "*":
				e, err := eA.Multiply(eB)
				if err != nil {
					return nil, errors.Wrap(err, "can't multiply expression")
				}
				expressionStack.Push(e)
			}
		default:
			linearExpression, err := variableResolver.ResolveConstant(e)
			if err != nil {
				v, err := variableResolver.ResolveVariable(e)
				if err != nil {
					return nil, errors.Wrap(err, "can't resolve variable")
				}
				t := NewTermFrom(v)
				e := NewExpressionFrom(t)
				linearExpression = e
			}
			expressionStack.Push(linearExpression)
		}
	}

	re, exists := expressionStack.Pop()
	if !exists {
		return nil, errors.New("can't find return expression")
	}
	returnExpression := re.(*Expression)
	return returnExpression, nil
}

func (cp *ConstraintParser) infixToPostfix(tokens []string) []string {
	const OPS = "-+/*^"
	s := arraystack.New() //int

	postFix := make([]string, 0, len(tokens))

	for _, token := range tokens {
		c := token[0]
		idx := strings.Index(OPS, token)

		switch {
		case idx != -1 && len(token) == 1:
			if s.Empty() {
				s.Push(idx)
			} else {
				for !s.Empty() {
					pi, _ := s.Peek()
					prec2 := pi.(int) / 2
					prec1 := idx / 2
					if prec2 > prec1 || (prec2 == prec1 && token != "^") {
						x, _ := s.Pop()
						i := x.(int)
						op := string(OPS[i])
						postFix = append(
							postFix,
							op,
						)
					} else {
						break
					}
				}
				s.Push(idx)
			}
		case c == '(':
			s.Push(-2)
		case c == ')':
			for {
				x, _ := s.Peek()
				i := x.(int)
				if i == -2 {
					break
				}
				s.Pop()
				op := string(OPS[i])
				postFix = append(postFix, op)
			}
			s.Pop()
		default:
			postFix = append(postFix, token)
		}
	}
	for !s.Empty() {
		x, _ := s.Pop()
		i := x.(int)
		op := string(OPS[i])
		postFix = append(
			postFix,
			op,
		)
	}
	return postFix
}

func (cp *ConstraintParser) tokenizeExpression(rawExpression string) []string {
	var sb strings.Builder
	tokens := []string{}

	for _, c := range []rune(rawExpression) {
		switch c {
		case '+', '-', '*', '/', '(', ')':
			if sb.Len() > 0 {
				tokens = append(tokens, sb.String())
				sb.Reset()
			}
			tokens = append(tokens, string(c))
		case ' ':
			// ignore space
		default:
			sb.WriteRune(c)
		}
	}
	if sb.Len() > 0 {
		tokens = append(tokens, sb.String())
	}
	return tokens
}
