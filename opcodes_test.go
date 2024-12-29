package main

import (
	"testing"
)

func TestLoadMemoryOperations(t *testing.T) {
	cpu := InitCPU()
	address := uint16(0x8000)

	// Test LD (HL),B
	cpu.Registers[RegB] = 0x42
	cpu.LoadMemory(address, RegB)

	if cpu.Memory[address] != 0x42 {
		t.Errorf("LoadMemory failed: expected memory[0x8000]=0x42, got 0x%02X", cpu.Memory[address])
	}

	// Test LD C,(HL)
	cpu.LoadFromMemory(RegC, address)

	if cpu.Registers[RegC] != 0x42 {
		t.Errorf("LoadFromMemory failed: expected C=0x42, got 0x%02X", cpu.Registers[RegC])
	}
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
	testCases := []struct {
		name     string
		program  []uint8
		setup    func(*CPU)
		validate func(*testing.T, *CPU)
	}{
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

			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
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
			}

			// Validate the final state
			tc.validate(t, cpu)
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

func TestLoadFromMemory(t *testing.T) {
	testCases := []struct {
		name     string
		reg      uint8
		opcode   uint8
		memValue uint8
	}{
		{"LD B,(HL)", RegB, 0x46, 0x42},
		{"LD C,(HL)", RegC, 0x4E, 0x43},
		{"LD D,(HL)", RegD, 0x56, 0x44},
		{"LD E,(HL)", RegE, 0x5E, 0x45},
		{"LD H,(HL)", RegH, 0x66, 0x46},
		{"LD L,(HL)", RegL, 0x6E, 0x47},
		{"LD A,(HL)", RegA, 0x7E, 0x48},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			cpu := InitCPU()

			// Setup HL to point to test memory location
			address := uint16(0x8000)
			cpu.Registers[RegH] = uint8(address >> 8)
			cpu.Registers[RegL] = uint8(address & 0xFF)
			cpu.Memory[address] = tc.memValue

			// Test via opcode
			cpu.ROM[0] = tc.opcode
			cpu.ParseNextOpcode()

			if cpu.Registers[tc.reg] != tc.memValue {
				t.Errorf("%s failed, expected 0x%02X, got 0x%02X",
					tc.name, tc.memValue, cpu.Registers[tc.reg])
			}

			if cpu.PC != 1 {
				t.Errorf("%s PC increment failed, expected 1, got %d", tc.name, cpu.PC)
			}
		})
	}
}
