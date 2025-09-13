package main

import (
	"fmt"
	action "github.com/DilemaFixer/GAction"
)

func main() {
	onTick, emitTick := action.NewEvent0()

	onTick.Subscribe(func() { fmt.Println("tick A") })
	onTick.Subscribe(func() { fmt.Println("tick B") })

	emitTick.Emit()
	emitTick.Emit()
}
