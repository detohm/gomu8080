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
		p.stax(&p.B, &p.C)
	case 0x03:
		p.inx(&p.B, &p.C)
	case 0x04:
		p.inr(&p.B)
	case 0x05:
		p.dcr(&p.B)
	case 0x06:
		p.mvi(&p.B)
	case 0x07:
		p.rlc()
	case 0x08:
		p.nop()
	case 0x09:
		p.dad(&p.B, &p.C)
	case 0x10:
		p.ldax(&p.B, &p.C)
	default:
		p.unimplemented()
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
	result16 := uint16(*reg) + 1
	*reg = byte(result16 & 0x00ff)

	p.SetZero(*reg)
	p.SetSign(*reg)
	p.SetParity(*reg)
	p.SetAuxiliaryCarry(result16)
}

// Decrease value of 8-bit register by 1
func (p *Processor) dcr(reg *byte) {
	p.dasm("DCR")
	result16 := uint16(*reg) - 1
	*reg = byte(result16 & 0xff)

	p.SetZero(*reg)
	p.SetSign(*reg)
	p.SetParity(*reg)
	p.SetAuxiliaryCarry(result16)
}

// Move immediate data
func (p *Processor) mvi(reg *byte) {
	p.dasm("MVI")
	*reg = p.mmu.Memory[p.PC]
	p.PC += 1
}

// Rotate accumulator left
func (p *Processor) rlc() {
	p.dasm("RLC")
	aux := p.A
	p.A = aux<<1 | aux>>7
	p.Carry = (aux >> 7) > 0
}

// Double Add - add specified register pair to HL
func (p *Processor) dad(msb *byte, lsb *byte) {
	p.dasm("DAD")
	adder := uint32(*msb)<<8 | uint32(*lsb)
	result := uint32(p.H)<<8 | uint32(p.L)
	result += adder
	p.H = byte(result >> 8)
	p.L = byte(result & 0xFF)
	p.Carry = result > 0xFFFF // greater than 2-byte

}

// Load Accumulator - load data from the provided address
func (p *Processor) ldax(msb *byte, lsb *byte) {
	p.dasm("LDAX")
	address := uint16(*msb)<<8 | uint16(*lsb)
	p.A = p.mmu.Memory[address]
}

// Store Accumulator - store data from A to the provided address
func (p *Processor) stax(msb *byte, lsb *byte) {
	p.dasm("STAX")
	address := uint16(*msb)<<8 | uint16(*lsb)
	p.mmu.Memory[address] = p.A
}

// Unimplemented
func (p *Processor) unimplemented() {
	p.dasm(fmt.Sprintf("%02x:UNIMPLEMENTED", p.mmu.Memory[p.PC]))
}

// Flags Helper
func (p *Processor) SetZero(result byte) {
	p.Zero = (result & 0xFF) == 0
}
func (p *Processor) SetSign(result byte) {
	p.Sign = (result & 0x80) == 0x80
}
func (p *Processor) SetParity(result byte) {

	oneCount := 0

	for i := 0; i < 8; i++ {

		if result&0x01 > 0 {
			oneCount += 1
		}
		result = result >> 1
	}

	p.Parity = oneCount%2 == 0
}
func (p *Processor) SetAuxiliaryCarry(result uint16) {
	// TODO - implement ac calculation
	p.AuxiliaryCarry = false
}

// DEBUGGER
// disassemble helper
func (p *Processor) dasm(opcode string) {

	if !p.DebugMode {
		return
	}

	fmt.Printf("%s\n", opcode)
}
