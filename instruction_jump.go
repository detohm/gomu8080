package gomu8080

/* Jump Instructions */

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
