package gomu8080

import (
	"fmt"
)

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
	SP uint16

	// program counter
	PC uint16

	// processor flags
	Sign           bool
	Zero           bool
	Parity         bool
	Carry          bool
	AuxiliaryCarry bool

	// MMU
	mmu *MMU

	// debug
	DebugMode bool
}

func NewProcessor(mmu *MMU, debugMode bool) *Processor {
	p := &Processor{}
	p.mmu = mmu
	p.DebugMode = debugMode
	return p
}

func (p *Processor) Run() {

	opcode := p.mmu.Memory[p.PC]
	p.PC += 1

	switch opcode {
	case 0x00:
		p.nop()
	case 0x01:
		p.lxi(&p.B, &p.C)
	case 0x02:
		p.unimplemented()
	case 0x03:
		p.inx(&p.B, &p.C)
	case 0x04:
		p.inr(&p.B)
	case 0x05:
		p.dcr(&p.B)
	}
}

// Instruction
func (p *Processor) nop() {
	p.dasm("NOP")
}

func (p *Processor) lxi(msb *byte, lsb *byte) {
	p.dasm("LXI") // TODO - identify which register?
	*lsb = p.mmu.Memory[p.PC+1]
	*msb = p.mmu.Memory[p.PC+2]
	p.PC += 2
}

// Increase value of register pair by 1
func (p *Processor) inx(msb *byte, lsb *byte) {
	p.dasm("INX")
	*lsb += 1
	if *lsb == 0 {
		*msb += 1
	}
}

// Increase value of 8-bit register by 1
func (p *Processor) inr(reg *byte) {
	p.dasm("INR")
	result := uint16(*reg) + 1
	*reg = byte(result & 0x00ff)

	// TODO - calculate flags
}

// Decrease value of 8-bit register by 1
func (p *Processor) dcr(reg *byte) {
	p.dasm("DCR")
	result := uint16(*reg) - 1
	*reg = byte(result & 0xff)

	// TODO - calculate flags
}

func (p *Processor) unimplemented() {
	p.dasm(fmt.Sprintf("%02x:UNIMPLEMENTED", p.mmu.Memory[p.PC]))
}

// DEBUGGER
// disassemble helper
func (p *Processor) dasm(opcode string) {

	if !p.DebugMode {
		return
	}
	fmt.Printf("%s\n", opcode)
}
