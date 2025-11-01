package heroesfactory

import (
	"fmt"
)

func GetHeroes(heroType string) (IHero, error) {
	switch heroType {
	case "Warrior":
		return newWarrior(), nil
	case "Mage":
		return newMage(), nil
	case "Archer":
		return newArcher(), nil
	default:
		return nil, fmt.Errorf("unknown hero type: %s", heroType)
	}
}
