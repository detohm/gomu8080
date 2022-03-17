package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/detohm/gomu8080"
	"github.com/hajimehoshi/ebiten/v2"
)

func main() {
	path := flag.String("path", "", "")       // TODO - add more detail
	debugMode := flag.Bool("debug", true, "") // TODO - add more detail
	isSpaceInvader := flag.Bool("spaceinvader", false, "")
	flag.Parse()

	mmu := gomu8080.NewMMU()
	p := gomu8080.NewProcessor(mmu, *debugMode)

	// test/console emulator
	if !*isSpaceInvader {
		bytes, err := os.ReadFile(*path)
		if err != nil {
			fmt.Println(err)
			return
		}

		mmu.Load(len(bytes), bytes, 0x0100)
		p.PC = 0x0100

		for !p.IsHalt {

			if p.DebugMode {
				fmt.Printf("PC=%04X OP=%02X %02X %02X | ",
					p.PC,
					mmu.Memory[p.PC],
					mmu.Memory[p.PC+1],
					mmu.Memory[p.PC+2]) // TODO - bound check
			}

			p.Run()
			if p.DebugMode {
				time.Sleep(5 * time.Millisecond)
			}
		}
		fmt.Println()
		return
	}

	// space invader emulator
	if *isSpaceInvader {

		files := []string{"invaders_h.rom", "invaders_g.rom", "invaders_f.rom", "invaders_e.rom"}
		pos := 0x0000
		for _, file := range files {
			bytes, _ := os.ReadFile(*path + file)
			mmu.Load(len(bytes), bytes, pos)
			pos += 0x0800
		}

		p.PC = 0x0000

		game := gomu8080.NewGame(p, mmu, 224, 256)
		ebiten.SetWindowSize(224*2, 256*2)
		ebiten.SetWindowTitle("Hello, World!")
		ebiten.SetFPSMode(ebiten.FPSModeVsyncOn)
		if err := ebiten.RunGame(game); err != nil {
			log.Fatal(err)
		}
		return
	}
}
