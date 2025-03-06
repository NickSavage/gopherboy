package main

import (
	"flag"
	"fmt"
	"log"
	"os"
)

type CPU struct {
	Registers []uint8
	Clock     uint8
	PC        uint16
	SP        uint16
	IME       uint16
	Flags     *Flags
	MaxCycles int // for testing

	Memory []uint8
	ROM    []uint8
	Halted bool
}

type Flags struct {
	value byte
	CPU   *CPU
}

const (
	FlagZ byte = 1 << 7 // Zero flag (Bit 7)
	FlagN byte = 1 << 6 // Subtract flag (Bit 6)
	FlagH byte = 1 << 5 // Half Carry flag (Bit 5)
	FlagC byte = 1 << 4 // Carry flag (Bit 4)
	// Bits 0-3 are unused and always 0
)

// Methods to get flag values
func (f *Flags) Z() bool { return f.value&FlagZ != 0 }
func (f *Flags) N() bool { return f.value&FlagN != 0 }
func (f *Flags) H() bool { return f.value&FlagH != 0 }
func (f *Flags) C() bool { return f.value&FlagC != 0 }

// Methods to set flag values
func (f *Flags) SetZ(value bool) { f.setBit(FlagZ, value) }
func (f *Flags) SetN(value bool) { f.setBit(FlagN, value) }
func (f *Flags) SetH(value bool) { f.setBit(FlagH, value) }
func (f *Flags) SetC(value bool) { f.setBit(FlagC, value) }

// Helper method for setting bits
func (f *Flags) setBit(bit byte, value bool) {
	if value {
		f.value |= bit
	} else {
		f.value &= ^bit
	}
	f.CPU.Registers[RegF] = f.value
}

// Get the raw byte value
func (f *Flags) Value() byte {
	return f.value
}

// Set the raw byte value
func (f *Flags) SetValue(value byte) {
	f.value = value & 0xF0 // Only upper 4 bits are used
	f.CPU.Registers[RegF] = f.value
}

// 8-bit register constants
const (
	RegA = iota // Accumulator
	RegF        // Flags
	RegB        // General purpose
	RegC        // General purpose
	RegD        // General purpose
	RegE        // General purpose
	RegH        // General purpose
	RegL        // General purpose
)

// 16-bit register pair constants
const (
	RegAF = iota // Accumulator & Flags
	RegBC        // BC pair
	RegDE        // DE pair
	RegHL        // HL pair
	RegSP        // Stack Pointer
	RegPC        // Program Counter
)

func InitCPU() *CPU {
	result := CPU{
		Memory:    make([]uint8, 65535),
		ROM:       make([]uint8, 32768),
		Registers: make([]uint8, 8),
		Halted:    false,
		SP:        0xFFFE,
		Flags:     &Flags{},
		PC:        0x0150,
	}
	result.Flags.CPU = &result
	result.Memory[0xFF43] = 0
	result.Memory[0xFF44] = 144
	return &result
}

// LoadROM loads a ROM file into the CPU's memory
func LoadROM(cpu *CPU, romFilePath string) error {
	romData, err := os.ReadFile(romFilePath)
	if err != nil {
		return fmt.Errorf("error reading ROM file: %v", err)
	}

	// Make sure we don't overflow the ROM buffer
	if len(romData) > len(cpu.ROM) {
		return fmt.Errorf("ROM file too large: %d bytes (max %d)", len(romData), len(cpu.ROM))
	}

	// Copy ROM data to CPU's memory
	copy(cpu.ROM, romData)

	// Also copy to memory starting at address 0
	copy(cpu.Memory, romData)

	return nil
}

// RunProgram executes the program loaded in the CPU's memory
func RunProgram(cpu *CPU, maxCycles int) {
	for i := 0; i < maxCycles && !cpu.Halted; i++ {
		cpu.ParseNextOpcode()
		cpu.Memory[0xFF44]++
		if cpu.Memory[0xFF44] == 154 {
			cpu.Memory[0xFF44] = 0
		}
		// Optional: Add delay or debug output here
		if i%1000 == 0 {
			log.Printf("Executed %d instructions, PC: 0x%04X", i, cpu.PC)
		}
	}

	log.Printf("Program execution stopped. PC: 0x%04X, Halted: %v", cpu.PC, cpu.Halted)

	// Dump memory contents to file
	if err := DumpMemoryToFile(cpu, "memory_dump.bin"); err != nil {
		log.Printf("Failed to dump memory: %v", err)
	} else {
		log.Printf("Memory dumped to memory_dump.hex")
	}
}

// DumpMemoryToFile writes the CPU memory contents to a binary file
func DumpMemoryToFile(cpu *CPU, filename string) error {
	// Open file for writing in binary mode
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("failed to create file: %v", err)
	}
	defer file.Close()

	// Write the entire memory buffer directly to the file
	_, err = file.Write(cpu.Memory)
	if err != nil {
		return fmt.Errorf("failed to write memory dump: %v", err)
	}

	return nil
}

func main() {
	// Parse command-line flags
	romFile := flag.String("rom", "", "Path to Game Boy ROM file")
	maxCycles := flag.Int("cycles", 2400, "Maximum number of CPU cycles to execute")
	debug := flag.Bool("debug", false, "Enable debug output")

	flag.Parse()

	// Check if ROM file was provided
	if *romFile == "" {
		log.Fatal("No ROM file specified. Use -rom flag to specify a Game Boy ROM file.")
	}

	// Initialize CPU
	cpu := InitCPU()

	// Load ROM file
	log.Printf("Loading ROM file: %s", *romFile)
	if err := LoadROM(cpu, *romFile); err != nil {
		log.Fatalf("Failed to load ROM: %v", err)
	}

	// Set debug level if needed
	if *debug {
		log.Printf("Debug mode enabled")
		// You can add more detailed debug setup here
	}

	// Run the program
	log.Printf("Starting program execution with max %d cycles", *maxCycles)
	RunProgram(cpu, *maxCycles)

	log.Printf("Emulation complete")
}
