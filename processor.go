package gomu8080

type Processor struct {
	// main register
	A byte
	F byte
	B byte
	C byte
	D byte
	E byte
	H byte
	L byte

	// stack register
	SP byte

	// program counter
	PC byte

	// processor flags
	Sign           bool
	Zero           bool
	Parity         bool
	Carry          bool
	AuxiliaryCarry bool
}
