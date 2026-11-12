package bytecode

type Op uint8

const (
	Push Op = iota
	Load
	Store
	Pop
	Unary
	Binary
	Compare
	Call
	JumpFalse
	Jump
	Return
)

type Instruction struct {
	Op  Op
	Arg any
}

type Code struct {
	Name         string
	Params       []string
	Instructions []Instruction
}
