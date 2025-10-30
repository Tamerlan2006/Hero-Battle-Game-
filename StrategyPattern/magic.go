package strategypattern

import "fmt"

type Magic struct {
}

func (m *Magic) Attack() {
	fmt.Println("Attacking with powerful magic spells!")
}