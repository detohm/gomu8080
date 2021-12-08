package gomu8080

import "fmt"

// Instruction
func (p *Processor) nop() {
	p.dasm("NOP")
}

// Increase value of 8-bit register by 1
func (p *Processor) inr(reg *byte) {
	p.dasm("INR")

	// result16 := uint16(*reg) + 1
	// result8 := byte(result16 & 0x00FF)

	// var addHalfCarryTable = []bool{false, false, true, false, true, false, true, true}

	// index := (((*reg & 0x88) >> 1) | ((0x01 & 0x88) >> 2) | ((result8 & 0x88) >> 3)) & 0x7
	// p.AuxiliaryCarry = addHalfCarryTable[index]

	// *reg = byte(result16 & 0x00ff)

	// p.SetZero(*reg)
	// p.SetSign(*reg)
	// p.SetParity(*reg)
	p.SetFlagsAdd(*reg, 1, 0, 3) // update only ac
	*reg += 1

}

// Decrease value of 8-bit register by 1
func (p *Processor) dcr(reg *byte) {
	p.dasm("DCR")

	// result16 := uint16(*reg) + (uint16(0xFFFF) ^ uint16(0x01)) + 0x01
	// result8 := byte(result16 & 0xFF)

	// var subHalfCarryTable = []bool{true, false, false, false, true, true, true, false}
	// index := (((*reg & 0x88) >> 1) | ((0x01 & 0x88) >> 2) | ((result8 & 0x88) >> 3)) & 0x7
	// p.AuxiliaryCarry = subHalfCarryTable[index]

	// *reg = byte(result16 & 0xff)

	// p.SetZero(*reg)
	// p.SetSign(*reg)
	// p.SetParity(*reg)
	p.SetFlagsSub(*reg, 1, 0, 3) // update ac only
	*reg -= 1

}

// Move Instruction - move data from src to dst
func (p *Processor) mov(dst *byte, src *byte) {
	p.dasm("MOV")
	*dst = *src
}

// Load Accumulator - load data from the provided address
func (p *Processor) ldax(msb *byte, lsb *byte) {
	p.dasm("LDAX")
	address := uint16(*msb)<<8 | uint16(*lsb)
	p.A = p.mmu.Memory[address]
}

// Load accumulator direct from the operand address to A
func (p *Processor) lda() {
	address := uint16(p.mmu.Memory[p.PC+1]) << 8
	address |= uint16(p.mmu.Memory[p.PC])

	p.dasm(fmt.Sprintf("LDA %04X", address))

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
	// lsb4 := p.A & 0x0F

	// if lsb4 > 0x09 || p.AuxiliaryCarry {
	// 	p.A += 0x06
	// 	p.AuxiliaryCarry = true
	// }

	// msb4 := p.A >> 4
	// if msb4 > 0x09 || p.Carry {
	// 	p.A += (0x06 << 4)
	// 	p.Carry = true
	// }

	carry := p.Carry
	add := uint8(0)
	if (p.A&0x0F) > 0x09 || p.AuxiliaryCarry {
		add += 0x06
	}

	if (p.A>>4) > 0x09 || ((p.A>>4) >= 9 && p.A&0x0F > 9) || p.Carry {
		add += 0x60
		carry = true
	}

	op1 := p.A
	result := uint16(op1) + uint16(add)
	p.A = byte(result & 0x00FF)

	var addHalfCarryTable = []bool{false, false, true, false, true, false, true, true}
	index := (((op1 & 0x88) >> 1) | ((add & 0x88) >> 2) | ((p.A & 0x88) >> 3)) & 0x7
	p.AuxiliaryCarry = addHalfCarryTable[index]

	p.SetZero(p.A)
	p.SetSign(p.A)
	p.SetParity(p.A)

	p.Carry = carry

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

// restart
func (p *Processor) rst(pos uint16) {
	p.dasm("RST")
	msb := byte(p.PC >> 8)
	lsb := byte(p.PC & 0x00FF)
	p.mmu.Memory[p.SP-1] = msb
	p.mmu.Memory[p.SP-2] = lsb
	p.SP -= 2

	address := uint16(pos << 3)
	p.PC = address
}

// get flags
func (p *Processor) getFlags() byte {
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

	// add unused flag for psw push-pop test
	if p.FlagBit1 {
		flags |= 0b00000010
	}
	if p.FlagBit3 {
		flags |= 0b00001000
	}
	if p.FlagBit5 {
		flags |= 0b00100000
	}

	return flags
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

// Enable Interuption
func (p *Processor) ei() {
	p.dasm("EI")
	// TODO - implementation
	p.IsInteruptsEnabled = true
}

// Disable Interuption
func (p *Processor) di() {
	p.dasm("DI")
	// TODO - implementation
	p.IsInteruptsEnabled = false
}

// Unimplemented
func (p *Processor) unimplemented() {
	p.dasm(fmt.Sprintf("%02x:UNIMPLEMENTED", p.mmu.Memory[p.PC]))
}
