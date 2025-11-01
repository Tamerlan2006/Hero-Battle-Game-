package heroesfactory

import (
	"fmt"
	observerpattern "week67/ObserverPattern"
	strategypattern "week67/StrategyPattern"
)

type Hero struct {
	heroType  string
	health    int
	strength  int
	armor     int
	strategy  strategypattern.Strategy
	observers []observerpattern.Observer
}

func (h *Hero) SetStrategy(strategy strategypattern.Strategy) {
	h.strategy = strategy
}

func (h *Hero) ExecuteAttack() {
	if h.strategy != nil {
		event := fmt.Sprintf("%s is attacking!", h.heroType)
		h.NotifyObservers(event)
		h.strategy.Attack()
	}
}

func (h *Hero) RegisterObserver(o observerpattern.Observer) {
	h.observers = append(h.observers, o)
}

func (h *Hero) UnregisterObserver(o observerpattern.Observer) {
	for i, obs := range h.observers {
		if obs == o {
			h.observers = append(h.observers[:i], h.observers[i+1:]...)
			break
		}
	}
}

func (h *Hero) NotifyObservers(event string) {
	for _, o := range h.observers {
		o.Update(event)
	}
}

// ================= GETTERS/SETTERS =================
func (h *Hero) SetType(heroType string)  { h.heroType = heroType }
func (h *Hero) GetType() string          { return h.heroType }
func (h *Hero) SetHealth(health int)     { h.health = health }
func (h *Hero) GetHealth() int           { return h.health }
func (h *Hero) SetStrength(strength int) { h.strength = strength }
func (h *Hero) GetStrength() int         { return h.strength }
func (h *Hero) SetArmor(armor int)       { h.armor = armor }
func (h *Hero) GetArmor() int            { return h.armor }
