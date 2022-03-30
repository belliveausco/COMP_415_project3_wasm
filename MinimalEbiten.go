package main

import (
	"embed"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
	"image/png"
	"image/color"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"log"
	"math/rand"
	"time"
	"fmt"
	"os"
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    48,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

var (
	mplusNormalFont font.Face
	mplusBigFont    font.Face
	jaKanjis        = []rune{}
)

//go:embed assets/*
var EmbeddedAssets embed.FS

const (
	GameWidth   = 700
	GameHeight  = 700
	PlayerSpeed = 5
)

type Sprite struct {
	pict *ebiten.Image
	xloc int
	yloc int
	dX   int
	dY   int
}

type Game struct {
	player  Sprite
	enemy   Sprite
	listOfenemies []Sprite
	// enemy slice of sprites
	score   int
	drawOps ebiten.DrawImageOptions
}

func (g *Game) Update() error {
	processPlayerInput(g)
	for i := 0; i < len(g.listOfenemies); i++ {
		if gotPeople(g.player, g.listOfenemies[i]) == false {
			g.listOfenemies = remove(g.listOfenemies, i)
			g.score++
		}
	}
	if len(g.listOfenemies) == 0 {
		fmt.Println("End of game")
		fmt.Println("Score is:",g.score)
		os.Exit(0)
	}
	return nil
}

// reference : https://github.com/jsantore/FirstGameDemo/blob/master/GmeEngineDemo.go
func gotPeople(player, listOfenemies Sprite) bool {
	enemyWidth, enemyHeight := listOfenemies.pict.Size()
	playerWidth, playerHeight := player.pict.Size()
	if player.xloc < listOfenemies.xloc+enemyWidth &&
		player.xloc+playerWidth > listOfenemies.xloc &&
		player.yloc < listOfenemies.yloc+enemyHeight &&
		player.yloc+playerHeight > listOfenemies.yloc {
		return false
	}
	return true
}
// reference https://stackoverflow.com/questions/37334119/how-to-delete-an-element-from-a-slice-in-golang
func remove(slice []Sprite, s int) []Sprite {
	return append(slice[:s], slice[s+1:]...)
}

func (g Game) Draw(screen *ebiten.Image) {
	g.drawOps.GeoM.Reset()
	g.drawOps.GeoM.Translate(float64(g.player.xloc), float64(g.player.yloc))
	screen.DrawImage(g.player.pict, &g.drawOps)
	for _, currentEnemy := range g.listOfenemies {
		g.drawOps.GeoM.Reset()
		g.drawOps.GeoM.Translate(float64(currentEnemy.xloc), float64(currentEnemy.yloc))
		screen.DrawImage(currentEnemy.pict, &g.drawOps)
	}
	text.Draw(screen, fmt.Sprint(g.score), mplusNormalFont, 100, 40, color.White)
}

func (g Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return GameWidth, GameHeight
}

// sets window size sets window title
// creates game starts 0, player gets nil value
// add player values with image loaded with an image
func main() {
	ebiten.SetWindowSize(GameWidth, GameHeight)
	ebiten.SetWindowTitle("Project 3")
	simpleGame := Game{score: 0}
	simpleGame.player = Sprite{
		pict: loadPNGImageFromEmbedded("bat.png"),
		xloc: 200,
		yloc: 300,
		dX:   0,
		dY:   0,
	}
	simpleGame.listOfenemies = make([]Sprite, 10)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < 10; i++ {
		simpleGame.enemy = Sprite{
			loadPNGImageFromEmbedded("stickfigure.png"),
			rand.Intn(650),
			rand.Intn(650),
			0,
			0,
		}
		simpleGame.listOfenemies[i] = simpleGame.enemy
	}
	if err := ebiten.RunGame(&simpleGame); err != nil {
		log.Fatal("Oh no! something terrible happened and the game crashed", err)
	}
}
//position gives uppers, pict.Size()
//pict.bounds(y) gives you height
//pict.bounds(x) gives you width to give you lower corner
func loadPNGImageFromEmbedded(name string) *ebiten.Image {
	pictNames, err := EmbeddedAssets.ReadDir("assets")
	if err != nil {
		log.Fatal("failed to read embedded dir ", pictNames, " ", err)
	}
	embeddedFile, err := EmbeddedAssets.Open("assets/" + name)
	if err != nil {
		log.Fatal("failed to load embedded image ", embeddedFile, err)
	}
	rawImage, err := png.Decode(embeddedFile)
	if err != nil {
		log.Fatal("failed to load embedded image ", name, err)
	}
	gameImage := ebiten.NewImageFromImage(rawImage)
	return gameImage
}

func processPlayerInput(theGame *Game) {
	if inpututil.IsKeyJustPressed(ebiten.KeyUp) {
		theGame.player.dY = -PlayerSpeed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyDown) {
		theGame.player.dY = PlayerSpeed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyUp) || inpututil.IsKeyJustReleased(ebiten.KeyDown) {
		theGame.player.dY = 0
	}
	theGame.player.yloc += theGame.player.dY
	if theGame.player.yloc <= 0 {
		theGame.player.dY = 0
		theGame.player.yloc = 0
	} else if theGame.player.yloc+theGame.player.pict.Bounds().Size().Y > GameHeight {
		theGame.player.dY = 0
		theGame.player.yloc = GameHeight - theGame.player.pict.Bounds().Size().Y
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyLeft) {
		theGame.player.dX = -PlayerSpeed
	} else if inpututil.IsKeyJustPressed(ebiten.KeyRight) {
		theGame.player.dX = PlayerSpeed
	} else if inpututil.IsKeyJustReleased(ebiten.KeyRight) || inpututil.IsKeyJustReleased(ebiten.KeyLeft) {
		theGame.player.dX = 0
	}
	theGame.player.xloc += theGame.player.dX
	if theGame.player.xloc <= 0 {
		theGame.player.dX = 0
		theGame.player.xloc = 0
	} else if theGame.player.xloc+theGame.player.pict.Bounds().Size().X > GameWidth {
		theGame.player.dX = 0
		theGame.player.xloc = GameWidth - theGame.player.pict.Bounds().Size().X
	}
}
