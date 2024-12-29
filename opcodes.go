package main

import "log"

func (cpu *CPU) ParseNextOpcode() {
	next := cpu.ROM[cpu.PC]
	log.Printf("Opcode: 0x%02X", next)
	switch next {
	case 0x02: // LD (BC), A
		cpu.LoadMemoryImmediate(cpu.GetBC(), cpu.Registers[RegA])
		cpu.PC += 1
	case 0x03: // INC BC
		cpu.IncrementU16Register(RegB, RegC)
		cpu.PC += 1
	case 0x04: // INC B
		cpu.Registers[RegB]++
		cpu.PC += 1
	case 0x05: // DEC B
		cpu.Registers[RegB]--
		cpu.PC += 1
	case 0x06: // LD B, u8
		cpu.LoadImmediate(RegB, cpu.ROM[cpu.PC+1])
		cpu.PC += 2
	case 0x0B: // DEC BC
		cpu.DecrementU16Register(RegB, RegC)
		cpu.PC += 1
	case 0x0C: // INC C
		cpu.Registers[RegC]++
		cpu.PC += 1
	case 0x0D: // DEC C
		cpu.Registers[RegC]--
		cpu.PC += 1
	case 0x0E: // LD C, u8
		cpu.LoadImmediate(RegC, cpu.ROM[cpu.PC+1])
		cpu.PC += 2
	case 0x12: // LD (DE), A
		cpu.LoadMemoryImmediate(cpu.GetDE(), cpu.Registers[RegA])
		cpu.PC += 1
	case 0x13: // INC DE
		cpu.IncrementU16Register(RegD, RegE)
		cpu.PC += 1
	case 0x14: // INC D
		cpu.Registers[RegD]++
		cpu.PC += 1
	case 0x15: // DEC D
		cpu.Registers[RegD]--
		cpu.PC += 1
	case 0x16: // LD D, u8
		cpu.LoadImmediate(RegD, cpu.ROM[cpu.PC+1])
		cpu.PC += 2
	case 0x1B: // DEC DE
		cpu.DecrementU16Register(RegD, RegE)
		cpu.PC += 1
	case 0x1C: // INC E
		cpu.Registers[RegE]++
		cpu.PC += 1
	case 0x1D: // DEC E
		cpu.Registers[RegE]--
		cpu.PC += 1
	case 0x1E: // LD E, u8
		cpu.LoadImmediate(RegE, cpu.ROM[cpu.PC+1])
		cpu.PC += 2
	case 0x22: //LD (HL+), A
		cpu.LoadMemoryImmediate(cpu.GetHL(), cpu.Registers[RegA])
		cpu.IncrementU16Register(RegH, RegL)
		cpu.PC += 1
	case 0x23: // INC HL
		cpu.IncrementU16Register(RegH, RegL)
		cpu.PC += 1
	case 0x24: // INC H
		cpu.Registers[RegH]++
		cpu.PC += 1
	case 0x25: // DEC H
		cpu.Registers[RegH]--
		cpu.PC += 1
	case 0x26: // LD H, u8
		cpu.LoadImmediate(RegH, cpu.ROM[cpu.PC+1])
		cpu.PC += 2
	case 0x2B: // DEC HL
		cpu.DecrementU16Register(RegH, RegL)
		cpu.PC += 1
	case 0x2C: // INC L
		cpu.Registers[RegL]++
		cpu.PC += 1
	case 0x2D: // DEC L
		cpu.Registers[RegL]--
		cpu.PC += 1
	case 0x2E: // LD L, u8
		cpu.LoadImmediate(RegL, cpu.ROM[cpu.PC+1])
		cpu.PC += 2
	case 0x32: //LD (HL-), A
		cpu.LoadMemoryImmediate(cpu.GetHL(), cpu.Registers[RegA])
		cpu.DecrementU16Register(RegH, RegL)
		cpu.PC += 1
	case 0x33: // INC SP
		cpu.SP += 1
		cpu.PC += 1
	case 0x34: // INC (HL)
		value := cpu.Memory[cpu.GetHL()]
		cpu.Memory[cpu.GetHL()] = value + 1
		cpu.PC += 1
	case 0x35: //DEC (HL)
		value := cpu.Memory[cpu.GetHL()]
		cpu.Memory[cpu.GetHL()] = value - 1
		cpu.PC += 1
	case 0x36: // LD (HL),u8
		cpu.LoadMemoryImmediate(cpu.GetHL(), cpu.ROM[cpu.PC+1])
		cpu.PC += 2
	case 0x3B: // DEC SP
		cpu.SP -= 1
		cpu.PC += 1
	case 0x3C: // INC A
		cpu.Registers[RegA]++
		cpu.PC += 1
	case 0x3D: // DEC A
		cpu.Registers[RegA]--
		cpu.PC += 1
	case 0x3E: // LD A, u8
		cpu.LoadImmediate(RegA, cpu.ROM[cpu.PC+1])
		cpu.PC += 2
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
	case 0x70: // LD (HL), B
		cpu.LoadMemory(cpu.GetHL(), RegB)
		cpu.PC++
	case 0x71: // LD (HL), C
		cpu.LoadMemory(cpu.GetHL(), RegC)
		cpu.PC++
	case 0x72: // LD (HL), D
		cpu.LoadMemory(cpu.GetHL(), RegD)
		cpu.PC++
	case 0x73: // LD (HL), E
		cpu.LoadMemory(cpu.GetHL(), RegE)
		cpu.PC++
	case 0x74: // LD (HL), H
		cpu.LoadMemory(cpu.GetHL(), RegH)
		cpu.PC++
	case 0x75: // LD (HL), L
		cpu.LoadMemory(cpu.GetHL(), RegL)
		cpu.PC++
	case 0x76: // HALT
		cpu.Halt()
		cpu.PC++
	case 0x77: // LD (HL), A
		cpu.LoadMemory(cpu.GetHL(), RegA)
		cpu.PC++
	case 0x78: // LD A, B
		cpu.LoadRegister(RegA, RegB)
		cpu.PC++
	case 0x79: // LD A, C
		cpu.LoadRegister(RegA, RegC)
		cpu.PC++
	case 0x7A: // LD A, D
		cpu.LoadRegister(RegA, RegD)
		cpu.PC++
	case 0x7B: // LD A, E
		cpu.LoadRegister(RegA, RegE)
		cpu.PC++
	case 0x7C: // LD A, H
		cpu.LoadRegister(RegA, RegH)
		cpu.PC++
	case 0x7D: // LD A, L
		cpu.LoadRegister(RegA, RegL)
		cpu.PC++
	case 0x7E: // LD A, (HL)
		cpu.LoadFromMemory(RegA, cpu.GetHL())
		cpu.PC++
	case 0x7F: // LD A, A
		cpu.LoadRegister(RegA, RegA)
		cpu.PC++
	}

}

func (cpu *CPU) GetHL() uint16 {
	hl := uint16(cpu.Registers[RegH])<<8 | uint16(cpu.Registers[RegL])
	return hl
}

func (cpu *CPU) GetBC() uint16 {
	bc := uint16(cpu.Registers[RegB])<<8 | uint16(cpu.Registers[RegC])
	return bc
}
func (cpu *CPU) GetDE() uint16 {
	de := uint16(cpu.Registers[RegD])<<8 | uint16(cpu.Registers[RegE])
	return de
}

func (cpu *CPU) LoadMemoryImmediate(address uint16, value uint8) {
	log.Printf("load mem immediate: 0x%02X -> 0x%04X", address, value)
	cpu.Memory[address] = value
}

func (cpu *CPU) LoadMemory(address uint16, reg uint8) {
	cpu.LoadMemoryImmediate(address, cpu.Registers[reg])
}

func (cpu *CPU) LoadFromMemory(reg uint8, address uint16) {
	cpu.Registers[reg] = cpu.Memory[address]
}

func (cpu *CPU) LoadRegister(dest uint8, source uint8) {
	cpu.Registers[dest] = cpu.Registers[source]
}

func (cpu *CPU) LoadImmediate(reg uint8, value uint8) {
	cpu.Registers[reg] = value
}

func (cpu *CPU) Halt() {
	cpu.Halted = true
}

func (cpu *CPU) IncrementU16Register(high uint8, low uint8) {
	value := uint16(cpu.Registers[high])<<8 | uint16(cpu.Registers[low])
	value++
	cpu.Registers[high] = uint8(value >> 8)
	cpu.Registers[low] = uint8(value & 0xFF)
}

func (cpu *CPU) DecrementU16Register(high uint8, low uint8) {
	value := uint16(cpu.Registers[high])<<8 | uint16(cpu.Registers[low])
	value--
	cpu.Registers[high] = uint8(value >> 8)
	cpu.Registers[low] = uint8(value & 0xFF)
}
