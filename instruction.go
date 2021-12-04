package gomu8080

import "fmt"

// Instruction
func (p *Processor) nop() {
	p.dasm("NOP")
}

// Increase value of 8-bit register by 1
func (p *Processor) inr(reg *byte) {
	p.dasm("INR")
	op1 := *reg
	result16 := uint16(*reg) + 1
	*reg = byte(result16 & 0x00ff)

	p.SetZero(*reg)
	p.SetSign(*reg)
	p.SetParity(*reg)
	p.SetAuxiliaryCarry(result16, op1, 1, true)
}

// Decrease value of 8-bit register by 1
func (p *Processor) dcr(reg *byte) {
	p.dasm("DCR")
	op1 := *reg
	result16 := uint16(*reg) - 1
	*reg = byte(result16 & 0xff)

	p.SetZero(*reg)
	p.SetSign(*reg)
	p.SetParity(*reg)
	p.SetAuxiliaryCarry(result16, op1, 0x01, false)
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
	lsb4 := p.A & 0x0F
	// msb4 := p.A >> 4

	if lsb4 > 0x09 || p.AuxiliaryCarry {
		p.A += 0x06
		p.AuxiliaryCarry = true
	}

	msb4 := p.A >> 4
	if msb4 > 0x09 || p.Carry {
		p.A += (0x06 << 4)
		p.Carry = true
	}

	p.SetZero(p.A)
	p.SetSign(p.A)
	p.SetParity(p.A)

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
	p.SP += 2

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
}

// Disable Interuption
func (p *Processor) di() {
	p.dasm("DI")
	// TODO - implementation
}

// Unimplemented
func (p *Processor) unimplemented() {
	p.dasm(fmt.Sprintf("%02x:UNIMPLEMENTED", p.mmu.Memory[p.PC]))
}
