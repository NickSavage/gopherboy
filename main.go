package main

import "log"

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
	}
	result.Flags.CPU = &result
	return &result
}

func main() {
	log.Printf("hello world")
}
