package heroesfactory

import "week67/StrategyPattern"

type Warrior struct {
	Hero
}

func newWarrior() IHero {
	return &Warrior{
		Hero: Hero{
			heroType: "Warrior",
			health:   150,
			strength: 100,
			armor:    80,
			strategy: &strategypattern.Melee{},
		},
	}
}
