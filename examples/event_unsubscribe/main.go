package main

import (
	"fmt"
	action "github.com/DilemaFixer/GAction"
)

func main() {
	onEvt, emitEvt := action.NewEvent1[int]()

	sub := onEvt.Subscribe(func(v int) { fmt.Println("sub1:", v) })
	onEvt.Subscribe(func(v int) { fmt.Println("sub2:", v) })

	emitEvt.Emit(10)
	sub.Unsubscribe()
	emitEvt.Emit(20)
}
