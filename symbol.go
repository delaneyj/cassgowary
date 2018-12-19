package cassgowary

type symbols []*symbol

type symbolType int

var nextSymbolID = 0

const (
	symbolInvalid symbolType = iota
	symbolExternal
	symbolSlack
	symbolError
	symbolDummy
)

type symbol struct {
	id   int
	kind symbolType
}

func newSymbol() *symbol {
	s := &symbol{
		id:   nextSymbolID,
		kind: symbolInvalid,
	}
	nextSymbolID++
	return s
}

func newSymbolFrom(t symbolType) *symbol {
	s := newSymbol()
	s.kind = t
	return s
}
