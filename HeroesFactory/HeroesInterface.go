package heroesfactory

type Heroes interface {
	setType(string)
	setHealth(int)
	setStrength(int)
	setArmor(int)
	getType() string
	getHealth() int
	getStrength() int
	getArmor() int
}
