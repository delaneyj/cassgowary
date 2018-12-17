package cassgowary

import (
	"log"

	"github.com/emirpasic/gods/maps/linkedhashmap"
)

type row struct {
	constant Float
	cells    *linkedhashmap.Map
}

func newRow() *row {
	return &row{
		cells: linkedhashmap.New(),
	}
}

func newRowWith(constant Float) *row {
	r := newRow()
	r.constant = constant
	return r
}

func newRowFrom(other *row) *row {
	r := &row{
		cells:    linkedhashmap.New(),
		constant: other.constant,
	}
	other.cells.Each(func(k, v interface{}) {
		r.cells.Put(k, v)
	})

	log.Printf("%+v", r)
	r.cells.Each(func(k, v interface{}) {
		log.Printf("cell k:%+v v:%+v", k, v)
	})
	return r
}

//Add a constant value to the row constant.
func (r *row) add(value Float) Float {
	r.constant += value
	return r.constant
}

// Insert a symbol into the row with a given coefficient.
// If the symbol already exists in the row, the coefficient will be
// added to the existing coefficient. If the resulting coefficient
// is zero, the symbol will be removed from the row
func (r *row) insertSymbol(s Symbol, coeffecient Float) {
	if x, exists := r.cells.Get(s); exists {
		existingCoefficient := x.(Float)
		coeffecient += existingCoefficient
	}

	if coeffecient.NearZero() {
		r.cells.Remove(s)
		return
	}

	r.cells.Put(s, coeffecient)
}

// Insert a symbol into the row with a given coefficient.
// If the symbol already exists in the row, the coefficient will be
// added to the existing coefficient. If the resulting coefficient
// is zero, the symbol will be removed from the row
func (r *row) insertSymbolDefault(s Symbol) {
	r.insertSymbol(s, 1)
}

//Insert a row into this row with a given coefficient.
//The constant and the cells of the other row will be multiplied by
//the coefficient and added to this row. Any cell with a resulting
//coefficient of zero will be removed from the row.
func (r *row) insertRow(other *row, coefficient Float) {
	r.constant += other.constant * coefficient

	other.cells.Each(func(k, v interface{}) {
		s := k.(Symbol)
		coeff := v.(Float) * coefficient

		var temp Float
		if x, exists := r.cells.Get(s); exists {
			temp = x.(Float)
		}
		temp += coeff
		r.cells.Put(s, temp)
		if temp.NearZero() {
			r.cells.Remove(s)
		}
	})
}

func (r *row) insertFromDefault(other *row) {
	r.insertRow(other, 1)
}

func (r *row) remove(s Symbol) {
	r.cells.Remove(s)
}

func (r *row) reverseSign() {
	r.constant *= -1
	r.cells.Each(func(k, v interface{}) {
		r.cells.Put(k, -v.(Float))
	})
}

// Solve the row for the given symbol.
// This method assumes the row is of the form a * x + b * y + c = 0
// and (assuming solve for x) will modify the row to represent the
// right hand side of x = -b/a * y - c / a. The target symbol will
// be removed from the row, and the constant and other cells will
// be multiplied by the negative inverse of the target coefficient.
// The given symbol *must* exist in the row.
func (r *row) solveFor(s Symbol) {
	v, _ := r.cells.Get(s)
	value := v.(Float)
	coeff := -1 / value
	r.cells.Remove(s)
	r.constant *= coeff

	r.cells.Each(func(k, v interface{}) {
		r.cells.Put(k, v.(Float)*coeff)
	})
}

//  Solve the row for the given symbols.
//  This method assumes the row is of the form x = b * y + c and will
//  solve the row such that y = x / b - c / b. The rhs symbol will be
//  removed from the row, the lhs added, and the result divided by the
//  negative inverse of the rhs coefficient.
//  The lhs symbol *must not* exist in the row, and the rhs symbol
//  must* exist in the row.
func (r *row) solveForSymbols(lhs, rhs Symbol) {
	r.insertSymbol(lhs, -1.0)
	r.solveFor(rhs)
}

// Get the coefficient for the given symbol.
// <p/>
// If the symbol does not exist in the row, zero will be returned.
func (r *row) coefficientFor(s Symbol) Float {
	log.Printf("%+v %d", r, r.cells.Size())
	r.cells.Each(func(k, v interface{}) {
		log.Printf("cell k:%+v v:%+v", k, v)
	})

	if v, exists := r.cells.Get(s); exists {
		f := v.(Float)
		return f
	}
	return 0
}

// Substitute a symbol with the data from another row.
// <p/>
// Given a row of the form a * x + b and a substitution of the
// form x = 3 * y + c the row will be updated to reflect the
// expression 3 * a * y + a * c + b.
// If the symbol does not exist in the row, this is a no-op.
func (r *row) substitute(s Symbol, other *row) {
	if c, exists := r.cells.Get(s); exists {
		coefficient := c.(Float)
		r.cells.Remove(s)
		r.insertRow(other, coefficient)
	}
}

// Test whether a row is composed of all dummy variables.
func (r *row) allDummies() bool {
	return r.cells.All(func(k, v interface{}) bool {
		return k.(Symbol) == SymbolDummy
	})
}
