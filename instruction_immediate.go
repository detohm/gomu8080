package gomu8080

import "fmt"

/* Immediate Instructions */

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

// Move immediate data
func (p *Processor) mvi(reg *byte) {
	p.dasm("MVI")
	*reg = p.mmu.Memory[p.PC]
	p.PC += 1
}

// Add Immediate to Accumulator
func (p *Processor) adi() {

	op1 := p.A
	op2 := p.mmu.Memory[p.PC]

	p.dasm(fmt.Sprintf("ADI %02X", op2))

	result := uint16(op1) + uint16(op2)
	p.A = byte(result & 0x00FF)

	p.SetSign(p.A)
	p.SetZero(p.A)
	p.SetAuxiliaryCarry(result, op1, op2, true)
	p.SetParity(p.A)
	p.Carry = result > 0xFF
	p.PC += 1
}

// Add Immediate to Accumulator with carry
func (p *Processor) aci() {
	op1 := p.A
	op2 := p.mmu.Memory[p.PC]

	p.dasm(fmt.Sprintf("ACI %02X", op2))

	result := uint16(op1) + uint16(op2)
	if p.Carry {
		result += 0x01
	}

	p.A = byte(result & 0x00FF)

	p.SetSign(p.A)
	p.SetZero(p.A)
	p.SetAuxiliaryCarry(result, op1, op2, true)
	p.SetParity(p.A)
	p.Carry = result > 0xFF
	p.PC += 1
}

// Subtract immediate from accumulator
func (p *Processor) sui() {

	op1 := p.A
	op2 := p.mmu.Memory[p.PC]

	p.dasm(fmt.Sprintf("SUI %02X", op2))

	result := uint16(op1) + uint16(^op2) + 0x1
	p.A = byte(result & 0x00FF)

	p.SetSign(p.A)
	p.SetZero(p.A)
	p.SetAuxiliaryCarry(result, op1, op2, false)
	p.SetParity(p.A)
	p.Carry = false
	if result <= 0x00FF {
		p.Carry = true
	}
	p.PC += 1
}

// Subtract immediate from accumulator with borrow
func (p *Processor) sbi() {

	op1 := p.A
	op2 := p.mmu.Memory[p.PC]

	p.dasm(fmt.Sprintf("SBI %02X", op2))

	result := uint16(op1) + uint16(^op2) + 0x1
	if p.Carry {
		result += ^uint16(0x01) + 0x1
	}

	p.A = byte(result & 0x00FF)
	p.SetSign(p.A)
	p.SetZero(p.A)
	p.SetAuxiliaryCarry(result, op1, op2, false)
	p.SetParity(p.A)
	p.Carry = false
	if result <= 0x00FF {
		p.Carry = true
	}
	p.PC += 1
}

// Logial AND immediate with accumulator
func (p *Processor) ani() {
	op1 := p.mmu.Memory[p.PC]
	p.dasm(fmt.Sprintf("ANI %02X", op1))
	p.A &= op1

	p.SetSign(p.A)
	p.SetZero(p.A)
	p.AuxiliaryCarry = false
	p.SetParity(p.A)
	p.Carry = false
	p.PC += 1
}

// Logical XOR immediate with accumulator
func (p *Processor) xri() {
	op1 := p.mmu.Memory[p.PC]
	p.dasm(fmt.Sprintf("XRI %02X", op1))
	p.A ^= op1

	p.SetSign(p.A)
	p.SetZero(p.A)
	p.AuxiliaryCarry = false
	p.SetParity(p.A)
	p.Carry = false
	p.PC += 1
}

// Logical OR immediate with accumulator
func (p *Processor) ori() {
	op1 := p.mmu.Memory[p.PC]
	p.dasm(fmt.Sprintf("ORI %02X", op1))
	p.A |= op1

	p.SetSign(p.A)
	p.SetZero(p.A)
	// TODO implement auxiliary carry
	p.SetParity(p.A)
	p.Carry = false
	p.PC += 1
}

// Compare immediate with accumulator
func (p *Processor) cpi() {
	op1 := p.mmu.Memory[p.PC]
	p.dasm(fmt.Sprintf("CPI %02X", op1))
	result := uint16(p.A) + uint16(^op1) + 0x01
	lsb := byte(result & 0x00FF)
	p.SetSign(lsb)
	p.SetZero(lsb)
	// TODO implement auxiliary carry
	p.SetParity(lsb)
	p.Carry = false
	if result <= 0xFF {
		p.Carry = true
	}
	p.PC += 1
}
