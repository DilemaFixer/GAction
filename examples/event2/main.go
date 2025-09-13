package main

import (
	"fmt"
	action "github.com/DilemaFixer/GAction"
)

func main() {
	onPair, emitPair := action.NewEvent2[string, int]()

	onPair.Subscribe(func(tag string, code int) {
		fmt.Printf("[%s] code=%d\n", tag, code)
	})

	emitPair.Emit("INFO", 200)
	emitPair.Emit("ERROR", 500)
}
