package main

import (
	"github.com/veandco/go-sdl2/sdl"
	"math"
	"math/rand"
	"time"
)

const (
	SCREEN_WIDTH  = 800
	SCREEN_HEIGHT = 600

	INITIAL_VELOCITY = 1
)

func main() {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Go Click Ball", sdl.WINDOWPOS_CENTERED, sdl.WINDOWPOS_CENTERED, SCREEN_WIDTH, SCREEN_HEIGHT, sdl.WINDOW_SHOWN|sdl.WINDOW_RESIZABLE)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	bmp, err := sdl.LoadBMP("ball.bmp")
	if err != nil {
		panic(err)
	}

	texture, err := renderer.CreateTextureFromSurface(bmp)
	if err != nil {
		panic(err)
	}
	bmp.Free()
	defer texture.Destroy()

	_, _, width, height, err := texture.Query()
	if err != nil {
		panic(err)
	}

	ball := &Ball{
		Box:      sdl.Rect{W: 100, H: 100},
		Velocity: sdl.Point{X: INITIAL_VELOCITY, Y: INITIAL_VELOCITY},
	}

GameLoop:
	for {
		for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
			switch event.(type) {
			case *sdl.QuitEvent:
				break GameLoop
			case *sdl.MouseButtonEvent:
				mouseEvent := event.(*sdl.MouseButtonEvent)
				if mouseEvent.Type != sdl.MOUSEBUTTONDOWN {
					break
				}

				mousePosition := sdl.Point{X: mouseEvent.X, Y: mouseEvent.Y}

				rect, _ := sdl.EnclosePoints([]sdl.Point{mousePosition}, &ball.Box)
				if rect.X != 0 && rect.Y != 0 && mouseEvent.Button == sdl.BUTTON_LEFT {

					absVel := math.Abs(float64(ball.Velocity.X))
					if absVel < 32 {
						incr := int32(absVel) + 1
						ball.Velocity.Y, ball.Velocity.X = incr, incr
						r := rand.New(rand.NewSource(time.Now().UnixNano()))
						if r.Int()%2 == 0 {
							ball.Velocity.X = -ball.Velocity.X
						}
						if r.Int()%2 == 0 {
							ball.Velocity.Y = -ball.Velocity.Y
						}
					}

				}
			}
		}

		ball.Box.X += ball.Velocity.X
		ball.Box.Y += ball.Velocity.Y

		if ball.Box.X < 0 || ball.Box.X+ball.Box.W >= SCREEN_WIDTH {
			ball.Velocity.X = -ball.Velocity.X
		}

		if ball.Box.Y < 0 || ball.Box.Y+ball.Box.H >= SCREEN_HEIGHT {
			ball.Velocity.Y = -ball.Velocity.Y
		}

		renderer.SetDrawColor(0, 255, 0, 255)
		renderer.Clear()

		//renderer.SetDrawColor(127, 0, 127, 0)
		//renderer.FillRect(&ball.Box)
		renderer.Copy(texture, &sdl.Rect{W: width, H: height}, &ball.Box)
		renderer.Present()

		sdl.Delay(10)
	}

}

type Ball struct {
	Box      sdl.Rect
	Velocity sdl.Point
}
