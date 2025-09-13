package main

import (
	"fmt"
	action "github.com/DilemaFixer/GAction"
)

func main() {
	d := action.NewDelegate0()
	d.Set(func() { fmt.Println("hello delegate0") })

	var inv action.Invoker0 = d
	inv.Invoke()
}
