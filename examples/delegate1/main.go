package main

import (
	"fmt"
	action "github.com/DilemaFixer/GAction"
)

func main() {
	d := action.NewDelegate1[int]()
	d.Set(func(v int) { fmt.Println("value:", v) })

	d.Invoke(42) // value: 42
}
