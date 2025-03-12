package gopherboy

import (
	rl "github.com/gen2brain/raylib-go/raylib"
)

func bgTileMapMode(cpu *CPU) uint8 {
	byte := cpu.Memory[0xFF40]
	result := ((byte & 0b00001000) >> 3) & 0b00001
	return result
}

func bgTileDataMode(cpu *CPU) uint8 {
	byte := cpu.Memory[0xFF40]
	result := ((byte & 0b00010000) >> 4)
	return result
}

func (cpu *CPU) BuildFrame() [][]uint32 {
	ly := cpu.Memory[0xFF44]

	fb := make([][]uint32, 144)
	for i := range fb {
		fb[i] = make([]uint32, 160)
	}

	var bgTileDataModeAddr uint16
	if bgTileDataMode(cpu) == 1 {
		bgTileDataModeAddr = 0x8000
	} else {
		bgTileDataModeAddr = 0x9000
	}
	var tileIndexAddr uint16
	if bgTileMapMode(cpu) == 0 {
		tileIndexAddr = 0x9800
	} else {
		tileIndexAddr = 0x9C00
	}

	for x := 0; x < 160; x++ {
		var addr uint16
		tileX := uint8(x / 8)
		tileY := uint8(ly / 8)
		tilePixelX := uint8(x % 8)
		tilePixelY := uint8(ly % 8)

		tileId := cpu.Memory[tileIndexAddr+uint16(tileY*32+tileX)]
		addr = bgTileDataModeAddr + uint16(tileId*16) + uint16(tilePixelY*2)

		pixel := interleaveTilePixel(cpu.Memory[addr], cpu.Memory[addr+1], 7-tilePixelX)
		fb[ly][x] = colourizePixel(int(pixel))
	}

	return fb
}
func colourizePixel(input int) uint32 {
	// The input is expected to be a value between 0 and 3
	// where 0 is white and 3 is black in the Game Boy's 2-bit color space.
	switch input {
	case 0: // White
		return 0xFFFFFFFF // ARGB for white
	case 1: // Light gray
		return 0xFFAAAAAA // ARGB for light gray
	case 2: // Dark gray
		return 0xFF555555 // ARGB for dark gray
	case 3: // Black
		return 0xFF000000 // ARGB for black
	default:
		// If the input is out of range, return a noticeable color (e.g., red)
		// to indicate an error or unexpected value.
		return 0xFFFF0000
	}
}
func interleaveTilePixel(low, high, index uint8) uint16 {
	result := uint16(((high>>index)&0x1)<<1) + uint16((low>>index)&0x1)
	return result
}

func buildFb(cpu *CPU, ly uint8) {
	bgTileMapModeAddr := bgTileMapMode(cpu)

	var tileIndexAddr uint16
	if bgTileMapModeAddr == 0 {
		tileIndexAddr = 0x9800
	} else {
		tileIndexAddr = 0x9C00
	}

	var (
		addr        uint16
		pixel       uint16
		colourPixel uint32
		tileX       uint8
		tileY       uint8
		tilePixelX  uint8
		tilePixelY  uint8
		tileIndex   uint16
		tileID      uint8
	)

	for x := uint8(0); x < 160; x++ {
		tileX = x / 8
		tileY = ly / 8
		tilePixelX = x % 8
		tilePixelY = ly % 8
		tileIndex = uint16(tileY)*32 + uint16(tileX)
		tileID = cpu.Memory[tileIndexAddr+tileIndex]

		if bgTileDataMode(cpu) == 1 {
			addr = 0x8000 + uint16(tileID)*16 + uint16(tilePixelY)*2
		} else {
			if tileID > 127 {
				addr = 0x8800
			} else {
				addr = 0x9000
			}
			addr += uint16(tileID)*16 + uint16(tilePixelY)*2
		}

		pixel = interleaveTilePixel(cpu.Memory[addr], cpu.Memory[addr+1], 7-tilePixelX)
		colourPixel = colourizePixel(int(pixel))
		cpu.Framebuffer[ly][x] = colourPixel
	}
}

func (cpu *CPU) RenderGameBoy() {
	// Approach 1: Create image data directly from framebuffer
	var imageData []rl.Color = make([]rl.Color, 160*144)

	for ly := uint8(0); ly < 144; ly++ {
		buildFb(cpu, ly)
		for x := uint8(0); x < 160; x++ {
			index := int(ly)*160 + int(x)
			colorValue := cpu.Framebuffer[ly][x]

			// Convert uint32 ARGB to Raylib Color
			imageData[index] = rl.Color{
				R: uint8((colorValue >> 16) & 0xFF),
				G: uint8((colorValue >> 8) & 0xFF),
				B: uint8(colorValue & 0xFF),
				A: uint8((colorValue >> 24) & 0xFF),
			}
		}
	}

	// Update the texture directly with the new color data
	rl.UpdateTexture(cpu.Texture, imageData)
}
