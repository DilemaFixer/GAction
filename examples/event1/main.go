package main

import (
	"fmt"
	action "github.com/DilemaFixer/GAction"
)

func main() {
	onMsg, emitMsg := action.NewEvent1[string]()

	onMsg.Subscribe(func(s string) { fmt.Println("A:", s) })
	onMsg.Subscribe(func(s string) { fmt.Println("B:", len(s)) })

	emitMsg.Emit("hello")
}
