package main

import (
	"fmt"
	action "github.com/DilemaFixer/GAction"
)

func main() {
	onLog, emitLog := action.NewEvent3[string, int, error]()

	onLog.Subscribe(func(tag string, code int, err error) {
		fmt.Println("log:", tag, code, err)
	})

	emitLog.Emit("AUTH", 401, nil)
}
