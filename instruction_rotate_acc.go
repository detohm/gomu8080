package gomu8080

/* Rotate Accumulator Instructions */

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
