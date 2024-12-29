package main

import "log"

type CPU struct {
	Registers []uint8
	Clock     uint8
	PC        uint16
	SP        uint16
	IME       uint16

	Memory []uint8
	ROM    []uint8
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

// Flag register bit positions
const (
	FlagZ = 7 // Zero flag
	FlagN = 6 // Subtract flag
	FlagH = 5 // Half carry flag
	FlagC = 4 // Carry flag
)

func InitCPU() *CPU {
	result := CPU{
		Memory:    make([]uint8, 65535),
		ROM:       make([]uint8, 32768),
		Registers: make([]uint8, 8),
	}
	return &result
}

func main() {
	log.Printf("hello world")
}
