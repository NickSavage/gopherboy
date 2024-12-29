package main

import "log"

func (cpu *CPU) ParseNextOpcode() {
	next := cpu.ROM[cpu.PC]
	switch next {
	case 0x40: // LD B,B
		cpu.LoadRegister(RegB, RegB)
		cpu.PC++
	case 0x41: // LD B,C
		cpu.LoadRegister(RegB, RegC)
		cpu.PC++
	case 0x42: // LD B,D
		cpu.LoadRegister(RegB, RegD)
		cpu.PC++
	case 0x43: // LD B,E
		cpu.LoadRegister(RegB, RegE)
		cpu.PC++
	case 0x44: // LD B,H
		cpu.LoadRegister(RegB, RegH)
		cpu.PC++
	case 0x45: // LD B,L
		cpu.LoadRegister(RegB, RegL)
		cpu.PC++
	case 0x46: // LD B, (HL):
		cpu.LoadFromMemory(RegB, cpu.GetHL())
		cpu.PC++
	case 0x47: // LD B,A
		cpu.LoadRegister(RegB, RegA)
		cpu.PC++
	case 0x48: // LD C,B
		cpu.LoadRegister(RegC, RegB)
		cpu.PC++
	case 0x49: // LD C,C
		cpu.LoadRegister(RegC, RegC)
		cpu.PC++
	case 0x4A: // LD C,D
		cpu.LoadRegister(RegC, RegD)
		cpu.PC++
	case 0x4B: // LD C,E
		cpu.LoadRegister(RegC, RegE)
		cpu.PC++
	case 0x4C: // LD C,H
		cpu.LoadRegister(RegC, RegH)
		cpu.PC++
	case 0x4D: // LD C,L
		cpu.LoadRegister(RegC, RegL)
		cpu.PC++
	case 0x4E: // LD C, (HL):
		cpu.LoadFromMemory(RegC, cpu.GetHL())
		cpu.PC++
	case 0x4F: // LD C,A
		cpu.LoadRegister(RegC, RegA)
		cpu.PC++
	case 0x50: // LD D,B
		cpu.LoadRegister(RegD, RegB)
		cpu.PC++
	case 0x51: // LD D,C
		cpu.LoadRegister(RegD, RegC)
		cpu.PC++
	case 0x52: // LD D,D
		cpu.LoadRegister(RegD, RegD)
		cpu.PC++
	case 0x53: // LD D,E
		cpu.LoadRegister(RegD, RegE)
		cpu.PC++
	case 0x54: // LD D,H
		cpu.LoadRegister(RegD, RegH)
		cpu.PC++
	case 0x55: // LD D,L
		cpu.LoadRegister(RegD, RegL)
		cpu.PC++
	case 0x56: // LD D,(HL):
		cpu.LoadFromMemory(RegD, cpu.GetHL())
		cpu.PC++
	case 0x57: // LD D,A
		cpu.LoadRegister(RegD, RegA)
		cpu.PC++
	case 0x58: // LD E,B
		cpu.LoadRegister(RegE, RegB)
		cpu.PC++
	case 0x59: // LD E,C
		cpu.LoadRegister(RegE, RegC)
		cpu.PC++
	case 0x5A: // LD E,D
		cpu.LoadRegister(RegE, RegD)
		cpu.PC++
	case 0x5B: // LD E,E
		cpu.LoadRegister(RegE, RegE)
		cpu.PC++
	case 0x5C: // LD E,H
		cpu.LoadRegister(RegE, RegH)
		cpu.PC++
	case 0x5D: // LD E,L
		cpu.LoadRegister(RegE, RegL)
		cpu.PC++
	case 0x5E: // LD E, (HL):
		cpu.LoadFromMemory(RegE, cpu.GetHL())
		cpu.PC++
	case 0x5F: // LD E,A
		cpu.LoadRegister(RegE, RegA)
		cpu.PC++
	case 0x60: // LD H,B
		cpu.LoadRegister(RegH, RegB)
		cpu.PC++
	case 0x61: // LD H,C
		cpu.LoadRegister(RegH, RegC)
		cpu.PC++
	case 0x62: // LD H,D
		cpu.LoadRegister(RegH, RegD)
		cpu.PC++
	case 0x63: // LD H,E
		cpu.LoadRegister(RegH, RegE)
		cpu.PC++
	case 0x64: // LD H,H
		cpu.LoadRegister(RegH, RegH)
		cpu.PC++
	case 0x65: // LD H,L
		cpu.LoadRegister(RegH, RegL)
		cpu.PC++
	case 0x66: // LD H,(HL):
		cpu.LoadFromMemory(RegH, cpu.GetHL())
		cpu.PC++
	case 0x67: // LD H,A
		cpu.LoadRegister(RegH, RegA)
		cpu.PC++
	case 0x68: // LD L,B
		cpu.LoadRegister(RegL, RegB)
		cpu.PC++
	case 0x69: // LD L,C
		cpu.LoadRegister(RegL, RegC)
		cpu.PC++
	case 0x6A: // LD L,D
		cpu.LoadRegister(RegL, RegD)
		cpu.PC++
	case 0x6B: // LD L,E
		cpu.LoadRegister(RegL, RegE)
		cpu.PC++
	case 0x6C: // LD L,H
		cpu.LoadRegister(RegL, RegH)
		cpu.PC++
	case 0x6D: // LD L,L
		cpu.LoadRegister(RegL, RegL)
		cpu.PC++
	case 0x6E: // LD L,(HL)
		cpu.LoadFromMemory(RegL, cpu.GetHL())
		cpu.PC++
	case 0x6F: // LD L,A
		cpu.LoadRegister(RegL, RegA)
		cpu.PC++
	}

}

func (cpu *CPU) GetHL() uint16 {
	return uint16(cpu.Registers[RegH])<<8 | uint16(cpu.Registers[RegL])
}

func (cpu *CPU) LoadMemory(address uint16, reg uint8) {
	log.Printf("address %v reg %v", address, reg)
	cpu.Memory[address] = cpu.Registers[reg]
}

func (cpu *CPU) LoadFromMemory(reg uint8, address uint16) {
	cpu.Registers[reg] = cpu.Memory[address]
}

func (cpu *CPU) LoadRegister(dest uint8, source uint8) {
	cpu.Registers[dest] = cpu.Registers[source]
}
