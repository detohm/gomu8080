package gomu8080

/* Register or Memory to Accumulator instructions */

// Add register or memory to accumulator
func (p *Processor) add(reg *byte) {
	p.dasm("ADD")
	op1 := p.A
	op2 := *reg

	// result := uint16(op1) + uint16(op2)
	// lsb := byte(result & 0x00FF)
	// p.A = lsb

	// var addHalfCarryTable = []bool{false, false, true, false, true, false, true, true}

	// index := (((op1 & 0x88) >> 1) | ((op2 & 0x88) >> 2) | ((lsb & 0x88) >> 3)) & 0x7
	// p.AuxiliaryCarry = addHalfCarryTable[index]

	// p.SetSign(lsb)
	// p.SetZero(lsb)
	// // p.SetAuxiliaryCarry(result, op1, op2, true)
	// p.SetParity(lsb)
	// p.Carry = result > 0xFF
	// panic("")
	p.SetFlagsAdd(op1, op2, 0, 1)
	p.A += op2
}

// Add register or memory to accumulator with carry
func (p *Processor) adc(reg *byte) {
	p.dasm("ADC")
	op1 := p.A
	op2 := *reg

	// result := uint16(op1) + uint16(op2)
	// if p.Carry {
	// 	result += 0x01
	// }
	// lsb := byte(result & 0x00FF)
	// p.A = lsb

	// var addHalfCarryTable = []bool{false, false, true, false, true, false, true, true}

	// index := (((op1 & 0x88) >> 1) | ((op2 & 0x88) >> 2) | ((lsb & 0x88) >> 3)) & 0x7
	// p.AuxiliaryCarry = addHalfCarryTable[index]

	// p.SetSign(lsb)
	// p.SetZero(lsb)
	// // p.SetAuxiliaryCarry(result, op1, op2, true)
	// p.SetParity(lsb)
	// p.Carry = result > 0xFF
	// panic("")
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
	// result := uint16(op1) + uint16(^op2) + 0x1
	// result := uint16(op1) - uint16(^op2)

	// p.A = byte(result & 0x00FF)

	// var subHalfCarryTable = []bool{true, false, false, false, true, true, true, false}
	// index := (((op1 & 0x88) >> 1) | ((op2 & 0x88) >> 2) | ((p.A & 0x88) >> 3)) & 0x7
	// p.AuxiliaryCarry = subHalfCarryTable[index]

	// p.SetSign(p.A)
	// p.SetZero(p.A)
	// // p.SetAuxiliaryCarry(result, op1, op2, false)
	// p.SetParity(p.A)
	// p.Carry = false
	// // if result <= 0x00FF {
	// // 	p.Carry = true
	// // }
	// if result&0x100 > 0 {
	// 	p.Carry = true
	// }
	// panic("")
	p.SetFlagsSub(op1, op2, 0, 1)
	p.A -= op2
}

// Subtract register or memory from accumulator with borrow
func (p *Processor) sbb(reg *byte) {
	p.dasm("SBB")
	op1 := p.A
	op2 := *reg
	// // result := uint16(op1) + uint16(^op2) + 0x1
	// // if p.Carry {
	// // 	result += ^uint16(0x01) + 0x1
	// // }
	// result := uint16(op1) - uint16(op2)
	// if p.Carry {
	// 	result -= 1
	// }

	// p.A = byte(result & 0x00FF)

	// var subHalfCarryTable = []bool{true, false, false, false, true, true, true, false}
	// index := (((op1 & 0x88) >> 1) | ((op2 & 0x88) >> 2) | ((p.A & 0x88) >> 3)) & 0x7
	// p.AuxiliaryCarry = subHalfCarryTable[index]

	// p.SetSign(p.A)
	// p.SetZero(p.A)
	// // p.SetAuxiliaryCarry(result, op1, op2, false)
	// p.SetParity(p.A)
	// p.Carry = false
	// // if result <= 0x00FF {
	// // 	p.Carry = true
	// // }
	// if result&0x100 > 0 {
	// 	p.Carry = true
	// }
	// panic("")
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
	// p.A &= *reg

	// p.SetSign(p.A)
	// p.SetZero(p.A)
	// p.AuxiliaryCarry = false
	// p.SetParity(p.A)
	// p.Carry = false
	// panic("")
	p.Carry = false
	p.AuxiliaryCarry = ((p.A | *reg) & 0x08) != 0 // TODO
	p.A &= *reg
	p.SetZSP(p.A)
}

// Logical XOR register or memory with accumulator
func (p *Processor) xra(reg *byte) {
	p.dasm("XRA")
	// p.A ^= *reg

	// p.SetSign(p.A)
	// p.SetZero(p.A)
	// p.AuxiliaryCarry = false
	// p.SetParity(p.A)
	// p.Carry = false
	// panic("")
	p.A ^= *reg
	p.SetFlagsAdd(p.A, 0, 0, 0)

}

// Logical OR register or memory with accumulator
func (p *Processor) ora(reg *byte) {
	p.dasm("ORA")
	// p.A |= *reg

	// p.SetSign(p.A)
	// p.SetZero(p.A)
	// // TODO implement auxiliary carry
	// p.SetParity(p.A)
	// p.Carry = false
	// panic("")
	p.A |= *reg
	p.SetFlagsAdd(p.A, 0, 0, 0)
}

// Compare register or memory with accumulator
func (p *Processor) cmp(reg *byte) {
	p.dasm("CMP")
	// result := uint16(p.A) + uint16(^*reg) + 0x01
	// lsb := byte(result & 0x00FF)
	// p.SetSign(lsb)
	// p.SetZero(lsb)
	// // TODO implement auxiliary carry
	// p.SetParity(lsb)
	// p.Carry = false
	// if result <= 0xFF {
	// 	p.Carry = true
	// }
	p.SetFlagsSub(p.A, *reg, 0, 1)
}
