package gomu8080

/* Register or Memory to Accumulator instructions */

// Add register or memory to accumulator
func (p *Processor) add(reg *byte) {
	p.dasm("ADD")
	op1 := p.A
	op2 := *reg
	p.SetFlagsAdd(op1, op2, 0, 1)
	p.A += op2
}

// Add register or memory to accumulator with carry
func (p *Processor) adc(reg *byte) {
	p.dasm("ADC")
	op1 := p.A
	op2 := *reg
	workValue := op2
	carry := uint8(0)
	if p.Carry {
		workValue += 1
		carry = 1
	}
	p.SetFlagsAdd(op1, op2, carry, 1)
	p.A += workValue
}

// Subtract register or memory from accumulator
func (p *Processor) sub(reg *byte) {
	p.dasm("SUB")
	op1 := p.A
	op2 := *reg
	p.SetFlagsSub(op1, op2, 0, 1)
	p.A -= op2
}

// Subtract register or memory from accumulator with borrow
func (p *Processor) sbb(reg *byte) {
	p.dasm("SBB")
	op1 := p.A
	op2 := *reg
	workValue := op2
	carry := uint8(0)
	if p.Carry {
		workValue += 1
		carry = 1
	}
	p.SetFlagsSub(op1, op2, carry, 1)
	p.A -= workValue
}

// Logical AND register or memory with accumulator
func (p *Processor) ana(reg *byte) {
	p.dasm("ANA")
	p.Carry = false
	p.AuxiliaryCarry = ((p.A | *reg) & 0x08) != 0 // TODO
	p.A &= *reg
	p.SetZSP(p.A)
}

// Logical XOR register or memory with accumulator
func (p *Processor) xra(reg *byte) {
	p.dasm("XRA")
	p.A ^= *reg
	p.SetFlagsAdd(p.A, 0, 0, 0)

}

// Logical OR register or memory with accumulator
func (p *Processor) ora(reg *byte) {
	p.dasm("ORA")
	p.A |= *reg
	p.SetFlagsAdd(p.A, 0, 0, 0)
}

// Compare register or memory with accumulator
func (p *Processor) cmp(reg *byte) {
	p.dasm("CMP")
	p.SetFlagsSub(p.A, *reg, 0, 1)
}
