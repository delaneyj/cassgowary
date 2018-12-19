package cassgowary

import (
	"github.com/emirpasic/gods/maps/linkedhashmap"

	"github.com/pkg/errors"
)

type tag struct {
	marker, other *symbol
}

type editInfo struct {
	tag        tag
	constraint *Constraint
	constant   Float
}

func newEditInfo(c *Constraint, t tag, constant Float) *editInfo {
	return &editInfo{
		constraint: c,
		tag:        t,
		constant:   constant,
	}
}

type Solver struct {
	cns                   *linkedhashmap.Map //[*Constraint]tag
	rows                  *linkedhashmap.Map //[Symbol]*row
	vars                  *linkedhashmap.Map //[*Variable]Symbol
	edits                 *linkedhashmap.Map //[*Variable]editInfo
	infeasibleRows        symbols
	objective, artificial *row
}

func NewSolver() *Solver {
	return &Solver{
		cns:            linkedhashmap.New(),
		rows:           linkedhashmap.New(),
		vars:           linkedhashmap.New(),
		edits:          linkedhashmap.New(),
		infeasibleRows: symbols{},
		objective:      newRow(),
		artificial:     nil,
	}
}

func (s *Solver) AddVariable(name string) {

}

// AddConstraint adds a constraint to the solver.
func (s *Solver) AddConstraint(c *Constraint) error {
	if _, exists := s.cns.Get(c); exists {
		return DuplicateConstraintErr(c)
	}

	var t tag
	r := s.createRow(c, t)
	subject := s.chooseSubject(r, t)

	switch {
	case subject.kind == symbolInvalid && r.allDummies():
		if r.constant.NearZero() {
			return UnsatisfiableConstraintErr(c)
		} else {
			subject = t.marker
		}

	case subject.kind == symbolInvalid:
		if added, err := s.addWithArtificialVariable(r); !added || err != nil {
			return UnsatisfiableConstraintErr(c)
		}

	default:
		r.solveFor(subject)
		s.substitute(subject, r)
		s.rows.Put(subject, r)
	}

	s.cns.Put(c, t)
	s.optimize(s.objective)

	return nil
}

func (s *Solver) removeConstraint(c *Constraint) error {
	t, exists := s.cns.Get(c)
	if !exists {
		return UnknownConstraintErr(c)
	}

	tag := t.(tag)
	s.cns.Remove(c)
	s.removeConstraintEffects(c, tag)

	if r, exists := s.rows.Get(tag.marker); exists {
		r.(*row).cells.Remove(tag.marker)
	} else {
		r := s.markerLeavingRow(tag.marker)
		if r == nil {
			return InternalSolverErr
		}

		//This looks wrong! changes made below
		//Symbol leaving = tag.marker;
		//rows.remove(tag.marker);

		leaving := newSymbol()
		if fs, _ := s.rows.Find(func(k, v interface{}) bool {
			return r == v
		}); fs != nil {
			leaving = fs.(*symbol)
		}

		if leaving != nil && leaving.kind == symbolInvalid {
			return InternalSolverErr
		}

		s.rows.Remove(leaving)
		r.solveForSymbols(leaving, tag.marker)
		s.substitute(tag.marker, r)
	}
	s.optimize(s.objective)
	return nil
}

func (s *Solver) removeConstraintEffects(c *Constraint, t tag) {
	if t.marker != nil && t.marker.kind == symbolError {
		s.removeMarkerEffects(t.marker, c.Strength.Float())
	} else if t.other != nil && t.other.kind == symbolError {
		s.removeMarkerEffects(t.other, c.Strength.Float())
	}
}

func (s *Solver) removeMarkerEffects(marker *symbol, strength Float) {
	if r, exists := s.rows.Get(marker); exists {
		s.objective.insertRow(r.(*row), -strength)
	} else {
		s.objective.insertSymbol(marker, -strength)
	}
}

func (s *Solver) markerLeavingRow(marker *symbol) *row {
	r1, r2 := FloatMax, FloatMax
	var first, second, third *row

	s.rows.Each(func(k, v interface{}) {
		symbol := k.(*symbol)
		candidate := v.(*row)

		c := candidate.coefficientFor(marker)
		if c == 0 {
			return
		}

		if symbol.kind == symbolExternal {
			third = candidate
		} else if c < 0 {
			r := -candidate.constant / c
			if r < r1 {
				r1 = r
				first = candidate
			}
		} else {
			if r := candidate.constant / c; r < r2 {
				r2 = r
				second = candidate
			}
		}
	})

	if first != nil {
		return first
	}
	if second != nil {
		return second
	}
	return third
}

func (s *Solver) HasConstraint(c Constraint) bool {
	_, exists := s.cns.Get(c)
	return exists
}

func (s *Solver) AddEditVariable(v *Variable, strength Strength) error {
	if _, exists := s.edits.Get(v); exists {
		return DuplicateEditVariableErr
	}

	strength = ClipStrength(strength)

	if strength == Required {
		return RequiredFailureErr
	}

	Terms := Terms{NewTermFrom(v)}
	c := NewConstraint(
		NewExpression(0, Terms...),
		OP_EQ,
		strength,
	)

	if err := s.AddConstraint(c); err != nil {
		return errors.Wrap(err, "can't add edit variable *Constraint")
	}

	t, _ := s.cns.Get(c)
	tag := t.(tag)
	info := newEditInfo(c, tag, 0)
	s.edits.Put(v, info)

	return nil
}

func (s *Solver) RemoveEditVariable(v *Variable) error {
	e, exists := s.edits.Get(v)
	edit := e.(*editInfo)
	if !exists {
		return UnknownEditVariableErr
	}

	if err := s.removeConstraint(edit.constraint); err != nil {
		return UnknownConstraintErr(edit.constraint)
	}

	s.edits.Remove(v)
	return nil
}

func (s *Solver) HasEditVariable(v *Variable) bool {
	_, exists := s.edits.Get(v)
	return exists
}

func (s *Solver) SuggestValue(v Variable, value Float) error {
	e, exists := s.edits.Get(v)
	edit := e.(*editInfo)
	if !exists {
		return UnknownEditVariableErr
	}

	delta := value - edit.constant
	edit.constant = value

	x, exists := s.rows.Get(edit.tag.marker)
	r := x.(*row)
	if exists {
		if r.add(-delta) < 0.0 {
			s.infeasibleRows = append(
				s.infeasibleRows,
				edit.tag.marker,
			)
		}
		return s.dualOptimize()
	}

	x, exists = s.rows.Get(edit.tag.other)
	r = x.(*row)
	if exists {
		if r.add(delta) < 0 {
			s.infeasibleRows = append(
				s.infeasibleRows,
				edit.tag.other,
			)
		}
		return s.dualOptimize()
	}

	s.rows.Each(func(k, v interface{}) {
		symbol := k.(*symbol)
		r := v.(*row)
		coefficient := r.coefficientFor(edit.tag.marker)
		if coefficient != 0.0 &&
			r.add(delta*coefficient) < 0.0 &&
			symbol != nil &&
			symbol.kind == symbolExternal {
			s.infeasibleRows = append(
				s.infeasibleRows,
				symbol,
			)
		}
	})

	return s.dualOptimize()
}

func (s *Solver) UpdateVariables() {
	s.vars.Each(func(k, v interface{}) {
		variable := k.(*Variable)
		symbol := v.(*symbol)

		if r, exists := s.rows.Get(symbol); exists {
			c := r.(*row).constant
			variable.Value = c
		} else {
			variable.Value = 0
		}
	})
}

// Create a new Row object for the given constraint.
//
// The Terms in the constraint will be converted to cells in the row.
// Any Term in the constraint with a coefficient of zero is ignored.
// This method uses the `getVarSymbol` method to get the symbol for
// the variables added to the row. If the symbol for a given cell
// variable is basic, the cell variable will be substituted with the
// basic row.
//
// The necessary slack and error variables will be added to the row.
// If the constant for the row is negative, the sign for the row
// will be inverted so the constant becomes positive.
//
// The tag will be updated with the marker and error symbols to use
// for tracking the movement of the constraint in the tableau.
func (s *Solver) createRow(c *Constraint, tag tag) *row {
	r := newRowWith(c.expression.Constant)
	for _, t := range c.expression.Terms {
		if !t.Coefficient.NearZero() {
			symbol := s.varSymbol(t.Variable)
			if otherRow, exists := s.rows.Get(symbol); exists {
				r.insertRow(otherRow.(*row), t.Coefficient)
			} else {
				r.insertSymbol(symbol, t.Coefficient)
			}
		}
	}

	switch c.Op {
	case OP_LE, OP_GE:
		coeff := Float(-1)
		if c.Op == OP_LE {
			coeff = 1
		}
		slack := newSymbolFrom(symbolSlack)
		tag.marker = slack
		r.insertSymbol(tag.marker, coeff)
		if c.Strength < Required {
			serror := newSymbolFrom(symbolError)
			tag.other = serror
			r.insertSymbol(serror, -coeff)
			s.objective.insertSymbol(serror, c.Strength.Float())
		}

	case OP_EQ:
		if c.Strength < Required {
			errPlus := newSymbolFrom(symbolError)
			errMinus := newSymbolFrom(symbolError)
			tag.marker = errPlus
			tag.other = errMinus
			r.insertSymbol(errPlus, -1) // v = eplus - eminus
			r.insertSymbol(errMinus, 1) // v - eplus + eminus = 0
			s.objective.insertSymbol(errPlus, c.Strength.Float())
			s.objective.insertSymbol(errMinus, c.Strength.Float())
		} else {
			dummy := newSymbolFrom(symbolDummy)
			tag.marker = dummy
			r.insertSymbol(dummy, 1)
		}
	}

	// Ensure the row as a positive constant.
	if r.constant < 0.0 {
		r.reverseSign()
	}

	return r
}

// Choose the subject for solving for the row
// This method will choose the best subject for using as the solve
// target for the row. An invalid symbol will be returned if there
// is no valid target.
// The symbols are chosen according to the following precedence:
// 1) The first symbol representing an external variable.
// 2) A negative slack or error tag variable.
// If a subject cannot be found, an invalid symbol will be returned.
func (s *Solver) chooseSubject(r *row, t tag) *symbol {
	if fk, _ := r.cells.Find(func(k, v interface{}) bool {
		return k.(*symbol).kind == symbolExternal
	}); fk != nil {
		return fk.(*symbol)
	}

	if t.marker != nil && (t.marker.kind == symbolSlack || t.marker.kind == symbolError) {
		if r.coefficientFor(t.marker) < 0.0 {
			return t.marker
		}
	}

	if t.other != nil && (t.other.kind == symbolSlack || t.other.kind == symbolError) {
		if r.coefficientFor(t.other) < 0.0 {
			return t.other
		}
	}

	return newSymbol()
}

// Add the row to the tableau using an artificial variable.
// This will return false if the constraint cannot be satisfied.
func (s *Solver) addWithArtificialVariable(r *row) (bool, error) {
	// Create and add the artificial variable to the tableau
	art := newSymbolFrom(symbolSlack)
	s.rows.Put(art, newRowFrom(r))
	s.artificial = newRowFrom(r)

	// Optimize the artificial objective. This is successful
	// only if the artificial objective is optimized to zero.
	if err := s.optimize(s.artificial); err != nil {
		return false, errors.Wrap(err, "can't optimize")
	}

	success := s.artificial.constant.NearZero()
	s.artificial = nil

	// If the artificial variable is basic, pivot the row so that
	// it becomes basic. If the row is constant, exit early.
	if x, exists := s.rows.Get(art); exists {
		/**this looks wrong!!!*/
		//rows.remove(rowptr);
		rowptr := x.(*row)
		s.rows = s.rows.Select(func(k, v interface{}) bool {
			return v.(*row) != rowptr
		})

		if rowptr.cells.Size() == 0 {
			return success, nil
		}

		entering := s.anyPivotableSymbol(rowptr)
		if entering.kind == symbolInvalid {
			return false, nil // unsatisfiable (will this ever happen?)
		}
		rowptr.solveForSymbols(art, entering)
		s.substitute(entering, rowptr)
		s.rows.Put(entering, rowptr)
	}

	// Remove the artificial variable from the tableau.
	s.rows.Each(func(k, v interface{}) {
		v.(*row).cells.Remove(art)
	})

	s.objective.cells.Remove(art)
	return success, nil
}

// Substitute the parametric symbol with the given row.
// This method will substitute all instances of the parametric symbol
// in the tableau and the objective function with the given row.
func (s *Solver) substitute(sym *symbol, r *row) {
	s.rows.Each(func(k, v interface{}) {
		ss := k.(*symbol)
		row := v.(*row)
		row.substitute(sym, r)

		if ss.kind != symbolExternal && row.constant < 0 {
			s.infeasibleRows = append(s.infeasibleRows, ss)
		}
	})

	s.objective.substitute(sym, r)

	if s.artificial != nil {
		s.artificial.substitute(sym, r)
	}
}

// Optimize the system for the given objective function.
// This method performs iterations of Phase 2 of the simplex method
// until the objective function reaches a minimum.
func (s *Solver) optimize(objective *row) error {
	for {
		entering := s.enteringSymbol(objective)
		if entering.kind == symbolInvalid {
			return nil
		}

		entry := s.leavingRow(entering)
		if entry == nil {
			return errors.New("The objective is unbounded.")
		}

		var leaving, entryKey *symbol
		{
			if fs, _ := s.rows.Find(func(k, v interface{}) bool {
				return v.(*row) == entry
			}); fs != nil {
				leaving = fs.(*symbol)
			}

			if fs, _ := s.rows.Find(func(k, v interface{}) bool {
				return v.(*row) == entry
			}); fs != nil {
				entryKey = fs.(*symbol)
			}
		}

		s.rows.Remove(entryKey)
		entry.solveForSymbols(leaving, entering)
		s.substitute(entering, entry)
		s.rows.Put(entering, entry)
	}
}

func (s *Solver) dualOptimize() error {
	for len(s.infeasibleRows) > 0 {
		lastIndex := len(s.infeasibleRows) - 1
		leaving := s.infeasibleRows[lastIndex]
		s.infeasibleRows = s.infeasibleRows[:lastIndex]

		if x, exists := s.rows.Get(leaving); exists {
			r := x.(*row)
			if r.constant < 0 {
				entering := s.dualEnteringSymbol(r)
				if entering.kind == symbolInvalid {
					return InternalSolverErr
				}
				s.rows.Remove(leaving)
				r.solveForSymbols(leaving, entering)
				s.substitute(entering, r)
				s.rows.Put(entering, r)
			}
		}
	}
	return nil
}

// Compute the entering variable for a pivot operation.
// This method will return first symbol in the objective function which
// is non-dummy and has a coefficient less than zero. If no symbol meets
// the criteria, it means the objective function is at a minimum, and an
// invalid symbol is returned.
func (s *Solver) enteringSymbol(objective *row) *symbol {
	foundSymbolRaw, _ := objective.cells.Find(func(k, v interface{}) bool {
		symbol := k.(*symbol)
		value := v.(Float)

		if symbol.kind != symbolDummy && value < 0 {
			return true
		}
		return false
	})
	if foundSymbolRaw != nil {
		return foundSymbolRaw.(*symbol)
	}

	return newSymbol()
}

func (s *Solver) dualEnteringSymbol(r *row) *symbol {
	entering, ratio := newSymbol(), FloatMax
	r.cells.Each(func(k, v interface{}) {
		if sym := k.(*symbol); sym.kind != symbolDummy {
			x, _ := r.cells.Get(sym)
			currentCell := x.(Float)
			if currentCell > 0.0 {
				coefficient := s.objective.coefficientFor(sym)
				r := coefficient / currentCell
				if r < ratio {
					ratio = r
					entering = sym
				}
			}
		}
	})

	return entering
}

// Get the first Slack or Error symbol in the row.
// If no such symbol is present, and Invalid symbol will be returned.
func (s *Solver) anyPivotableSymbol(r *row) *symbol {
	if fs, _ := r.cells.Find(func(k, v interface{}) bool {
		sym := k.(*symbol)
		return sym.kind == symbolSlack || sym.kind == symbolError
	}); fs != nil {
		return fs.(*symbol)
	}
	return newSymbol()
}

// Compute the row which holds the exit symbol for a pivot.
// This documentation is copied from the C++ version and is outdated
// This method will return an iterator to the row in the row map
// which holds the exit symbol. If no appropriate exit symbol is
// found, the end() iterator will be returned. This indicates that
// the objective function is unbounded.
func (s *Solver) leavingRow(entering *symbol) *row {
	ratio := FloatMax
	var r *row

	s.rows.Each(func(k, v interface{}) {
		sym := k.(*symbol)

		if sym.kind != symbolExternal {
			candidate := v.(*row)

			t := candidate.coefficientFor(entering)
			if t < 0 {
				if tr := -candidate.constant / t; tr < ratio {
					ratio = tr
					r = candidate
				}
			}
		}
	})
	return r
}

// Get the symbol for the given variable.
// If a symbol does not exist for the variable, one will be created.
func (s *Solver) varSymbol(v *Variable) *symbol {
	if x, exists := s.vars.Get(v); exists {
		return x.(*symbol)
	}

	symbol := newSymbolFrom(symbolExternal)
	s.vars.Put(v, symbol)
	return symbol
}
