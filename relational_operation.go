package cassgowary

type RelationalOperator int

const (
	OP_LE RelationalOperator = iota
	OP_GE
	OP_EQ
)

var OperationNames = map[RelationalOperator]string{
	OP_LE: "LEQ",
	OP_GE: "GEQ",
	OP_EQ: "EQ",
}

var OperationFromString = map[string]RelationalOperator{
	"LEQ": OP_LE,
	"GEQ": OP_GE,
	"EQ":  OP_EQ,
}
