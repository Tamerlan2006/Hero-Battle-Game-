package heroesfactory

type Warrior struct {
	Hero
}

func newWarrior() IHero {
	return &Warrior{
		Warrior: Warrior{
			heroType: "Warrior",
			health:   100,
			strength: 90,
			armor:    100,
		},
	}
}
