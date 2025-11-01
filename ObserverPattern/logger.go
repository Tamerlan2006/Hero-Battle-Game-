package observerpattern

import "fmt"

type Logger struct{}

func (l *Logger) Update(event string) {
	fmt.Println("[LOG]:", event)
}
