package main

import (
	"fmt"

	"github.com/detohm/gomu8080"
)

func main() {
	mmu := gomu8080.NewMMU()
	p := gomu8080.NewProcessor(mmu, true)
	fmt.Printf("%+v", p)
}
