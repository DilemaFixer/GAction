package main

import (
	"fmt"
	action "github.com/DilemaFixer/GAction"
)

func main() {
	onA, emitA := action.NewEvent0()
	onB, emitB := action.NewEvent0()

	// проброс события A -> B
	onA.Subscribe(func() { emitB.Emit() })

	onB.Subscribe(func() { fmt.Println("B triggered") })

	emitA.Emit() // вызовет B
}
