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
	p.SetFlagsAdd(op1, op2, 0, 1) // update both c and ac
	p.A += op2
	p.PC += 1
}

// Add Immediate to Accumulator with carry
func (p *Processor) aci() {
	op1 := p.A
	op2 := p.mmu.Memory[p.PC]
	workValue := op2
	carry := uint8(0)
	if p.Carry {
		workValue += 1
		carry = 1
	}

	p.SetFlagsAdd(op1, op2, carry, 1) // set both c and ac
	p.A += workValue
	p.PC += 1
}

// Subtract immediate from accumulator
func (p *Processor) sui() {

	op1 := p.A
	op2 := p.mmu.Memory[p.PC]
	p.SetFlagsSub(op1, op2, 0, 1) // update both c and ac
	p.A -= op2
	p.PC += 1
}

// Subtract immediate from accumulator with borrow
func (p *Processor) sbi() {

	op1 := p.A
	op2 := p.mmu.Memory[p.PC]

	p.dasm(fmt.Sprintf("SBI %02X", op2))

	workValue := op2
	carry := uint8(0)
	if p.Carry {
		workValue += 1
		carry = 1
	}
	p.SetFlagsSub(op1, op2, carry, 1) // update both c and ac
	p.A -= workValue
	p.PC += 1
}

// Logial AND immediate with accumulator
func (p *Processor) ani() {
	op1 := p.mmu.Memory[p.PC]
	p.dasm(fmt.Sprintf("ANI %02X", op1))

	p.Carry = false
	p.AuxiliaryCarry = ((p.A | op1) & 0x08) != 0 // TODO
	p.A &= op1
	p.SetZSP(p.A)
	p.PC += 1
}

// Logical XOR immediate with accumulator
func (p *Processor) xri() {
	op1 := p.mmu.Memory[p.PC]
	p.dasm(fmt.Sprintf("XRI %02X", op1))

	p.A ^= op1
	p.SetFlagsAdd(p.A, 0, 0, 0) // reset c and ac
	p.PC += 1
}

// Logical OR immediate with accumulator
func (p *Processor) ori() {
	op1 := p.mmu.Memory[p.PC]
	p.dasm(fmt.Sprintf("ORI %02X", op1))

	p.A |= op1
	p.SetFlagsAdd(p.A, 0, 0, 0) // reset c and ac
	p.PC += 1
}

// Compare immediate with accumulator
func (p *Processor) cpi() {
	op1 := p.mmu.Memory[p.PC]
	p.dasm(fmt.Sprintf("CPI %02X", op1))

	p.SetFlagsSub(p.A, op1, 0, 1) // affects both c and ac
	p.PC += 1
}
