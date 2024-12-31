package main

import "log"

func (cpu *CPU) ParseNextOpcode() {
	next := cpu.ROM[cpu.PC]
	log.Printf("Opcode: 0x%02X", next)
	switch next {
	case 0x01: // LD BC, u16
		cpu.Registers[RegB] = cpu.ROM[cpu.PC+2]
		cpu.Registers[RegC] = cpu.ROM[cpu.PC+1]
		cpu.PC += 3
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
	case 0x07: // RLCA
		cpu.Registers[RegA] = (cpu.Registers[RegA] << 1) | (cpu.Registers[RegA] >> 7)
		cpu.Flags.SetC((cpu.Registers[RegA] & 0x01) != 0)
		cpu.Flags.SetZ(false)
		cpu.Flags.SetN(false)
		cpu.Flags.SetH(false)
		cpu.PC += 1
	case 0x08: // LD (u16), SP
		address := uint16(cpu.ROM[cpu.PC+1]) | (uint16(cpu.ROM[cpu.PC+2]) << 8)
		cpu.Memory[address] = uint8(cpu.SP & 0xFF) // Store low byte
		cpu.Memory[address+1] = uint8(cpu.SP >> 8) // Store high byte
		cpu.PC += 3
	case 0x09: // ADD HL, BC
		cpu.AddU16Registers(RegH, RegL, RegB, RegC)
		cpu.PC += 1
	case 0x0A: // LD A, (BC)
		cpu.Registers[RegA] = cpu.Memory[cpu.GetBC()]
		cpu.PC += 1
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
	case 0x11: // LD DE, u16
		cpu.Registers[RegD] = cpu.ROM[cpu.PC+2]
		cpu.Registers[RegE] = cpu.ROM[cpu.PC+1]
		cpu.PC += 3
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
	case 0x19: // ADD HL, DE
		cpu.AddU16Registers(RegH, RegL, RegD, RegE)
		cpu.PC += 1
	case 0x1A: // LD A, (DE)
		cpu.Registers[RegA] = cpu.Memory[cpu.GetDE()]
		cpu.PC += 1
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
	case 0x21: // LD HL, u16
		cpu.Registers[RegH] = cpu.ROM[cpu.PC+2]
		cpu.Registers[RegL] = cpu.ROM[cpu.PC+1]
		cpu.PC += 3
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
	case 0x29: // ADD HL, HL
		cpu.AddU16Registers(RegH, RegL, RegH, RegL)
		cpu.PC += 1
	case 0x2A: // LD A, (HL+)
		cpu.Registers[RegA] = cpu.Memory[cpu.GetHL()]
		cpu.IncrementU16Register(RegH, RegL)
		cpu.PC += 1
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
	case 0x31: // LD SP, u16
		value := uint16(cpu.ROM[cpu.PC+1]) | (uint16(cpu.ROM[cpu.PC+2]) << 8)
		cpu.SP = value
		cpu.PC += 3
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
	case 0x39: // ADD HL, SP
		cpu.AddU16Register(RegH, RegL, cpu.SP)
		cpu.PC += 1
	case 0x3A: // LD A, (HL-)
		cpu.Registers[RegA] = cpu.Memory[cpu.GetHL()]
		cpu.DecrementU16Register(RegH, RegL)
		cpu.PC += 1
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
	case 0x80: // ADD A, B
		cpu.AddU8Register(RegB)
		cpu.PC++
	case 0x81: // ADD A, C
		cpu.AddU8Register(RegC)
		cpu.PC++
	case 0x82: // ADD A, D
		cpu.AddU8Register(RegD)
		cpu.PC++
	case 0x83: // ADD A, E
		cpu.AddU8Register(RegE)
		cpu.PC++
	case 0x84: // ADD A, H
		cpu.AddU8Register(RegH)
		cpu.PC++
	case 0x85: // ADD A, L
		cpu.AddU8Register(RegL)
		cpu.PC++
	case 0x86: // ADD A, (HL)
		cpu.AddU8(cpu.Memory[cpu.GetHL()])
		cpu.PC++
	case 0x87: // ADD A, A
		cpu.AddU8Register(RegA)
		cpu.PC++
	case 0x88: // ADC A, B
		cpu.AdcU8Register(RegB)
		cpu.PC++
	case 0x89: // ADC A, C
		cpu.AdcU8Register(RegC)
		cpu.PC++
	case 0x8A: // ADC A, D
		cpu.AdcU8Register(RegD)
		cpu.PC++
	case 0x8B: // ADC A, E
		cpu.AdcU8Register(RegE)
		cpu.PC++
	case 0x8C: // ADC A, H
		cpu.AdcU8Register(RegH)
		cpu.PC++
	case 0x8D: // ADC A, L
		cpu.AdcU8Register(RegL)
		cpu.PC++
	case 0x8E: // ADC A, (HL)
		cpu.AdcU8(cpu.Memory[cpu.GetHL()])
		cpu.PC++
	case 0x8F: // ADC A, A
		cpu.AdcU8Register(RegA)
		cpu.PC++
	case 0x90: // SUB A, B
		cpu.SubU8Register(RegB)
		cpu.PC++
	case 0x91: // SUB A, C
		cpu.SubU8Register(RegC)
		cpu.PC++
	case 0x92: // SUB A, D
		cpu.SubU8Register(RegD)
		cpu.PC++
	case 0x93: // SUB A, E
		cpu.SubU8Register(RegE)
		cpu.PC++
	case 0x94: // SUB A, H
		cpu.SubU8Register(RegH)
		cpu.PC++
	case 0x95: // SUB A, L
		cpu.SubU8Register(RegL)
		cpu.PC++
	case 0x96: // SUB A, (HL)
		cpu.SubU8(cpu.Memory[cpu.GetHL()])
		cpu.PC++
	case 0x97: // SUB A, A
		cpu.SubU8Register(RegA)
		cpu.PC++
	case 0x98: // SBC A, B
		cpu.SbcU8Register(RegB)
		cpu.PC++
	case 0x99: // SBC A, C
		cpu.SbcU8Register(RegC)
		cpu.PC++
	case 0x9A: // SBC A, D
		cpu.SbcU8Register(RegD)
		cpu.PC++
	case 0x9B: // SBC A, E
		cpu.SbcU8Register(RegE)
		cpu.PC++
	case 0x9C: // SBC A, H
		cpu.SbcU8Register(RegH)
		cpu.PC++
	case 0x9D: // SBC A, L
		cpu.SbcU8Register(RegL)
		cpu.PC++
	case 0x9E: // SBC A, (HL)
		cpu.SbcU8(cpu.Memory[cpu.GetHL()])
		cpu.PC++
	case 0x9F: // SBC A, A
		cpu.SbcU8Register(RegA)
		cpu.PC++
	case 0xA0: // AND A, B
		cpu.AndU8Register(RegB)
		cpu.PC++
	case 0xA1: // AND A, C
		cpu.AndU8Register(RegC)
		cpu.PC++
	case 0xA2: // AND A, D
		cpu.AndU8Register(RegD)
		cpu.PC++
	case 0xA3: // AND A, E
		cpu.AndU8Register(RegE)
		cpu.PC++
	case 0xA4: // AND A, H
		cpu.AndU8Register(RegH)
		cpu.PC++
	case 0xA5: // AND A, L
		cpu.AndU8Register(RegL)
		cpu.PC++
	case 0xA6: // AND A, (HL)
		cpu.AndU8(cpu.Memory[cpu.GetHL()])
		cpu.PC++
	case 0xA7: // AND A, A
		cpu.AndU8Register(RegA)
		cpu.PC++
	case 0xA8: // XOR A, B
		cpu.XorU8Register(RegB)
		cpu.PC++
	case 0xA9: // XOR A, C
		cpu.XorU8Register(RegC)
		cpu.PC++
	case 0xAA: // XOR A, D
		cpu.XorU8Register(RegD)
		cpu.PC++
	case 0xAB: // XOR A, E
		cpu.XorU8Register(RegE)
		cpu.PC++
	case 0xAC: // XOR A, H
		cpu.XorU8Register(RegH)
		cpu.PC++
	case 0xAD: // XOR A, L
		cpu.XorU8Register(RegL)
		cpu.PC++
	case 0xAE: // XOR A, (HL)
		cpu.XorU8(cpu.Memory[cpu.GetHL()])
		cpu.PC++
	case 0xAF: // XOR A, A
		cpu.XorU8Register(RegA)
		cpu.PC++
	case 0xB0: // OR A, B
		cpu.OrU8Register(RegB)
		cpu.PC++
	case 0xB1: // OR A, C
		cpu.OrU8Register(RegC)
		cpu.PC++
	case 0xB2: // OR A, D
		cpu.OrU8Register(RegD)
		cpu.PC++
	case 0xB3: // OR A, E
		cpu.OrU8Register(RegE)
		cpu.PC++
	case 0xB4: // OR A, H
		cpu.OrU8Register(RegH)
		cpu.PC++
	case 0xB5: // OR A, L
		cpu.OrU8Register(RegL)
		cpu.PC++
	case 0xB6: // OR A, (HL)
		cpu.OrU8(cpu.Memory[cpu.GetHL()])
		cpu.PC++
	case 0xB7: // OR A, A
		cpu.OrU8Register(RegA)
		cpu.PC++
	case 0xB8: // CP A, B
		cpu.CpU8Register(RegB)
		cpu.PC++
	case 0xB9: // CP A, C
		cpu.CpU8Register(RegC)
		cpu.PC++
	case 0xBA: // CP A, D
		cpu.CpU8Register(RegD)
		cpu.PC++
	case 0xBB: // CP A, E
		cpu.CpU8Register(RegE)
		cpu.PC++
	case 0xBC: // CP A, H
		cpu.CpU8Register(RegH)
		cpu.PC++
	case 0xBD: // CP A, L
		cpu.CpU8Register(RegL)
		cpu.PC++
	case 0xBE: // CP A, (HL)
		cpu.CpU8(cpu.Memory[cpu.GetHL()])
		cpu.PC++
	case 0xBF: // CP A, A
		cpu.CpU8Register(RegA)
		cpu.PC++
	case 0xC6: // ADD A, u8
		cpu.AddU8(cpu.ROM[cpu.PC+1])
		cpu.PC += 2
	case 0xCE: // ADC A, u8
		cpu.AdcU8(cpu.ROM[cpu.PC+1])
		cpu.PC += 2
	case 0xD6: // SUB A, u8
		cpu.SubU8(cpu.ROM[cpu.PC+1])
		cpu.PC += 2
	case 0xE6: // AND A, u8
		cpu.AndU8(cpu.ROM[cpu.PC+1])
		cpu.PC += 2
	case 0xF6: // OR A, u8
		cpu.OrU8(cpu.ROM[cpu.PC+1])
		cpu.PC += 2
	case 0xDE: // SBC A, u8
		cpu.SbcU8(cpu.ROM[cpu.PC+1])
		cpu.PC += 2
	case 0xEE: // XOR A, u8
		cpu.XorU8(cpu.ROM[cpu.PC+1])
		cpu.PC += 2
	case 0xFE: // CP A, u8
		cpu.CpU8(cpu.ROM[cpu.PC+1])
		cpu.PC += 2
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

func (cpu *CPU) LoadImmediateU16(regHigh uint8, regLow uint8, value uint16) {
	cpu.Registers[regHigh] = uint8(value >> 8)
	cpu.Registers[regLow] = uint8(value & 0xFF)
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

func (cpu *CPU) AddU16Register(high uint8, low uint8, value uint16) {
	currentValue := uint16(cpu.Registers[high])<<8 | uint16(cpu.Registers[low])
	result := currentValue + value

	// Set carry flag if result overflows 16 bits
	cpu.Flags.SetC(result < currentValue)

	// Half carry is set if bit 11 overflows
	cpu.Flags.SetH((currentValue&0x0FFF)+(value&0x0FFF) > 0x0FFF)

	cpu.Registers[high] = uint8(result >> 8)
	cpu.Registers[low] = uint8(result & 0xFF)

	cpu.Flags.SetN(false)
}

func (cpu *CPU) AddU16Registers(high1 uint8, low1 uint8, high2 uint8, low2 uint8) {
	value1 := uint16(cpu.Registers[high1])<<8 | uint16(cpu.Registers[low1])
	value2 := uint16(cpu.Registers[high2])<<8 | uint16(cpu.Registers[low2])

	result := value1 + value2

	// Set carry flag if result overflows 16 bits
	cpu.Flags.SetC(result < value1) // Changed from result > 0xFFFF

	// Half carry is set if bit 11 overflows
	cpu.Flags.SetH((value1&0x0FFF)+(value2&0x0FFF) > 0x0FFF)

	// Store result in destination registers
	cpu.Registers[high1] = uint8(result >> 8)
	cpu.Registers[low1] = uint8(result & 0xFF)

	// N flag is always reset
	cpu.Flags.SetN(false)
}

func (cpu *CPU) AddU8Register(reg uint8) {
	cpu.AddU8(cpu.Registers[reg])
}

func (cpu *CPU) AddU8(value uint8) {
	a := cpu.Registers[RegA]
	b := value
	result := uint16(a) + uint16(b) // Use uint16 to catch overflow

	// Half carry occurs when there's a carry from bit 3 to bit 4
	halfCarry := (a&0x0F)+(b&0x0F) > 0x0F

	finalResult := uint8(result & 0xFF)

	// Set flags
	cpu.Flags.SetZ(finalResult == 0)
	cpu.Flags.SetN(false)
	cpu.Flags.SetH(halfCarry)
	cpu.Flags.SetC(result > 0xFF)

	cpu.Registers[RegA] = finalResult
}

func (cpu *CPU) AdcU8(value uint8) {
	a := cpu.Registers[RegA]
	b := value
	carry := uint8(0)
	if cpu.Flags.C() {
		carry = 1
	}

	result := uint16(a) + uint16(b) + uint16(carry)

	// Half carry needs to account for carry flag too
	cpu.Flags.SetH((a&0x0F)+(b&0x0F)+carry > 0x0F)
	cpu.Flags.SetC(result > 0xFF)
	cpu.Flags.SetN(false)

	finalResult := uint8(result & 0xFF)
	cpu.Registers[RegA] = finalResult
	cpu.Flags.SetZ(finalResult == 0)

}

func (cpu *CPU) AdcU8Register(reg uint8) {
	cpu.AdcU8(cpu.Registers[reg])
}

func (cpu *CPU) SubU8Register(reg uint8) {
	cpu.SubU8(cpu.Registers[reg])
}

func (cpu *CPU) SbcU8Register(reg uint8) {
	cpu.SbcU8(cpu.Registers[reg])
}
func (cpu *CPU) SubU8(value uint8) {
	a := cpu.Registers[RegA]
	b := value

	result := a - b
	cpu.Flags.SetZ(result == 0)
	cpu.Flags.SetN(true) // N is always set for subtraction
	cpu.Flags.SetH((a & 0xF) < (b & 0xF))
	cpu.Flags.SetC(a < b)

	cpu.Registers[RegA] = result
}

func (cpu *CPU) SbcU8(value uint8) {
	a := cpu.Registers[RegA]
	b := value
	carry := uint8(0)
	if cpu.Flags.C() {
		carry = 1
	}

	// Calculate result including carry
	result := a - b - carry

	// Half carry occurs when borrowing from bit 4
	halfCarry := (a & 0x0F) < ((b & 0x0F) + carry)

	// Carry occurs if result would be negative
	carryOut := uint16(a) < uint16(b)+uint16(carry)

	cpu.Flags.SetZ(result == 0)
	cpu.Flags.SetN(true) // N is always set for subtraction
	cpu.Flags.SetH(halfCarry)
	cpu.Flags.SetC(carryOut)

	cpu.Registers[RegA] = result
}
func (cpu *CPU) CpU8Register(reg uint8) {
	cpu.CpU8(cpu.Registers[reg])
}

func (cpu *CPU) CpU8(value uint8) {
	a := cpu.Registers[RegA]
	b := value

	result := a - b
	cpu.Flags.SetZ(result == 0)
	cpu.Flags.SetN(true) // N is always set for subtraction
	cpu.Flags.SetH((a & 0xF) < (b & 0xF))
	cpu.Flags.SetC(a < b)
}
func (cpu *CPU) AndU8Register(reg uint8) {
	cpu.AndU8(cpu.Registers[reg])
}

func (cpu *CPU) AndU8(value uint8) {
	a := cpu.Registers[RegA]
	b := value

	result := a & b
	cpu.Registers[RegA] = result
	cpu.Flags.SetZ(result == 0)
	cpu.Flags.SetN(false) // N is always set for subtraction
	cpu.Flags.SetH(true)
	cpu.Flags.SetC(false)
}

func (cpu *CPU) XorU8Register(reg uint8) {
	cpu.XorU8(cpu.Registers[reg])
}

func (cpu *CPU) XorU8(value uint8) {
	a := cpu.Registers[RegA]
	b := value
	result := a ^ b

	cpu.Registers[RegA] = result
	cpu.Flags.SetZ(result == 0)
	cpu.Flags.SetN(false)
	cpu.Flags.SetH(false)
	cpu.Flags.SetC(false)
}

func (cpu *CPU) OrU8Register(reg uint8) {
	cpu.OrU8(cpu.Registers[reg])
}

func (cpu *CPU) OrU8(value uint8) {
	a := cpu.Registers[RegA]
	b := value
	result := a | b
	cpu.Registers[RegA] = result
	cpu.Flags.SetZ(result == 0)
	cpu.Flags.SetN(false)
	cpu.Flags.SetH(false)
	cpu.Flags.SetC(false)
}
