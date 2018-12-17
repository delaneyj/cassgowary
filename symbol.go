package cassgowary

type Symbol int

type Symbols []Symbol

const (
	SymbolInvalid Symbol = iota
	SymbolExternal
	SymbolSlack
	SymbolError
	SymbolDummy
)
