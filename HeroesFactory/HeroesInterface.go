package heroesfactory

import (
	observerpattern "week67/ObserverPattern"
	strategypattern "week67/StrategyPattern"
)

type IHero interface {
	SetType(string)
	SetHealth(int)
	SetStrength(int)
	SetArmor(int)
	GetType() string
	GetHealth() int
	GetStrength() int
	GetArmor() int

	SetStrategy(strategypattern.Strategy)
	ExecuteAttack()

	RegisterObserver(observerpattern.Observer)
	UnregisterObserver(observerpattern.Observer)
	NotifyObservers(string)
}
