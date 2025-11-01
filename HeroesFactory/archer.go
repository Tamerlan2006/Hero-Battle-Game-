package heroesfactory

import strategypattern "week67/StrategyPattern"

type Archer struct {
	Hero
}

func newArcher() IHero {
	return &Archer{
		Hero: Hero{
			heroType: "Archer",
			health:   80,
			strength: 120,
			armor:    60,
			strategy: &strategypattern.Ranged{},
		},
	}
}
