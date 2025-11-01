package gameui

import (
	"fmt"
	"math/rand"
	"time"

	heroesfactory "week67/HeroesFactory"
	observerpattern "week67/ObserverPattern"
	strategypattern "week67/StrategyPattern"

	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

type Game struct {
	leftHero  heroesfactory.IHero
	rightHero heroesfactory.IHero
	state     string // "menu", "fight", "end"
	winner    string

	leftHealth  int
	rightHealth int

	logger    *observerpattern.Logger
	announcer *observerpattern.Announcer
}

func NewGame() *Game {
	return &Game{
		state:     "menu",
		logger:    &observerpattern.Logger{},
		announcer: &observerpattern.Announcer{},
	}
}

func (g *Game) Update() error {
	switch g.state {
	case "menu":
		if ebiten.IsKeyPressed(ebiten.Key1) {
			g.leftHero, _ = heroesfactory.GetHeroes("Warrior")
		}
		if ebiten.IsKeyPressed(ebiten.Key2) {
			g.leftHero, _ = heroesfactory.GetHeroes("Mage")
		}
		if ebiten.IsKeyPressed(ebiten.Key3) {
			g.leftHero, _ = heroesfactory.GetHeroes("Archer")
		}
		if ebiten.IsKeyPressed(ebiten.KeyEnter) && g.leftHero != nil {
			rightOptions := []string{"Warrior", "Mage", "Archer"}
			rand.Seed(time.Now().UnixNano())
			randHero := rightOptions[rand.Intn(len(rightOptions))]
			g.rightHero, _ = heroesfactory.GetHeroes(randHero)

			g.leftHero.RegisterObserver(g.logger)
			g.leftHero.RegisterObserver(g.announcer)
			g.rightHero.RegisterObserver(g.logger)
			g.rightHero.RegisterObserver(g.announcer)

			g.leftHealth = g.leftHero.GetHealth()
			g.rightHealth = g.rightHero.GetHealth()
			g.state = "fight"
		}

	case "fight":
		if g.leftHealth <= 0 || g.rightHealth <= 0 {
			if g.leftHealth <= 0 {
				g.winner = g.rightHero.GetType()
			} else {
				g.winner = g.leftHero.GetType()
			}
			g.state = "end"
		} else if ebiten.IsKeyPressed(ebiten.KeySpace) {
			g.performAttack()
		}

	case "end":
		if ebiten.IsKeyPressed(ebiten.KeyR) {
			g.state = "menu"
			g.leftHero = nil
			g.rightHero = nil
		}
	}

	return nil
}

func (g *Game) performAttack() {
	leftStrats := []strategypattern.Strategy{
		&strategypattern.Melee{},
		&strategypattern.Magic{},
		&strategypattern.Ranged{},
	}
	rightStrats := []strategypattern.Strategy{
		&strategypattern.Melee{},
		&strategypattern.Magic{},
		&strategypattern.Ranged{},
	}

	rand.Seed(time.Now().UnixNano())
	g.leftHero.SetStrategy(leftStrats[rand.Intn(len(leftStrats))])
	g.rightHero.SetStrategy(rightStrats[rand.Intn(len(rightStrats))])

	g.leftHero.ExecuteAttack()
	g.rightHero.ExecuteAttack()

	// Damage calculation with hero strength and armor
	leftDmg := rand.Intn(20) + g.leftHero.GetStrength()/10
	rightDmg := rand.Intn(20) + g.rightHero.GetStrength()/10
	g.leftHealth -= rightDmg
	g.rightHealth -= leftDmg
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 10, 25, 255})

	switch g.state {
	case "menu":
		msg := "=== Hero Battle Game ===\n\nChoose your hero:\n[1] Warrior\n[2] Mage\n[3] Archer\n\nPress [ENTER] to start battle"
		ebitenutil.DebugPrint(screen, msg)

	case "fight":
		g.drawArena(screen)
		ebitenutil.DebugPrintAt(screen, "Press SPACE to attack", 300, 20)

	case "end":
		msg := fmt.Sprintf("Winner: %s\n\nPress [R] to restart", g.winner)
		ebitenutil.DebugPrint(screen, msg)
	}
}

func (g *Game) drawArena(screen *ebiten.Image) {
	vector.FillRect(screen, 100, 300, 100, 150, color.RGBA{255, 0, 0, 255}, false)
	vector.FillRect(screen, 600, 300, 100, 150, color.RGBA{0, 0, 255, 255}, false)

	vector.FillRect(screen, 100, 250, float32(g.leftHealth)*2, 15, color.RGBA{255, 100, 100, 255}, false)
	vector.FillRect(screen, 500, 250, float32(g.rightHealth)*2, 15, color.RGBA{100, 100, 255, 255}, false)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return 800, 600
}
