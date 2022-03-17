package gomu8080

import (
	"image"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	processor *Processor
	mmu       *MMU
	width     int
	height    int

	// external hardware
	ShiftReg       uint16
	ShiftRegOffset uint8

	// flipflop for display
	IsFullDraw bool
	LastDraw   int64

	// Game specific data io

	// Dips switch
	dip4 bool // self-test-request read at power up
	dip3 bool // 00 = 3 ships  10 = 5 ships
	dip5 bool // 01 = 4 ships  11 = 6 ships
	dip6 bool // extra ship at 1500, 1 = extra ship at 1000
	dip7 bool // Coin info displayed in demo screen 0=ON
}

func NewGame(pcs *Processor, mmu *MMU, width int, height int) *Game {
	game := Game{}
	game.processor = pcs
	game.mmu = mmu
	game.width = width
	game.height = height

	game.dip4 = true

	return &game
}

func (g *Game) Update() error {

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	instructionPerFrame := 1000

	img := image.NewRGBA(image.Rect(0, 0, g.height, g.width))

	for i := 0; i < instructionPerFrame/2; i++ {
		g.Process()
	}

	if g.processor.IsInteruptsEnabled {
		g.processor.IsInteruptsEnabled = false

		g.processor.SP -= 2
		g.mmu.Memory[g.processor.SP] = byte(g.processor.PC & 0xFF)
		g.mmu.Memory[g.processor.SP+1] = byte(g.processor.PC >> 8)

		g.processor.PC = 0x0008
	}
	g.render(img, true)

	for i := 0; i < instructionPerFrame/2; i++ {
		g.Process()
	}

	if g.processor.IsInteruptsEnabled {
		g.processor.IsInteruptsEnabled = false

		g.processor.SP -= 2
		g.mmu.Memory[g.processor.SP] = byte(g.processor.PC & 0xFF)
		g.mmu.Memory[g.processor.SP+1] = byte(g.processor.PC >> 8)
		g.processor.PC = 0x0010
	}
	g.render(img, false)

	ebiImg := ebiten.NewImage(g.height, g.width)
	ebiImg.ReplacePixels(img.Pix)
	opts := &ebiten.DrawImageOptions{}
	opts.GeoM.Rotate(-math.Pi / 2)
	opts.GeoM.Translate(0, float64(g.height))
	screen.DrawImage(ebiImg, opts)

}

func (g *Game) render(img *image.RGBA, isTop bool) {

	start := 0x2400
	startPix := 0
	if !isTop {
		start = 0x3200
		startPix = 0xE00 * 8
	}
	// 224 * 256 -> 28bytes * 256
	// halve it for separated rendering top and bottom
	for i := 0; i < 14*256; i++ {

		value := g.mmu.Memory[start+i]
		for bit := 0; bit < 8; bit++ {
			color := uint8(0)
			if (value>>uint32(bit))&0x01 > 0 {
				color = uint8(255)
			}
			pos := startPix + i*8 + bit
			img.Pix[4*pos] = color
			img.Pix[4*pos+1] = color
			img.Pix[4*pos+2] = color
			img.Pix[4*pos+3] = color
		}
	}
}

func (g *Game) Process() {
	opcode := g.mmu.Memory[g.processor.PC]
	if opcode == 0xDB { // IN
		port := g.mmu.Memory[g.processor.PC+1]
		// fmt.Printf("IN %02X\n", port)
		switch port {
		case 0:
			data := uint8(0b10001110)

			// bit 0 - self test
			if g.dip4 {
				data |= 0x1
			}
			// bit 4 - fire
			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				data |= (0x1 << 4)
			}
			// bit 5 - left
			if ebiten.IsKeyPressed(ebiten.KeyLeft) {
				data |= (0x1 << 5)
			}
			// bit 6 - right
			if ebiten.IsKeyPressed(ebiten.KeyRight) {
				data |= (0x1 << 6)
			}
			g.processor.A = data

		case 1:
			data := uint8(0b00001000)

			// bit 0 - deposit a credit
			if ebiten.IsKeyPressed(ebiten.KeyEnter) {
				data |= 0x1
			}

			// bit 1 - 2p start
			if ebiten.IsKeyPressed(ebiten.KeyO) {
				data |= (0x1 << 1)
			}
			// bit 2 - 1p start
			if ebiten.IsKeyPressed(ebiten.KeyP) {
				data |= (0x1 << 2)
			}
			// bit 4 - fire (1p)
			if ebiten.IsKeyPressed(ebiten.KeySpace) {
				data |= (0x1 << 4)
			}
			// bit 5 - left (1p)
			if ebiten.IsKeyPressed(ebiten.KeyLeft) {
				data |= (0x1 << 5)
			}
			// bit 6 - right (1p)
			if ebiten.IsKeyPressed(ebiten.KeyRight) {
				data |= (0x1 << 6)
			}
			g.processor.A = data
			// fmt.Printf("%02X\n", g.processor.A)

		case 2:
			data := uint8(0)
			// bit 0 - ship
			if g.dip3 {
				data |= 0x1
			}
			// bit 1 - ship
			if g.dip5 {
				data |= (0x1 << 1)
			}
			// bit 2 - tilt
			if ebiten.IsKeyPressed(ebiten.KeyT) {
				data |= (0x1 << 2)
			}
			// bit 3 - extra ship
			if g.dip6 {
				data |= (0x1 << 3)
			}
			// bit 4 - fire (2p)
			if ebiten.IsKeyPressed(ebiten.KeyW) {
				data |= (0x1 << 4)
			}
			// bit 5 - left (2p)
			if ebiten.IsKeyPressed(ebiten.KeyA) {
				data |= (0x1 << 5)
			}
			// bit 6 - right (2p)
			if ebiten.IsKeyPressed(ebiten.KeyD) {
				data |= (0x1 << 6)
			}
			// bit 7 - coin info displayed on screen
			if g.dip7 {
				data |= (0x1 << 7)
			}
			g.processor.A = data

		case 3:
			// load data from external hardware back to processor
			// for shifted data
			offset := 8 - g.ShiftRegOffset
			g.processor.A = uint8(g.ShiftRegOffset >> offset)
		default:
			panic("unimplemented")
		}
		g.processor.PC += 2
		return
	}
	if opcode == 0xD3 { // OUT
		port := g.mmu.Memory[g.processor.PC+1]
		// fmt.Printf("OUT %02X", port)
		switch port {
		case 2:
			// out from processor to external hardware
			// to set offset for shifting the data
			g.ShiftRegOffset = g.processor.A
		case 3:
		case 4:
			// out from processor to external hardware
			// to shift the data
			g.ShiftReg = (uint16(g.processor.A) << 8) | g.ShiftReg>>8
		case 5:
		case 6:
		default:
			panic("unimplemented")
		}
		g.processor.PC += 2
		return
	}
	// normal case
	g.processor.Run()
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (screenWidth, screenHeight int) {
	return g.width, g.height
}
