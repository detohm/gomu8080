package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/detohm/gomu8080"
)

func main() {
	path := flag.String("path", "", "")       // TODO - add more detail
	debugMode := flag.Bool("debug", true, "") // TODO - add more detail
	flag.Parse()

	mmu := gomu8080.NewMMU()
	p := gomu8080.NewProcessor(mmu, *debugMode)

	bytes, err := os.ReadFile(*path)
	if err != nil {
		fmt.Println(err)
		return
	}

	// load rom into byte slice
	mmu.Load(len(bytes), bytes, 0x0100)

	p.PC = 0x0100

	for !p.IsHalt {

		if p.DebugMode {
			fmt.Printf("pc: %X ", p.PC)
		}

		p.Run()
		if p.DebugMode {
			time.Sleep(5 * time.Millisecond)
		}
	}
}
