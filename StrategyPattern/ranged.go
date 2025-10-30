package strategypattern

import "fmt"

type Ranged struct {
}

func (r *Ranged) Attack() {
	fmt.Println("Attacking from a distance with ranged weapons!")
}