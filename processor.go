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

	/* 8x */
	// add
	case 0x80:
		p.add(&p.B)
	case 0x81:
		p.add(&p.C)
	case 0x82:
		p.add(&p.D)
	case 0x83:
		p.add(&p.E)
	case 0x84:
		p.add(&p.H)
	case 0x85:
		p.add(&p.L)
	case 0x86:
		p.add(&p.mmu.Memory[address])
	case 0x87:
		p.add(&p.A)
	// adc
	case 0x88:
		p.adc(&p.B)
	case 0x89:
		p.adc(&p.C)
	case 0x8A:
		p.adc(&p.D)
	case 0x8B:
		p.adc(&p.E)
	case 0x8C:
		p.adc(&p.H)
	case 0x8D:
		p.adc(&p.L)
	case 0x8E:
		p.adc(&p.mmu.Memory[address])
	case 0x8F:
		p.adc(&p.A)

	/* 9x */
	// sub
	case 0x90:
		p.sub(&p.B)
	case 0x91:
		p.sub(&p.C)
	case 0x92:
		p.sub(&p.D)
	case 0x93:
		p.sub(&p.E)
	case 0x94:
		p.sub(&p.H)
	case 0x95:
		p.sub(&p.L)
	case 0x96:
		p.sub(&p.mmu.Memory[address])
	case 0x97:
		p.sub(&p.A)
		// sbb
	case 0x98:
		p.sbb(&p.B)
	case 0x99:
		p.sbb(&p.C)
	case 0x9A:
		p.sbb(&p.D)
	case 0x9B:
		p.sbb(&p.E)
	case 0x9C:
		p.sbb(&p.H)
	case 0x9D:
		p.sbb(&p.L)
	case 0x9E:
		p.sbb(&p.mmu.Memory[address])
	case 0x9F:
		p.sbb(&p.A)

	/* Ax */
	// ana
	case 0xA0:
		p.ana(&p.B)
	case 0xA1:
		p.ana(&p.C)
	case 0xA2:
		p.ana(&p.D)
	case 0xA3:
		p.ana(&p.E)
	case 0xA4:
		p.ana(&p.H)
	case 0xA5:
		p.ana(&p.L)
	case 0xA6:
		p.ana(&p.mmu.Memory[address])
	case 0xA7:
		p.ana(&p.A)
	// xra
	case 0xA8:
		p.xra(&p.B)
	case 0xA9:
		p.xra(&p.C)
	case 0xAA:
		p.xra(&p.D)
	case 0xAB:
		p.xra(&p.E)
	case 0xAC:
		p.xra(&p.H)
	case 0xAD:
		p.xra(&p.L)
	case 0xAE:
		p.xra(&p.mmu.Memory[address])
	case 0xAF:
		p.xra(&p.A)

	/* Bx */
	// ora
	case 0xB0:
		p.ora(&p.B)
	case 0xB1:
		p.ora(&p.C)
	case 0xB2:
		p.ora(&p.D)
	case 0xB3:
		p.ora(&p.E)
	case 0xB4:
		p.ora(&p.H)
	case 0xB5:
		p.ora(&p.L)
	case 0xB6:
		p.ora(&p.mmu.Memory[address])
	case 0xB7:
		p.ora(&p.A)
	// cmp
	case 0xB8:
		p.cmp(&p.B)
	case 0xB9:
		p.cmp(&p.C)
	case 0xBA:
		p.cmp(&p.D)
	case 0xBB:
		p.cmp(&p.E)
	case 0xBC:
		p.cmp(&p.H)
	case 0xBD:
		p.cmp(&p.L)
	case 0xBE:
		p.cmp(&p.mmu.Memory[address])
	case 0xBF:
		p.cmp(&p.A)

	/* Cx */
	case 0xC0:
		p.rnz()
	case 0xC1:
		p.pop(&p.B, &p.C)
	case 0xC2:
		p.jnz()
	case 0xC3:
		p.jmp()
	case 0xC4:
		p.cnz()
	case 0xC5:
		p.push(&p.B, &p.C)
	case 0xC6:
		p.adi()
	case 0xC7:
		p.rst(0)
	case 0xC8:
		p.rz()
	case 0xC9:
		p.ret()
	case 0xCA:
		p.jz()
	case 0xCB:
		p.jmp()
	case 0xCC:
		p.cz()
	case 0xCD:
		p.call()
	case 0xCE:
		p.aci()
	case 0xCF:
		p.rst(1)

	/* Dx */
	case 0xD0:
		p.rnc()
	case 0xD1:
		p.pop(&p.D, &p.E)
	case 0xD2:
		p.jnc()
	case 0xD3:
		p.out()
	case 0xD4:
		p.cnc()
	case 0xD5:
		p.push(&p.D, &p.E)
	case 0xD6:
		p.sui()
	case 0xD7:
		p.rst(2)
	case 0xD8:
		p.rc()
	case 0xD9:
		p.ret()
	case 0xDA:
		p.jc()
	case 0xDB:
		p.in()
	case 0xDC:
		p.cc()
	case 0xDD:
		p.call()
	case 0xDE:
		p.sbi()
	case 0xDF:
		p.rst(3)

	/* Ex */
	case 0xE0:
		p.rpo()
	case 0xE1:
		p.pop(&p.H, &p.L)
	case 0xE2:
		p.jpo()
	case 0xE3:
		p.xthl()
	case 0xE4:
		p.cpo()
	case 0xE5:
		p.push(&p.H, &p.L)
	case 0xE6:
		p.ani()
	case 0xE7:
		p.rst(4)
	case 0xE8:
		p.rpe()
	case 0xE9:
		p.pchl()
	case 0xEA:
		p.jpe()
	case 0xEB:
		p.xchg()
	case 0xEC:
		p.cpe()
	case 0xED:
		p.call()
	case 0xEE:
		p.xri()
	case 0xEF:
		p.rst(5)

	/* Fx */
	case 0xF0:
		p.rp()
	case 0xF1:
		p.popPSW()
	case 0xF2:
		p.jp()
	case 0xF3:
		// DI
	case 0xF4:
		p.cp()
	case 0xF5:
		p.pushPSW()
	case 0xF6:
		p.ori()
	case 0xF7:
		p.rst(6)
	case 0xF8:
		p.rm()
	case 0xF9:
		p.sphl()
	case 0xFA:
		p.jm()
	case 0xFB:
		// EI
	case 0xFC:
		p.cm()
	case 0xFD:
		p.call()
	case 0xFE:
		p.cpi()
	case 0xFF:
		p.rst(7)

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

// Exchange register pair HL <-> DE
func (p *Processor) xchg() {
	p.dasm("XCHG")
	tlsb := p.L
	tmsb := p.H
	p.L = p.E
	p.H = p.D
	p.E = tlsb
	p.D = tmsb
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

// Add register or memory to accumulator
func (p *Processor) add(reg *byte) {
	p.dasm("ADD")
	result := uint16(p.A) + uint16(*reg)
	lsb := byte(result & 0x00FF)
	p.A = lsb

	p.SetSign(lsb)
	p.SetZero(lsb)
	// TODO implement auxiliary carry
	p.SetParity(lsb)
	p.Carry = result > 0xFF
}

// Add register or memory to accumulator with carry
func (p *Processor) adc(reg *byte) {
	p.dasm("ADC")
	result := uint16(p.A) + uint16(*reg)
	if p.Carry {
		result += 0x01
	}
	lsb := byte(result & 0x00FF)
	p.A = lsb

	p.SetSign(lsb)
	p.SetZero(lsb)
	// TODO implement auxiliary carry
	p.SetParity(lsb)
	p.Carry = result > 0xFF
}

// Add Immediate to Accumulator
func (p *Processor) adi() {
	p.dasm("ADI")
	result := uint16(p.A) + uint16(p.mmu.Memory[p.PC])
	lsb := byte(result & 0x00FF)

	p.SetSign(lsb)
	p.SetZero(lsb)
	// TODO implement auxiliary carry
	p.SetParity(lsb)
	p.Carry = result > 0xFF
}

// Add Immediate to Accumulator with carry
func (p *Processor) aci() {
	p.dasm("ACI")
	result := uint16(p.A) + uint16(p.mmu.Memory[p.PC])
	if p.Carry {
		result += 0x01
	}

	lsb := byte(result & 0x00FF)

	p.SetSign(lsb)
	p.SetZero(lsb)
	// TODO implement auxiliary carry
	p.SetParity(lsb)
	p.Carry = result > 0xFF
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

// Subtract register or memory from accumulator
func (p *Processor) sub(reg *byte) {
	p.dasm("SUB")
	result := uint16(p.A) + uint16(^*reg) + 0x1

	lsb := byte(result & 0x00FF)
	p.SetSign(lsb)
	p.SetZero(lsb)
	// TODO implement auxiliary carry
	p.SetParity(lsb)
	if result <= 0x00FF {
		p.Carry = true
	}
}

// Subtract register or memory from accumulator with borrow
func (p *Processor) sbb(reg *byte) {
	p.dasm("SBB")
	result := uint16(p.A) + uint16(^*reg) + 0x1
	if p.Carry {
		result += 0x01
	}

	lsb := byte(result & 0x00FF)
	p.SetSign(lsb)
	p.SetZero(lsb)
	// TODO implement auxiliary carry
	p.SetParity(lsb)
	if result <= 0x00FF {
		p.Carry = true
	}
}

// Subtract immediate from accumulator
func (p *Processor) sui() {
	p.dasm("SUI")
	result := uint16(p.A) + uint16(^p.mmu.Memory[p.PC]) + 0x1
	lsb := byte(result & 0x00FF)
	p.SetSign(lsb)
	p.SetZero(lsb)
	// TODO implement auxiliary carry
	p.SetParity(lsb)
	if result <= 0x00FF {
		p.Carry = true
	}
	p.PC += 1
}

// Subtract immediate from accumulator with borrow
func (p *Processor) sbi() {
	p.dasm("SBI")
	result := uint16(p.A) + uint16(^p.mmu.Memory[p.PC]) + 0x1
	if p.Carry {
		result += 0x01
	}

	lsb := byte(result & 0x00FF)
	p.SetSign(lsb)
	p.SetZero(lsb)
	// TODO implement auxiliary carry
	p.SetParity(lsb)
	if result <= 0x00FF {
		p.Carry = true
	}
	p.PC += 1
}

// Logical AND register or memory with accumulator
func (p *Processor) ana(reg *byte) {
	p.dasm("ANA")
	p.A &= *reg

	p.SetSign(p.A)
	p.SetZero(p.A)
	p.AuxiliaryCarry = false
	p.SetParity(p.A)
	p.Carry = false
}

// Logial AND immediate with accumulator
func (p *Processor) ani() {
	p.dasm("ANI")
	p.A &= p.mmu.Memory[p.PC]

	p.SetSign(p.A)
	p.SetZero(p.A)
	p.AuxiliaryCarry = false
	p.SetParity(p.A)
	p.Carry = false
	p.PC += 1
}

// Logical XOR register or memory with accumulator
func (p *Processor) xra(reg *byte) {
	p.dasm("XRA")
	p.A ^= *reg

	p.SetSign(p.A)
	p.SetZero(p.A)
	p.AuxiliaryCarry = false
	p.SetParity(p.A)
	p.Carry = false
}

// Logical XOR immediate with accumulator
func (p *Processor) xri() {
	p.dasm("XRI")
	p.A ^= p.mmu.Memory[p.PC]

	p.SetSign(p.A)
	p.SetZero(p.A)
	p.AuxiliaryCarry = false
	p.SetParity(p.A)
	p.Carry = false
	p.PC += 1
}

// Logical OR register or memory with accumulator
func (p *Processor) ora(reg *byte) {
	p.dasm("ORA")
	p.A |= *reg

	p.SetSign(p.A)
	p.SetZero(p.A)
	// TODO implement auxiliary carry
	p.SetParity(p.A)
	p.Carry = false
}

// Logical OR immediate with accumulator
func (p *Processor) ori() {
	p.dasm("ORI")
	p.A |= p.mmu.Memory[p.PC]

	p.SetSign(p.A)
	p.SetZero(p.A)
	// TODO implement auxiliary carry
	p.SetParity(p.A)
	p.Carry = false
	p.PC += 1
}

// Compare register or memory with accumulator
func (p *Processor) cmp(reg *byte) {
	p.dasm("CMP")
	result := uint16(p.A) + uint16(^*reg) + 0x01
	lsb := byte(result & 0x00FF)
	p.SetSign(lsb)
	p.SetZero(lsb)
	// TODO implement auxiliary carry
	p.SetParity(lsb)
	if result <= 0xFF {
		p.Carry = true
	}
}

// Compare immediate with accumulator
func (p *Processor) cpi() {
	p.dasm("CPI")
	result := uint16(p.A) + uint16(^p.mmu.Memory[p.PC]) + 0x01
	lsb := byte(result & 0x00FF)
	p.SetSign(lsb)
	p.SetZero(lsb)
	// TODO implement auxiliary carry
	p.SetParity(lsb)
	if result <= 0xFF {
		p.Carry = true
	}
	p.PC += 1
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

// Load SP from HL
func (p *Processor) sphl() {
	p.dasm("SPHL")
	p.SP = (uint16(p.H) << 8) | uint16(p.L)
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
	p.dasm("STC")
	p.Carry = true
}

// Halt
func (p *Processor) hlt() {
	p.dasm("HLT")
	p.IsHalt = true
}

/* subroutine instruction */
// Internal Call subroutine
func (p *Processor) intCall() {

	address := (uint16(p.PC+1) << 8) | uint16(p.PC)
	returnAddress := p.PC + 2

	// push return address into stack
	// msb
	p.mmu.Memory[p.SP-1] = byte(returnAddress >> 8)
	// lsb
	p.mmu.Memory[p.SP-2] = byte(returnAddress & 0x00FF)

	// move stack pointer downward as "push"
	p.SP -= 2

	// set next instruction fetch to the call address
	p.PC = address
}

// Return from subroutine
func (p *Processor) intRet() {
	p.dasm("RET")
	p.PC = uint16(p.mmu.Memory[p.SP+1]) << 8
	p.PC |= uint16(p.mmu.Memory[p.SP])

	p.SP += 2
}

// Call instruction
func (p *Processor) call() {
	p.dasm("CALL")
	p.intCall()
}

// Call if not zero
func (p *Processor) cnz() {
	p.dasm("CNZ")
	if !p.Zero {
		p.intCall()
	} else {
		p.PC += 2
	}
}

// Call if zero
func (p *Processor) cz() {
	p.dasm("CZ")
	if p.Zero {
		p.intCall()
	} else {
		p.PC += 2
	}
}

// Call if not carry
func (p *Processor) cnc() {
	p.dasm("CNC")
	if !p.Carry {
		p.intCall()
	} else {
		p.PC += 2
	}
}

// Call if carry
func (p *Processor) cc() {
	p.dasm("CC")
	if p.Carry {
		p.intCall()
	} else {
		p.PC += 2
	}
}

// Call if parity odd (zero)
func (p *Processor) cpo() {
	p.dasm("CPO")
	if !p.Parity {
		p.intCall()
	} else {
		p.PC += 2
	}
}

// Call if parity even (one)
func (p *Processor) cpe() {
	p.dasm("CPE")
	if p.Parity {
		p.intCall()
	} else {
		p.PC += 2
	}
}

// Call if plus (sign-zero)
func (p *Processor) cp() {
	p.dasm("CP")
	if !p.Sign {
		p.intCall()
	} else {
		p.PC += 2
	}
}

// Call if minus (sign-one)
func (p *Processor) cm() {
	p.dasm("CM")
	if p.Sign {
		p.intCall()
	} else {
		p.PC += 2
	}
}

// Return instruction
func (p *Processor) ret() {
	p.dasm("RET")
	p.intRet()
}

// return if not zero
func (p *Processor) rnz() {
	p.dasm("RNZ")
	if !p.Zero {
		p.intRet()
	}
}

// return if zero
func (p *Processor) rz() {
	p.dasm("RZ")
	if p.Zero {
		p.intRet()
	}
}

// return if not carry
func (p *Processor) rnc() {
	p.dasm("RNC")
	if !p.Carry {
		p.intRet()
	}
}

// return if carry
func (p *Processor) rc() {
	p.dasm("RC")
	if p.Carry {
		p.intRet()
	}
}

// return if parity odd (zero)
func (p *Processor) rpo() {
	p.dasm("RPO")
	if !p.Parity {
		p.intRet()
	}
}

// return if parity even (one)
func (p *Processor) rpe() {
	p.dasm("RPE")
	if p.Parity {
		p.intRet()
	}
}

// return if plus(sign-zero)
func (p *Processor) rp() {
	p.dasm("RP")
	if !p.Sign {
		p.intRet()
	}
}

// return if minus(sign-one)
func (p *Processor) rm() {
	p.dasm("RM")
	if p.Sign {
		p.intRet()
	}
}

// restart
func (p *Processor) rst(pos uint16) {
	p.dasm("RST")
	msb := byte(p.PC >> 8)
	lsb := byte(p.PC & 0x00FF)
	p.mmu.Memory[p.SP-1] = msb
	p.mmu.Memory[p.SP-2] = lsb
	p.SP += 2

	address := uint16(pos << 3)
	p.PC = address
}

/* STACK Instruction */
// push data onto stack
func (p *Processor) push(msb *byte, lsb *byte) {
	p.dasm("PUSH")
	p.mmu.Memory[p.SP-1] = *msb
	p.mmu.Memory[p.SP-2] = *lsb
	p.SP -= 2
}

// pop data off stack (register pair)
func (p *Processor) pop(msb *byte, lsb *byte) {
	p.dasm("POP")
	*lsb = p.mmu.Memory[p.SP]
	*msb = p.mmu.Memory[p.SP+1]
	p.SP += 2
}

// push PSW onto stack
func (p *Processor) pushPSW() {
	p.dasm("PUSH PSW")
	p.mmu.Memory[p.SP-1] = p.A

	// assemble flags byte
	flags := byte(0)
	if p.Carry {
		flags |= 0x00000001
	}
	if p.Parity {
		flags |= 0b00000100
	}
	if p.AuxiliaryCarry {
		flags |= 0b00010000
	}
	if p.Zero {
		flags |= 0b01000000
	}
	if p.Sign {
		flags |= 0b10000000
	}
	p.mmu.Memory[p.SP-2] = flags
	p.SP -= 2
}

// pop data off stack to PSW
func (p *Processor) popPSW() {
	p.dasm("POP PSW")
	flags := p.mmu.Memory[p.SP]
	p.A = p.mmu.Memory[p.SP+1]
	p.SP += 2

	//restore flags
	p.Carry = false
	p.Parity = false
	p.AuxiliaryCarry = false
	p.Zero = false
	p.Sign = false
	if flags&0b00000001 > 0 {
		p.Carry = true
	}
	if flags&0b00000100 > 0 {
		p.Parity = true
	}
	if flags&0b00010000 > 0 {
		p.AuxiliaryCarry = true
	}
	if flags&0b01000000 > 0 {
		p.Zero = true
	}
	if flags&0b10000000 > 0 {
		p.Sign = true
	}
}

// Exchange stack HL <-> mem[stack pointer]
func (p *Processor) xthl() {
	p.dasm("XTHL")
	tlsb := p.mmu.Memory[p.SP]
	tmsb := p.mmu.Memory[p.SP+1]
	p.mmu.Memory[p.SP] = p.L
	p.mmu.Memory[p.SP+1] = p.H
	p.L = tlsb
	p.H = tmsb
}

/* JUMP Instruction */
// internal jump
func (p *Processor) intJmp() {
	lsb := p.mmu.Memory[p.PC]
	msb := p.mmu.Memory[p.PC+1]

	p.PC = (uint16(msb) << 8) | uint16(lsb)
}

// load program counter
func (p *Processor) pchl() {
	p.dasm("PCHL")
	p.PC = (uint16(p.H) << 8) | uint16(p.L)
}

// jmp instruction
func (p *Processor) jmp() {
	p.dasm("JMP")
	p.intJmp()
}

// jump if not zero
func (p *Processor) jnz() {
	p.dasm("JNZ")
	if !p.Zero {
		p.intJmp()
	} else {
		p.PC += 2
	}
}

// jump if zero
func (p *Processor) jz() {
	p.dasm("JZ")
	if p.Zero {
		p.intJmp()
	} else {
		p.PC += 2
	}
}

// jump if not carry
func (p *Processor) jnc() {
	p.dasm("JNC")
	if !p.Carry {
		p.intJmp()
	} else {
		p.PC += 2
	}
}

// jump if carry
func (p *Processor) jc() {
	p.dasm("JC")
	if p.Carry {
		p.intJmp()
	} else {
		p.PC += 2
	}
}

// jump if parity odd (zero)
func (p *Processor) jpo() {
	p.dasm("JPO")
	if !p.Parity {
		p.intJmp()
	} else {
		p.PC += 2
	}
}

// jump if parity even (one)
func (p *Processor) jpe() {
	p.dasm("JPE")
	if p.Parity {
		p.intJmp()
	} else {
		p.PC += 2
	}
}

// jump if plus (sign-zero)
func (p *Processor) jp() {
	p.dasm("JP")
	if !p.Sign {
		p.intJmp()
	} else {
		p.PC += 2
	}
}

// jump if minus (sign-one)
func (p *Processor) jm() {
	p.dasm("JM")
	if p.Sign {
		p.intJmp()
	} else {
		p.PC += 2
	}
}

// in - read from specified input device to accumulator
// TODO - implementation
func (p *Processor) in() {
	p.dasm("IN")
	p.PC += 1
}

// out - send accumulator's content to the specified output device
// TODO - implementation
func (p *Processor) out() {
	p.dasm("OUT")
	p.PC += 1
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
