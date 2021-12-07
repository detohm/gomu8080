package gomu8080

import "fmt"

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
	// fmt.Println()
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
