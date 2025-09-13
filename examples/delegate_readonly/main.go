package main

import (
	"fmt"

	action "github.com/DilemaFixer/GAction"
)

type Worker struct {
	OnDone action.Invoker1[int]
}

func NewWorker() *Worker {
	d := action.NewDelegate1[int]()
	d.Set(func(n int) { fmt.Println("done with", n) })
	return &Worker{OnDone: d}
}

func main() {
	w := NewWorker()
	w.OnDone.Invoke(100)
}
