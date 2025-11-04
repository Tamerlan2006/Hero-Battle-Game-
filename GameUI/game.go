package gameui

import (
	"bytes"
	"fmt"
	"image"
	"image/color"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/audio"
	"github.com/hajimehoshi/ebiten/v2/audio/vorbis"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"

	heroesfactory "week67/HeroesFactory"
	observerpattern "week67/ObserverPattern"
	strategypattern "week67/StrategyPattern"
)

const (
	screenW = 800
	screenH = 600
	frameW  = 64
	frameH  = 128
)

type animState int

const (
	animIdle animState = iota
	animRun
	animAttack
)


type Animation struct {
	sheet       *ebiten.Image
	frameCount  int
	frameWidth  int
	frameHeight int
	frameTime   time.Duration
	current     int
	timer       time.Time
}

func NewAnimation(sheet *ebiten.Image, frameCount int, frameTime time.Duration) *Animation {
	return &Animation{
		sheet:       sheet,
		frameCount:  frameCount,
		frameWidth:  frameW,
		frameHeight: frameH,
		frameTime:   frameTime,
		timer:       time.Now(),
	}
}

func (a *Animation) Update() {
	if time.Since(a.timer) >= a.frameTime {
		a.current = (a.current + 1) % a.frameCount
		a.timer = time.Now()
	}
}

func (a *Animation) Draw(screen *ebiten.Image, opts *ebiten.DrawImageOptions) {
	sx := a.current * a.frameWidth
	sub := a.sheet.SubImage(image.Rect(sx, 0, sx+a.frameWidth, a.frameHeight)).(*ebiten.Image)
	screen.DrawImage(sub, opts)
}

func (a *Animation) Reset() {
	a.current = 0
	a.timer = time.Now()
}

type Game struct {
	leftHero  heroesfactory.IHero
	rightHero heroesfactory.IHero

	state string
	turn  string

	winner          string
	flawlessVictory bool
	fightTextTimer  time.Time
	shakeOffset     float64
	shakeTime       time.Time

	audioCtx    *audio.Context
	fightSnd    *audio.Player
	punchSnd    *audio.Player
	magicSnd    *audio.Player
	fatalitySnd *audio.Player

	arenaImg    *ebiten.Image
	leftSheet   *ebiten.Image
	rightSheet  *ebiten.Image

	leftAnim    *Animation
	rightAnim   *Animation
	leftState   animState
	rightState  animState
	runTimer    time.Time
	prevTime    time.Time

	logger    *observerpattern.Logger
	announcer *observerpattern.Announcer
}

func NewGame() *Game {
	rand.Seed(time.Now().UnixNano())
	g := &Game{
		state:     "menu",
		logger:    &observerpattern.Logger{},
		announcer: &observerpattern.Announcer{},
		prevTime:  time.Now(),
	}

	g.audioCtx = audio.NewContext(44100)
	g.loadSounds()
	g.loadArena()
	g.loadHeroSheets()

	g.leftAnim = NewAnimation(g.leftSheet, 2, 300*time.Millisecond)
	g.rightAnim = NewAnimation(g.rightSheet, 2, 300*time.Millisecond)

	g.leftState = animIdle
	g.rightState = animIdle

	return g
}

func (g *Game) loadSounds() {
	load := func(path string) *audio.Player {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil
		}
		stream, err := vorbis.DecodeWithSampleRate(44100, bytes.NewReader(data))
		if err != nil {
			return nil
		}
		p, err := g.audioCtx.NewPlayer(stream)
		if err != nil {
			return nil
		}
		return p
	}
	g.fightSnd = load("GameUI/assets/sounds/fight.ogg")
	g.punchSnd = load("GameUI/assets/sounds/punch.ogg")
	g.magicSnd = load("GameUI/assets/sounds/magic.ogg")
	g.fatalitySnd = load("GameUI/assets/sounds/fatality.ogg")
}

func (g *Game) play(p *audio.Player) {
	if p == nil {
		return
	}
	p.Rewind()
	p.Play()
}

func (g *Game) loadArena() {
	img, _, err := ebitenutil.NewImageFromFile("GameUI/assets/arena.png")
	if err != nil {
		img = ebiten.NewImage(screenW, screenH)
		img.Fill(color.RGBA{20, 20, 40, 255})
	}
	g.arenaImg = img
}

func (g *Game) loadSprite(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		img = ebiten.NewImage(frameW*2, frameH)
		img.Fill(color.RGBA{120, 120, 120, 255})
	}
	return img
}

func (g *Game) loadHeroSheets() {
	g.leftSheet = g.loadSprite("GameUI/assets/heroes/warrior.png")
	g.rightSheet = g.loadSprite("GameUI/assets/heroes/mage.png")
}

func (g *Game) Update() error {
	now := time.Now()
	g.prevTime = now 

	switch g.state {
	case "menu":
		g.updateMenu()
	case "fight":
		g.updateFight()
	case "end":
		if inpututil.IsKeyJustPressed(ebiten.KeyR) {
			g.reset()
		}
	}

	g.leftAnim.Update()
	g.rightAnim.Update()

	if time.Since(g.shakeTime) < 200*time.Millisecond {
		g.shakeOffset = rand.Float64()*4 - 2
	} else {
		g.shakeOffset = 0
	}

	g.checkRunAndAttack()
	return nil
}

func (g *Game) updateMenu() {
	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		g.leftHero, _ = heroesfactory.GetHeroes("Warrior")
	}
	if inpututil.IsKeyJustPressed(ebiten.Key2) {
		g.leftHero, _ = heroesfactory.GetHeroes("Mage")
	}
	if inpututil.IsKeyJustPressed(ebiten.Key3) {
		g.leftHero, _ = heroesfactory.GetHeroes("Archer")
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) && g.leftHero != nil {
		opts := []string{"Warrior", "Mage", "Archer"}
		r := opts[rand.Intn(len(opts))]
		g.rightHero, _ = heroesfactory.GetHeroes(r)

		g.clearObservers()
		g.registerObservers()

		g.state = "fight"
		g.turn = "player"
		g.fightTextTimer = time.Now()
		g.play(g.fightSnd)
	}
}

func (g *Game) updateFight() {
	if g.leftHero.GetHealth() <= 0 || g.rightHero.GetHealth() <= 0 {
		g.endFight()
		return
	}

	if g.turn == "player" {
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			g.leftHero.SetStrategy(&strategypattern.Melee{})
			g.startRun(true)
		} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
			g.leftHero.SetStrategy(&strategypattern.Ranged{})
			g.startRun(true)
		} else if inpututil.IsKeyJustPressed(ebiten.Key3) {
			g.leftHero.SetStrategy(&strategypattern.Magic{})
			g.startRun(true)
		}
	} else {
		g.startAIAttack()
	}
}

func (g *Game) startRun(isPlayer bool) {
	if isPlayer {
		g.leftState = animRun
		g.runTimer = time.Now()
		g.leftAnim = NewAnimation(g.leftSheet, 2, 120*time.Millisecond) // run: 2 кадра
	} else {
		g.rightState = animRun
		g.runTimer = time.Now()
		g.rightAnim = NewAnimation(g.rightSheet, 2, 120*time.Millisecond)
	}
}

func (g *Game) startAIAttack() {
	time.Sleep(600 * time.Millisecond)
	strats := []strategypattern.Strategy{
		&strategypattern.Melee{},
		&strategypattern.Ranged{},
		&strategypattern.Magic{},
	}
	g.rightHero.SetStrategy(strats[rand.Intn(len(strats))])
	g.startRun(false)
}

func (g *Game) checkRunAndAttack() {
	// run → attack
	if g.leftState == animRun && time.Since(g.runTimer) > 400*time.Millisecond {
		g.leftState = animAttack
		g.leftAnim = NewAnimation(g.leftSheet, 1, 200*time.Millisecond)
		g.leftAnim.current = 1 // attack = кадр 1
		g.performAttack(g.leftHero, g.rightHero)
	}
	if g.rightState == animRun && time.Since(g.runTimer) > 400*time.Millisecond {
		g.rightState = animAttack
		g.rightAnim = NewAnimation(g.rightSheet, 1, 200*time.Millisecond)
		g.rightAnim.current = 1
		g.performAttack(g.rightHero, g.leftHero)
	}

	if g.leftState == animAttack && time.Since(g.runTimer) > 600*time.Millisecond {
		g.leftState = animIdle
		g.leftAnim = NewAnimation(g.leftSheet, 1, 300*time.Millisecond)
		g.leftAnim.current = 0
		if g.turn == "player" {
			g.turn = "ai"
		}
	}
	if g.rightState == animAttack && time.Since(g.runTimer) > 600*time.Millisecond {
		g.rightState = animIdle
		g.rightAnim = NewAnimation(g.rightSheet, 1, 300*time.Millisecond)
		g.rightAnim.current = 0
		if g.turn == "ai" {
			g.turn = "player"
		}
	}
}

func (g *Game) performAttack(att, def heroesfactory.IHero) {
	att.ExecuteAttack()

	base := rand.Intn(20) + att.GetStrength()/10
	bonus := 0
	switch att.GetType() {
	case "Warrior":
		bonus = att.GetStrength() / 5
	case "Archer":
		bonus = att.GetStrength()/8 - def.GetArmor()/2
	case "Mage":
		bonus = att.GetStrength() / 4
	}
	dmg := max(1, base+bonus-def.GetArmor()/10)
	def.SetHealth(def.GetHealth() - dmg)

	g.shakeTime = time.Now()
	g.play(g.punchSnd)
}

func (g *Game) endFight() {
	if g.leftHero.GetHealth() <= 0 {
		g.winner = g.rightHero.GetType()
	} else {
		g.winner = g.leftHero.GetType()
		g.flawlessVictory = true
	}
	g.state = "end"
	g.play(g.fatalitySnd)
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{10, 10, 25, 255})

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(screenW)/float64(g.arenaImg.Bounds().Dx()),
		float64(screenH)/float64(g.arenaImg.Bounds().Dy()))
	screen.DrawImage(g.arenaImg, op)

	shakeX := g.shakeOffset
	shakeY := g.shakeOffset

	leftOpts := &ebiten.DrawImageOptions{}
	leftOpts.GeoM.Scale(2, 2)
	leftOpts.GeoM.Translate(150+shakeX, 320+shakeY)
	g.leftAnim.Draw(screen, leftOpts)

	rightOpts := &ebiten.DrawImageOptions{}
	rightOpts.GeoM.Scale(-2, 2)
	rightOpts.GeoM.Translate(650+shakeX, 320+shakeY)
	g.rightAnim.Draw(screen, rightOpts)

	g.drawHP(screen, g.leftHero, 100, 50, true)
	g.drawHP(screen, g.rightHero, 500, 50, false)

	switch g.state {
	case "menu":
		ebitenutil.DebugPrint(screen,
			"=== MORTAL KOMBAT GO ===\n\n[1] Warrior\n[2] Mage\n[3] Archer\n\n[ENTER] FIGHT!")
	case "fight":
		if time.Since(g.fightTextTimer) < 1500*time.Millisecond {
			ebitenutil.DebugPrintAt(screen, "ROUND 1  FIGHT!", 300, 200)
		}
		ebitenutil.DebugPrintAt(screen,
			"Choose attack:  [1] Melee   [2] Ranged   [3] Magic", 200, 560)
	case "end":
		msg := fmt.Sprintf("FATALITY!  %s WINS!", g.winner)
		if g.flawlessVictory {
			msg = "FLAWLESS VICTORY!\n" + msg
		}
		ebitenutil.DebugPrint(screen, msg+"\n\n[R] Restart")
	}
}

func (g *Game) drawHP(screen *ebiten.Image, h heroesfactory.IHero, x, y int, _ bool) {
	if h == nil {
		return
	}
	label := fmt.Sprintf("%s  HP:%d", h.GetType(), h.GetHealth())
	ebitenutil.DebugPrintAt(screen, label, x, y-20)

	maxHP := 150
	fill := float32(h.GetHealth()) / float32(maxHP) * 200
	col := color.RGBA{0, 255, 0, 255}
	if fill < 0.4 {
		col = color.RGBA{255, 0, 0, 255}
	}
	vector.FillRect(screen, float32(x), float32(y), fill, 20, col, false)
}

func (g *Game) Layout(_, _ int) (int, int) { return screenW, screenH }

func (g *Game) clearObservers() {
	if g.leftHero != nil {
		g.leftHero.UnregisterObserver(g.logger)
		g.leftHero.UnregisterObserver(g.announcer)
	}
	if g.rightHero != nil {
		g.rightHero.UnregisterObserver(g.logger)
		g.rightHero.UnregisterObserver(g.announcer)
	}
}

func (g *Game) registerObservers() {
	g.leftHero.RegisterObserver(g.logger)
	g.leftHero.RegisterObserver(g.announcer)
	g.rightHero.RegisterObserver(g.logger)
	g.rightHero.RegisterObserver(g.announcer)
}

func (g *Game) reset() {
	g.clearObservers()
	g.leftHero = nil
	g.rightHero = nil
	g.state = "menu"
	g.winner = ""
	g.flawlessVictory = false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}