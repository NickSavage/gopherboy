package main

import (
	"log"
)

func (cpu *CPU) Bit(bit uint8, value uint8) {
	cpu.Flags.SetZ(value&(1<<bit) == 0)
	cpu.Flags.SetN(false)
	cpu.Flags.SetH(true)
}

func RLC(value byte) (result uint8, flags uint8) {
	highBit := (value & 0x80) >> 7
	result = (value << 1) | highBit
	flags = 0
	if result == 0 {
		flags |= 0x80 // Set Z flag
	}
	if highBit == 1 {
		flags |= 0x10 // Set C flag
	}
	// N and H flags are always reset
	return result, flags
}

func RL(value uint8, carryFlag bool) (result uint8, flags uint8) {
	highBit := (value & 0x80) >> 7

	result = (value << 1)
	if carryFlag {
		result |= 0x01
	}

	flags = 0
	if result == 0 {
		flags |= 0x80 // Set Z flag
	}
	if highBit == 1 {
		flags |= 0x10 // Set C flag
	}
	return result, flags
}

func RR(value uint8, carryFlag bool) (result uint8, flags uint8) {
	lowBit := value & 0x01

	result = (value >> 1)
	if carryFlag {
		result |= 0x80
	}

	flags = 0
	if result == 0 {
		flags |= 0x80 // Set Z flag
	}
	if lowBit == 1 {
		flags |= 0x10 // Set C flag
	}
	return result, flags
}
func RRC(value byte) (result byte, flags byte) {
	lowBit := value & 0x01
	result = (value >> 1) | (lowBit << 7)
	flags = 0
	if result == 0 {
		flags |= 0x80 // Set Z flag
	}
	if lowBit == 1 {
		flags |= 0x10 // Set C flag
	}
	return result, flags
}

func SLA(value uint8) (result uint8, flags uint8) {
	highBit := (value & 0x80) >> 7

	result = value << 1 // bit 0 becomes 0 automatically

	flags = 0
	if result == 0 {
		flags |= 0x80 // Set Z flag
	}
	if highBit == 1 {
		flags |= 0x10 // Set C flag
	}
	return result, flags
}

func SRA(value uint8) (result uint8, flags uint8) {
	lowBit := value & 0x01

	result = (value >> 1) | (value & 0x80) // keep bit 7 the same

	flags = 0
	if result == 0 {
		flags |= 0x80 // Set Z flag
	}
	if lowBit == 1 {
		flags |= 0x10 // Set C flag
	}
	return result, flags
}

func SRL(value uint8) (result uint8, flags uint8) {
	lowBit := value & 0x01

	result = value >> 1 // bit 7 becomes 0 automatically

	flags = 0
	if result == 0 {
		flags |= 0x80 // Set Z flag
	}
	if lowBit == 1 {
		flags |= 0x10 // Set C flag
	}
	return result, flags
}

func Res(bit uint8, value uint8) uint8 {
	mask := ^(uint8(1) << bit)
	return value & mask
}

func Set(bit uint8, value uint8) uint8 {
	return value | (1 << bit)
}

func (cpu *CPU) ReadMemory(address uint16) uint8 {
	// if cpu.DMAActive {
	// 	log.Printf("DMA active, reading from 0x%04X", address)
	// 	if address < 0xFF80 || address > 0xFFFE {
	// 		return 0xFF // Return 0xFF for non-HRAM memory during DMA
	// 	}
	// 	return cpu.Memory[address]
	// }
	return cpu.Memory[address]
}

func (cpu *CPU) ParseNextCBOpcode() {
	next := cpu.ReadMemory(cpu.PC)

	// fmt.Printf("CB Opcode: 0x%02X PC: 0x%04X SP: 0x%04X A: 0x%02X B: 0x%02X C: 0x%02X D: 0x%02X E: 0x%02X H: 0x%02X L: 0x%02X Flags: Z:%t N:%t H:%t C:%t\n",
	// 	next, cpu.PC, cpu.SP,
	// 	cpu.Registers[RegA], cpu.Registers[RegB], cpu.Registers[RegC],
	// 	cpu.Registers[RegD], cpu.Registers[RegE], cpu.Registers[RegH], cpu.Registers[RegL],
	// 	cpu.Flags.Z(), cpu.Flags.N(), cpu.Flags.H(), cpu.Flags.C())
	switch next {
	case 0x00: // RLC B
		result, flags := RLC(cpu.Registers[RegB])
		cpu.Registers[RegB] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x01: // RLC C
		result, flags := RLC(cpu.Registers[RegC])
		cpu.Registers[RegC] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x02: // RLC D
		result, flags := RLC(cpu.Registers[RegD])
		cpu.Registers[RegD] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x03: // RLC E
		result, flags := RLC(cpu.Registers[RegE])
		cpu.Registers[RegE] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x04: // RLC H
		result, flags := RLC(cpu.Registers[RegH])
		cpu.Registers[RegH] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x05: // RLC L
		result, flags := RLC(cpu.Registers[RegL])
		cpu.Registers[RegL] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x06: // RLC (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		result, flags := RLC(value)
		cpu.Memory[cpu.GetHL()] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 16
	case 0x07: // RLC A
		result, flags := RLC(cpu.Registers[RegA])
		cpu.Registers[RegA] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x08: // RRC B
		result, flags := RRC(cpu.Registers[RegB])
		cpu.Registers[RegB] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x09: // RRC C
		result, flags := RRC(cpu.Registers[RegC])
		cpu.Registers[RegC] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x0A: // RRC D
		result, flags := RRC(cpu.Registers[RegD])
		cpu.Registers[RegD] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x0B: // RRC E
		result, flags := RRC(cpu.Registers[RegE])
		cpu.Registers[RegE] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x0C: // RRC H
		result, flags := RRC(cpu.Registers[RegH])
		cpu.Registers[RegH] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x0D: // RRC L
		result, flags := RRC(cpu.Registers[RegL])
		cpu.Registers[RegL] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x0E: // RRC (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		result, flags := RRC(value)
		cpu.Memory[cpu.GetHL()] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 16
	case 0x0F: // RRC A
		result, flags := RRC(cpu.Registers[RegA])
		cpu.Registers[RegA] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x10: // RL B
		result, flags := RL(cpu.Registers[RegB], cpu.Flags.C())
		cpu.Registers[RegB] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x11: // RL C
		result, flags := RL(cpu.Registers[RegC], cpu.Flags.C())
		cpu.Registers[RegC] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x12: // RL D
		result, flags := RL(cpu.Registers[RegD], cpu.Flags.C())
		cpu.Registers[RegD] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x13: // RL E
		result, flags := RL(cpu.Registers[RegE], cpu.Flags.C())
		cpu.Registers[RegE] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x14: // RL H
		result, flags := RL(cpu.Registers[RegH], cpu.Flags.C())
		cpu.Registers[RegH] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x15: // RL L
		result, flags := RL(cpu.Registers[RegL], cpu.Flags.C())
		cpu.Registers[RegL] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x16: // RL (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		result, flags := RL(value, cpu.Flags.C())
		cpu.Memory[cpu.GetHL()] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 16
	case 0x17: // RL A
		result, flags := RL(cpu.Registers[RegA], cpu.Flags.C())
		cpu.Registers[RegA] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x18: // RR B
		result, flags := RR(cpu.Registers[RegB], cpu.Flags.C())
		cpu.Registers[RegB] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x19: // RR C
		result, flags := RR(cpu.Registers[RegC], cpu.Flags.C())
		cpu.Registers[RegC] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x1A: // RR D
		result, flags := RR(cpu.Registers[RegD], cpu.Flags.C())
		cpu.Registers[RegD] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x1B: // RR E
		result, flags := RR(cpu.Registers[RegE], cpu.Flags.C())
		cpu.Registers[RegE] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x1C: // RR H
		result, flags := RR(cpu.Registers[RegH], cpu.Flags.C())
		cpu.Registers[RegH] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x1D: // RR L
		result, flags := RR(cpu.Registers[RegL], cpu.Flags.C())
		cpu.Registers[RegL] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x1E: // RR (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		result, flags := RR(value, cpu.Flags.C())
		cpu.Memory[cpu.GetHL()] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 16
	case 0x1F: // RR A
		result, flags := RR(cpu.Registers[RegA], cpu.Flags.C())
		cpu.Registers[RegA] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x20: // SLA B
		result, flags := SLA(cpu.Registers[RegB])
		cpu.Registers[RegB] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x21: // SLA C
		result, flags := SLA(cpu.Registers[RegC])
		cpu.Registers[RegC] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x22: // SLA D
		result, flags := SLA(cpu.Registers[RegD])
		cpu.Registers[RegD] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x23: // SLA E
		result, flags := SLA(cpu.Registers[RegE])
		cpu.Registers[RegE] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x24: // SLA H
		result, flags := SLA(cpu.Registers[RegH])
		cpu.Registers[RegH] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x25: // SLA L
		result, flags := SLA(cpu.Registers[RegL])
		cpu.Registers[RegL] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x26: // SLA (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		result, flags := SLA(value)
		cpu.Memory[cpu.GetHL()] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 16
	case 0x27: // SLA A
		result, flags := SLA(cpu.Registers[RegA])
		cpu.Registers[RegA] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x28: // SRA B
		result, flags := SRA(cpu.Registers[RegB])
		cpu.Registers[RegB] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x29: // SRA C
		result, flags := SRA(cpu.Registers[RegC])
		cpu.Registers[RegC] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x2A: // SRA D
		result, flags := SRA(cpu.Registers[RegD])
		cpu.Registers[RegD] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x2B: // SRA E
		result, flags := SRA(cpu.Registers[RegE])
		cpu.Registers[RegE] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x2C: // SRA H
		result, flags := SRA(cpu.Registers[RegH])
		cpu.Registers[RegH] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x2D: // SRA L
		result, flags := SRA(cpu.Registers[RegL])
		cpu.Registers[RegL] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x2E: // SRA (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		result, flags := SRA(value)
		cpu.Memory[cpu.GetHL()] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 16
	case 0x2F: // SRA A
		result, flags := SRA(cpu.Registers[RegA])
		cpu.Registers[RegA] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x30: // SWAP B
		cpu.Registers[RegB] = cpu.Swap(cpu.Registers[RegB])
		cpu.Clock += 8
	case 0x31: // SWAP C
		cpu.Registers[RegC] = cpu.Swap(cpu.Registers[RegC])
		cpu.Clock += 8
	case 0x32: // SWAP D
		cpu.Registers[RegD] = cpu.Swap(cpu.Registers[RegD])
		cpu.Clock += 8
	case 0x33: // SWAP E
		cpu.Registers[RegE] = cpu.Swap(cpu.Registers[RegE])
		cpu.Clock += 8
	case 0x34: // SWAP H
		cpu.Registers[RegH] = cpu.Swap(cpu.Registers[RegH])
		cpu.Clock += 8
	case 0x35: // SWAP L
		cpu.Registers[RegL] = cpu.Swap(cpu.Registers[RegL])
		cpu.Clock += 8
	case 0x36: // SWAP (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		value = cpu.Swap(value)
		cpu.Memory[cpu.GetHL()] = value
		cpu.Clock += 16
	case 0x37: // SWAP A
		cpu.Registers[RegA] = cpu.Swap(cpu.Registers[RegA])
		cpu.Clock += 8
	case 0x38: // SRL B
		result, flags := SRL(cpu.Registers[RegB])
		cpu.Registers[RegB] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x39: // SRL C
		result, flags := SRL(cpu.Registers[RegC])
		cpu.Registers[RegC] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x3A: // SRL D
		result, flags := SRL(cpu.Registers[RegD])
		cpu.Registers[RegD] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x3B: // SRL E
		result, flags := SRL(cpu.Registers[RegE])
		cpu.Registers[RegE] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x3C: // SRL H
		result, flags := SRL(cpu.Registers[RegH])
		cpu.Registers[RegH] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x3D: // SRL L
		result, flags := SRL(cpu.Registers[RegL])
		cpu.Registers[RegL] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x3E: // SRL (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		result, flags := SRL(value)
		cpu.Memory[cpu.GetHL()] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 16
	case 0x3F: // SRL A
		result, flags := SRL(cpu.Registers[RegA])
		cpu.Registers[RegA] = result
		cpu.Flags.SetValue(flags)
		cpu.Clock += 8
	case 0x40: // BIT 0, B
		cpu.Bit(0, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0x41: // BIT 0, C
		cpu.Bit(0, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0x42: // BIT 0, D
		cpu.Bit(0, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0x43: // BIT 0, E
		cpu.Bit(0, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0x44: // BIT 0, H
		cpu.Bit(0, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0x45: // BIT 0, L
		cpu.Bit(0, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0x46: // BIT 0, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		cpu.Bit(0, value)
		cpu.Clock += 16
	case 0x47: // BIT 0, A
		cpu.Bit(0, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0x48: // BIT 1, B
		cpu.Bit(1, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0x49: // BIT 1, C
		cpu.Bit(1, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0x4A: // BIT 1, D
		cpu.Bit(1, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0x4B: // BIT 1, E
		cpu.Bit(1, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0x4C: // BIT 1, H
		cpu.Bit(1, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0x4D: // BIT 1, L
		cpu.Bit(1, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0x4E: // BIT 1, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		cpu.Bit(1, value)
		cpu.Clock += 16
	case 0x4F: // BIT 1, A
		cpu.Bit(1, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0x50: // BIT 2, B
		cpu.Bit(2, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0x51: // BIT 2, C
		cpu.Bit(2, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0x52: // BIT 2, D
		cpu.Bit(2, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0x53: // BIT 2, E
		cpu.Bit(2, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0x54: // BIT 2, H
		cpu.Bit(2, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0x55: // BIT 2, L
		cpu.Bit(2, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0x56: // BIT 2, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		cpu.Bit(2, value)
		cpu.Clock += 16
	case 0x57: // BIT 2, A
		cpu.Bit(2, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0x58: // BIT 3, B
		cpu.Bit(3, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0x59: // BIT 3, C
		cpu.Bit(3, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0x5A: // BIT 3, D
		cpu.Bit(3, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0x5B: // BIT 3, E
		cpu.Bit(3, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0x5C: // BIT 3, H
		cpu.Bit(3, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0x5D: // BIT 3, L
		cpu.Bit(3, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0x5E: // BIT 3, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		cpu.Bit(3, value)
		cpu.Clock += 16
	case 0x5F: // BIT 3, A
		cpu.Bit(3, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0x60: // BIT 4, B
		cpu.Bit(4, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0x61: // BIT 4, C
		cpu.Bit(4, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0x62: // BIT 4, D
		cpu.Bit(4, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0x63: // BIT 4, E
		cpu.Bit(4, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0x64: // BIT 4, H
		cpu.Bit(4, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0x65: // BIT 4, L
		cpu.Bit(4, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0x66: // BIT 4, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		cpu.Bit(4, value)
		cpu.Clock += 16
	case 0x67: // BIT 4, A
		cpu.Bit(4, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0x68: // BIT 5, B
		cpu.Bit(5, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0x69: // BIT 5, C
		cpu.Bit(5, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0x6A: // BIT 5, D
		cpu.Bit(5, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0x6B: // BIT 5, E
		cpu.Bit(5, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0x6C: // BIT 5, H
		cpu.Bit(5, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0x6D: // BIT 5, L
		cpu.Bit(5, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0x6E: // BIT 5, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		cpu.Bit(5, value)
		cpu.Clock += 16
	case 0x6F: // BIT 5, A
		cpu.Bit(5, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0x70: // BIT 6, B
		cpu.Bit(6, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0x71: // BIT 6, C
		cpu.Bit(6, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0x72: // BIT 6, D
		cpu.Bit(6, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0x73: // BIT 6, E
		cpu.Bit(6, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0x74: // BIT 6, H
		cpu.Bit(6, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0x75: // BIT 6, L
		cpu.Bit(6, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0x76: // BIT 6, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		cpu.Bit(6, value)
		cpu.Clock += 16
	case 0x77: // BIT 6, A
		cpu.Bit(6, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0x78: // BIT 7, B
		cpu.Bit(7, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0x79: // BIT 7, C
		cpu.Bit(7, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0x7A: // BIT 7, D
		cpu.Bit(7, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0x7B: // BIT 7, E
		cpu.Bit(7, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0x7C: // BIT 7, H
		cpu.Bit(7, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0x7D: // BIT 7, L
		cpu.Bit(7, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0x7E: // BIT 7, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		cpu.Bit(7, value)
		cpu.Clock += 16
	case 0x7F: // BIT 7, A
		cpu.Bit(7, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0x80: // RES 0, B
		cpu.Registers[RegB] = Res(0, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0x81: // RES 0, C
		cpu.Registers[RegC] = Res(0, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0x82: // RES 0, D
		cpu.Registers[RegD] = Res(0, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0x83: // RES 0, E
		cpu.Registers[RegE] = Res(0, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0x84: // RES 0, H
		cpu.Registers[RegH] = Res(0, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0x85: // RES 0, L
		cpu.Registers[RegL] = Res(0, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0x86: // RES 0, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		value = Res(0, value)
		cpu.Memory[cpu.GetHL()] = value
		cpu.Clock += 16
	case 0x87: // RES 0, A
		cpu.Registers[RegA] = Res(0, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0x88: // RES 1, B
		cpu.Registers[RegB] = Res(1, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0x89: // RES 1, C
		cpu.Registers[RegC] = Res(1, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0x8A: // RES 1, D
		cpu.Registers[RegD] = Res(1, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0x8B: // RES 1, E
		cpu.Registers[RegE] = Res(1, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0x8C: // RES 1, H
		cpu.Registers[RegH] = Res(1, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0x8D: // RES 1, L
		cpu.Registers[RegL] = Res(1, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0x8E: // RES 1, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		value = Res(1, value)
		cpu.Memory[cpu.GetHL()] = value
		cpu.Clock += 16
	case 0x8F: // RES 1, A
		cpu.Registers[RegA] = Res(1, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0x90: // RES 2, B
		cpu.Registers[RegB] = Res(2, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0x91: // RES 2, C
		cpu.Registers[RegC] = Res(2, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0x92: // RES 2, D
		cpu.Registers[RegD] = Res(2, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0x93: // RES 2, E
		cpu.Registers[RegE] = Res(2, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0x94: // RES 2, H
		cpu.Registers[RegH] = Res(2, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0x95: // RES 2, L
		cpu.Registers[RegL] = Res(2, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0x96: // RES 2, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		value = Res(2, value)
		cpu.Memory[cpu.GetHL()] = value
		cpu.Clock += 16
	case 0x97: // RES 2, A
		cpu.Registers[RegA] = Res(2, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0x98: // RES 3, B
		cpu.Registers[RegB] = Res(3, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0x99: // RES 3, C
		cpu.Registers[RegC] = Res(3, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0x9A: // RES 3, D
		cpu.Registers[RegD] = Res(3, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0x9B: // RES 3, E
		cpu.Registers[RegE] = Res(3, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0x9C: // RES 3, H
		cpu.Registers[RegH] = Res(3, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0x9D: // RES 3, L
		cpu.Registers[RegL] = Res(3, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0x9E: // RES 3, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		value = Res(3, value)
		cpu.Memory[cpu.GetHL()] = value
		cpu.Clock += 16
	case 0x9F: // RES 3, A
		cpu.Registers[RegA] = Res(3, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0xA0: // RES 4, B
		cpu.Registers[RegB] = Res(4, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0xA1: // RES 4, C
		cpu.Registers[RegC] = Res(4, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0xA2: // RES 4, D
		cpu.Registers[RegD] = Res(4, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0xA3: // RES 4, E
		cpu.Registers[RegE] = Res(4, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0xA4: // RES 4, H
		cpu.Registers[RegH] = Res(4, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0xA5: // RES 4, L
		cpu.Registers[RegL] = Res(4, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0xA6: // RES 4, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		value = Res(4, value)
		cpu.Memory[cpu.GetHL()] = value
		cpu.Clock += 16
	case 0xA7: // RES 4, A
		cpu.Registers[RegA] = Res(4, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0xA8: // RES 5, B
		cpu.Registers[RegB] = Res(5, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0xA9: // RES 5, C
		cpu.Registers[RegC] = Res(5, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0xAA: // RES 5, D
		cpu.Registers[RegD] = Res(5, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0xAB: // RES 5, E
		cpu.Registers[RegE] = Res(5, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0xAC: // RES 5, H
		cpu.Registers[RegH] = Res(5, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0xAD: // RES 5, L
		cpu.Registers[RegL] = Res(5, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0xAE: // RES 5, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		value = Res(5, value)
		cpu.Memory[cpu.GetHL()] = value
		cpu.Clock += 16
	case 0xAF: // RES 5, A
		cpu.Registers[RegA] = Res(5, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0xB0: // RES 6, B
		cpu.Registers[RegB] = Res(6, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0xB1: // RES 6, C
		cpu.Registers[RegC] = Res(6, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0xB2: // RES 6, D
		cpu.Registers[RegD] = Res(6, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0xB3: // RES 6, E
		cpu.Registers[RegE] = Res(6, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0xB4: // RES 6, H
		cpu.Registers[RegH] = Res(6, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0xB5: // RES 6, L
		cpu.Registers[RegL] = Res(6, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0xB6: // RES 6, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		value = Res(6, value)
		cpu.Memory[cpu.GetHL()] = value
		cpu.Clock += 16
	case 0xB7: // RES 6, A
		cpu.Registers[RegA] = Res(6, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0xB8: // RES 7, B
		cpu.Registers[RegB] = Res(7, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0xB9: // RES 7, C
		cpu.Registers[RegC] = Res(7, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0xBA: // RES 7, D
		cpu.Registers[RegD] = Res(7, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0xBB: // RES 7, E
		cpu.Registers[RegE] = Res(7, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0xBC: // RES 7, H
		cpu.Registers[RegH] = Res(7, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0xBD: // RES 7, L
		cpu.Registers[RegL] = Res(7, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0xBE: // RES 7, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		value = Res(7, value)
		cpu.Memory[cpu.GetHL()] = value
		cpu.Clock += 16
	case 0xBF: // RES 7, A
		cpu.Registers[RegA] = Res(7, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0xC0: // SET 0, B
		cpu.Registers[RegB] = Set(0, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0xC1: // SET 0, C
		cpu.Registers[RegC] = Set(0, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0xC2: // SET 0, D
		cpu.Registers[RegD] = Set(0, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0xC3: // SET 0, E
		cpu.Registers[RegE] = Set(0, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0xC4: // SET 0, H
		cpu.Registers[RegH] = Set(0, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0xC5: // SET 0, L
		cpu.Registers[RegL] = Set(0, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0xC6: // SET 0, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		value = Set(0, value)
		cpu.Memory[cpu.GetHL()] = value
		cpu.Clock += 16
	case 0xC7: // SET 0, A
		cpu.Registers[RegA] = Set(0, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0xC8: // SET 1, B
		cpu.Registers[RegB] = Set(1, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0xC9: // SET 1, C
		cpu.Registers[RegC] = Set(1, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0xCA: // SET 1, D
		cpu.Registers[RegD] = Set(1, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0xCB: // SET 1, E
		cpu.Registers[RegE] = Set(1, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0xCC: // SET 1, H
		cpu.Registers[RegH] = Set(1, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0xCD: // SET 1, L
		cpu.Registers[RegL] = Set(1, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0xCE: // SET 1, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		value = Set(1, value)
		cpu.Memory[cpu.GetHL()] = value
		cpu.Clock += 16
	case 0xCF: // SET 1, A
		cpu.Registers[RegA] = Set(1, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0xD0: // SET 2, B
		cpu.Registers[RegB] = Set(2, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0xD1: // SET 2, C
		cpu.Registers[RegC] = Set(2, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0xD2: // SET 2, D
		cpu.Registers[RegD] = Set(2, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0xD3: // SET 2, E
		cpu.Registers[RegE] = Set(2, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0xD4: // SET 2, H
		cpu.Registers[RegH] = Set(2, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0xD5: // SET 2, L
		cpu.Registers[RegL] = Set(2, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0xD6: // SET 2, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		value = Set(2, value)
		cpu.Memory[cpu.GetHL()] = value
		cpu.Clock += 16
	case 0xD7: // SET 2, A
		cpu.Registers[RegA] = Set(2, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0xD8: // SET 3, B
		cpu.Registers[RegB] = Set(3, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0xD9: // SET 3, C
		cpu.Registers[RegC] = Set(3, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0xDA: // SET 3, D
		cpu.Registers[RegD] = Set(3, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0xDB: // SET 3, E
		cpu.Registers[RegE] = Set(3, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0xDC: // SET 3, H
		cpu.Registers[RegH] = Set(3, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0xDD: // SET 3, L
		cpu.Registers[RegL] = Set(3, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0xDE: // SET 3, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		value = Set(3, value)
		cpu.Memory[cpu.GetHL()] = value
		cpu.Clock += 16
	case 0xDF: // SET 3, A
		cpu.Registers[RegA] = Set(3, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0xE0: // SET 4, B
		cpu.Registers[RegB] = Set(4, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0xE1: // SET 4, C
		cpu.Registers[RegC] = Set(4, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0xE2: // SET 4, D
		cpu.Registers[RegD] = Set(4, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0xE3: // SET 4, E
		cpu.Registers[RegE] = Set(4, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0xE4: // SET 4, H
		cpu.Registers[RegH] = Set(4, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0xE5: // SET 4, L
		cpu.Registers[RegL] = Set(4, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0xE6: // SET 4, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		value = Set(4, value)
		cpu.Memory[cpu.GetHL()] = value
		cpu.Clock += 16
	case 0xE7: // SET 4, A
		cpu.Registers[RegA] = Set(4, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0xE8: // SET 5, B
		cpu.Registers[RegB] = Set(5, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0xE9: // SET 5, C
		cpu.Registers[RegC] = Set(5, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0xEA: // SET 5, D
		cpu.Registers[RegD] = Set(5, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0xEB: // SET 5, E
		cpu.Registers[RegE] = Set(5, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0xEC: // SET 5, H
		cpu.Registers[RegH] = Set(5, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0xED: // SET 5, L
		cpu.Registers[RegL] = Set(5, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0xEE: // SET 5, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		value = Set(5, value)
		cpu.Memory[cpu.GetHL()] = value
		cpu.Clock += 16
	case 0xEF: // SET 5, A
		cpu.Registers[RegA] = Set(5, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0xF0: // SET 6, B
		cpu.Registers[RegB] = Set(6, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0xF1: // SET 6, C
		cpu.Registers[RegC] = Set(6, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0xF2: // SET 6, D
		cpu.Registers[RegD] = Set(6, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0xF3: // SET 6, E
		cpu.Registers[RegE] = Set(6, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0xF4: // SET 6, H
		cpu.Registers[RegH] = Set(6, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0xF5: // SET 6, L
		cpu.Registers[RegL] = Set(6, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0xF6: // SET 6, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		value = Set(6, value)
		cpu.Memory[cpu.GetHL()] = value
		cpu.Clock += 16
	case 0xF7: // SET 6, A
		cpu.Registers[RegA] = Set(6, cpu.Registers[RegA])
		cpu.Clock += 8
	case 0xF8: // SET 7, B
		cpu.Registers[RegB] = Set(7, cpu.Registers[RegB])
		cpu.Clock += 8
	case 0xF9: // SET 7, C
		cpu.Registers[RegC] = Set(7, cpu.Registers[RegC])
		cpu.Clock += 8
	case 0xFA: // SET 7, D
		cpu.Registers[RegD] = Set(7, cpu.Registers[RegD])
		cpu.Clock += 8
	case 0xFB: // SET 7, E
		cpu.Registers[RegE] = Set(7, cpu.Registers[RegE])
		cpu.Clock += 8
	case 0xFC: // SET 7, H
		cpu.Registers[RegH] = Set(7, cpu.Registers[RegH])
		cpu.Clock += 8
	case 0xFD: // SET 7, L
		cpu.Registers[RegL] = Set(7, cpu.Registers[RegL])
		cpu.Clock += 8
	case 0xFE: // SET 7, (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		value = Set(7, value)
		cpu.Memory[cpu.GetHL()] = value
		cpu.Clock += 16
	case 0xFF: // SET 7, A
		cpu.Registers[RegA] = Set(7, cpu.Registers[RegA])
		cpu.Clock += 8
	default:
		log.Fatalf("Unknown CB opcode: 0x%02X", next)
	}
	cpu.PC += 1
}

func (cpu *CPU) ParseNextOpcode() {
	next := cpu.ReadMemory(cpu.PC)
	// fmt.Printf("Opcode: 0x%02X 0x%02X 0x%02X PC: 0x%04X SP: 0x%04X A: 0x%02X B: 0x%02X C: 0x%02X D: 0x%02X E: 0x%02X H: 0x%02X L: 0x%02X Flags: Z:%t N:%t H:%t C:%t\n",
	// 	next, cpu.Memory[cpu.PC+1], cpu.Memory[cpu.PC+2], cpu.PC, cpu.SP,
	// 	cpu.Registers[RegA], cpu.Registers[RegB], cpu.Registers[RegC],
	// 	cpu.Registers[RegD], cpu.Registers[RegE], cpu.Registers[RegH], cpu.Registers[RegL],
	// 	cpu.Flags.Z(), cpu.Flags.N(), cpu.Flags.H(), cpu.Flags.C())

	switch next {
	case 0x00: // NOP
		cpu.PC++
		cpu.Clock += 4
	case 0x01: // LD BC, u16
		cpu.Registers[RegB] = cpu.ReadMemory(cpu.PC + 2)
		cpu.Registers[RegC] = cpu.ReadMemory(cpu.PC + 1)
		cpu.PC += 3
		cpu.Clock += 12
	case 0x02: // LD (BC), A
		cpu.LoadMemoryImmediate(cpu.GetBC(), cpu.Registers[RegA])
		cpu.PC += 1
		cpu.Clock += 8
	case 0x03: // INC BC
		cpu.IncrementU16Register(RegB, RegC)
		cpu.PC += 1
		cpu.Clock += 8
	case 0x04: // INC B
		flags := uint8(0)
		cpu.Registers[RegB], flags = cpu.IncrementU8(cpu.Registers[RegB])
		cpu.Flags.SetValue(flags)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x05: // DEC B
		flags := uint8(0)
		cpu.Registers[RegB], flags = cpu.DecrementU8(cpu.Registers[RegB])
		cpu.Flags.SetValue(flags)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x06: // LD B, u8
		cpu.LoadImmediate(RegB, cpu.ReadMemory(cpu.PC+1))
		cpu.PC += 2
		cpu.Clock += 8
	case 0x07: // RLCA
		cpu.Registers[RegA] = (cpu.Registers[RegA] << 1) | (cpu.Registers[RegA] >> 7)
		cpu.Flags.SetC((cpu.Registers[RegA] & 0x01) != 0)
		cpu.Flags.SetZ(false)
		cpu.Flags.SetN(false)
		cpu.Flags.SetH(false)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x08: // LD (u16), SP
		address := uint16(cpu.ReadMemory(cpu.PC+1)) | (uint16(cpu.ReadMemory(cpu.PC+2)) << 8)
		cpu.Memory[address] = uint8(cpu.SP & 0xFF) // Store low byte
		cpu.Memory[address+1] = uint8(cpu.SP >> 8) // Store high byte
		cpu.PC += 3
		cpu.Clock += 20
	case 0x09: // ADD HL, BC
		cpu.AddU16Registers(RegH, RegL, RegB, RegC)
		cpu.PC += 1
		cpu.Clock += 8
	case 0x0A: // LD A, (BC)
		cpu.Registers[RegA] = cpu.ReadMemory(cpu.GetBC())
		cpu.PC += 1
		cpu.Clock += 8
	case 0x0B: // DEC BC
		cpu.DecrementU16Register(RegB, RegC)
		cpu.PC += 1
		cpu.Clock += 8
	case 0x0C: // INC C
		flags := uint8(0)
		cpu.Registers[RegC], flags = cpu.IncrementU8(cpu.Registers[RegC])
		cpu.Flags.SetValue(flags)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x0D: // DEC C
		flags := uint8(0)
		cpu.Registers[RegC], flags = cpu.DecrementU8(cpu.Registers[RegC])
		cpu.Flags.SetValue(flags)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x0E: // LD C, u8
		cpu.LoadImmediate(RegC, cpu.ReadMemory(cpu.PC+1))
		cpu.PC += 2
		cpu.Clock += 8
	case 0x0F: // RRCA
		oldBit0 := cpu.Registers[RegA] & 0x01
		cpu.Registers[RegA] = (cpu.Registers[RegA] >> 1) | (cpu.Registers[RegA] << 7)
		if oldBit0 != 0 {
			cpu.Flags.SetC(true)
		} else {
			cpu.Flags.SetC(false)
		}
		cpu.Flags.SetZ(false)
		cpu.Flags.SetN(false)
		cpu.Flags.SetH(false)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x10: // STOP
		cpu.Halt()
		cpu.PC += 1
		cpu.Clock += 4
	case 0x11: // LD DE, u16
		cpu.Registers[RegD] = cpu.ReadMemory(cpu.PC + 2)
		cpu.Registers[RegE] = cpu.ReadMemory(cpu.PC + 1)
		cpu.PC += 3
		cpu.Clock += 12
	case 0x12: // LD (DE), A
		cpu.LoadMemoryImmediate(cpu.GetDE(), cpu.Registers[RegA])
		cpu.PC += 1
		cpu.Clock += 8
	case 0x13: // INC DE
		cpu.IncrementU16Register(RegD, RegE)
		cpu.PC += 1
		cpu.Clock += 8
	case 0x14: // INC D
		flags := uint8(0)
		cpu.Registers[RegD], flags = cpu.IncrementU8(cpu.Registers[RegD])
		cpu.Flags.SetValue(flags)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x15: // DEC D
		flags := uint8(0)
		cpu.Registers[RegD], flags = cpu.DecrementU8(cpu.Registers[RegD])
		cpu.Flags.SetValue(flags)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x16: // LD D, u8
		cpu.LoadImmediate(RegD, cpu.ReadMemory(cpu.PC+1))
		cpu.PC += 2
		cpu.Clock += 8
	case 0x17: // RLA

		// Save original bit 7 before rotation
		oldBit7 := (cpu.Registers[RegA] & 0x80) != 0

		// Perform rotation
		carryBit := uint8(0)
		if cpu.Flags.C() {
			carryBit = 1
		}
		cpu.Registers[RegA] = (cpu.Registers[RegA] << 1) | carryBit

		// Use the saved bit 7 to set carry flag
		cpu.Flags.SetC(oldBit7)
		cpu.Flags.SetZ(false)
		cpu.Flags.SetN(false)
		cpu.Flags.SetH(false)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x18: // JR i8
		offset := int8(cpu.Memory[cpu.PC+1])
		cpu.PC += 2
		cpu.PC += uint16(offset)
		cpu.Clock += 12
	case 0x19: // ADD HL, DE
		cpu.AddU16Registers(RegH, RegL, RegD, RegE)
		cpu.PC += 1
		cpu.Clock += 8
	case 0x1A: // LD A, (DE)
		cpu.Registers[RegA] = cpu.ReadMemory(cpu.GetDE())
		cpu.PC += 1
		cpu.Clock += 8
	case 0x1B: // DEC DE
		cpu.DecrementU16Register(RegD, RegE)
		cpu.PC += 1
		cpu.Clock += 8
	case 0x1C: // INC E
		flags := uint8(0)
		cpu.Registers[RegE], flags = cpu.IncrementU8(cpu.Registers[RegE])
		cpu.Flags.SetValue(flags)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x1D: // DEC E
		flags := uint8(0)
		cpu.Registers[RegE], flags = cpu.DecrementU8(cpu.Registers[RegE])
		cpu.Flags.SetValue(flags)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x1E: // LD E, u8
		cpu.LoadImmediate(RegE, cpu.ReadMemory(cpu.PC+1))
		cpu.PC += 2
		cpu.Clock += 8
	case 0x1F: // RLA
		oldBit0 := cpu.Registers[RegA] & 0x01
		var oldCarry uint8 = 0
		if cpu.Flags.C() {
			oldCarry = 1
		}
		cpu.Registers[RegA] = (cpu.Registers[RegA] >> 1) | (oldCarry << 7)
		cpu.Flags.SetC(oldBit0 == 1)
		cpu.Flags.SetZ(false)
		cpu.Flags.SetN(false)
		cpu.Flags.SetH(false)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x20: // JR NZ, i8
		if !cpu.Flags.Z() {
			offset := int8(cpu.ReadMemory(cpu.PC + 1)) // Treat as signed byte
			cpu.PC += 2                                // Move past opcode and offset
			cpu.PC += uint16(offset)                   // Add the offset
		} else {
			cpu.PC += 2 // Skip the instruction without jumping
		}
		cpu.Clock += 8
	case 0x21: // LD HL, u16
		cpu.Registers[RegH] = cpu.ReadMemory(cpu.PC + 2)
		cpu.Registers[RegL] = cpu.ReadMemory(cpu.PC + 1)
		cpu.PC += 3
		cpu.Clock += 12
	case 0x22: //LD (HL+), A
		cpu.LoadMemoryImmediate(cpu.GetHL(), cpu.Registers[RegA])
		cpu.IncrementU16Register(RegH, RegL)
		cpu.PC += 1
		cpu.Clock += 8
	case 0x23: // INC HL
		cpu.IncrementU16Register(RegH, RegL)
		cpu.PC += 1
		cpu.Clock += 8
	case 0x24: // INC H
		flags := uint8(0)
		cpu.Registers[RegH], flags = cpu.IncrementU8(cpu.Registers[RegH])
		cpu.Flags.SetValue(flags)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x25: // DEC H
		flags := uint8(0)
		cpu.Registers[RegH], flags = cpu.DecrementU8(cpu.Registers[RegH])
		cpu.Flags.SetValue(flags)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x26: // LD H, u8
		cpu.LoadImmediate(RegH, cpu.ReadMemory(cpu.PC+1))
		cpu.PC += 2
		cpu.Clock += 8
	case 0x27: // DAA
		a := cpu.Registers[RegA]
		var adjust uint8 = 0

		if cpu.Flags.H() || (!cpu.Flags.N() && (a&0x0F) > 9) {
			adjust |= 0x06
		}

		if cpu.Flags.C() || (!cpu.Flags.N() && a > 0x99) {
			adjust |= 0x60
			cpu.Flags.SetC(true)
		} else {
			cpu.Flags.SetC(false)
		}

		if cpu.Flags.N() {
			cpu.Registers[RegA] -= adjust
		} else {
			cpu.Registers[RegA] += adjust
		}

		cpu.Flags.SetZ(cpu.Registers[RegA] == 0)
		cpu.Flags.SetH(false)
		cpu.PC++
		cpu.Clock += 4
	case 0x28: // JR Z, i8
		if cpu.Flags.Z() {
			offset := int8(cpu.ReadMemory(cpu.PC + 1)) // Treat as signed byte
			cpu.PC += 2                                // Move past opcode and offset
			cpu.PC += uint16(offset)                   // Add the offset
		} else {
			cpu.PC += 2 // Skip the instruction without jumping
		}
		cpu.Clock += 8
	case 0x29: // ADD HL, HL
		cpu.AddU16Registers(RegH, RegL, RegH, RegL)
		cpu.PC += 1
		cpu.Clock += 8
	case 0x2A: // LD A, (HL+)
		cpu.Registers[RegA] = cpu.ReadMemory(cpu.GetHL())
		cpu.IncrementU16Register(RegH, RegL)
		cpu.PC += 1
		cpu.Clock += 8
	case 0x2B: // DEC HL
		cpu.DecrementU16Register(RegH, RegL)
		cpu.PC += 1
		cpu.Clock += 8
	case 0x2C: // INC L
		flags := uint8(0)
		cpu.Registers[RegL], flags = cpu.IncrementU8(cpu.Registers[RegL])
		cpu.Flags.SetValue(flags)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x2D: // DEC L
		flags := uint8(0)
		cpu.Registers[RegL], flags = cpu.DecrementU8(cpu.Registers[RegL])
		cpu.Flags.SetValue(flags)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x2E: // LD L, u8
		cpu.LoadImmediate(RegL, cpu.ReadMemory(cpu.PC+1))
		cpu.PC += 2
		cpu.Clock += 8
	case 0x2F: // CPL
		cpu.Registers[RegA] = ^cpu.Registers[RegA]
		cpu.Flags.SetN(true)
		cpu.Flags.SetH(true)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x30: // JR NC, i8
		if !cpu.Flags.C() {
			offset := int8(cpu.ReadMemory(cpu.PC + 1)) // Treat as signed byte
			cpu.PC += 2                                // Move past opcode and offset
			cpu.PC += uint16(offset)                   // Add the offset
		} else {
			cpu.PC += 2 // Skip the instruction without jumping
		}
		cpu.Clock += 8
	case 0x31: // LD SP, u16
		value := uint16(cpu.ReadMemory(cpu.PC+1)) | (uint16(cpu.ReadMemory(cpu.PC+2)) << 8)
		cpu.SP = value
		cpu.PC += 3
		cpu.Clock += 12
	case 0x32: //LD (HL-), A
		cpu.LoadMemoryImmediate(cpu.GetHL(), cpu.Registers[RegA])
		cpu.DecrementU16Register(RegH, RegL)
		cpu.PC += 1
		cpu.Clock += 8
	case 0x33: // INC SP
		cpu.SP += 1
		cpu.PC += 1
		cpu.Clock += 4
	case 0x34: // INC (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		flags := uint8(0)
		cpu.Memory[cpu.GetHL()], flags = cpu.IncrementU8(value)
		cpu.Flags.SetValue(flags)
		cpu.PC += 1
		cpu.Clock += 12
	case 0x35: //DEC (HL)
		value := cpu.ReadMemory(cpu.GetHL())
		flags := uint8(0)
		cpu.Memory[cpu.GetHL()], flags = cpu.DecrementU8(value)
		cpu.Flags.SetValue(flags)
		cpu.PC += 1
		cpu.Clock += 12
		cpu.Flags.SetZ(cpu.Memory[cpu.GetHL()] == 0)
		cpu.Flags.SetN(true)
	case 0x36: // LD (HL),u8
		cpu.LoadMemoryImmediate(cpu.GetHL(), cpu.ReadMemory(cpu.PC+1))
		cpu.PC += 2
		cpu.Clock += 12
	case 0x37: // SCF
		cpu.Flags.SetH(false)
		cpu.Flags.SetC(true)
		cpu.Flags.SetN(false)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x38: // JR C, i8
		if cpu.Flags.C() {
			offset := int8(cpu.ReadMemory(cpu.PC + 1)) // Treat as signed byte
			cpu.PC += 2                                // Move past opcode and offset
			cpu.PC += uint16(offset)                   // Add the offset
		} else {
			cpu.PC += 2 // Skip the instruction without jumping
		}
		cpu.Clock += 8
	case 0x39: // ADD HL, SP
		cpu.AddU16Register(RegH, RegL, cpu.SP)
		cpu.PC += 1
		cpu.Clock += 8
	case 0x3A: // LD A, (HL-)
		cpu.Registers[RegA] = cpu.ReadMemory(cpu.GetHL())
		cpu.DecrementU16Register(RegH, RegL)
		cpu.PC += 1
		cpu.Clock += 8
	case 0x3B: // DEC SP
		cpu.SP -= 1
		cpu.PC += 1
		cpu.Clock += 4
	case 0x3C: // INC A
		flags := uint8(0)
		cpu.Registers[RegA], flags = cpu.IncrementU8(cpu.Registers[RegA])
		cpu.Flags.SetValue(flags)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x3D: // DEC A
		flags := uint8(0)
		cpu.Registers[RegA], flags = cpu.DecrementU8(cpu.Registers[RegA])
		cpu.Flags.SetValue(flags)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x3E: // LD A, u8
		cpu.LoadImmediate(RegA, cpu.ReadMemory(cpu.PC+1))
		cpu.PC += 2
		cpu.Clock += 8
	case 0x3F: // CCF
		cpu.Flags.SetH(false)
		cpu.Flags.SetC(!cpu.Flags.C())
		cpu.Flags.SetN(false)
		cpu.PC += 1
		cpu.Clock += 4
	case 0x40: // LD B,B
		cpu.LoadRegister(RegB, RegB)
		cpu.PC++
		cpu.Clock += 4
	case 0x41: // LD B,C
		cpu.LoadRegister(RegB, RegC)
		cpu.PC++
		cpu.Clock += 4
	case 0x42: // LD B,D
		cpu.LoadRegister(RegB, RegD)
		cpu.PC++
		cpu.Clock += 4
	case 0x43: // LD B,E
		cpu.LoadRegister(RegB, RegE)
		cpu.PC++
		cpu.Clock += 4
	case 0x44: // LD B,H
		cpu.LoadRegister(RegB, RegH)
		cpu.PC++
		cpu.Clock += 4
	case 0x45: // LD B,L
		cpu.LoadRegister(RegB, RegL)
		cpu.PC++
		cpu.Clock += 4
	case 0x46: // LD B, (HL):
		cpu.LoadFromMemory(RegB, cpu.GetHL())
		cpu.PC++
		cpu.Clock += 8
	case 0x47: // LD B,A
		cpu.LoadRegister(RegB, RegA)
		cpu.PC++
		cpu.Clock += 4
	case 0x48: // LD C,B
		cpu.LoadRegister(RegC, RegB)
		cpu.PC++
		cpu.Clock += 4
	case 0x49: // LD C,C
		cpu.LoadRegister(RegC, RegC)
		cpu.PC++
		cpu.Clock += 4
	case 0x4A: // LD C,D
		cpu.LoadRegister(RegC, RegD)
		cpu.PC++
		cpu.Clock += 4
	case 0x4B: // LD C,E
		cpu.LoadRegister(RegC, RegE)
		cpu.PC++
		cpu.Clock += 4
	case 0x4C: // LD C,H
		cpu.LoadRegister(RegC, RegH)
		cpu.PC++
		cpu.Clock += 4
	case 0x4D: // LD C,L
		cpu.LoadRegister(RegC, RegL)
		cpu.PC++
		cpu.Clock += 4
	case 0x4E: // LD C, (HL):
		cpu.LoadFromMemory(RegC, cpu.GetHL())
		cpu.PC++
		cpu.Clock += 8
	case 0x4F: // LD C,A
		cpu.LoadRegister(RegC, RegA)
		cpu.PC++
	case 0x50: // LD D,B
		cpu.LoadRegister(RegD, RegB)
		cpu.PC++
		cpu.Clock += 4
	case 0x51: // LD D,C
		cpu.LoadRegister(RegD, RegC)
		cpu.PC++
		cpu.Clock += 4
	case 0x52: // LD D,D
		cpu.LoadRegister(RegD, RegD)
		cpu.PC++
		cpu.Clock += 4
	case 0x53: // LD D,E
		cpu.LoadRegister(RegD, RegE)
		cpu.PC++
		cpu.Clock += 4
	case 0x54: // LD D,H
		cpu.LoadRegister(RegD, RegH)
		cpu.PC++
		cpu.Clock += 4
	case 0x55: // LD D,L
		cpu.LoadRegister(RegD, RegL)
		cpu.PC++
		cpu.Clock += 4
	case 0x56: // LD D,(HL):
		cpu.LoadFromMemory(RegD, cpu.GetHL())
		cpu.PC++
		cpu.Clock += 8
	case 0x57: // LD D,A
		cpu.LoadRegister(RegD, RegA)
		cpu.PC++
		cpu.Clock += 4
	case 0x58: // LD E,B
		cpu.LoadRegister(RegE, RegB)
		cpu.PC++
		cpu.Clock += 4
	case 0x59: // LD E,C
		cpu.LoadRegister(RegE, RegC)
		cpu.PC++
		cpu.Clock += 4
	case 0x5A: // LD E,D
		cpu.LoadRegister(RegE, RegD)
		cpu.PC++
	case 0x5B: // LD E,E
		cpu.LoadRegister(RegE, RegE)
		cpu.PC++
		cpu.Clock += 4
	case 0x5C: // LD E,H
		cpu.LoadRegister(RegE, RegH)
		cpu.PC++
		cpu.Clock += 4
	case 0x5D: // LD E,L
		cpu.LoadRegister(RegE, RegL)
		cpu.PC++
		cpu.Clock += 4
	case 0x5E: // LD E, (HL):
		cpu.LoadFromMemory(RegE, cpu.GetHL())
		cpu.PC++
		cpu.Clock += 8
	case 0x5F: // LD E,A
		cpu.LoadRegister(RegE, RegA)
		cpu.PC++
		cpu.Clock += 4
	case 0x60: // LD H,B
		cpu.LoadRegister(RegH, RegB)
		cpu.PC++
		cpu.Clock += 4
	case 0x61: // LD H,C
		cpu.LoadRegister(RegH, RegC)
		cpu.PC++
		cpu.Clock += 4
	case 0x62: // LD H,D
		cpu.LoadRegister(RegH, RegD)
		cpu.PC++
		cpu.Clock += 4
	case 0x63: // LD H,E
		cpu.LoadRegister(RegH, RegE)
		cpu.PC++
		cpu.Clock += 4
	case 0x64: // LD H,H
		cpu.LoadRegister(RegH, RegH)
		cpu.PC++
		cpu.Clock += 4
	case 0x65: // LD H,L
		cpu.LoadRegister(RegH, RegL)
		cpu.PC++
		cpu.Clock += 4
	case 0x66: // LD H,(HL):
		cpu.LoadFromMemory(RegH, cpu.GetHL())
		cpu.PC++
		cpu.Clock += 8
	case 0x67: // LD H,A
		cpu.LoadRegister(RegH, RegA)
		cpu.PC++
		cpu.Clock += 4
	case 0x68: // LD L,B
		cpu.LoadRegister(RegL, RegB)
		cpu.PC++
		cpu.Clock += 4
	case 0x69: // LD L,C
		cpu.LoadRegister(RegL, RegC)
		cpu.PC++
		cpu.Clock += 4
	case 0x6A: // LD L,D
		cpu.LoadRegister(RegL, RegD)
		cpu.PC++
		cpu.Clock += 4
	case 0x6B: // LD L,E
		cpu.LoadRegister(RegL, RegE)
		cpu.PC++
		cpu.Clock += 4
	case 0x6C: // LD L,H
		cpu.LoadRegister(RegL, RegH)
		cpu.PC++
		cpu.Clock += 4
	case 0x6D: // LD L,L
		cpu.LoadRegister(RegL, RegL)
		cpu.PC++
		cpu.Clock += 4
	case 0x6E: // LD L,(HL)
		cpu.LoadFromMemory(RegL, cpu.GetHL())
		cpu.PC++
		cpu.Clock += 8
	case 0x6F: // LD L,A
		cpu.LoadRegister(RegL, RegA)
		cpu.PC++
		cpu.Clock += 4
	case 0x70: // LD (HL), B
		cpu.LoadMemory(cpu.GetHL(), RegB)
		cpu.PC++
		cpu.Clock += 8
	case 0x71: // LD (HL), C
		cpu.LoadMemory(cpu.GetHL(), RegC)
		cpu.PC++
		cpu.Clock += 8
	case 0x72: // LD (HL), D
		cpu.LoadMemory(cpu.GetHL(), RegD)
		cpu.PC++
		cpu.Clock += 8
	case 0x73: // LD (HL), E
		cpu.LoadMemory(cpu.GetHL(), RegE)
		cpu.PC++
		cpu.Clock += 8
	case 0x74: // LD (HL), H
		cpu.LoadMemory(cpu.GetHL(), RegH)
		cpu.PC++
		cpu.Clock += 8
	case 0x75: // LD (HL), L
		cpu.LoadMemory(cpu.GetHL(), RegL)
		cpu.PC++
		cpu.Clock += 8
	case 0x76: // HALT
		cpu.Halt()
		cpu.PC++
		cpu.Clock += 4
	case 0x77: // LD (HL), A
		cpu.LoadMemory(cpu.GetHL(), RegA)
		cpu.PC++
		cpu.Clock += 8
	case 0x78: // LD A, B
		cpu.LoadRegister(RegA, RegB)
		cpu.PC++
		cpu.Clock += 4
	case 0x79: // LD A, C
		cpu.LoadRegister(RegA, RegC)
		cpu.PC++
		cpu.Clock += 4
	case 0x7A: // LD A, D
		cpu.LoadRegister(RegA, RegD)
		cpu.PC++
		cpu.Clock += 4
	case 0x7B: // LD A, E
		cpu.LoadRegister(RegA, RegE)
		cpu.PC++
		cpu.Clock += 4
	case 0x7C: // LD A, H
		cpu.LoadRegister(RegA, RegH)
		cpu.PC++
		cpu.Clock += 4
	case 0x7D: // LD A, L
		cpu.LoadRegister(RegA, RegL)
		cpu.PC++
		cpu.Clock += 4
	case 0x7E: // LD A, (HL)
		cpu.LoadFromMemory(RegA, cpu.GetHL())
		cpu.PC++
		cpu.Clock += 8
	case 0x7F: // LD A, A
		cpu.LoadRegister(RegA, RegA)
		cpu.PC++
		cpu.Clock += 4
	case 0x80: // ADD A, B
		cpu.AddU8Register(RegB)
		cpu.PC++
		cpu.Clock += 4
	case 0x81: // ADD A, C
		cpu.AddU8Register(RegC)
		cpu.PC++
		cpu.Clock += 4
	case 0x82: // ADD A, D
		cpu.AddU8Register(RegD)
		cpu.PC++
		cpu.Clock += 4
	case 0x83: // ADD A, E
		cpu.AddU8Register(RegE)
		cpu.PC++
		cpu.Clock += 4
	case 0x84: // ADD A, H
		cpu.AddU8Register(RegH)
		cpu.PC++
		cpu.Clock += 4
	case 0x85: // ADD A, L
		cpu.AddU8Register(RegL)
		cpu.PC++
		cpu.Clock += 4
	case 0x86: // ADD A, (HL)
		cpu.AddU8(cpu.ReadMemory(cpu.GetHL()))
		cpu.PC++
		cpu.Clock += 8
	case 0x87: // ADD A, A
		cpu.AddU8Register(RegA)
		cpu.PC++
		cpu.Clock += 4
	case 0x88: // ADC A, B
		cpu.AdcU8Register(RegB)
		cpu.PC++
		cpu.Clock += 4
	case 0x89: // ADC A, C
		cpu.AdcU8Register(RegC)
		cpu.PC++
		cpu.Clock += 4
	case 0x8A: // ADC A, D
		cpu.AdcU8Register(RegD)
		cpu.PC++
		cpu.Clock += 4
	case 0x8B: // ADC A, E
		cpu.AdcU8Register(RegE)
		cpu.PC++
		cpu.Clock += 4
	case 0x8C: // ADC A, H
		cpu.AdcU8Register(RegH)
		cpu.PC++
		cpu.Clock += 4
	case 0x8D: // ADC A, L
		cpu.AdcU8Register(RegL)
		cpu.PC++
		cpu.Clock += 4
	case 0x8E: // ADC A, (HL)
		cpu.AdcU8(cpu.ReadMemory(cpu.GetHL()))
		cpu.PC++
		cpu.Clock += 8
	case 0x8F: // ADC A, A
		cpu.AdcU8Register(RegA)
		cpu.PC++
		cpu.Clock += 4
	case 0x90: // SUB A, B
		cpu.SubU8Register(RegB)
		cpu.PC++
		cpu.Clock += 4
	case 0x91: // SUB A, C
		cpu.SubU8Register(RegC)
		cpu.PC++
		cpu.Clock += 4
	case 0x92: // SUB A, D
		cpu.SubU8Register(RegD)
		cpu.PC++
		cpu.Clock += 4
	case 0x93: // SUB A, E
		cpu.SubU8Register(RegE)
		cpu.PC++
		cpu.Clock += 4
	case 0x94: // SUB A, H
		cpu.SubU8Register(RegH)
		cpu.PC++
		cpu.Clock += 4
	case 0x95: // SUB A, L
		cpu.SubU8Register(RegL)
		cpu.PC++
		cpu.Clock += 4
	case 0x96: // SUB A, (HL)
		cpu.SubU8(cpu.ReadMemory(cpu.GetHL()))
		cpu.PC++
		cpu.Clock += 8
	case 0x97: // SUB A, A
		cpu.SubU8Register(RegA)
		cpu.PC++
		cpu.Clock += 4
	case 0x98: // SBC A, B
		cpu.SbcU8Register(RegB)
		cpu.PC++
		cpu.Clock += 4
	case 0x99: // SBC A, C
		cpu.SbcU8Register(RegC)
		cpu.PC++
		cpu.Clock += 4
	case 0x9A: // SBC A, D
		cpu.SbcU8Register(RegD)
		cpu.PC++
		cpu.Clock += 4
	case 0x9B: // SBC A, E
		cpu.SbcU8Register(RegE)
		cpu.PC++
		cpu.Clock += 4
	case 0x9C: // SBC A, H
		cpu.SbcU8Register(RegH)
		cpu.PC++
		cpu.Clock += 4
	case 0x9D: // SBC A, L
		cpu.SbcU8Register(RegL)
		cpu.PC++
		cpu.Clock += 4
	case 0x9E: // SBC A, (HL)
		cpu.SbcU8(cpu.ReadMemory(cpu.GetHL()))
		cpu.PC++
		cpu.Clock += 8
	case 0x9F: // SBC A, A
		cpu.SbcU8Register(RegA)
		cpu.PC++
		cpu.Clock += 4
	case 0xA0: // AND A, B
		cpu.AndU8Register(RegB)
		cpu.PC++
		cpu.Clock += 4
	case 0xA1: // AND A, C
		cpu.AndU8Register(RegC)
		cpu.PC++
		cpu.Clock += 4
	case 0xA2: // AND A, D
		cpu.AndU8Register(RegD)
		cpu.PC++
		cpu.Clock += 4
	case 0xA3: // AND A, E
		cpu.AndU8Register(RegE)
		cpu.PC++
		cpu.Clock += 4
	case 0xA4: // AND A, H
		cpu.AndU8Register(RegH)
		cpu.PC++
		cpu.Clock += 4
	case 0xA5: // AND A, L
		cpu.AndU8Register(RegL)
		cpu.PC++
		cpu.Clock += 4
	case 0xA6: // AND A, (HL)
		cpu.AndU8(cpu.ReadMemory(cpu.GetHL()))
		cpu.PC++
		cpu.Clock += 8
	case 0xA7: // AND A, A
		cpu.AndU8Register(RegA)
		cpu.PC++
		cpu.Clock += 4
	case 0xA8: // XOR A, B
		cpu.XorU8Register(RegB)
		cpu.PC++
		cpu.Clock += 4
	case 0xA9: // XOR A, C
		cpu.XorU8Register(RegC)
		cpu.PC++
		cpu.Clock += 4
	case 0xAA: // XOR A, D
		cpu.XorU8Register(RegD)
		cpu.PC++
		cpu.Clock += 4
	case 0xAB: // XOR A, E
		cpu.XorU8Register(RegE)
		cpu.PC++
		cpu.Clock += 4
	case 0xAC: // XOR A, H
		cpu.XorU8Register(RegH)
		cpu.PC++
		cpu.Clock += 4
	case 0xAD: // XOR A, L
		cpu.XorU8Register(RegL)
		cpu.PC++
		cpu.Clock += 4
	case 0xAE: // XOR A, (HL)
		cpu.XorU8(cpu.ReadMemory(cpu.GetHL()))
		cpu.PC++
		cpu.Clock += 8
	case 0xAF: // XOR A, A
		cpu.XorU8Register(RegA)
		cpu.PC++
		cpu.Clock += 4
	case 0xB0: // OR A, B
		cpu.OrU8Register(RegB)
		cpu.PC++
		cpu.Clock += 4
	case 0xB1: // OR A, C
		cpu.OrU8Register(RegC)
		cpu.PC++
		cpu.Clock += 4
	case 0xB2: // OR A, D
		cpu.OrU8Register(RegD)
		cpu.PC++
		cpu.Clock += 4
	case 0xB3: // OR A, E
		cpu.OrU8Register(RegE)
		cpu.PC++
		cpu.Clock += 4
	case 0xB4: // OR A, H
		cpu.OrU8Register(RegH)
		cpu.PC++
		cpu.Clock += 4
	case 0xB5: // OR A, L
		cpu.OrU8Register(RegL)
		cpu.PC++
		cpu.Clock += 4
	case 0xB6: // OR A, (HL)
		cpu.OrU8(cpu.ReadMemory(cpu.GetHL()))
		cpu.PC++
		cpu.Clock += 8
	case 0xB7: // OR A, A
		cpu.OrU8Register(RegA)
		cpu.PC++
		cpu.Clock += 4
	case 0xB8: // CP A, B
		cpu.CpU8Register(RegB)
		cpu.PC++
		cpu.Clock += 4
	case 0xB9: // CP A, C
		cpu.CpU8Register(RegC)
		cpu.PC++
		cpu.Clock += 4
	case 0xBA: // CP A, D
		cpu.CpU8Register(RegD)
		cpu.PC++
		cpu.Clock += 4
	case 0xBB: // CP A, E
		cpu.CpU8Register(RegE)
		cpu.PC++
		cpu.Clock += 4
	case 0xBC: // CP A, H
		cpu.CpU8Register(RegH)
		cpu.PC++
		cpu.Clock += 4
	case 0xBD: // CP A, L
		cpu.CpU8Register(RegL)
		cpu.PC++
		cpu.Clock += 4
	case 0xBE: // CP A, (HL)
		cpu.CpU8(cpu.ReadMemory(cpu.GetHL()))
		cpu.PC++
		cpu.Clock += 8
	case 0xBF: // CP A, A
		cpu.CpU8Register(RegA)
		cpu.PC++
		cpu.Clock += 4
	case 0xC0: // RET NZ
		if !cpu.Flags.Z() {
			low := cpu.ReadMemory(cpu.SP)
			cpu.SP++
			high := cpu.ReadMemory(cpu.SP)
			cpu.SP++
			cpu.PC = uint16(high)<<8 | uint16(low)
			cpu.Clock += 20
		} else {
			cpu.PC++
			cpu.Clock += 8
		}
	case 0xC1: // POP BC
		cpu.PopU16(RegB, RegC)
		cpu.PC++
		cpu.Clock += 12
	case 0xC2: // JP NZ, u16
		low := cpu.ReadMemory(cpu.PC + 1)
		high := cpu.ReadMemory(cpu.PC + 2)
		cpu.PC += 3
		if !cpu.Flags.Z() {
			cpu.PC = uint16(high)<<8 | uint16(low)
			cpu.Clock += 16
		} else {
			cpu.Clock += 12
		}
	case 0xC3: // JP u16
		low := cpu.ReadMemory(cpu.PC + 1)
		high := cpu.ReadMemory(cpu.PC + 2)
		cpu.PC = uint16(high)<<8 | uint16(low)
		cpu.Clock += 16
	case 0xC4: // CALL NZ, u16
		if !cpu.Flags.Z() {
			newPC := uint16(cpu.ReadMemory(cpu.PC+2))<<8 | uint16(cpu.ReadMemory(cpu.PC+1))
			cpu.PC += 3
			high := uint8(cpu.PC >> 8)
			low := uint8(cpu.PC & 0xFF)
			cpu.SP--
			cpu.Memory[cpu.SP] = high
			cpu.SP--
			cpu.Memory[cpu.SP] = low
			cpu.PC = newPC
			cpu.Clock += 24
		} else {
			cpu.PC += 3
			cpu.Clock += 12
		}
	case 0xC5: // PUSH BC
		cpu.PushU16(RegB, RegC)
		cpu.PC++
		cpu.Clock += 16
	case 0xC6: // ADD A, u8
		cpu.AddU8(cpu.ReadMemory(cpu.PC + 1))
		cpu.PC += 2
		cpu.Clock += 8
	case 0xC7: // RST 00H
		cpu.PC++
		high := uint8(cpu.PC >> 8)
		low := uint8(cpu.PC & 0xFF)
		cpu.SP--
		cpu.Memory[cpu.SP] = high
		cpu.SP--
		cpu.Memory[cpu.SP] = low
		cpu.PC = 0x0000
		cpu.Clock += 16
	case 0xC8: // RET Z
		if cpu.Flags.Z() {
			low := cpu.ReadMemory(cpu.SP)
			cpu.SP++
			high := cpu.ReadMemory(cpu.SP)
			cpu.SP++
			cpu.PC = uint16(high)<<8 | uint16(low)
			cpu.Clock += 20
		} else {
			cpu.PC++
			cpu.Clock += 8
		}
	case 0xC9: // RET
		low := cpu.ReadMemory(cpu.SP)
		cpu.SP++
		high := cpu.ReadMemory(cpu.SP)
		cpu.SP++
		cpu.PC = uint16(high)<<8 | uint16(low)
		cpu.Clock += 16
	case 0xCA: // JP Z, u16
		low := cpu.ReadMemory(cpu.PC + 1)
		high := cpu.ReadMemory(cpu.PC + 2)
		cpu.PC += 3
		if cpu.Flags.Z() {
			cpu.PC = uint16(high)<<8 | uint16(low)
			cpu.Clock += 16
		} else {
			cpu.Clock += 12
		}
	case 0xCB: // CB prefix
		cpu.PC++
		cpu.Clock += 4
		cpu.ParseNextCBOpcode()
	case 0xCC: // CALL Z, u16
		if cpu.Flags.Z() {
			newPC := uint16(cpu.ReadMemory(cpu.PC+2))<<8 | uint16(cpu.ReadMemory(cpu.PC+1))
			cpu.PC += 3
			high := uint8(cpu.PC >> 8)
			low := uint8(cpu.PC & 0xFF)
			cpu.SP--
			cpu.Memory[cpu.SP] = high
			cpu.SP--
			cpu.Memory[cpu.SP] = low
			cpu.PC = newPC
			cpu.Clock += 24
		} else {
			cpu.PC += 3
			cpu.Clock += 12
		}
	case 0xCD: // CALL u16
		newPC := uint16(cpu.ReadMemory(cpu.PC+2))<<8 | uint16(cpu.ReadMemory(cpu.PC+1))
		cpu.PC += 3
		high := uint8(cpu.PC >> 8)
		low := uint8(cpu.PC & 0xFF)
		cpu.SP--
		cpu.Memory[cpu.SP] = high
		cpu.SP--
		cpu.Memory[cpu.SP] = low
		cpu.PC = newPC
		cpu.Clock += 24
	case 0xCE: // ADC A, u8
		cpu.AdcU8(cpu.ReadMemory(cpu.PC + 1))
		cpu.PC += 2
		cpu.Clock += 8
	case 0xCF: // RST 08H
		cpu.PC++
		high := uint8(cpu.PC >> 8)
		low := uint8(cpu.PC & 0xFF)
		cpu.SP--
		cpu.Memory[cpu.SP] = high
		cpu.SP--
		cpu.Memory[cpu.SP] = low
		cpu.PC = 0x0008
		cpu.Clock += 16
	case 0xD0: // RET NC
		if !cpu.Flags.C() {
			low := cpu.ReadMemory(cpu.SP)
			cpu.SP++
			high := cpu.ReadMemory(cpu.SP)
			cpu.SP++
			cpu.PC = uint16(high)<<8 | uint16(low)
			cpu.Clock += 20
		} else {
			cpu.PC++
			cpu.Clock += 8
		}
	case 0xD1: // POP DE
		cpu.PopU16(RegD, RegE)
		cpu.PC++
		cpu.Clock += 12
	case 0xD2: // JP NC, u16
		low := cpu.ReadMemory(cpu.PC + 1)
		high := cpu.ReadMemory(cpu.PC + 2)
		cpu.PC += 3
		if !cpu.Flags.C() {
			cpu.PC = uint16(high)<<8 | uint16(low)
			cpu.Clock += 16
		} else {
			cpu.Clock += 12
		}
	case 0xD4: // CALL NC, u16
		if !cpu.Flags.C() {
			newPC := uint16(cpu.ReadMemory(cpu.PC+2))<<8 | uint16(cpu.ReadMemory(cpu.PC+1))
			cpu.PC += 3
			high := uint8(cpu.PC >> 8)
			low := uint8(cpu.PC & 0xFF)
			cpu.SP--
			cpu.Memory[cpu.SP] = high
			cpu.SP--
			cpu.Memory[cpu.SP] = low
			cpu.PC = newPC
			cpu.Clock += 24
		} else {
			cpu.PC += 3
			cpu.Clock += 12
		}
	case 0xD5: // PUSH DE
		cpu.PushU16(RegD, RegE)
		cpu.PC++
		cpu.Clock += 16
	case 0xD6: // SUB A, u8
		cpu.SubU8(cpu.ReadMemory(cpu.PC + 1))
		cpu.PC += 2
		cpu.Clock += 8
	case 0xD7: // RST 10H
		cpu.PC++
		high := uint8(cpu.PC >> 8)
		low := uint8(cpu.PC & 0xFF)
		cpu.SP--
		cpu.Memory[cpu.SP] = high
		cpu.SP--
		cpu.Memory[cpu.SP] = low
		cpu.PC = 0x0010
		cpu.Clock += 16
	case 0xD8: // RET C
		if cpu.Flags.C() {
			low := cpu.ReadMemory(cpu.SP)
			cpu.SP++
			high := cpu.ReadMemory(cpu.SP)
			cpu.SP++
			cpu.PC = uint16(high)<<8 | uint16(low)
			cpu.Clock += 20
		} else {
			cpu.PC++
			cpu.Clock += 8
		}
	case 0xD9: // RETI
		low := cpu.ReadMemory(cpu.SP)
		cpu.SP++
		high := cpu.ReadMemory(cpu.SP)
		cpu.SP++
		cpu.PC = uint16(high)<<8 | uint16(low)
		cpu.IME = 1
		cpu.Clock += 16
	case 0xDA: // JP C, u16
		low := cpu.ReadMemory(cpu.PC + 1)
		high := cpu.ReadMemory(cpu.PC + 2)
		cpu.PC += 3
		if cpu.Flags.C() {
			log.Printf("jumping to 0x%04X", uint16(high)<<8|uint16(low))
			cpu.PC = uint16(high)<<8 | uint16(low)
			cpu.Clock += 16
		} else {
			cpu.Clock += 12
		}
	case 0xDC: // CALL C, u16
		if cpu.Flags.C() {
			newPC := uint16(cpu.ReadMemory(cpu.PC+2))<<8 | uint16(cpu.ReadMemory(cpu.PC+1))
			cpu.PC += 3
			high := uint8(cpu.PC >> 8)
			low := uint8(cpu.PC & 0xFF)
			cpu.SP--
			cpu.Memory[cpu.SP] = high
			cpu.SP--
			cpu.Memory[cpu.SP] = low
			cpu.PC = newPC
			cpu.Clock += 24
		} else {
			cpu.PC += 3
			cpu.Clock += 12
		}
	case 0xDE: // SBC A, u8
		cpu.SbcU8(cpu.ReadMemory(cpu.PC + 1))
		cpu.PC += 2
		cpu.Clock += 8
	case 0xDF: // RST 18H
		cpu.PC++
		high := uint8(cpu.PC >> 8)
		low := uint8(cpu.PC & 0xFF)
		cpu.SP--
		cpu.Memory[cpu.SP] = high
		cpu.SP--
		cpu.Memory[cpu.SP] = low
		cpu.PC = 0x0018
		cpu.Clock += 16
	case 0xE0: // LD (0xFF00 + u8), A
		address := uint16(0xFF00 + uint16(cpu.ReadMemory(cpu.PC+1)))
		cpu.Memory[address] = cpu.Registers[RegA]
		cpu.PC += 2
		cpu.Clock += 12
		if address == 0xFF46 {
			log.Printf("DMA active")
			cpu.DMAActive = true
			cpu.DMASourceBase = uint16(cpu.Registers[RegA]) << 8
			cpu.DMACycles = 160
		}
		if address == 0xFF50 {
			copy(cpu.Memory[0x0000:0x0150], cpu.ROM[0x0000:0x0150])
		}
	case 0xE1: // POP HL
		cpu.PopU16(RegH, RegL)
		cpu.PC++
		cpu.Clock += 12
	case 0xE2: // LD (C), A
		cpu.Memory[0xFF00+uint16(cpu.Registers[RegC])] = cpu.Registers[RegA]
		cpu.PC += 1
		cpu.Clock += 8
	case 0xE5: // PUSH HL
		cpu.PushU16(RegH, RegL)
		cpu.PC++
		cpu.Clock += 16
	case 0xE6: // AND A, u8
		cpu.AndU8(cpu.ReadMemory(cpu.PC + 1))
		cpu.PC += 2
		cpu.Clock += 8
	case 0xE7: // RST 20H
		cpu.PC++
		high := uint8(cpu.PC >> 8)
		low := uint8(cpu.PC & 0xFF)
		cpu.SP--
		cpu.Memory[cpu.SP] = high
		cpu.SP--
		cpu.Memory[cpu.SP] = low
		cpu.PC = 0x0020
		cpu.Clock += 16
	case 0xE8: // ADD SP, u8
		offset := uint16(int8(cpu.ReadMemory(cpu.PC + 1)))

		// Check carry from bit 7
		if ((cpu.SP & 0xFF) + (offset & 0xFF)) > 0xFF {
			cpu.Flags.SetC(true)
		} else {
			cpu.Flags.SetC(false)
		}

		// Check half carry from bit 3
		if ((cpu.SP & 0xF) + (offset & 0xF)) > 0xF {
			cpu.Flags.SetH(true)
		} else {
			cpu.Flags.SetH(false)
		}
		cpu.SP += offset
		cpu.Flags.SetN(false)
		cpu.Flags.SetZ(false)
		cpu.PC += 2
		cpu.Clock += 16
	case 0xE9: // JP (HL)
		cpu.PC = cpu.GetHL()
		cpu.Clock += 4
	case 0xEA: // LD (u16), A
		address := uint16(cpu.ReadMemory(cpu.PC+2))<<8 | uint16(cpu.ReadMemory(cpu.PC+1))
		cpu.Memory[address] = cpu.Registers[RegA]
		cpu.PC += 3
		cpu.Clock += 16
	case 0xEE: // XOR A, u8
		cpu.XorU8(cpu.ReadMemory(cpu.PC + 1))
		cpu.PC += 2
		cpu.Clock += 8
	case 0xEF: // RST 28H
		cpu.PC++
		high := uint8(cpu.PC >> 8)
		low := uint8(cpu.PC & 0xFF)
		cpu.SP--
		cpu.Memory[cpu.SP] = high
		cpu.SP--
		cpu.Memory[cpu.SP] = low
		cpu.PC = 0x0028
		cpu.Clock += 16
	case 0xF0: // LD A, (0xFF00 + u8)
		n := cpu.ReadMemory(cpu.PC + 1)
		cpu.Registers[RegA] = cpu.ReadMemory(0xFF00 + uint16(n))
		cpu.PC += 2
		cpu.Clock += 12
	case 0xF1: // POP AF
		cpu.PopU16(RegA, RegF)
		cpu.PC++
		cpu.Clock += 12
	case 0xF2: // LD A, (C)
		cpu.Registers[RegA] = cpu.ReadMemory(0xFF00 + uint16(cpu.Registers[RegC]))
		cpu.PC++
		cpu.Clock += 8
	case 0xF3: // DI
		cpu.IME = 0
		cpu.PC++
		cpu.Clock += 4
	case 0xF5: // PUSH AF
		cpu.PushU16(RegA, RegF)
		cpu.PC++
		cpu.Clock += 16
	case 0xF6: // OR A, u8
		cpu.OrU8(cpu.ReadMemory(cpu.PC + 1))
		cpu.PC += 2
		cpu.Clock += 8
	case 0xF7: // RST 30H
		cpu.PC++
		high := uint8(cpu.PC >> 8)
		low := uint8(cpu.PC & 0xFF)
		cpu.SP--
		cpu.Memory[cpu.SP] = high
		cpu.SP--
		cpu.Memory[cpu.SP] = low
		cpu.PC = 0x0030
		cpu.Clock += 16
	case 0xF8: // LD HL, SP + s8
		immediate := int8(cpu.ReadMemory(cpu.PC + 1))
		value := uint16(int32(cpu.SP) + int32(immediate))

		// Set HL to result
		cpu.LoadImmediateU16(RegH, RegL, value)

		// Flags are set based on the lower byte addition only
		cpu.Flags.SetH((cpu.SP&0x0F + uint16(immediate&0x0F)) > 0x0F)
		cpu.Flags.SetC((cpu.SP&0xFF + uint16(immediate)&0xFF) > 0xFF)
		cpu.Flags.SetZ(false)
		cpu.Flags.SetN(false)

		cpu.PC += 2
		cpu.Clock += 12
	case 0xF9: // LD SP, HL
		cpu.SP = cpu.GetHL()
		cpu.PC += 1
		cpu.Clock += 8
	case 0xFA: // LD A, (0xFF00 + C)
		high := cpu.ReadMemory(cpu.PC + 2)
		low := cpu.ReadMemory(cpu.PC + 1)
		cpu.Registers[RegA] = cpu.ReadMemory(uint16(high)<<8 | uint16(low))
		cpu.PC += 3
		cpu.Clock += 16
	case 0xFB: // EI
		cpu.IME = 1
		cpu.PC++
		cpu.Clock += 4
	case 0xFE: // CP A, u8
		cpu.CpU8(cpu.ReadMemory(cpu.PC + 1))
		cpu.PC += 2
		cpu.Clock += 8
	case 0xFF: // RST 00H
		cpu.PC++
		high := uint8(cpu.PC >> 8)
		low := uint8(cpu.PC & 0xFF)
		cpu.SP--
		cpu.Memory[cpu.SP] = high
		cpu.SP--
		cpu.Memory[cpu.SP] = low
		cpu.PC = 0x0038
		cpu.Clock += 16
	default:
		log.Fatalf("Unknown opcode: 0x%02X", cpu.Memory[cpu.PC])
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

func (cpu *CPU) GetAF() uint16 {
	af := uint16(cpu.Registers[RegA])<<8 | uint16(cpu.Registers[RegF])
	return af
}

func (cpu *CPU) LoadMemoryImmediate(address uint16, value uint8) {
	cpu.Memory[address] = value
}

func (cpu *CPU) LoadMemory(address uint16, reg uint8) {
	cpu.LoadMemoryImmediate(address, cpu.Registers[reg])
}

func (cpu *CPU) LoadFromMemory(reg uint8, address uint16) {
	cpu.Registers[reg] = cpu.ReadMemory(address)
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

func (cpu *CPU) IncrementU8(value uint8) (uint8, uint8) {
	flags := uint8(0)

	if (value & 0x0F) == 0x0F {
		flags |= 1 << 5
	}

	value++

	// Zero flag
	if value == 0 {
		flags |= 1 << 7
	}

	// preserve c
	if cpu.Flags.C() {
		flags |= 1 << 4
	}

	return value, flags
}

func (cpu *CPU) DecrementU16Register(high uint8, low uint8) {
	value := uint16(cpu.Registers[high])<<8 | uint16(cpu.Registers[low])
	value--
	cpu.Registers[high] = uint8(value >> 8)
	cpu.Registers[low] = uint8(value & 0xFF)
}

func (cpu *CPU) DecrementU8(value uint8) (uint8, uint8) {
	flags := uint8(0)

	// Check half carry BEFORE decrementing
	// Half carry occurs if lower nibble is 0
	if (value & 0x0F) == 0 {
		flags |= 1 << 5
	}

	value--

	// Zero flag
	if value == 0 {
		flags |= 1 << 7
	}

	// n always set
	flags |= 1 << 6

	// preserve c
	if cpu.Flags.C() {
		flags |= 1 << 4
	}

	return value, flags
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

	// N flag is always reset
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

func (cpu *CPU) PushU16(high, low uint8) {
	cpu.SP--
	cpu.Memory[cpu.SP] = cpu.Registers[high]
	cpu.SP--
	cpu.Memory[cpu.SP] = cpu.Registers[low]
}

func (cpu *CPU) PopU16(high, low uint8) {
	// Read low byte first
	cpu.Registers[low] = cpu.ReadMemory(cpu.SP)
	cpu.SP++

	// Read high byte
	cpu.Registers[high] = cpu.ReadMemory(cpu.SP)
	cpu.SP++

	// If we're popping into AF, we need to handle the flags register specially
	if low == RegF {
		cpu.Flags.SetValue(cpu.Registers[low]) // Use the value we just popped into F
	}
}

// Example implementation
func (cpu *CPU) Swap(value uint8) uint8 {
	// Extract upper and lower nibbles
	upper := (value & 0xF0) >> 4 // Upper 4 bits shifted right
	lower := value & 0x0F        // Lower 4 bits

	// Combine them in swapped order
	result := (lower << 4) | upper

	// Set flags
	cpu.Flags.SetZ(result == 0)
	cpu.Flags.SetN(false)
	cpu.Flags.SetH(false)
	cpu.Flags.SetC(false)

	return result
}
