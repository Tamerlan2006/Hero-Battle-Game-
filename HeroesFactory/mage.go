package heroesfactory

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
		},
	}
}