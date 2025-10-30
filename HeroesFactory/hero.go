package heroesfactory

import "week67/StrategyPattern"

type Hero struct {
	heroType string
	health   int
	strength int
	armor    int
	strategy strategypattern.Strategy
}

func (h *Hero) setStrategy(strategy strategypattern.Strategy) {
	h.strategy = strategy
}

func (h *Hero) executeAttack() {
	if h.strategy != nil {
		h.strategy.Attack()
	}
}

func (h *Hero) setType(heroType string) {
	h.heroType = heroType
}

func (h *Hero) getType() string {
	return h.heroType
}

func (h *Hero) setHealth(health int) {
	h.health = health
}

func (h *Hero) getHealth() int {
	return h.health
}

func (h *Hero) setStrength(strength int) {
	h.strength = strength
}

func (h *Hero) getStrength() int {
	return h.strength
}

func (h *Hero) setArmor(armor int) {
	h.armor = armor
}

func (h *Hero) getArmor() int {
	return h.armor
}
