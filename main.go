package main

import (
	"image"
	"image/draw"
	"log"
	"os"
	"runtime/pprof"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
)

const (
	FrameDuration  time.Duration = 33 * time.Millisecond
	PlayerMovement int           = 2
	BallMovement   int           = 3
)

type Direction uint

const (
	UpLeft Direction = iota
	DownLeft
	DownRight
	UpRight
)

type Side uint

const (
	Left Side = iota
	Right
)

type Action uint

const (
	Up Action = iota
	Down
	Nothing
)

type Player struct {
	Side   Side
	Height int
	Score  int
	Action Action
}

func (player Player) Render(img *image.RGBA) {
	var x int
	if player.Side == Left {
		x = 15
	} else {
		x = 738
	}
	rect := image.Rect(x, player.Height, x+15, player.Height+100)
	draw.Draw(img, rect, image.White, image.Point{}, draw.Over)
}

func (player *Player) Update(img *image.RGBA) {
	var x int
	if player.Side == Left {
		x = 15
	} else {
		x = 738
	}
	switch player.Action {
	case Up:
		if player.Height == 0 {
			break
		}
		draw.Draw(img, image.Rect(x, player.Height+100-PlayerMovement, x+15, player.Height+100), image.Black, image.Point{}, draw.Over)
		draw.Draw(img, image.Rect(x, player.Height-PlayerMovement, x+15, player.Height), image.White, image.Point{}, draw.Over)
		player.Height -= PlayerMovement
	case Down:
		if player.Height == 412 {
			break
		}
		draw.Draw(img, image.Rect(x, player.Height+100-PlayerMovement, x+15, player.Height+100), image.White, image.Point{}, draw.Over)
		draw.Draw(img, image.Rect(x, player.Height-PlayerMovement, x+15, player.Height), image.Black, image.Point{}, draw.Over)
		player.Height += PlayerMovement
	}
	player.Action = Nothing
}

type Ball struct {
	Direction Direction
	Position  image.Point
}

func (ball Ball) Render(img *image.RGBA) {
	rect := image.Rect(ball.Position.X, ball.Position.Y, ball.Position.X+8, ball.Position.Y+8)
	draw.Draw(img, rect, image.White, image.Point{}, draw.Over)
}

func (ball *Ball) Update(img *image.RGBA) {
	rect := image.Rect(ball.Position.X, ball.Position.Y, ball.Position.X+8, ball.Position.Y+8)
	draw.Draw(img, rect, image.Black, image.Point{}, draw.Over)
	if ball.Position.Y == 0 && ball.Direction == UpLeft {
		ball.Direction = DownLeft
	} else if ball.Position.Y == 0 && ball.Direction == UpRight {
		ball.Direction = DownRight
	}

	if ball.Position.Y == 508 && ball.Direction == DownLeft {
		ball.Direction = UpLeft
	} else if ball.Position.Y == 508 && ball.Direction == DownRight {
		ball.Direction = UpRight
	}

	if ball.Position.X == 30 && ball.Direction == UpLeft {
		ball.Direction = UpRight
	}
	if ball.Position.X == 30 && ball.Direction == DownLeft {
		ball.Direction = DownRight
	}

	if ball.Position.X == 730 && ball.Direction == UpRight {
		ball.Direction = UpLeft
	}
	if ball.Position.X == 730 && ball.Direction == DownRight {
		ball.Direction = DownLeft
	}

	switch ball.Direction {
	case UpLeft:
		ball.Position.Y -= PlayerMovement
		ball.Position.X -= PlayerMovement
	case DownLeft:
		ball.Position.Y += PlayerMovement
		ball.Position.X -= PlayerMovement
	case DownRight:
		ball.Position.Y += PlayerMovement
		ball.Position.X += PlayerMovement
	case UpRight:
		ball.Position.Y -= PlayerMovement
		ball.Position.X += PlayerMovement
	}
	rect = image.Rect(ball.Position.X, ball.Position.Y, ball.Position.X+8, ball.Position.Y+8)
	draw.Draw(img, rect, image.White, image.Point{}, draw.Over)
}

type Game struct {
	P1   Player
	P2   Player
	Ball Ball
	img  *image.RGBA
}

func (game *Game) Update() {
	game.P1.Update(game.img)
	game.P2.Update(game.img)
	game.Ball.Update(game.img)
}

func (game Game) GetImg() image.Image {
	return game.img
}

func NewGame() Game {
	rect := image.Rect(0, 0, 768, 512)
	img := image.NewRGBA(rect)
	draw.Draw(img, rect, image.Black, image.Point{}, draw.Src)
	game := Game{
		P1: Player{
			Side:   Left,
			Height: 206,
			Score:  0,
		},
		P2: Player{
			Side:   Right,
			Height: 206,
			Score:  0,
		},
		Ball: Ball{
			Direction: Direction(time.Now().Nanosecond() % 4),
			Position: image.Point{
				X: 380,
				Y: 252,
			},
		},
		img: img,
	}
	game.P1.Render(game.img)
	game.P2.Render(game.img)
	game.Ball.Render(game.img)
	return game
}

func main() {
	f, err := os.Create("prof")
	if err != nil {
		log.Fatal(err)
	}
	_ = pprof.StartCPUProfile(f)
	defer pprof.StopCPUProfile()
	game := NewGame()

	a := app.New()
	w := a.NewWindow("Hello World")
	w.SetFixedSize(true)
	w.Resize(fyne.NewSize(768, 512))
	c := canvas.NewImageFromImage(game.GetImg())
	c.FillMode = canvas.ImageFillOriginal
	w.SetContent(c)
	go func() {
		time.Sleep(2 * time.Second)
		for i := 0; i < 1000; i++ {
			time.Sleep(FrameDuration)
			game.Update()
			canvas.Refresh(c)
		}
		w.Close()
	}()
	w.ShowAndRun()
}
