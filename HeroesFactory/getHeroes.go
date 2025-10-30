package heroesfactory

import "fmt"

func GetHeroes(heroType string) (IHero, error) {
	if heroType == "Warrior" {
		return newWarrior(), nil
	} else if heroType == "Mage" {
		return newMage(), nil
	} else if heroType == "Archer" {
		return newArcher(), nil
	}
	return nil, fmt.Errorf("unknown hero type: %s", heroType)
}