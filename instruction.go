package gomu8080

import "fmt"

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
	op1 := p.A
	op2 := *reg

	result := uint16(op1) + uint16(op2)
	lsb := byte(result & 0x00FF)
	p.A = lsb

	p.SetSign(lsb)
	p.SetZero(lsb)
	p.SetAuxiliaryCarry(result, op1, op2, true)
	p.SetParity(lsb)
	p.Carry = result > 0xFF
}

// Add register or memory to accumulator with carry
func (p *Processor) adc(reg *byte) {
	p.dasm("ADC")
	op1 := p.A
	op2 := *reg

	result := uint16(op1) + uint16(op2)
	if p.Carry {
		result += 0x01
	}
	lsb := byte(result & 0x00FF)
	p.A = lsb

	p.SetSign(lsb)
	p.SetZero(lsb)
	p.SetAuxiliaryCarry(result, op1, op2, true)
	p.SetParity(lsb)
	p.Carry = result > 0xFF
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
	op1 := p.A
	op2 := *reg
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
}

// Subtract register or memory from accumulator with borrow
func (p *Processor) sbb(reg *byte) {
	p.dasm("SBB")
	op1 := p.A
	op2 := *reg
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
	p.Carry = false
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
	p.Carry = false
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

/* subroutine instruction */
// Internal Call subroutine
func (p *Processor) intCall() {

	address := (uint16(p.mmu.Memory[p.PC+1]) << 8)
	address |= uint16(p.mmu.Memory[p.PC])
	returnAddress := p.PC + 2

	// Inject emulated CP/M routines
	if address == 0x0000 {
		p.IsHalt = true
		return
	}
	if address == 0x0005 {
		if p.C == 0x02 {
			p.BdosConsoleOutput()
		}
		if p.C == 0x09 {
			p.BdosWriteStr()
		}
		p.PC += 2
		return
	}

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

// for debugging purpose
/*
  Emulate BDOS in CP/M for message output routine
  C_WRITESTR - Output string
  C=9, DE=address of string
*/
func (p *Processor) BdosWriteStr() {

	address := (uint16(p.D) << 8) | uint16(p.E)
	for p.mmu.Memory[address] != '$' {
		fmt.Printf("%c", p.mmu.Memory[address])
		address += 1
	}
	fmt.Println()
}

/*
  Emulate BDOS in CP/M for character output routine
  C_WRITE - Output character
  C=2, E=ascii character
*/
func (p *Processor) BdosConsoleOutput() {
	fmt.Printf("%c", p.E)
}

// Return from subroutine
func (p *Processor) intRet() {
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

// push PSW onto stack
func (p *Processor) pushPSW() {
	p.dasm("PUSH PSW")
	p.mmu.Memory[p.SP-1] = p.A
	p.mmu.Memory[p.SP-2] = p.getFlags()
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

	address := (uint16(msb) << 8) | uint16(lsb)

	// inject warm boot from CP/M
	if address == 0x0000 {
		p.IsHalt = true
		return
	}

	p.PC = address
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
