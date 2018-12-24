package cassgowary

import (
	"fmt"
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestBenchmarkTestAddingLotsOfConstraints(t *testing.T) {
	solver := NewSolver()
	vr := &benchmarkVariableResolver{
		solver:    solver,
		variables: nodeMap{},
	}

	cp := NewConstraintParser()
	c, err := cp.ParseConstraint("variable0 == 100", vr)
	assert.NoError(t, err)
	solver.AddConstraint(c)

	getVariableName := func(number int) string {
		return fmt.Sprintf("getVariable:%d", number)
	}

	runs := 3000
	for i := 1; i < runs; i++ {
		constraintString := fmt.Sprintf(
			"%s == 100 + %s",
			getVariableName(i),
			getVariableName(i-1),
		)
		constraint, err := cp.ParseConstraint(constraintString, vr)
		assert.NoError(t, err)
		start := time.Now()
		solver.AddConstraint(constraint)
		log.Printf("%d, %s", i, time.Since(start))
	}
}

type benchmarkVariableResolver struct {
	solver    *Solver
	variables nodeMap
}

func (vr *benchmarkVariableResolver) ResolveVariable(variableName string) (*Variable, error) {
	if v, exists := vr.variables[variableName]; exists {
		return v, nil
	}

	v := NewVariable(variableName)
	vr.variables[variableName] = v
	return v, nil
}

func (vr *benchmarkVariableResolver) ResolveConstant(name string) (*Expression, error) {
	f, err := strconv.ParseFloat(name, 64)
	if err != nil {
		return nil, err
	}

	e := NewExpression(Float(f))
	return e, nil
}