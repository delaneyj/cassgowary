package tests

const TestBasicsEpsilon = 1.0e-8

// func TestSimpleNew(t *testing.T) {
// 	solver := cassgowary.NewSolver()
// 	x := cassgowary.NewVariable("x")
// 	e := x.AddFloat(2)
// 	c := e.EqualsFloat(20.0)
// 	solver.AddConstraint(c)
// 	solver.UpdateVariables()
// 	assert.InDelta(t, 18.0, x.Value, TestBasicsEpsilon)
// }

// func TestSimple0(t *testing.T) {
// 	solver := cassgowary.NewSolver()
// 	x := cassgowary.NewVariable("x")
// 	y := cassgowary.NewVariable("y")
// 	solver.AddConstraint(x.EqualsFloat(20))
// 	solver.AddConstraint(x.AddFloat(2).Equals(y.AddFloat(10)))
// 	solver.UpdateVariables()
// 	assert.Equal(t, 12, y.Value, TestBasicsEpsilon)
// 	assert.Equal(t, 20, x.Value, TestBasicsEpsilon)
// }

//     @Test
//     public void simple1() throws DuplicateConstraintException, UnsatisfiableConstraintException {
//         Variable x = new Variable("x");
//         Variable y = new Variable("y");
//         Solver solver = new Solver();
//         solver.addConstraint(Symbolics.equals(x, y));
//         solver.updateVariables();
//         assertEquals(x.getValue(), y.getValue(), EPSILON);
//     }

//     @Test
//     public void casso1() throws DuplicateConstraintException, UnsatisfiableConstraintException {
//         Variable x = new Variable("x");
//         Variable y = new Variable("y");
//         Solver solver = new Solver();

//         solver.addConstraint(Symbolics.lessThanOrEqualTo(x, y));
//         solver.addConstraint(Symbolics.equals(y, Symbolics.add(x, 3.0)));
//         solver.addConstraint(Symbolics.equals(x, 10.0).setStrength(Strength.WEAK));
//         solver.addConstraint(Symbolics.equals(y, 10.0).setStrength(Strength.WEAK));

//         solver.updateVariables();

//         if (Math.abs(x.getValue() - 10.0) < EPSILON) {
//             assertEquals(10, x.getValue(), EPSILON);
//             assertEquals(13, y.getValue(), EPSILON);
//         } else {
//             assertEquals(7, x.getValue(), EPSILON);
//             assertEquals(10, y.getValue(), EPSILON);
//         }
//     }

//     @Test
//     public void addDelete1() throws DuplicateConstraintException, UnsatisfiableConstraintException, UnknownConstraintException {
//         Variable x = new Variable("x");
//         Solver solver = new Solver();

//         solver.addConstraint(Symbolics.lessThanOrEqualTo(x, 100).setStrength(Strength.WEAK));

//         solver.updateVariables();
//         assertEquals(100, x.getValue(), EPSILON);

//         Constraint c10 = Symbolics.lessThanOrEqualTo(x, 10.0);
//         Constraint c20 = Symbolics.lessThanOrEqualTo(x, 20.0);

//         solver.addConstraint(c10);
//         solver.addConstraint(c20);

//         solver.updateVariables();

//         assertEquals(10, x.getValue(), EPSILON);

//         solver.removeConstraint(c10);

//         solver.updateVariables();

//         assertEquals(20, x.getValue(), EPSILON);

//         solver.removeConstraint(c20);
//         solver.updateVariables();

//         assertEquals(100, x.getValue(), EPSILON);

//         Constraint c10again = Symbolics.lessThanOrEqualTo(x, 10.0);

//         solver.addConstraint(c10again);
//         solver.addConstraint(c10);
//         solver.updateVariables();

//         assertEquals(10, x.getValue(), EPSILON);

//         solver.removeConstraint(c10);
//         solver.updateVariables();
//         assertEquals(10, x.getValue(), EPSILON);

//         solver.removeConstraint(c10again);
//         solver.updateVariables();
//         assertEquals(100, x.getValue(), EPSILON);
//     }

//     @Test
//     public void addDelete2() throws DuplicateConstraintException, UnsatisfiableConstraintException, UnknownConstraintException {
//         Variable x = new Variable("x");
//         Variable y = new Variable("y");
//         Solver solver = new Solver();

//         solver.addConstraint(Symbolics.equals(x, 100).setStrength(Strength.WEAK));
//         solver.addConstraint(Symbolics.equals(y, 120).setStrength(Strength.STRONG));

//         Constraint c10 = Symbolics.lessThanOrEqualTo(x, 10.0);
//         Constraint c20 = Symbolics.lessThanOrEqualTo(x, 20.0);

//         solver.addConstraint(c10);
//         solver.addConstraint(c20);
//         solver.updateVariables();

//         assertEquals(10, x.getValue(), EPSILON);
//         assertEquals(120, y.getValue(), EPSILON);

//         solver.removeConstraint(c10);
//         solver.updateVariables();

//         assertEquals(20, x.getValue(), EPSILON);
//         assertEquals(120, y.getValue(), EPSILON);

//         Constraint cxy = Symbolics.equals(Symbolics.multiply(x, 2.0), y);
//         solver.addConstraint(cxy);
//         solver.updateVariables();

//         assertEquals(20, x.getValue(), EPSILON);
//         assertEquals(40, y.getValue(), EPSILON);

//         solver.removeConstraint(c20);
//         solver.updateVariables();

//         assertEquals(60, x.getValue(), EPSILON);
//         assertEquals(120, y.getValue(), EPSILON);

//         solver.removeConstraint(cxy);
//         solver.updateVariables();

//         assertEquals(100, x.getValue(), EPSILON);
//         assertEquals(120, y.getValue(), EPSILON);
//     }

//     @Test(expected = UnsatisfiableConstraintException.class)
//     public void inconsistent1() throws InternalError, DuplicateConstraintException, UnsatisfiableConstraintException {
//         Variable x = new Variable("x");
//         Solver solver = new Solver();

//         solver.addConstraint(Symbolics.equals(x, 10.0));
//         solver.addConstraint(Symbolics.equals(x, 5.0));

//         solver.updateVariables();
//     }

//     @Test(expected = UnsatisfiableConstraintException.class)
//     public void inconsistent2() throws DuplicateConstraintException, UnsatisfiableConstraintException {
//         Variable x = new Variable("x");
//         Solver solver = new Solver();

//         solver.addConstraint(Symbolics.greaterThanOrEqualTo(x, 10.0));
//         solver.addConstraint(Symbolics.lessThanOrEqualTo(x, 5.0));
//         solver.updateVariables();
//     }

//     @Test(expected = UnsatisfiableConstraintException.class)
//     public void inconsistent3() throws DuplicateConstraintException, UnsatisfiableConstraintException {

//         Variable w = new Variable("w");
//         Variable x = new Variable("x");
//         Variable y = new Variable("y");
//         Variable z = new Variable("z");
//         Solver solver = new Solver();

//         solver.addConstraint(Symbolics.greaterThanOrEqualTo(w, 10.0));
//         solver.addConstraint(Symbolics.greaterThanOrEqualTo(x, w));
//         solver.addConstraint(Symbolics.greaterThanOrEqualTo(y, x));
//         solver.addConstraint(Symbolics.greaterThanOrEqualTo(z, y));
//         solver.addConstraint(Symbolics.greaterThanOrEqualTo(z, 8.0));
//         solver.addConstraint(Symbolics.lessThanOrEqualTo(z, 4.0));
//         solver.updateVariables();
//     }

// }
