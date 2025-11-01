package observerpattern

type Observer interface {
	Update(event string)
}
