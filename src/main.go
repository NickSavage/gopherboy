package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/veandco/go-sdl2/sdl"
)

type CPU struct {
	Registers     []uint8
	Clock         uint16
	PC            uint16
	SP            uint16
	IME           uint16
	Flags         *Flags
	MaxCycles     int // for testing
	DMAActive     bool
	DMACycles     int
	DMASourceBase uint16

	Memory []uint8
	ROM    []uint8
	Halted bool

	Framebuffer [][]uint32
	Window      *sdl.Window
	Renderer    *sdl.Renderer
	Texture     *sdl.Texture
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
		Memory:    make([]uint8, 65536),
		ROM:       make([]uint8, 32768),
		Registers: make([]uint8, 8),
		Halted:    false,
		SP:        0xFFFE,
		Flags:     &Flags{},
		PC:        0x0000,
	}
	result.Flags.CPU = &result
	result.Memory[0xFF43] = 0
	result.Memory[0xFF44] = 0xFF
	return &result
}

func (cpu *CPU) CheckError() error {
	// Check if registers B through L all contain 0x42
	if cpu.Registers[RegB] == 0x42 &&
		cpu.Registers[RegC] == 0x42 &&
		cpu.Registers[RegD] == 0x42 &&
		cpu.Registers[RegE] == 0x42 &&
		cpu.Registers[RegH] == 0x42 &&
		cpu.Registers[RegL] == 0x42 {

		// Check if current opcode is LD B, B (0x40) or we're in an infinite JR loop (0x18 0x00)
		opcode := cpu.Memory[cpu.PC]
		nextByte := cpu.Memory[cpu.PC+1]

		// Check for the infinite JR loop (JR 0 - jump to self)
		if opcode == 0x18 && nextByte == 0x00 {
			return fmt.Errorf("test failure detected: infinite JR loop after setting registers to 0x42")
		}

		// Check for LD B, B opcode
		if opcode == 0x40 {
			// Check serial port for evidence of sending 0x42 six times
			// In a real implementation, you would need to track serial port usage
			// This is a simplified version that just checks if a specific memory pattern exists

			// For GameBoy, 0xFF01 is the Serial Transfer Data Register (SB)
			// 0xFF02 is the Serial Transfer Control Register (SC)
			// We can check if these have been used to send 0x42 repeatedly

			// Since we don't have a proper way to track the 6 serial transfers in this code,
			// we'll just check if the current SB register contains 0x42
			if cpu.Memory[0xFF01] == 0x42 {
				return fmt.Errorf("test failure detected: registers set to 0x42 and serial port used with 0x42")
			}
		}
	}

	return nil // No error detected
}

func LoadBoot(cpu *CPU, bootFilePath string) error {
	bootData, err := os.ReadFile(bootFilePath)
	if err != nil {
		return fmt.Errorf("error reading boot file: %v", err)
	}
	copy(cpu.Memory[0x0000:0x0000+len(bootData)], bootData)
	return nil
}

// LoadROM loads a ROM file into the CPU's memory
func LoadROM(cpu *CPU, romFilePath string) error {
	romData, err := os.ReadFile(romFilePath)
	if err != nil {
		return fmt.Errorf("error reading ROM file: %v", err)
	}

	if len(romData) > len(cpu.ROM) {
		return fmt.Errorf("ROM file too large: %d bytes (max %d)", len(romData), len(cpu.ROM))
	}

	copy(cpu.ROM, romData)

	copy(cpu.Memory, romData)

	return nil
}

func (cpu *CPU) Exit() {
	if cpu.Texture != nil {
		cpu.Texture.Destroy()
	}
	if cpu.Renderer != nil {
		cpu.Renderer.Destroy()
	}
	if cpu.Window != nil {
		cpu.Window.Destroy()
	}
	sdl.Quit()

	// Dump memory contents to file
	if err := DumpMemoryToFile(cpu, "memory_dump.bin"); err != nil {
		log.Printf("Failed to dump memory: %v", err)
	} else {
		log.Printf("Memory dumped to memory_dump.bin")
	}
	os.Exit(0)
}

func (cpu *CPU) HandleKeyboard() {
	for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
		switch event.(type) {
		case *sdl.QuitEvent:
			cpu.Exit()
		case *sdl.KeyboardEvent:
			keyEvent := event.(*sdl.KeyboardEvent)
			if keyEvent.Type == sdl.KEYDOWN {
				if keyEvent.Keysym.Sym == sdl.K_ESCAPE {
					cpu.Exit()
				}
			}
		}
	}
}

func (cpu *CPU) RequestVBlank() {
	cpu.Memory[0xFF0F] |= 1 << 0
}

func (cpu *CPU) RequestStatInterrupt() {
	cpu.Memory[0xFF0F] |= 1 << 1
}

func (cpu *CPU) HandleInterrupts() {
	cpu.Halted = false
	// vblank
	// Check each interrupt type

	ie := cpu.Memory[0xFFFF]

	if_ := cpu.Memory[0xFF0F]

	// VBlank
	if (ie&0x01 != 0) && (if_&0x01 != 0) {
		cpu.IME = 0 // Disable interrupts
		// Clear the interrupt flag
		cpu.Memory[0xFF0F] &= ^uint8(0x01)
		// Push current PC to stack
		// Jump to interrupt handler
		high := uint8(cpu.PC >> 8)
		low := uint8(cpu.PC & 0xFF)
		cpu.SP--
		cpu.Memory[cpu.SP] = high
		cpu.SP--
		cpu.Memory[cpu.SP] = low
		cpu.PC = 0x0040
		return
	}

	// STAT
	if (ie&0x02 != 0) && (if_&0x02 != 0) {
		cpu.IME = 0
		cpu.Memory[0xFF0F] &= ^uint8(0x02)

		high := uint8(cpu.PC >> 8)
		low := uint8(cpu.PC & 0xFF)
		cpu.SP--
		cpu.Memory[cpu.SP] = high
		cpu.SP--
		cpu.Memory[cpu.SP] = low
		cpu.PC = 0x0048
		return
	}
}

// RunProgram executes the program loaded in the CPU's memory
func RunProgram(cpu *CPU, maxCycles int) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		log.Fatalf("Failed to initialize SDL: %v", err)
	}

	window, err := sdl.CreateWindow("Gopherboy", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED, 160*4, 144*4, sdl.WINDOW_SHOWN)
	if err != nil {
		log.Fatalf("Failed to create window: %v", err)
	}
	cpu.Window = window

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		log.Fatalf("Failed to create renderer: %v", err)
	}
	cpu.Renderer = renderer

	texture, err := renderer.CreateTexture(sdl.PIXELFORMAT_RGBA8888, sdl.TEXTUREACCESS_STREAMING, 160, 144)
	if err != nil {
		log.Fatalf("Failed to create texture: %v", err)
	}
	cpu.Texture = texture

	// Set up framebuffer
	cpu.Framebuffer = make([][]uint32, 144)
	for i := range cpu.Framebuffer {
		cpu.Framebuffer[i] = make([]uint32, 160)
	}

	// Set the logical size to maintain aspect ratio
	renderer.SetLogicalSize(160, 144)

	start := time.Now()

	for i := 0; i < maxCycles; i++ {
		// cycleStart := time.Now()

		cpu.HandleKeyboard()
		cpu.HandleInterrupts()
		if !cpu.Halted {
			cpu.ParseNextOpcode()
		}

		if cpu.DMAActive {
			if cpu.DMACycles > 0 {
				cpu.Memory[0xFE00+uint16(160-cpu.DMACycles)] = cpu.Memory[cpu.DMASourceBase+uint16(160-cpu.DMACycles)]
				cpu.DMACycles--
			} else {
				cpu.DMAActive = false
			}
		}

		// super cludge, just want to make sure there is a little delay
		// otherwise it loops at a constant rate
		if i%100 != 0 {
			cpu.Memory[0xFF44]++
		}
		if cpu.Memory[0xFF44] == 144 {
			cpu.RequestVBlank()
		}
		// if cpu.Memory[0xFF44] == cpu.Memory[0xFF45] {
		// 	cpu.RequestStatInterrupt()
		// }
		if cpu.Memory[0xFF44] == 154 {
			cpu.RenderGameBoy()

			// Clear the renderer
			cpu.Renderer.Clear()

			// Copy the texture to the renderer
			cpu.Renderer.Copy(cpu.Texture, nil, nil)

			// Present the renderer
			cpu.Renderer.Present()

			cpu.Clock = cpu.Clock % 114
			cpu.Memory[0xFF44] = 0
		}
		// elapsed := time.Since(start)
		// cycleTime := time.Since(cycleStart)
		// avgCycleTime := elapsed / time.Duration(i+1)
		// log.Printf("Executed %d instructions, PC: 0x%04X, clock: %d, Avg cycle: %v, Current cycle: %v",
		// 	i, cpu.PC, cpu.Clock, avgCycleTime, cycleTime)
		err := cpu.CheckError()
		if err != nil {
			log.Printf("Test has failed")
			break
		}

		if cpu.Clock >= 456 {
		}

	}

	totalTime := time.Since(start)
	avgCycleTime := totalTime / time.Duration(maxCycles)
	log.Printf("Program execution stopped. PC: 0x%04X, Halted: %v", cpu.PC, cpu.Halted)
	log.Printf("Total execution time: %v, Average cycle time: %v", totalTime, avgCycleTime)

	// Dump memory contents to file
	if err := DumpMemoryToFile(cpu, "memory_dump.bin"); err != nil {
		log.Printf("Failed to dump memory: %v", err)
	} else {
		log.Printf("Memory dumped to memory_dump.hex")
	}

	if cpu.Texture != nil {
		cpu.Texture.Destroy()
	}
	if cpu.Renderer != nil {
		cpu.Renderer.Destroy()
	}
	if cpu.Window != nil {
		cpu.Window.Destroy()
	}
	sdl.Quit()
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
	maxCycles := flag.Int("cycles", 5000000, "Maximum number of CPU cycles to execute")
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
	if err := LoadBoot(cpu, "boot.gb"); err != nil {
		log.Fatalf("Failed to load boot ROM: %v", err)
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
