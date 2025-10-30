package strategypattern

import "fmt"

type Melee struct {
}

func (m *Melee) Attack() {
	fmt.Println("Attacking from close range with melee weapons!")
}
