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

	// processor state
	IsHalt bool
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

	// TODO - validate address bound
	address := (uint16(p.H) << 8) | uint16(p.L)

	switch opcode {
	/* 0x */
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
	case 0x0A:
		p.ldax(&p.B, &p.C)
	case 0x0B:
		p.dcx(&p.B, &p.C)
	case 0x0C:
		p.inr(&p.C)
	case 0x0D:
		p.dcr(&p.C)
	case 0x0E:
		p.mvi(&p.C)
	case 0x0F:
		p.rrc()
	/* 1x */
	case 0x10:
		p.nop()
	case 0x11:
		p.lxi(&p.H, &p.L)
	case 0x12:
		p.stax(&p.D, &p.E)
	case 0x13:
		p.inx(&p.D, &p.E)
	case 0x14:
		p.inr(&p.D)
	case 0x15:
		p.dcr(&p.D)
	case 0x16:
		p.mvi(&p.D)
	case 0x17:
		p.ral()
	case 0x18:
		p.nop()
	case 0x19:
		p.dad(&p.D, &p.E)
	case 0x1A:
		p.ldax(&p.D, &p.E)
	case 0x1B:
		p.dcx(&p.D, &p.E)
	case 0x1C:
		p.inr(&p.E)
	case 0x1D:
		p.dcr(&p.E)
	case 0x1E:
		p.mvi(&p.E)
	case 0x1F:
		p.rar()

	/* 2x */
	case 0x20:
		p.nop()
	case 0x21:
		p.lxi(&p.H, &p.L)
	case 0x22:
		p.shld()
	case 0x23:
		p.inx(&p.H, &p.L)
	case 0x24:
		p.inr(&p.H)
	case 0x25:
		p.dcr(&p.H)
	case 0x26:
		p.mvi(&p.H)
	case 0x27:
		p.daa()
	case 0x28:
		p.nop()
	case 0x29:
		p.dad(&p.H, &p.L)
	case 0x2A:
		p.lhld()
	case 0x2B:
		p.dcx(&p.H, &p.L)
	case 0x2C:
		p.inr(&p.L)
	case 0x2D:
		p.dcr(&p.L)
	case 0x2E:
		p.mvi(&p.L)
	case 0x2F:
		p.cma()

	/* 3x */
	case 0x30:
		p.nop()
	case 0x31:
		p.lxi16(&p.SP)
	case 0x32:
		p.sta()
	case 0x33:
		p.inx16(&p.SP)
	case 0x34:
		p.inr(&p.mmu.Memory[address])
	case 0x35:
		p.dcr(&p.mmu.Memory[address])
	case 0x36:
		p.mvi(&p.mmu.Memory[address])
	case 0x37:
		p.stc()
	case 0x38:
		p.nop()
	case 0x39:
		p.dad16(&p.SP)
	case 0x3A:
		p.lda()
	case 0x3B:
		p.dcx16(&p.SP)
	case 0x3C:
		p.inr(&p.A)
	case 0x3D:
		p.dcr(&p.A)
	case 0x3E:
		p.mvi(&p.A)
	case 0x3F:
		p.cmc()

	/* 4x */
	// B - destination
	case 0x40:
		p.mov(&p.B, &p.B)
	case 0x41:
		p.mov(&p.B, &p.C)
	case 0x42:
		p.mov(&p.B, &p.D)
	case 0x43:
		p.mov(&p.B, &p.E)
	case 0x44:
		p.mov(&p.B, &p.H)
	case 0x45:
		p.mov(&p.B, &p.L)
	case 0x46:
		p.mov(&p.B, &p.mmu.Memory[address])
	case 0x47:
		p.mov(&p.B, &p.A)
	// C - destination
	case 0x48:
		p.mov(&p.C, &p.B)
	case 0x49:
		p.mov(&p.C, &p.C)
	case 0x4A:
		p.mov(&p.C, &p.D)
	case 0x4B:
		p.mov(&p.C, &p.E)
	case 0x4C:
		p.mov(&p.C, &p.H)
	case 0x4D:
		p.mov(&p.C, &p.L)
	case 0x4E:
		p.mov(&p.C, &p.mmu.Memory[address])
	case 0x4F:
		p.mov(&p.C, &p.A)

	/* 5x */
	// D - destination
	case 0x50:
		p.mov(&p.D, &p.B)
	case 0x51:
		p.mov(&p.D, &p.C)
	case 0x52:
		p.mov(&p.D, &p.D)
	case 0x53:
		p.mov(&p.D, &p.E)
	case 0x54:
		p.mov(&p.D, &p.H)
	case 0x55:
		p.mov(&p.D, &p.L)
	case 0x56:
		p.mov(&p.D, &p.mmu.Memory[address])
	case 0x57:
		p.mov(&p.D, &p.A)
	// E - destination
	case 0x58:
		p.mov(&p.E, &p.B)
	case 0x59:
		p.mov(&p.E, &p.C)
	case 0x5A:
		p.mov(&p.E, &p.D)
	case 0x5B:
		p.mov(&p.E, &p.E)
	case 0x5C:
		p.mov(&p.E, &p.H)
	case 0x5D:
		p.mov(&p.E, &p.L)
	case 0x5E:
		p.mov(&p.E, &p.mmu.Memory[address])
	case 0x5F:
		p.mov(&p.E, &p.A)

	/* 6x */
	// H - destination
	case 0x60:
		p.mov(&p.H, &p.B)
	case 0x61:
		p.mov(&p.H, &p.C)
	case 0x62:
		p.mov(&p.H, &p.D)
	case 0x63:
		p.mov(&p.H, &p.E)
	case 0x64:
		p.mov(&p.H, &p.H)
	case 0x65:
		p.mov(&p.H, &p.L)
	case 0x66:
		p.mov(&p.H, &p.mmu.Memory[address])
	case 0x67:
		p.mov(&p.H, &p.A)
	// L - destination
	case 0x68:
		p.mov(&p.L, &p.B)
	case 0x69:
		p.mov(&p.L, &p.C)
	case 0x6A:
		p.mov(&p.L, &p.D)
	case 0x6B:
		p.mov(&p.L, &p.E)
	case 0x6C:
		p.mov(&p.L, &p.H)
	case 0x6D:
		p.mov(&p.L, &p.L)
	case 0x6E:
		p.mov(&p.L, &p.mmu.Memory[address])
	case 0x6F:
		p.mov(&p.L, &p.A)

	/* 7x */
	// M - destination
	case 0x70:
		p.mov(&p.mmu.Memory[address], &p.B)
	case 0x71:
		p.mov(&p.mmu.Memory[address], &p.C)
	case 0x72:
		p.mov(&p.mmu.Memory[address], &p.D)
	case 0x73:
		p.mov(&p.mmu.Memory[address], &p.E)
	case 0x74:
		p.mov(&p.mmu.Memory[address], &p.H)
	case 0x75:
		p.mov(&p.mmu.Memory[address], &p.L)
	// HALT
	case 0x76:
		p.hlt()
	case 0x77:
		p.mov(&p.mmu.Memory[address], &p.A)
	// A - destination
	case 0x78:
		p.mov(&p.A, &p.B)
	case 0x79:
		p.mov(&p.A, &p.C)
	case 0x7A:
		p.mov(&p.A, &p.D)
	case 0x7B:
		p.mov(&p.A, &p.E)
	case 0x7C:
		p.mov(&p.A, &p.H)
	case 0x7D:
		p.mov(&p.A, &p.L)
	case 0x7E:
		p.mov(&p.A, &p.mmu.Memory[address])
	case 0x7F:
		p.mov(&p.A, &p.A)

	default:
		p.unimplemented()
	}
}

// Instruction
func (p *Processor) nop() {
	p.dasm("NOP")
}

// Load register pair immediate
func (p *Processor) lxi(msb *byte, lsb *byte) {
	p.dasm("LXI") // TODO - identify which register?
	*lsb = p.mmu.Memory[p.PC]
	*msb = p.mmu.Memory[p.PC+1]
	p.PC += 2
}

// Load register pair immediate (16-bit)
func (p *Processor) lxi16(reg *uint16) {
	p.dasm("LXI")
	lsb := p.mmu.Memory[p.PC]
	msb := p.mmu.Memory[p.PC+1]
	*reg = uint16(msb)<<8 | uint16(lsb)
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

// Increase value of register pair by 1 (16-bit input)
func (p *Processor) inx16(reg *uint16) {
	p.dasm("INX")
	*reg += 1
}

// Decrease value of register pair by 1
func (p *Processor) dcx(msb *byte, lsb *byte) {
	p.dasm("DCX")
	*lsb -= 1
	if *lsb == 0xFF {
		*msb -= 1
	}
}

// Decrease value of register pair by 1 (16-bit)
func (p *Processor) dcx16(reg *uint16) {
	p.dasm("DCX")
	*reg -= 1
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

// Move Instruction - move data from src to dst
func (p *Processor) mov(dst *byte, src *byte) {
	p.dasm("MOV")
	*dst = *src
}

// Rotate accumulator left
func (p *Processor) rlc() {
	p.dasm("RLC")
	aux := p.A
	p.A = aux<<1 | aux>>7
	p.Carry = (aux >> 7) > 0
}

// Rotate accumulator right
func (p *Processor) rrc() {
	p.dasm("RRC")
	aux := p.A
	p.A = aux>>1 | ((aux << 7) & 0x80)
	p.Carry = aux&0x01 > 0
}

// Rorate accumulator left through carry
func (p *Processor) ral() {
	p.dasm("RAL")
	aux := p.A
	p.A = aux << 1

	if p.Carry {
		p.A = p.A | 0x01
	}
	p.Carry = (aux >> 7) > 0
}

// Rotate accumulator right through carry
func (p *Processor) rar() {
	p.dasm("RAR")
	aux := p.A
	p.A = aux >> 1
	if p.Carry {
		p.A += 0x80 // set 1 to most significant bit
	}
	p.Carry = aux&0x01 > 0
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

// Double Add - add specified register pair to HL (16-bit operand)
func (p *Processor) dad16(reg *uint16) {
	p.dasm("DAD")
	adder := uint32(*reg)
	HL := uint32(p.H)<<8 | uint32(p.L)
	HL += adder
	p.H = byte(HL >> 8)
	p.L = byte(HL & 0xFF)
	p.Carry = HL > 0xFFFF
}

// Load Accumulator - load data from the provided address
func (p *Processor) ldax(msb *byte, lsb *byte) {
	p.dasm("LDAX")
	address := uint16(*msb)<<8 | uint16(*lsb)
	p.A = p.mmu.Memory[address]
}

// Load accumulator direct from the operand address to A
func (p *Processor) lda() {
	p.dasm("LDA")
	address := uint16(p.mmu.Memory[p.PC+1]) << 8
	address |= uint16(p.mmu.Memory[p.PC])
	p.A = p.mmu.Memory[address]
	p.PC += 2
}

// Store Accumulator - store data from A to the provided address
func (p *Processor) stax(msb *byte, lsb *byte) {
	p.dasm("STAX")
	address := uint16(*msb)<<8 | uint16(*lsb)
	p.mmu.Memory[address] = p.A
}

// Store accumulator direct from A to the operand address
func (p *Processor) sta() {
	p.dasm("STA")
	address := uint16(p.mmu.Memory[p.PC+1]) << 8
	address |= uint16(p.mmu.Memory[p.PC])
	p.mmu.Memory[address] = p.A
	p.PC += 2
}

// Store H and L direct to the provided address and the next one
func (p *Processor) shld() {
	p.dasm("SHLD")
	lsb := uint16(p.mmu.Memory[p.PC])
	msb := uint16(p.mmu.Memory[p.PC+1])
	address := (msb << 8) | lsb

	if address < 65535 {
		p.mmu.Memory[address] = p.L
		p.mmu.Memory[address+1] = p.H
	}
	p.PC += 2
}

// Load H and L direct from the provided address and the next one
func (p *Processor) lhld() {
	p.dasm("LHLD")
	lsb := uint16(p.mmu.Memory[p.PC])
	msb := uint16(p.mmu.Memory[p.PC+1])
	address := (msb << 8) | lsb
	if address < 65535 {
		p.L = p.mmu.Memory[address]
		p.H = p.mmu.Memory[address+1]
	}
	p.PC += 2
}

// Decimal Adjust Accumulator
func (p *Processor) daa() {
	p.dasm("DAA")
	// divide 4-bit
	lsb4 := p.A & 0x0F
	msb4 := p.A >> 4

	if lsb4 > 0x09 || p.AuxiliaryCarry {
		// result := uint16(lsb4) + uint16(0x06)
		p.A += 0x06
	}

	if msb4 > 0x09 || p.Carry {
		p.A += (0x06 << 4)
	}
	// TODO - implement logic to update flags
}

// Complement Accumulator
func (p *Processor) cma() {
	p.dasm("CMA")
	p.A = 0xFF ^ p.A
}

// Complement carry
func (p *Processor) cmc() {
	p.dasm("CMC")
	p.Carry = !p.Carry
}

// Set Carry
func (p *Processor) stc() {
	p.Carry = true
}

// Halt
func (p *Processor) hlt() {
	p.IsHalt = true
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
