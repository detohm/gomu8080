package gomu8080

/* Register Pair Instructions */

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

// Load SP from HL
func (p *Processor) sphl() {
	p.dasm("SPHL")
	p.SP = (uint16(p.H) << 8) | uint16(p.L)
}
