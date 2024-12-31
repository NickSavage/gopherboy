package main

import (
	"log"
	"testing"
)

type TestProgram struct {
	name     string
	program  []uint8
	setup    func(*CPU)
	validate func(*testing.T, *CPU)
}

func RunTestProgram(t *testing.T, tc TestProgram) {
	cpu := InitCPU()

	// Load the program into ROM
	copy(cpu.ROM, tc.program)

	// Run setup if provided
	if tc.setup != nil {
		tc.setup(cpu)
	}

	// Execute the program
	for i := 0; i < len(tc.program); i++ {
		cpu.ParseNextOpcode()
		log.Printf("a: %d\n", cpu.Registers[RegA])
	}

	// Validate the final state
	tc.validate(t, cpu)
}
func TestLoadMemoryOpcodes(t *testing.T) {
	cpu := InitCPU()
	address := uint16(0x8000)
	cpu.Registers[RegH] = uint8(address >> 8)
	cpu.Registers[RegL] = uint8(address & 0xFF)

	// Define test cases for all load memory operations in the 0x7* range
	testCases := []struct {
		name     string
		srcReg   uint8
		opcode   uint8
		expected uint8
	}{
		{"LD (HL),B", RegB, 0x70, 0x42},
		{"LD (HL),C", RegC, 0x71, 0x33},
		{"LD (HL),D", RegD, 0x72, 0x22},
		{"LD (HL),E", RegE, 0x73, 0x11},
		{"LD (HL),A", RegA, 0x77, 0xFF},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu = InitCPU() // Reset CPU

			// Set up HL with the expected values for the test case
			cpu.Registers[RegH] = 0x00
			cpu.Registers[RegL] = 0x01

			// Test both direct memory load and opcode execution

			// Test 1: Direct memory load
			cpu.Registers[tc.srcReg] = tc.expected
			cpu.LoadMemory(address, tc.srcReg)
			if cpu.Memory[address] != tc.expected {
				t.Errorf("%s direct load failed, expected 0x%02X, got 0x%02X",
					tc.name, tc.expected, cpu.Memory[address])
			}

			// Test 2: Opcode execution
			cpu = InitCPU() // Reset CPU
			cpu.Registers[tc.srcReg] = tc.expected
			cpu.Registers[RegH] = uint8(address >> 8)
			cpu.Registers[RegL] = uint8(address & 0xFF)
			cpu.Registers[tc.srcReg] = tc.expected

			cpu.ROM[0] = tc.opcode
			cpu.ParseNextOpcode()

			if cpu.Memory[address] != tc.expected {
				t.Errorf("%s opcode 0x%02X failed, expected 0x%02X, got 0x%02X",
					tc.name, tc.opcode, tc.expected, cpu.Memory[address])
			}

			if cpu.PC != 1 {
				t.Errorf("%s PC increment failed, expected 1, got %d", tc.name, cpu.PC)
			}
		})
	}
}

func TestExecuteProgram(t *testing.T) {
	testCases := []TestProgram{
		{
			name: "Load series",
			program: []uint8{
				0x41, // LD B,C
				0x42, // LD B,D
				0x43, // LD B,E
			},
			setup: func(cpu *CPU) {
				cpu.Registers[RegC] = 0x12
				cpu.Registers[RegD] = 0x34
				cpu.Registers[RegE] = 0x56
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 3 {
					t.Errorf("Expected PC to be 3, got %d", cpu.PC)
				}
				if cpu.Registers[RegB] != 0x56 {
					t.Errorf("Expected final B value to be 0x56, got 0x%02X", cpu.Registers[RegB])
				}
			},
		},
		{
			name: "Load from memory",
			program: []uint8{
				0x46, // LD B,(HL)
				0x41, // LD B,C
				0x46, // LD B,(HL)
			},
			setup: func(cpu *CPU) {
				// Set up HL to point to 0x8000
				cpu.Registers[RegH] = 0x80
				cpu.Registers[RegL] = 0x00
				// Set value at (HL)
				cpu.Memory[0x8000] = 0x42
				cpu.Registers[RegC] = 0x24
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 3 {
					t.Errorf("Expected PC to be 3, got %d", cpu.PC)
				}
				if cpu.Registers[RegB] != 0x42 {
					t.Errorf("Expected final B value to be 0x42, got 0x%02X", cpu.Registers[RegB])
				}
			},
		},
		{
			name: "Load immediate 1",
			program: []uint8{
				0x0E, 0x01, // LD C, u8
				0x1E, 0x02, // LD E, u8
				0x2E, 0x03, // LD L, u8
				0x3E, 0x04, // LD A, u8
			},
			setup: func(cpu *CPU) {
				// Set up HL to point to 0x8000
				cpu.Registers[RegH] = 0x80
				cpu.Registers[RegL] = 0x00
				// Set value at (HL)
				cpu.Memory[0x8000] = 0x42
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 8 {
					t.Errorf("Expected PC to be 8, got %d", cpu.PC)
				}
				if cpu.Registers[RegC] != 0x01 {
					t.Errorf("Expected final C value to be 0x01, got 0x%02X", cpu.Registers[RegB])
				}
				if cpu.Registers[RegE] != 0x02 {
					t.Errorf("Expected final E value to be 0x01, got 0x%02X", cpu.Registers[RegB])
				}
				if cpu.Registers[RegL] != 0x03 {
					t.Errorf("Expected final L value to be 0x01, got 0x%02X", cpu.Registers[RegB])
				}
				if cpu.Registers[RegA] != 0x04 {
					t.Errorf("Expected final A value to be 0x01, got 0x%02X", cpu.Registers[RegB])
				}

			},
		},
		{
			name: "Load immediate 1",
			program: []uint8{
				0x06, 0x01, // LD B, u8
				0x16, 0x02, // LD D, u8
				0x26, 0x03, // LD H, u8
				0x36, 0x04, // LD (HL), u8
			},
			setup: func(cpu *CPU) {
				// Set up HL to point to 0x8000
				cpu.Registers[RegH] = 0x80
				cpu.Registers[RegL] = 0x00
				// Set value at (HL)
				cpu.Memory[0x8000] = 0x42
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 8 {
					t.Errorf("Expected PC to be 8, got %d", cpu.PC)
				}
				if cpu.Registers[RegB] != 0x01 {
					t.Errorf("Expected final B value to be 0x01, got 0x%02X", cpu.Registers[RegB])
				}
				if cpu.Registers[RegD] != 0x02 {
					t.Errorf("Expected final D value to be 0x01, got 0x%02X", cpu.Registers[RegB])
				}
				if cpu.Registers[RegH] != 0x03 {
					t.Errorf("Expected final H value to be 0x01, got 0x%02X", cpu.Registers[RegB])
				}
				if cpu.Memory[0x0300] != 0x04 {
					t.Errorf("Expected final 0x0300 value to be 0x01, got 0x%02X", cpu.Registers[RegB])
				}

			},
		},
		{
			name: "Load u16 regs",
			program: []uint8{
				0x3E, 0x01, // LD A, u8
				0x06, 0x04, // LD B, u8
				0x0E, 0x04, // LD C, u8
				0x02,
				0x16, 0x05, // LD D, u8
				0x1E, 0x05, // LD E, u8
				0x12,       // LD (DE), A
				0x26, 0x06, // LD H, u8
				0x2E, 0x06, // LD L, u8
				0x22, // LD (HL+), A
				0x32, // LD (HL-), A
				0x32, // LD (HL-), A

			},
			setup: func(cpu *CPU) {
				// Set up HL to point to 0x8000
				cpu.Registers[RegH] = 0x80
				cpu.Registers[RegL] = 0x00
				// Set value at (HL)
				cpu.Memory[0x8000] = 0x42
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 19 {
					t.Errorf("Expected PC to be 19, got %d", cpu.PC)
				}
				if cpu.Memory[0x0404] != 0x01 {
					t.Errorf("Expected memory at 0x0404 to be 0x01, got %02X", cpu.Memory[0x0404])
				}
				if cpu.Memory[0x0505] != 0x01 {
					t.Errorf("Expected memory at 0x0404 to be 0x01, got %02X", cpu.Memory[0x0404])
				}
				if cpu.Memory[0x0606] != 0x01 {
					t.Errorf("Expected memory at 0x0404 to be 0x01, got %02X", cpu.Memory[0x0404])
				}
				if cpu.Registers[RegL] != 0x05 {
					t.Errorf("Expected Register L to be 0x05, got %02X", cpu.Registers[RegL])
				}

			},
		},
		{
			name: "increment u16 regs",
			program: []uint8{
				0x3E, 0x01, // LD A, u8
				0x47,       // LD B, A
				0x0E, 0x04, // LD C, u8
				0x57,       // LD D, A
				0x1E, 0x05, // LD E, u8
				0x67,       // LD H, A
				0x2E, 0x06, // LD L, u8
				0x03, // INC BC
				0x13, // INC DE
				0x23, // INC HL
				0x33, // INC SP
			},
			setup: func(cpu *CPU) {
				// Set up HL to point to 0x8000
				cpu.Registers[RegH] = 0x80
				cpu.Registers[RegL] = 0x00
				// Set value at (HL)
				cpu.Memory[0x8000] = 0x42
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 15 {
					t.Errorf("Expected PC to be 15, got %d", cpu.PC)
				}
				if cpu.Registers[RegC] != 0x05 {
					t.Errorf("Expected Register C to be 0x05, got %02X", cpu.Registers[RegC])
				}
				if cpu.Registers[RegE] != 0x06 {
					t.Errorf("Expected Register E to be 0x06, got %02X", cpu.Registers[RegE])
				}
				if cpu.Registers[RegL] != 0x07 {
					t.Errorf("Expected Register L to be 0x07, got %02X", cpu.Registers[RegL])
				}
				if cpu.SP != 0xFFFF {
					t.Errorf("Expected SP to be 0xFFFF, got %04X", cpu.SP)
				}

			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			RunTestProgram(t, tc)

		})
	}
}

func TestLoadRegister(t *testing.T) {
	cpu := InitCPU()

	// Define test cases for all register combinations
	testCases := []struct {
		name     string
		destReg  uint8
		srcReg   uint8
		initVal  uint8
		expected uint8
		opcode   uint8
	}{
		{"LD B,B", RegB, RegB, 0x42, 0x42, 0x40},
		{"LD B,C", RegB, RegC, 0x42, 0x42, 0x41},
		{"LD B,D", RegB, RegD, 0x42, 0x42, 0x42},
		{"LD B,E", RegB, RegE, 0x42, 0x42, 0x43},
		{"LD B,H", RegB, RegH, 0x42, 0x42, 0x44},
		{"LD B,L", RegB, RegL, 0x42, 0x42, 0x45},
		{"LD B,A", RegB, RegA, 0x42, 0x42, 0x47},

		{"LD C,B", RegC, RegB, 0x42, 0x42, 0x48},
		{"LD C,C", RegC, RegC, 0x42, 0x42, 0x49},
		{"LD C,D", RegC, RegD, 0x42, 0x42, 0x4A},
		{"LD C,E", RegC, RegE, 0x42, 0x42, 0x4B},
		{"LD C,H", RegC, RegH, 0x42, 0x42, 0x4C},
		{"LD C,L", RegC, RegL, 0x42, 0x42, 0x4D},
		{"LD C,A", RegC, RegA, 0x42, 0x42, 0x4F},

		// D register loads (0x50-0x57)
		{"LD D,B", RegD, RegB, 0x42, 0x42, 0x50},
		{"LD D,C", RegD, RegC, 0x42, 0x42, 0x51},
		{"LD D,D", RegD, RegD, 0x42, 0x42, 0x52},
		{"LD D,E", RegD, RegE, 0x42, 0x42, 0x53},
		{"LD D,H", RegD, RegH, 0x42, 0x42, 0x54},
		{"LD D,L", RegD, RegL, 0x42, 0x42, 0x55},
		{"LD D,A", RegD, RegA, 0x42, 0x42, 0x57},

		// E register loads (0x58-0x5F)
		{"LD E,B", RegE, RegB, 0x42, 0x42, 0x58},
		{"LD E,C", RegE, RegC, 0x42, 0x42, 0x59},
		{"LD E,D", RegE, RegD, 0x42, 0x42, 0x5A},
		{"LD E,E", RegE, RegE, 0x42, 0x42, 0x5B},
		{"LD E,H", RegE, RegH, 0x42, 0x42, 0x5C},
		{"LD E,L", RegE, RegL, 0x42, 0x42, 0x5D},
		{"LD E,A", RegE, RegA, 0x42, 0x42, 0x5F},

		// H register loads (0x60-0x67)
		{"LD H,B", RegH, RegB, 0x42, 0x42, 0x60},
		{"LD H,C", RegH, RegC, 0x42, 0x42, 0x61},
		{"LD H,D", RegH, RegD, 0x42, 0x42, 0x62},
		{"LD H,E", RegH, RegE, 0x42, 0x42, 0x63},
		{"LD H,H", RegH, RegH, 0x42, 0x42, 0x64},
		{"LD H,L", RegH, RegL, 0x42, 0x42, 0x65},
		{"LD H,A", RegH, RegA, 0x42, 0x42, 0x67},

		// L register loads (0x68-0x6F)
		{"LD L,B", RegL, RegB, 0x42, 0x42, 0x68},
		{"LD L,C", RegL, RegC, 0x42, 0x42, 0x69},
		{"LD L,D", RegL, RegD, 0x42, 0x42, 0x6A},
		{"LD L,E", RegL, RegE, 0x42, 0x42, 0x6B},
		{"LD L,H", RegL, RegH, 0x42, 0x42, 0x6C},
		{"LD L,L", RegL, RegL, 0x42, 0x42, 0x6D},
		{"LD L,A", RegL, RegA, 0x42, 0x42, 0x6F},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu = InitCPU()

			// Test both direct register load and opcode execution

			// Test 1: Direct register load
			cpu.Registers[tc.srcReg] = tc.initVal
			cpu.LoadRegister(tc.destReg, tc.srcReg)
			if cpu.Registers[tc.destReg] != tc.expected {
				t.Errorf("%s direct load failed, expected 0x%02X, got 0x%02X",
					tc.name, tc.expected, cpu.Registers[tc.destReg])
			}

			// Test 2: Opcode execution
			cpu = InitCPU() // Reset CPU
			cpu.Registers[tc.srcReg] = tc.initVal
			cpu.ROM[0] = tc.opcode
			cpu.ParseNextOpcode()

			if cpu.Registers[tc.destReg] != tc.expected {
				t.Errorf("%s opcode 0x%02X failed, expected 0x%02X, got 0x%02X",
					tc.name, tc.opcode, tc.expected, cpu.Registers[tc.destReg])
			}

			if cpu.PC != 1 {
				t.Errorf("%s PC increment failed, expected 1, got %d", tc.name, cpu.PC)
			}
		})
	}
}

func TestIncrementRegisters(t *testing.T) {
	testCases := []TestProgram{
		{
			name: "Increment 8-bit registers",
			program: []uint8{
				0x3E, 0x41, // LD A,u8
				0x06, 0x41, // LD B,u8
				0x0E, 0x41, // LD C,u8
				0x16, 0x41, // LD D,u8
				0x1E, 0x41, // LD E,u8
				0x26, 0x41, // LD H,u8
				0x2E, 0x41, // LD L,u8
				0x04, // INC B
				0x0C, // INC C
				0x14, // INC D
				0x1C, // INC E
				0x24, // INC H
				0x2C, // INC L
				0x3C, // INC A
			},
			setup: func(cpu *CPU) {
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 21 {
					t.Errorf("Expected PC to be 21, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x42 {
					t.Errorf("Expected Register A to be 0x42, got %02X", cpu.Registers[RegA])
				}
				if cpu.Registers[RegB] != 0x42 {
					t.Errorf("Expected Register B to be 0x42, got %02X", cpu.Registers[RegB])
				}
				if cpu.Registers[RegC] != 0x42 {
					t.Errorf("Expected Register C to be 0x42, got %02X", cpu.Registers[RegC])
				}
				if cpu.Registers[RegD] != 0x42 {
					t.Errorf("Expected Register D to be 0x42, got %02X", cpu.Registers[RegD])
				}
				if cpu.Registers[RegE] != 0x42 {
					t.Errorf("Expected Register E to be 0x42, got %02X", cpu.Registers[RegE])
				}
				if cpu.Registers[RegH] != 0x42 {
					t.Errorf("Expected Register H to be 0x42, got %02X", cpu.Registers[RegH])
				}
				if cpu.Registers[RegL] != 0x42 {
					t.Errorf("Expected Register L to be 0x42, got %02X", cpu.Registers[RegL])
				}
			},
		},
		{
			name: "Increment (HL)",
			program: []uint8{
				0x26, 0x80, // LD H,u8
				0x2E, 0x00, // LD L,u8
				0x36, 0x41, // LD (HL),u8
				0x34, // INC (HL)
			},
			setup: func(cpu *CPU) {
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 7 {
					t.Errorf("Expected PC to be 7, got %d", cpu.PC)
				}
				if cpu.Memory[0x8000] != 0x42 {
					t.Errorf("Expected memory at (HL) to be 0x42, got %02X", cpu.Memory[0x8000])
				}
				// Verify HL wasn't modified
				if cpu.Registers[RegH] != 0x80 {
					t.Errorf("Expected Register H to be 0x80, got %02X", cpu.Registers[RegH])
				}
				if cpu.Registers[RegL] != 0x00 {
					t.Errorf("Expected Register L to be 0x00, got %02X", cpu.Registers[RegL])
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			RunTestProgram(t, tc)
		})
	}
}
func TestDecrementRegisters(t *testing.T) {
	testCases := []TestProgram{
		{
			name: "Decrement 16-bit registers",
			program: []uint8{
				0x06, 0x80, // LD B,u8
				0x16, 0x80, // LD B,u8
				0x26, 0x80, // LD B,u8
				0x0B, // DEC BC
				0x1B, // DEC DE
				0x2B, // DEC HL
				0x3B, // DEC SP
			},
			setup: func(cpu *CPU) {
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 10 {
					t.Errorf("Expected PC to be 10, got %d", cpu.PC)
				}
				if cpu.Registers[RegC] != 0xFF {
					t.Errorf("Expected final C value to be 0xFF, got 0x%02X", cpu.Registers[RegC])
				}
				if cpu.Registers[RegE] != 0xFF {
					t.Errorf("Expected final E value to be 0xFF, got 0x%02X", cpu.Registers[RegE])
				}
				if cpu.Registers[RegL] != 0xFF {
					t.Errorf("Expected final L value to be 0xFF, got 0x%02X", cpu.Registers[RegL])
				}
				if cpu.SP != 0xFFFD {
					t.Errorf("Expected final SP value to be 0xFFFD, got 0x%02X", cpu.SP)
				}
			},
		},
		{
			name: "Decrement 8-bit registers",
			program: []uint8{
				0x3E, 0x42, // LD A,u8
				0x06, 0x42, // LD B,u8
				0x0E, 0x42, // LD C,u8
				0x16, 0x42, // LD D,u8
				0x1E, 0x42, // LD E,u8
				0x26, 0x42, // LD H,u8
				0x2E, 0x42, // LD L,u8
				0x05, // DEC B
				0x0D, // DEC C
				0x15, // DEC D
				0x1D, // DEC E
				0x25, // DEC H
				0x2D, // DEC L
				0x3D, // DEC A
			},
			setup: func(cpu *CPU) {
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 21 {
					t.Errorf("Expected PC to be 21, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x41 {
					t.Errorf("Expected Register A to be 0x41, got %02X", cpu.Registers[RegA])
				}
				if cpu.Registers[RegB] != 0x41 {
					t.Errorf("Expected Register B to be 0x41, got %02X", cpu.Registers[RegB])
				}
				if cpu.Registers[RegC] != 0x41 {
					t.Errorf("Expected Register C to be 0x41, got %02X", cpu.Registers[RegC])
				}
				if cpu.Registers[RegD] != 0x41 {
					t.Errorf("Expected Register D to be 0x41, got %02X", cpu.Registers[RegD])
				}
				if cpu.Registers[RegE] != 0x41 {
					t.Errorf("Expected Register E to be 0x41, got %02X", cpu.Registers[RegE])
				}
				if cpu.Registers[RegH] != 0x41 {
					t.Errorf("Expected Register H to be 0x41, got %02X", cpu.Registers[RegH])
				}
				if cpu.Registers[RegL] != 0x41 {
					t.Errorf("Expected Register L to be 0x41, got %02X", cpu.Registers[RegL])
				}
			},
		},
		{
			name: "Decrement (HL)",
			program: []uint8{
				0x26, 0x80, // LD H,u8
				0x2E, 0x00, // LD L,u8
				0x36, 0x42, // LD (HL),u8
				0x35, // DEC (HL)
			},
			setup: func(cpu *CPU) {
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 7 {
					t.Errorf("Expected PC to be 7, got %d", cpu.PC)
				}
				if cpu.Memory[0x8000] != 0x41 {
					t.Errorf("Expected memory at (HL) to be 0x41, got %02X", cpu.Memory[0x8000])
				}
				// Verify HL wasn't modified
				if cpu.Registers[RegH] != 0x80 {
					t.Errorf("Expected Register H to be 0x80, got %02X", cpu.Registers[RegH])
				}
				if cpu.Registers[RegL] != 0x00 {
					t.Errorf("Expected Register L to be 0x00, got %02X", cpu.Registers[RegL])
				}
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			RunTestProgram(t, tc)

		})
	}
}

func TestLDu16SP(t *testing.T) {
	testCases := []TestProgram{
		{
			name: "LD (u16), SP",
			program: []uint8{
				0x31, 0x34, 0x12, // LD SP, 0x1234
				0x08, 0x00, 0x80, // LD (0x8000), SP
			},
			setup: func(cpu *CPU) {
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 6 {
					t.Errorf("Expected PC to be 6, got %d", cpu.PC)
				}
				// Check that SP was stored correctly in little-endian format
				if cpu.Memory[0x8000] != 0x34 {
					t.Errorf("Expected memory at 0x8000 (low byte) to be 0x34, got %02X", cpu.Memory[0x8000])
				}
				if cpu.Memory[0x8001] != 0x12 {
					t.Errorf("Expected memory at 0x8001 (high byte) to be 0x12, got %02X", cpu.Memory[0x8001])
				}
				// Verify SP wasn't modified
				if cpu.SP != 0x1234 {
					t.Errorf("Expected SP to be 0x1234, got %04X", cpu.SP)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			RunTestProgram(t, tc)
		})
	}
}
func TestLoadAFromRegisters(t *testing.T) {
	testCases := []TestProgram{
		{
			name: "Load A from register pairs",
			program: []uint8{
				0x06, 0x80, // LD B,0x80
				0x0E, 0x00, // LD C,0x00
				0x16, 0x80, // LD D,0x80
				0x1E, 0x01, // LD E,0x01
				0x26, 0x80, // LD H,0x80
				0x2E, 0x02, // LD L,0x02
				0x3E, 0x00, // LD A,0x00
				0x0A, // LD A,(BC)
				0x1A, // LD A,(DE)
				0x2A, // LD A,(HL+)
				0x3A, // LD A,(HL-)
			},
			setup: func(cpu *CPU) {
				// Set up test values in memory
				cpu.Memory[0x8000] = 0x42 // Value at (BC)
				cpu.Memory[0x8001] = 0x43 // Value at (DE)
				cpu.Memory[0x8002] = 0x44 // First value at (HL)
				cpu.Memory[0x8003] = 0x45 // Second value at (HL)
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 18 {
					t.Errorf("Expected PC to be 18, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x45 {
					t.Errorf("Expected Register A to be 0x45, got %02X", cpu.Registers[RegA])
				}
				// Check if HL was incremented and then decremented correctly
				if cpu.Registers[RegH] != 0x80 {
					t.Errorf("Expected Register H to be 0x80, got %02X", cpu.Registers[RegH])
				}
				if cpu.Registers[RegL] != 0x02 {
					t.Errorf("Expected Register L to be 0x02, got %02X", cpu.Registers[RegL])
				}
			},
		},
		{
			name: "Load A from HL with increment/decrement",
			program: []uint8{
				0x26, 0x80, // LD H,0x80
				0x2E, 0x00, // LD L,0x00
				0x2A, // LD A,(HL+)
				0x2A, // LD A,(HL+)
				0x2A, // LD A,(HL+)
				0x3A, // LD A,(HL-)
				0x3A, // LD A,(HL-)
			},
			setup: func(cpu *CPU) {
				// Set up sequential values in memory
				cpu.Memory[0x8000] = 0x41
				cpu.Memory[0x8001] = 0x42
				cpu.Memory[0x8002] = 0x43
				cpu.Memory[0x8003] = 0x44
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 9 {
					t.Errorf("Expected PC to be 9, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x43 {
					t.Errorf("Expected Register A to be 0x43, got %02X", cpu.Registers[RegA])
				}
				// Check if HL was modified correctly after all operations
				if cpu.Registers[RegH] != 0x80 {
					t.Errorf("Expected Register H to be 0x80, got %02X", cpu.Registers[RegH])
				}
				if cpu.Registers[RegL] != 0x01 {
					t.Errorf("Expected Register L to be 0x01, got %02X", cpu.Registers[RegL])
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			RunTestProgram(t, tc)
		})
	}
}

func TestRotateInstructions(t *testing.T) {
	testCases := []TestProgram{
		{
			name: "RLCA basic rotation",
			program: []uint8{
				0x3E, 0x85, // LD A, 0x85 (10000101)
				0x07, // RLCA
			},
			setup: func(cpu *CPU) {
				cpu.Flags.SetValue(0x00) // Clear all flags
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 3 {
					t.Errorf("Expected PC to be 3, got %d", cpu.PC)
				}
				// After RLCA, 0x85 (10000101) becomes 0x0B (00001011)
				// Because bit 7 (1) goes into both carry and bit 0
				if cpu.Registers[RegA] != 0x0B {
					t.Errorf("Expected Register A to be 0x0B, got %02X", cpu.Registers[RegA])
				}
				// Carry flag should be set (1)
				if !cpu.Flags.C() {
					t.Error("Expected Carry flag to be set")
				}
				// Other flags should be reset
				if cpu.Flags.Z() {
					t.Error("Expected Zero flag to be reset")
				}
				if cpu.Flags.N() {
					t.Error("Expected N flag to be reset")
				}
				if cpu.Flags.H() {
					t.Error("Expected H flag to be reset")
				}
			},
		},
		{
			name: "RLCA no carry case",
			program: []uint8{
				0x3E, 0x42, // LD A, 0x42 (01000010)
				0x07, // RLCA
			},
			setup: func(cpu *CPU) {
				cpu.Flags.SetValue(0xF0) // Set all flags initially
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 3 {
					t.Errorf("Expected PC to be 3, got %d", cpu.PC)
				}
				// After RLCA, 0x42 (01000010) becomes 0x84 (10000100)
				if cpu.Registers[RegA] != 0x84 {
					t.Errorf("Expected Register A to be 0x84, got %02X", cpu.Registers[RegA])
				}
				// Carry flag should be reset (0)
				if cpu.Flags.C() {
					t.Error("Expected Carry flag to be reset")
				}
				// Other flags should be reset
				if cpu.Flags.Z() {
					t.Error("Expected Zero flag to be reset")
				}
				if cpu.Flags.N() {
					t.Error("Expected N flag to be reset")
				}
				if cpu.Flags.H() {
					t.Error("Expected H flag to be reset")
				}
			},
		},
		{
			name: "RLCA multiple rotations",
			program: []uint8{
				0x3E, 0xFF, // LD A, 0xFF
				0x07, // RLCA
				0x07, // RLCA
				0x07, // RLCA
			},
			setup: func(cpu *CPU) {
				cpu.Flags.SetValue(0x00)
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				// After 3 RLCAs, 0xFF should still be 0xFF
				if cpu.Registers[RegA] != 0xFF {
					t.Errorf("Expected Register A to be 0xFF, got %02X", cpu.Registers[RegA])
				}
				// Carry flag should be set (1)
				if !cpu.Flags.C() {
					t.Error("Expected Carry flag to be set")
				}
				// Other flags should be reset
				if cpu.Flags.Z() {
					t.Error("Expected Zero flag to be reset")
				}
				if cpu.Flags.N() {
					t.Error("Expected N flag to be reset")
				}
				if cpu.Flags.H() {
					t.Error("Expected H flag to be reset")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			RunTestProgram(t, tc)
		})
	}
}

func TestLoadRP16Immediate(t *testing.T) {
	testCases := []TestProgram{
		{
			name: "LD BC,u16",
			program: []uint8{
				0x01, 0x34, 0x12, // LD BC,0x1234
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 3 {
					t.Errorf("Expected PC to be 3, got %d", cpu.PC)
				}
				bc := uint16(cpu.Registers[RegB])<<8 | uint16(cpu.Registers[RegC])
				if bc != 0x1234 {
					t.Errorf("Expected BC to be 0x1234, got %04X", bc)
				}
				// Verify individual registers
				if cpu.Registers[RegB] != 0x12 {
					t.Errorf("Expected B to be 0x12, got %02X", cpu.Registers[RegB])
				}
				if cpu.Registers[RegC] != 0x34 {
					t.Errorf("Expected C to be 0x34, got %02X", cpu.Registers[RegC])
				}
			},
		},
		{
			name: "LD DE,u16",
			program: []uint8{
				0x11, 0x78, 0x56, // LD DE,0x5678
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 3 {
					t.Errorf("Expected PC to be 3, got %d", cpu.PC)
				}
				de := uint16(cpu.Registers[RegD])<<8 | uint16(cpu.Registers[RegE])
				if de != 0x5678 {
					t.Errorf("Expected DE to be 0x5678, got %04X", de)
				}
				// Verify individual registers
				if cpu.Registers[RegD] != 0x56 {
					t.Errorf("Expected D to be 0x56, got %02X", cpu.Registers[RegD])
				}
				if cpu.Registers[RegE] != 0x78 {
					t.Errorf("Expected E to be 0x78, got %02X", cpu.Registers[RegE])
				}
			},
		},
		{
			name: "LD HL,u16",
			program: []uint8{
				0x21, 0xBC, 0x9A, // LD HL,0x9ABC
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 3 {
					t.Errorf("Expected PC to be 3, got %d", cpu.PC)
				}
				hl := uint16(cpu.Registers[RegH])<<8 | uint16(cpu.Registers[RegL])
				if hl != 0x9ABC {
					t.Errorf("Expected HL to be 0x9ABC, got %04X", hl)
				}
				// Verify individual registers
				if cpu.Registers[RegH] != 0x9A {
					t.Errorf("Expected H to be 0x9A, got %02X", cpu.Registers[RegH])
				}
				if cpu.Registers[RegL] != 0xBC {
					t.Errorf("Expected L to be 0xBC, got %02X", cpu.Registers[RegL])
				}
			},
		},
		{
			name: "LD SP,u16",
			program: []uint8{
				0x31, 0xEF, 0xCD, // LD SP,0xCDEF
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 3 {
					t.Errorf("Expected PC to be 3, got %d", cpu.PC)
				}
				if cpu.SP != 0xCDEF {
					t.Errorf("Expected SP to be 0xCDEF, got %04X", cpu.SP)
				}
			},
		},
		{
			name: "Multiple LD rr,u16",
			program: []uint8{
				0x01, 0x34, 0x12, // LD BC,0x1234
				0x11, 0x78, 0x56, // LD DE,0x5678
				0x21, 0xBC, 0x9A, // LD HL,0x9ABC
				0x31, 0xEF, 0xCD, // LD SP,0xCDEF
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 12 {
					t.Errorf("Expected PC to be 12, got %d", cpu.PC)
				}

				bc := uint16(cpu.Registers[RegB])<<8 | uint16(cpu.Registers[RegC])
				if bc != 0x1234 {
					t.Errorf("Expected BC to be 0x1234, got %04X", bc)
				}

				de := uint16(cpu.Registers[RegD])<<8 | uint16(cpu.Registers[RegE])
				if de != 0x5678 {
					t.Errorf("Expected DE to be 0x5678, got %04X", de)
				}

				hl := uint16(cpu.Registers[RegH])<<8 | uint16(cpu.Registers[RegL])
				if hl != 0x9ABC {
					t.Errorf("Expected HL to be 0x9ABC, got %04X", hl)
				}

				if cpu.SP != 0xCDEF {
					t.Errorf("Expected SP to be 0xCDEF, got %04X", cpu.SP)
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			RunTestProgram(t, tc)
		})
	}
}

func TestAddRP(t *testing.T) {
	testCases := []TestProgram{
		{
			name: "ADD HL,BC",
			program: []uint8{
				0x21, 0x00, 0x10, // LD HL,0x1000
				0x01, 0xFF, 0x0F, // LD BC,0x0FFF
				0x09, // ADD HL,BC
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 7 {
					t.Errorf("Expected PC to be 7, got %d", cpu.PC)
				}

				// Result should be 0x1FFF
				hl := uint16(cpu.Registers[RegH])<<8 | uint16(cpu.Registers[RegL])
				if hl != 0x1FFF {
					t.Errorf("Expected HL to be 0x1FFF, got %04X", hl)
				}
				bc := uint16(cpu.Registers[RegB])<<8 | uint16(cpu.Registers[RegC])
				if bc != 0x0FFF {
					t.Errorf("Expected BC to be 0x0FFF, got %04X", bc)
				}

				// N flag should be reset
				if cpu.Flags.N() {
					t.Error("Expected N flag to be reset")
				}
			},
		},
		{
			name: "ADD HL,DE",
			program: []uint8{
				0x21, 0xFF, 0xFF, // LD HL,0xFFFF
				0x11, 0x01, 0x00, // LD DE,0x0001
				0x19, // ADD HL,DE
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 7 {
					t.Errorf("Expected PC to be 7, got %d", cpu.PC)
				}

				// Result should be 0x0000 with carry
				hl := uint16(cpu.Registers[RegH])<<8 | uint16(cpu.Registers[RegL])
				if hl != 0x0000 {
					t.Errorf("Expected HL to be 0x0000, got %04X", hl)
				}

				// Carry flag should be set
				if !cpu.Flags.C() {
					t.Error("Expected Carry flag to be set")
				}

				// N flag should be reset
				if cpu.Flags.N() {
					t.Error("Expected N flag to be reset")
				}
			},
		},
		{
			name: "ADD HL,HL",
			program: []uint8{
				0x21, 0xF0, 0x0F, // LD HL,0x0FF0
				0x29, // ADD HL,HL
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 4 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}

				// Result should be 0x1FE0
				hl := uint16(cpu.Registers[RegH])<<8 | uint16(cpu.Registers[RegL])
				if hl != 0x1FE0 {
					t.Errorf("Expected HL to be 0x1FE0, got %04X", hl)
				}

				// N flag should be reset
				if cpu.Flags.N() {
					t.Error("Expected N flag to be reset")
				}
			},
		},
		{
			name: "ADD HL,SP",
			program: []uint8{
				0x21, 0x00, 0x80, // LD HL,0x8000
				0x31, 0x00, 0x80, // LD SP,0x8000
				0x39, // ADD HL,SP
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 7 {
					t.Errorf("Expected PC to be 7, got %d", cpu.PC)
				}

				// Result should be 0x0000
				hl := uint16(cpu.Registers[RegH])<<8 | uint16(cpu.Registers[RegL])
				if hl != 0x0000 {
					t.Errorf("Expected HL to be 0x0000, got %04X", hl)
				}

				// Carry flag should be set since 0x8000 + 0x8000 = 0x10000
				if !cpu.Flags.C() {
					t.Error("Expected Carry flag to be set")
				}

				// N flag should be reset
				if cpu.Flags.N() {
					t.Error("Expected N flag to be reset")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			RunTestProgram(t, tc)
		})
	}
}

func TestAddOpcodes(t *testing.T) {
	testCases := []TestProgram{
		{
			name: "ADD A,B",
			program: []uint8{
				0x3E, 0x42, // LD A,0x42
				0x06, 0x24, // LD B,0x24
				0x80, // ADD A,B
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x66 {
					t.Errorf("Expected A to be 0x66, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected all flags to be reset")
				}
			},
		},
		{
			name: "ADD A,C",
			program: []uint8{
				0x3E, 0xFF, // LD A,0xFF
				0x0E, 0x01, // LD C,0x01
				0x81, // ADD A,C
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x00 {
					t.Errorf("Expected A to be 0x00, got %02X", cpu.Registers[RegA])
				}
				if !cpu.Flags.Z() || cpu.Flags.N() || !cpu.Flags.H() || !cpu.Flags.C() {
					t.Error("Expected Z, H, and C flags to be set, N to be reset")
				}
			},
		},
		{
			name: "ADD A,D",
			program: []uint8{
				0x3E, 0x0F, // LD A,0x0F
				0x16, 0x01, // LD D,0x01
				0x82, // ADD A,D
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x10 {
					t.Errorf("Expected A to be 0x10, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected H flag to be set, others reset")
				}
			},
		},
		{
			name: "ADD A,E",
			program: []uint8{
				0x3E, 0x80, // LD A,0x80
				0x1E, 0x80, // LD E,0x80
				0x83, // ADD A,E
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x00 {
					t.Errorf("Expected A to be 0x00, got %02X", cpu.Registers[RegA])
				}
				if !cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || !cpu.Flags.C() {
					t.Error("Expected Z and C flags to be set, N and H reset")
				}
			},
		},
		{
			name: "ADD A,H",
			program: []uint8{
				0x3E, 0x2F, // LD A,0x2F
				0x26, 0x11, // LD H,0x11
				0x84, // ADD A,H
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x40 {
					t.Errorf("Expected A to be 0x40, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Errorf("Expected H flag to be set, all else to be reset, got %v", cpu.Flags)
				}
			},
		},
		{
			name: "ADD A,L",
			program: []uint8{
				0x3E, 0x11, // LD A,0x11
				0x2E, 0x22, // LD L,0x22
				0x85, // ADD A,L
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x33 {
					t.Errorf("Expected A to be 0x33, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected all flags to be reset")
				}
			},
		},
		{
			name: "ADD A,(HL)",
			program: []uint8{
				0x21, 0x00, 0x80, // LD HL,0x8000
				0x3E, 0x42, // LD A,0x42
				0x86, // ADD A,(HL)
			},
			setup: func(cpu *CPU) {
				cpu.Memory[0x8000] = 0x24 // Set value at (HL)
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 6 {
					t.Errorf("Expected PC to be 6, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x66 {
					t.Errorf("Expected A to be 0x66, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected all flags to be reset")
				}
			},
		},
		{
			name: "ADD A,A",
			program: []uint8{
				0x3E, 0x44, // LD A,0x44
				0x87, // ADD A,A
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 3 {
					t.Errorf("Expected PC to be 3, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x88 {
					t.Errorf("Expected A to be 0x88, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected all flags to be reset")
				}
			},
		},
		{
			name: "ADD A,u8",
			program: []uint8{
				0x3E, 0x42, // LD A,0x42
				0xC6, 0x24, // ADD A,0x24
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 4 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x66 {
					t.Errorf("Expected A to be 0x66, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected all flags to be reset")
				}
			},
		},
		{
			name: "ADD A,u8 with carry and half carry",
			program: []uint8{
				0x3E, 0xFF, // LD A,0xFF
				0xC6, 0x01, // ADD A,0x01
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 4 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x00 {
					t.Errorf("Expected A to be 0x00, got %02X", cpu.Registers[RegA])
				}
				if !cpu.Flags.Z() || cpu.Flags.N() || !cpu.Flags.H() || !cpu.Flags.C() {
					t.Error("Expected Z, H, and C flags to be set, N to be reset")
				}
			},
		},
		{
			name: "ADD A,u8 with half carry only",
			program: []uint8{
				0x3E, 0x0F, // LD A,0x0F
				0xC6, 0x01, // ADD A,0x01
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 4 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x10 {
					t.Errorf("Expected A to be 0x10, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected H flag to be set, others reset")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			RunTestProgram(t, tc)
		})
	}
}

func TestADCOpcodes(t *testing.T) {
	testCases := []TestProgram{
		{
			name: "ADC A,B - Basic addition with no carry",
			program: []uint8{
				0x3E, 0x42, // LD A,0x42
				0x06, 0x24, // LD B,0x24
				0x88, // ADC A,B
			},
			setup: func(cpu *CPU) {
				cpu.Flags.SetValue(0) // Clear carry flag
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x66 {
					t.Errorf("Expected A to be 0x66, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected all flags to be reset")
				}
			},
		},
		{
			name: "ADC A,B - Addition with carry set",
			program: []uint8{
				0x3E, 0x42, // LD A,0x42
				0x06, 0x24, // LD B,0x24
				0x88, // ADC A,B
			},
			setup: func(cpu *CPU) {
				cpu.Flags.SetC(true) // Set carry flag
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x67 { // 0x42 + 0x24 + 1
					t.Errorf("Expected A to be 0x67, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected all flags to be reset")
				}
			},
		},
		{
			name: "ADC A,C - Carry and half-carry generation",
			program: []uint8{
				0x3E, 0xFF, // LD A,0xFF
				0x0E, 0x01, // LD C,0x01
				0x89, // ADC A,C
			},
			setup: func(cpu *CPU) {
				cpu.Flags.SetC(true) // Set carry flag
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x01 { // 0xFF + 0x01 + 1 = 0x101 (overflow)
					t.Errorf("Expected A to be 0x01, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || !cpu.Flags.H() || !cpu.Flags.C() {
					t.Error("Expected H and C flags to be set, Z and N to be reset")
				}
			},
		},
		{
			name: "ADC A,(HL) - Memory operation",
			program: []uint8{
				0x21, 0x00, 0x80, // LD HL,0x8000
				0x3E, 0x42, // LD A,0x42
				0x8E, // ADC A,(HL)
			},
			setup: func(cpu *CPU) {
				cpu.Memory[0x8000] = 0x24 // Set value at (HL)
				cpu.Flags.SetC(true)      // Set carry flag
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 6 {
					t.Errorf("Expected PC to be 6, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x67 { // 0x42 + 0x24 + 1
					t.Errorf("Expected A to be 0x67, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected all flags to be reset")
				}
			},
		},
		{
			name: "ADC A,u8 - Immediate value",
			program: []uint8{
				0x3E, 0x42, // LD A,0x42
				0xCE, 0x24, // ADC A,0x24
			},
			setup: func(cpu *CPU) {
				cpu.Flags.SetC(true) // Set carry flag
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 4 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x67 { // 0x42 + 0x24 + 1
					t.Errorf("Expected A to be 0x67, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected all flags to be reset")
				}
			},
		},
		{
			name: "ADC A,A - Add A to itself with carry",
			program: []uint8{
				0x3E, 0x80, // LD A,0x80
				0x8F, // ADC A,A (was incorrectly using 0x87 which is ADD A,A)
			},
			setup: func(cpu *CPU) {
				cpu.Flags.SetC(true) // Set carry flag
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 3 {
					t.Errorf("Expected PC to be 3, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x01 { // 0x80 + 0x80 + 1 = 0x101 (overflow)
					t.Errorf("Expected A to be 0x01, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || !cpu.Flags.C() {
					t.Error("Expected C flag to be set, others reset")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			RunTestProgram(t, tc)
		})
	}
}
func TestSubOpcodes(t *testing.T) {
	testCases := []TestProgram{
		{
			name: "SUB B - Basic subtraction",
			program: []uint8{
				0x3E, 0x42, // LD A,0x42
				0x06, 0x24, // LD B,0x24
				0x90, // SUB B
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x1E {
					t.Errorf("Expected A to be 0x1E, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || !cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Errorf("Expected N flag to be set, others reset, got %v", cpu.Flags)
				}
			},
		},
		{
			name: "SUB C - Subtraction with half carry",
			program: []uint8{
				0x3E, 0x20, // LD A,0x20
				0x0E, 0x11, // LD C,0x11
				0x91, // SUB C
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x0F {
					t.Errorf("Expected A to be 0x0F, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || !cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected N and H flags to be set, others reset")
				}
			},
		},
		{
			name: "SUB D - Subtraction with carry",
			program: []uint8{
				0x3E, 0x20, // LD A,0x20
				0x16, 0x30, // LD D,0x30
				0x92, // SUB D
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0xF0 { // 0x20 - 0x30 = -0x10 = 0xF0
					t.Errorf("Expected A to be 0xF0, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || !cpu.Flags.N() || cpu.Flags.H() || !cpu.Flags.C() {
					t.Error("Expected N, H, and C flags to be set, Z reset")
				}
			},
		},
		{
			name: "SUB (HL) - Memory operation",
			program: []uint8{
				0x21, 0x00, 0x80, // LD HL,0x8000
				0x3E, 0x42, // LD A,0x42
				0x96, // SUB (HL)
			},
			setup: func(cpu *CPU) {
				cpu.Memory[0x8000] = 0x24 // Set value at (HL)
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 6 {
					t.Errorf("Expected PC to be 6, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x1E { // 0x42 - 0x24
					t.Errorf("Expected A to be 0x1E, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || !cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected N and H flags to be set, others reset")
					t.Errorf("flags: z %v, n %v, h %v, c %v", cpu.Flags.Z(), cpu.Flags.N(), cpu.Flags.H(), cpu.Flags.C())

				}
			},
		},
		{
			name: "SUB A - Subtraction from self",
			program: []uint8{
				0x3E, 0x42, // LD A,0x42
				0x97, // SUB A
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 3 {
					t.Errorf("Expected PC to be 3, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x00 {
					t.Errorf("Expected A to be 0x00, got %02X", cpu.Registers[RegA])
				}
				if !cpu.Flags.Z() || !cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected Z and N flags to be set, H and C reset")
				}
			},
		},
		{
			name: "SUB u8 - Immediate value",
			program: []uint8{
				0x3E, 0x42, // LD A,0x42
				0xD6, 0x24, // SUB 0x24
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 4 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x1E {
					t.Errorf("Expected A to be 0x1E, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || !cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected N and H flags to be set, others reset")
				}
			},
		},
		{
			name: "SUB - Zero result",
			program: []uint8{
				0x3E, 0x24, // LD A,0x24
				0xD6, 0x24, // SUB 0x24
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 4 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x00 {
					t.Errorf("Expected A to be 0x00, got %02X", cpu.Registers[RegA])
				}
				if !cpu.Flags.Z() || !cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected Z and N flags to be set, H and C reset")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			RunTestProgram(t, tc)
		})
	}
}
func TestSBCOpcodes(t *testing.T) {
	testCases := []TestProgram{
		{
			name: "SBC A,B - Basic subtraction with no carry",
			program: []uint8{
				0x3E, 0x42, // LD A,0x42
				0x06, 0x24, // LD B,0x24
				0x98, // SBC B
			},
			setup: func(cpu *CPU) {
				cpu.Flags.SetC(false) // Clear carry flag
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x1E { // 0x42 - 0x24 - 0
					t.Errorf("Expected A to be 0x1E, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || !cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected N and H flags to be set, others reset")
				}
			},
		},
		{
			name: "SBC A,B - Subtraction with carry set",
			program: []uint8{
				0x3E, 0x42, // LD A,0x42
				0x06, 0x24, // LD B,0x24
				0x98, // SBC B
			},
			setup: func(cpu *CPU) {
				cpu.Flags.SetC(true) // Set carry flag
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x1D { // 0x42 - 0x24 - 1
					t.Errorf("Expected A to be 0x1D, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || !cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected N and H flags to be set, others reset")
				}
			},
		},
		{
			name: "SBC A,(HL) - Memory operation with carry",
			program: []uint8{
				0x21, 0x00, 0x80, // LD HL,0x8000
				0x3E, 0x42, // LD A,0x42
				0x9E, // SBC (HL)
			},
			setup: func(cpu *CPU) {
				cpu.Memory[0x8000] = 0x24 // Set value at (HL)
				cpu.Flags.SetC(true)      // Set carry flag
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 6 {
					t.Errorf("Expected PC to be 6, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x1D { // 0x42 - 0x24 - 1
					t.Errorf("Expected A to be 0x1D, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || !cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected N and H flags to be set, others reset")
				}
			},
		},
		{
			name: "SBC A,A - Subtract A from itself with carry",
			program: []uint8{
				0x3E, 0x42, // LD A,0x42
				0x9F, // SBC A
			},
			setup: func(cpu *CPU) {
				cpu.Flags.SetC(true) // Set carry flag
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 3 {
					t.Errorf("Expected PC to be 3, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0xFF { // 0x42 - 0x42 - 1 = -1 = 0xFF
					t.Errorf("Expected A to be 0xFF, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || !cpu.Flags.N() || !cpu.Flags.H() || !cpu.Flags.C() {
					t.Error("Expected N, H, and C flags to be set, Z reset")
				}
			},
		},
		{
			name: "SBC A,u8 - Immediate value with carry",
			program: []uint8{
				0x3E, 0x42, // LD A,0x42
				0xDE, 0x24, // SBC 0x24
			},
			setup: func(cpu *CPU) {
				cpu.Flags.SetC(true) // Set carry flag
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 4 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x1D { // 0x42 - 0x24 - 1
					t.Errorf("Expected A to be 0x1D, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || !cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected N and H flags to be set, others reset")
				}
			},
		},
		{
			name: "SBC A,u8 - Borrow and half-borrow test",
			program: []uint8{
				0x3E, 0x10, // LD A,0x10
				0xDE, 0x20, // SBC 0x20
			},
			setup: func(cpu *CPU) {
				cpu.Flags.SetC(true) // Set carry flag
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 4 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0xEF { // 0x10 - 0x20 - 1 = -0x11 = 0xEF
					t.Errorf("Expected A to be 0xEF, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || !cpu.Flags.N() || !cpu.Flags.H() || !cpu.Flags.C() {
					t.Error("Expected N, H, and C flags to be set, Z reset")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			RunTestProgram(t, tc)
		})
	}
}
func TestCPOpcodes(t *testing.T) {
	testCases := []TestProgram{
		{
			name: "CP B - Equal values",
			program: []uint8{
				0x3E, 0x42, // LD A,0x42
				0x06, 0x42, // LD B,0x42
				0xB8, // CP B
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x42 { // A should be unchanged
					t.Errorf("Expected A to be 0x42, got %02X", cpu.Registers[RegA])
				}
				if !cpu.Flags.Z() || !cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected Z and N flags to be set, H and C reset")
				}
			},
		},
		{
			name: "CP C - A greater than C",
			program: []uint8{
				0x3E, 0x42, // LD A,0x42
				0x0E, 0x24, // LD C,0x24
				0xB9, // CP C
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x42 { // A should be unchanged
					t.Errorf("Expected A to be 0x42, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || !cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected N and H flags to be set, Z and C reset")
				}
			},
		},
		{
			name: "CP D - A less than D",
			program: []uint8{
				0x3E, 0x20, // LD A,0x20
				0x16, 0x30, // LD D,0x30
				0xBA, // CP D
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x20 { // A should be unchanged
					t.Errorf("Expected A to be 0x20, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || !cpu.Flags.N() || cpu.Flags.H() || !cpu.Flags.C() {
					t.Error("Expected N and C flags to be set, Z and H reset")
				}
			},
		},
		{
			name: "CP (HL) - Memory comparison",
			program: []uint8{
				0x21, 0x00, 0x80, // LD HL,0x8000
				0x3E, 0x42, // LD A,0x42
				0xBE, // CP (HL)
			},
			setup: func(cpu *CPU) {
				cpu.Memory[0x8000] = 0x42 // Set value at (HL)
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 6 {
					t.Errorf("Expected PC to be 6, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x42 { // A should be unchanged
					t.Errorf("Expected A to be 0x42, got %02X", cpu.Registers[RegA])
				}
				if !cpu.Flags.Z() || !cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected Z and N flags to be set, H and C reset")
				}
			},
		},
		{
			name: "CP A - Compare with self",
			program: []uint8{
				0x3E, 0x42, // LD A,0x42
				0xBF, // CP A
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 3 {
					t.Errorf("Expected PC to be 3, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x42 { // A should be unchanged
					t.Errorf("Expected A to be 0x42, got %02X", cpu.Registers[RegA])
				}
				if !cpu.Flags.Z() || !cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected Z and N flags to be set, H and C reset")
				}
			},
		},
		{
			name: "CP u8 - Immediate value comparison",
			program: []uint8{
				0x3E, 0x42, // LD A,0x42
				0xFE, 0x42, // CP 0x42
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 4 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x42 { // A should be unchanged
					t.Errorf("Expected A to be 0x42, got %02X", cpu.Registers[RegA])
				}
				if !cpu.Flags.Z() || !cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected Z and N flags to be set, H and C reset")
				}
			},
		},
		{
			name: "CP - Half carry test",
			program: []uint8{
				0x3E, 0x20, // LD A,0x20
				0xFE, 0x11, // CP 0x11
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 4 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x20 { // A should be unchanged
					t.Errorf("Expected A to be 0x20, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || !cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected N and H flags to be set, Z and C reset")
				}
			},
		},
		{
			name: "CP - Zero and carry test",
			program: []uint8{
				0x3E, 0x10, // LD A,0x10
				0xFE, 0x20, // CP 0x20
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 4 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x10 { // A should be unchanged
					t.Errorf("Expected A to be 0x10, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || !cpu.Flags.N() || cpu.Flags.H() || !cpu.Flags.C() {
					t.Error("Expected N and C flags to be set, Z and H reset")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			RunTestProgram(t, tc)
		})
	}
}
func TestANDOpcodes(t *testing.T) {
	testCases := []TestProgram{
		{
			name: "AND B - Basic AND operation",
			program: []uint8{
				0x3E, 0xF0, // LD A,0xF0
				0x06, 0x0F, // LD B,0x0F
				0xA0, // AND B
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x00 {
					t.Errorf("Expected A to be 0x00, got %02X", cpu.Registers[RegA])
				}
				if !cpu.Flags.Z() || cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected Z and H flags to be set, N and C reset")
				}
			},
		},
		{
			name: "AND C - Non-zero result",
			program: []uint8{
				0x3E, 0xFF, // LD A,0xFF
				0x0E, 0x0F, // LD C,0x0F
				0xA1, // AND C
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x0F {
					t.Errorf("Expected A to be 0x0F, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected H flag to be set, others reset")
				}
			},
		},
		{
			name: "AND D - All ones",
			program: []uint8{
				0x3E, 0xFF, // LD A,0xFF
				0x16, 0xFF, // LD D,0xFF
				0xA2, // AND D
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0xFF {
					t.Errorf("Expected A to be 0xFF, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected H flag to be set, others reset")
				}
			},
		},
		{
			name: "AND (HL) - Memory operation",
			program: []uint8{
				0x21, 0x00, 0x80, // LD HL,0x8000
				0x3E, 0xFF, // LD A,0xFF
				0xA6, // AND (HL)
			},
			setup: func(cpu *CPU) {
				cpu.Memory[0x8000] = 0x0F // Set value at (HL)
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 6 {
					t.Errorf("Expected PC to be 6, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x0F {
					t.Errorf("Expected A to be 0x0F, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected H flag to be set, others reset")
				}
			},
		},
		{
			name: "AND A - Self operation",
			program: []uint8{
				0x3E, 0xFF, // LD A,0xFF
				0xA7, // AND A
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 3 {
					t.Errorf("Expected PC to be 3, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0xFF {
					t.Errorf("Expected A to be 0xFF, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected H flag to be set, others reset")
				}
			},
		},
		{
			name: "AND u8 - Immediate value",
			program: []uint8{
				0x3E, 0xFF, // LD A,0xFF
				0xE6, 0x0F, // AND 0x0F
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 4 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x0F {
					t.Errorf("Expected A to be 0x0F, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected H flag to be set, others reset")
				}
			},
		},
		{
			name: "AND - Zero result",
			program: []uint8{
				0x3E, 0xF0, // LD A,0xF0
				0xE6, 0x0F, // AND 0x0F
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 4 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x00 {
					t.Errorf("Expected A to be 0x00, got %02X", cpu.Registers[RegA])
				}
				if !cpu.Flags.Z() || cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected Z and H flags to be set, N and C reset")
				}
			},
		},
		{
			name: "AND - Alternating bits",
			program: []uint8{
				0x3E, 0xAA, // LD A,0xAA
				0xE6, 0x55, // AND 0x55
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 4 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x00 {
					t.Errorf("Expected A to be 0x00, got %02X", cpu.Registers[RegA])
				}
				if !cpu.Flags.Z() || cpu.Flags.N() || !cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected Z and H flags to be set, N and C reset")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			RunTestProgram(t, tc)
		})
	}
}
func TestXOROpcodes(t *testing.T) {
	testCases := []TestProgram{
		{
			name: "XOR B - Basic XOR operation",
			program: []uint8{
				0x3E, 0xF0, // LD A,0xF0
				0x06, 0xFF, // LD B,0xFF
				0xA8, // XOR B
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x0F {
					t.Errorf("Expected A to be 0x0F, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected all flags to be reset")
				}
			},
		},
		{
			name: "XOR C - Zero result",
			program: []uint8{
				0x3E, 0xFF, // LD A,0xFF
				0x0E, 0xFF, // LD C,0xFF
				0xA9, // XOR C
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x00 {
					t.Errorf("Expected A to be 0x00, got %02X", cpu.Registers[RegA])
				}
				if !cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected Z flag to be set, others reset")
				}
			},
		},
		{
			name: "XOR (HL) - Memory operation",
			program: []uint8{
				0x21, 0x00, 0x80, // LD HL,0x8000
				0x3E, 0xFF, // LD A,0xFF
				0xAE, // XOR (HL)
			},
			setup: func(cpu *CPU) {
				cpu.Memory[0x8000] = 0x0F // Set value at (HL)
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 6 {
					t.Errorf("Expected PC to be 6, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0xF0 {
					t.Errorf("Expected A to be 0xF0, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected all flags to be reset")
				}
			},
		},
		{
			name: "XOR A - Self operation (zero)",
			program: []uint8{
				0x3E, 0xFF, // LD A,0xFF
				0xAF, // XOR A
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 3 {
					t.Errorf("Expected PC to be 3, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x00 {
					t.Errorf("Expected A to be 0x00, got %02X", cpu.Registers[RegA])
				}
				if !cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected Z flag to be set, others reset")
				}
			},
		},
		{
			name: "XOR u8 - Immediate value",
			program: []uint8{
				0x3E, 0xFF, // LD A,0xFF
				0xEE, 0x0F, // XOR 0x0F
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 4 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0xF0 {
					t.Errorf("Expected A to be 0xF0, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected all flags to be reset")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			RunTestProgram(t, tc)
		})
	}
}

func TestOROpcodes(t *testing.T) {
	testCases := []TestProgram{
		{
			name: "OR B - Basic OR operation",
			program: []uint8{
				0x3E, 0xF0, // LD A,0xF0
				0x06, 0x0F, // LD B,0x0F
				0xB0, // OR B
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0xFF {
					t.Errorf("Expected A to be 0xFF, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected all flags to be reset")
				}
			},
		},
		{
			name: "OR C - Zero input",
			program: []uint8{
				0x3E, 0x00, // LD A,0x00
				0x0E, 0x00, // LD C,0x00
				0xB1, // OR C
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 5 {
					t.Errorf("Expected PC to be 5, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x00 {
					t.Errorf("Expected A to be 0x00, got %02X", cpu.Registers[RegA])
				}
				if !cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected Z flag to be set, others reset")
				}
			},
		},
		{
			name: "OR (HL) - Memory operation",
			program: []uint8{
				0x21, 0x00, 0x80, // LD HL,0x8000
				0x3E, 0xF0, // LD A,0xF0
				0xB6, // OR (HL)
			},
			setup: func(cpu *CPU) {
				cpu.Memory[0x8000] = 0x0F // Set value at (HL)
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 6 {
					t.Errorf("Expected PC to be 6, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0xFF {
					t.Errorf("Expected A to be 0xFF, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected all flags to be reset")
				}
			},
		},
		{
			name: "OR A - Self operation (no change)",
			program: []uint8{
				0x3E, 0xFF, // LD A,0xFF
				0xB7, // OR A
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 3 {
					t.Errorf("Expected PC to be 3, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0xFF {
					t.Errorf("Expected A to be 0xFF, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected all flags to be reset")
				}
			},
		},
		{
			name: "OR u8 - Immediate value",
			program: []uint8{
				0x3E, 0xF0, // LD A,0xF0
				0xF6, 0x0F, // OR 0x0F
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 4 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0xFF {
					t.Errorf("Expected A to be 0xFF, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected all flags to be reset")
				}
			},
		},
		{
			name: "OR - Alternating bits",
			program: []uint8{
				0x3E, 0xAA, // LD A,0xAA
				0xF6, 0x55, // OR 0x55
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 4 {
					t.Errorf("Expected PC to be 4, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0xFF {
					t.Errorf("Expected A to be 0xFF, got %02X", cpu.Registers[RegA])
				}
				if cpu.Flags.Z() || cpu.Flags.N() || cpu.Flags.H() || cpu.Flags.C() {
					t.Error("Expected all flags to be reset")
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			RunTestProgram(t, tc)
		})
	}
}
func TestPushPopOpcodes(t *testing.T) {
	testCases := []TestProgram{
		{
			name: "PUSH/POP BC",
			program: []uint8{
				0x01, 0x34, 0x12, // LD BC,0x1234
				0xC5,             // PUSH BC
				0x01, 0x00, 0x00, // LD BC,0x0000
				0xC1, // POP BC
			},
			setup: func(cpu *CPU) {
				cpu.SP = 0xFFFE // Initialize stack pointer
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 8 {
					t.Errorf("Expected PC to be 8, got %d", cpu.PC)
				}
				bc := uint16(cpu.Registers[RegB])<<8 | uint16(cpu.Registers[RegC])
				if bc != 0x1234 {
					t.Errorf("Expected BC to be 0x1234, got %04X", bc)
				}
				if cpu.SP != 0xFFFE {
					t.Errorf("Expected SP to be 0xFFFE, got %04X", cpu.SP)
				}
			},
		},
		{
			name: "PUSH/POP DE",
			program: []uint8{
				0x11, 0x78, 0x56, // LD DE,0x5678
				0xD5,             // PUSH DE
				0x11, 0x00, 0x00, // LD DE,0x0000
				0xD1, // POP DE
			},
			setup: func(cpu *CPU) {
				cpu.SP = 0xFFFE
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 8 {
					t.Errorf("Expected PC to be 8, got %d", cpu.PC)
				}
				de := uint16(cpu.Registers[RegD])<<8 | uint16(cpu.Registers[RegE])
				if de != 0x5678 {
					t.Errorf("Expected DE to be 0x5678, got %04X", de)
				}
				if cpu.SP != 0xFFFE {
					t.Errorf("Expected SP to be 0xFFFE, got %04X", cpu.SP)
				}
			},
		},
		{
			name: "PUSH/POP HL",
			program: []uint8{
				0x21, 0xBC, 0x9A, // LD HL,0x9ABC
				0xE5,             // PUSH HL
				0x21, 0x00, 0x00, // LD HL,0x0000
				0xE1, // POP HL
			},
			setup: func(cpu *CPU) {
				cpu.SP = 0xFFFE
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 8 {
					t.Errorf("Expected PC to be 8, got %d", cpu.PC)
				}
				hl := uint16(cpu.Registers[RegH])<<8 | uint16(cpu.Registers[RegL])
				if hl != 0x9ABC {
					t.Errorf("Expected HL to be 0x9ABC, got %04X", hl)
				}
				if cpu.SP != 0xFFFE {
					t.Errorf("Expected SP to be 0xFFFE, got %04X", cpu.SP)
				}
			},
		},
		{
			name: "PUSH/POP AF",
			program: []uint8{
				0x3E, 0xF0, // LD A,0xF0
				0x37,       // SCF (Set Carry Flag)
				0x27,       // DAA (Set H Flag)
				0xF5,       // PUSH AF
				0x3E, 0x00, // LD A,0x00
				0x3F, // CCF (Clear Carry Flag)
				0xF1, // POP AF
			},
			setup: func(cpu *CPU) {
				cpu.SP = 0xFFFE
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 9 {
					t.Errorf("Expected PC to be 9, got %d", cpu.PC)
				}
				if cpu.Registers[RegA] != 0x50 {
					t.Errorf("Expected A to be 0x50, got %02X", cpu.Registers[RegA])
				}
				if !cpu.Flags.C() {
					t.Error("Expected Carry flag to be set")
				}
				if cpu.SP != 0xFFFE {
					t.Errorf("Expected SP to be 0xFFFE, got %04X", cpu.SP)
				}
			},
		},
		{
			name: "Multiple PUSH/POP sequence",
			program: []uint8{
				0x01, 0x34, 0x12, // LD BC,0x1234
				0x11, 0x78, 0x56, // LD DE,0x5678
				0x21, 0xBC, 0x9A, // LD HL,0x9ABC
				0xC5,             // PUSH BC
				0xD5,             // PUSH DE
				0xE5,             // PUSH HL
				0x01, 0x00, 0x00, // LD BC,0x0000
				0x11, 0x00, 0x00, // LD DE,0x0000
				0x21, 0x00, 0x00, // LD HL,0x0000
				0xE1, // POP HL
				0xD1, // POP DE
				0xC1, // POP BC
			},
			setup: func(cpu *CPU) {
				cpu.SP = 0xFFFE
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 24 {
					t.Errorf("Expected PC to be 24, got %d", cpu.PC)
				}
				bc := uint16(cpu.Registers[RegB])<<8 | uint16(cpu.Registers[RegC])
				if bc != 0x1234 {
					t.Errorf("Expected BC to be 0x1234, got %04X", bc)
				}
				de := uint16(cpu.Registers[RegD])<<8 | uint16(cpu.Registers[RegE])
				if de != 0x5678 {
					t.Errorf("Expected DE to be 0x5678, got %04X", de)
				}
				hl := uint16(cpu.Registers[RegH])<<8 | uint16(cpu.Registers[RegL])
				if hl != 0x9ABC {
					t.Errorf("Expected HL to be 0x9ABC, got %04X", hl)
				}
				if cpu.SP != 0xFFFE {
					t.Errorf("Expected SP to be 0xFFFE, got %04X", cpu.SP)
				}
			},
		},
		{
			name: "Stack Pointer Manipulation",
			program: []uint8{
				0x01, 0x34, 0x12, // LD BC,0x1234
				0x11, 0x78, 0x56, // LD DE,0x5678
				0xC5, // PUSH BC
				0xD5, // PUSH DE
				0xD1, // POP DE
				0xC1, // POP BC
			},
			setup: func(cpu *CPU) {
				cpu.SP = 0xFFFE
			},
			validate: func(t *testing.T, cpu *CPU) {
				if cpu.PC != 10 {
					t.Errorf("Expected PC to be 10, got %d", cpu.PC)
				}
				// Check if stack pointer returned to original position
				if cpu.SP != 0xFFFE {
					t.Errorf("Expected SP to be 0xFFFE, got %04X", cpu.SP)
				}
				// Verify memory contents were correctly written and read
				if cpu.Memory[0xFFFC] != 0x34 || cpu.Memory[0xFFFD] != 0x12 {
					t.Errorf("Expected memory at 0xFFFC-0xFFFD to be 0x1234, got %02X%02X",
						cpu.Memory[0xFFFD], cpu.Memory[0xFFFC])
				}
				if cpu.Memory[0xFFFA] != 0x78 || cpu.Memory[0xFFFB] != 0x56 {
					t.Errorf("Expected memory at 0xFFFA-0xFFFB to be 0x5678, got %02X%02X",
						cpu.Memory[0xFFFB], cpu.Memory[0xFFFA])
				}
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			RunTestProgram(t, tc)
		})
	}
}
