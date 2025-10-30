package heroesfactory

import "week67/StrategyPattern"

type Mage struct {
	Hero
}

func newMage() IHero {
	return &Mage{
		Hero: Hero{
			heroType: "Mage",
			health:   70,
			strength: 150,
			armor:    40,
			strategy: &strategypattern.Magic{},
		},
	}
}