package observerpattern

import "fmt"

type Announcer struct{}

func (a *Announcer) Update(event string) {
	fmt.Println("[ANNOUNCER]:", event)
}
