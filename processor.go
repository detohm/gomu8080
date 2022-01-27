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
	// unused flag (must as for cpu pre test check in psw push pop)
	FlagBit1 bool
	FlagBit3 bool
	FlagBit5 bool

	// MMU
	mmu *MMU

	// debug
	DebugMode bool

	// processor state
	IsHalt bool

	// enable interupt
	IsInteruptsEnabled bool

	// pre calculation for zsp flags
	ZSP [0x100]uint8
}

func NewProcessor(mmu *MMU, debugMode bool) *Processor {
	p := &Processor{}
	p.mmu = mmu
	p.DebugMode = debugMode
	initZSPTable(p)
	// p.FlagBit1 = true
	return p
}

func initZSPTable(p *Processor) {

	for i := 0; i <= 0xFF; i++ {
		if i == 0 {
			p.ZSP[i] |= 0x01
		}
		if i>>7&0x1 > 0 {
			p.ZSP[i] |= 0x02
		}

		oneCount := 0
		for b := 0; b < 8; b++ {

			if (i>>b)&0x01 == 0x01 {
				oneCount += 1
			}
		}
		if oneCount%2 == 0 {
			p.ZSP[i] |= 0x04
		}
	}
}

func (p *Processor) SetZSP(value uint8) {
	p.Zero = (p.ZSP[value]>>0)&0x01 > 0
	p.Sign = (p.ZSP[value]>>1)&0x01 > 0
	p.Parity = (p.ZSP[value]>>2)&0x01 > 0
}

func (p *Processor) SetFlagsAdd(op1 uint8, op2 uint8, carry uint8, mode uint8) {
	result := op1 + op2 + carry
	switch mode {
	case 0: // reset
		p.Carry = false
		p.AuxiliaryCarry = false
	case 1: // both carry and auxCarry
		p.Carry = p.GetCarry(op1, op2, carry, 8)
		p.AuxiliaryCarry = p.GetCarry(op1, op2, carry, 4)
	case 2: // only carry
		p.Carry = p.GetCarry(op1, op2, carry, 8)
	case 3: // only auxCarry
		p.AuxiliaryCarry = p.GetCarry(op1, op2, carry, 4)
	}
	p.SetZSP(result)
}

func (p *Processor) SetFlagsSub(op1 uint8, op2 uint8, carry uint8, mode uint8) {
	p.SetFlagsAdd(op1, uint8(0xFF)^op2, 1-carry, mode)

	switch mode {
	case 1:
		p.Carry = !p.Carry
	}
}

func (p *Processor) GetCarry(op1 uint8, op2 uint8, carry uint8, bit uint8) bool {
	result := uint16(op1) + uint16(op2) + uint16(carry)
	newCarry := result ^ uint16(op1) ^ uint16(op2)
	return (newCarry & uint16(0x1<<bit)) != 0
}

func (p *Processor) Run() {

	opcode := p.mmu.Memory[p.PC]

	// TODO - validate address bound
	address := (uint16(p.H) << 8) | uint16(p.L)

	p.PC += 1

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
		p.lxi(&p.D, &p.E)
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
		p.di()
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
		p.ei()
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

	if p.DebugMode {
		p.PrintStatus()
	}
}

// debug print status
func (p *Processor) PrintStatus() {
	fmt.Printf("(A=%02X,H=%02X%02X,B=%02X%02X,D=%02X%02X,SP=%04X,PC=%04X,FLAG=%08b)\n",
		p.A,
		p.H,
		p.L,
		p.B,
		p.C,
		p.D,
		p.E,
		p.SP,
		p.PC,
		p.getFlags())
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
func (p *Processor) SetAuxiliaryCarry(result uint16, op1 uint8, op2 uint8, isAdd bool) {

	p.AuxiliaryCarry = false
	carryBits := result ^ uint16(op1) ^ uint16(op2)
	p.AuxiliaryCarry = carryBits&(0x01<<4) > 0

	// TODO - confirm on subtraction arithmetic

}

// DEBUGGER
// disassemble helper
// TODO - add operands detail
func (p *Processor) dasm(opcode string) {

	if !p.DebugMode {
		return
	}
	fmt.Printf("%s ", opcode)

}
