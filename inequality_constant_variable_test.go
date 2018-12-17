package cassgowary

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

const ConstantVariableTestEpsilon = 1.0e-8

func TestFloatLessThanEqualTo(t *testing.T) {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	solver := NewSolver()
	x := NewVariable("x")
	f := Float(100)
	c := f.LessThanOrEqualToVariable(x)
	log.Printf("%+v", c)
	solver.AddConstraint(c)
	solver.UpdateVariables()
	assert.True(t, 100 <= x.Value)
	c2 := x.EqualsFloat(110)
	log.Printf("%+v", c2)
	solver.AddConstraint(c2)
	log.Printf("%+v", solver)
	log.Printf("%+v", solver.cns.Size())
	solver.UpdateVariables()
	assert.InDelta(t, 110, x.Value.Raw(), ConstantVariableTestEpsilon)
}

//     @Test(expected = UnsatisfiableConstraintException.class)
//     public void lessThanEqualToUnsatisfiable() throws DuplicateConstraintException, UnsatisfiableConstraintException {
//         Variable x = new Variable("x");
//         Solver solver = new Solver();
//         solver.addConstraint(Symbolics.lessThanOrEqualTo(100, x));
//         solver.updateVariables();
//         assertTrue(x.getValue() <= 100);
//         solver.addConstraint(Symbolics.equals(x, 10));
//         solver.updateVariables();
//     }

//     @Test
//     public void greaterThanEqualTo() throws DuplicateConstraintException, UnsatisfiableConstraintException {
//         Variable x = new Variable("x");
//         Solver solver = new Solver();
//         solver.addConstraint(Symbolics.greaterThanOrEqualTo(100, x));
//         solver.updateVariables();
//         assertTrue(100 >= x.getValue());
//         solver.addConstraint(Symbolics.equals(x, 90));
//         solver.updateVariables();
//         assertEquals(x.getValue(), 90, EPSILON);
//     }

//     @Test(expected = UnsatisfiableConstraintException.class)
//     public void greaterThanEqualToUnsatisfiable() throws DuplicateConstraintException, UnsatisfiableConstraintException {
//         Variable x = new Variable("x");
//         Solver solver = new Solver();
//         solver.addConstraint(Symbolics.greaterThanOrEqualTo(100, x));
//         solver.updateVariables();
//         assertTrue(100 >= x.getValue());
//         solver.addConstraint(Symbolics.equals(x, 110));
//         solver.updateVariables();
//     }
// }
