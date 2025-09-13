package main

import (
	"fmt"

	action "github.com/DilemaFixer/GAction"
)

type Service struct {
	OnReady action.Event0
}

func NewService() (*Service, action.Emitter0) {
	evt, em := action.NewEvent0()
	return &Service{OnReady: evt}, em
}

func main() {
	s, fire := NewService()

	s.OnReady.Subscribe(func() { fmt.Println("service ready!") })

	fire.Emit()
}
