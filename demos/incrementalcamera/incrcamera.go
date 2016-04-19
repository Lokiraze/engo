package main

import (
	"image"
	"image/color"
	"log"

	"engo.io/ecs"
	"engo.io/engo"
	"time"
)

type Game struct{}

var (
	scrollSpeed float32 = 700
	worldWidth  float32 = 800
	worldHeight float32 = 800
)

// generateBackground creates a background of green tiles - might not be the most efficient way to do this
func generateBackground() *ecs.Entity {
	rect := image.Rect(0, 0, int(worldWidth), int(worldHeight))
	img := image.NewNRGBA(rect)
	c1 := color.RGBA{102, 153, 0, 255}
	c2 := color.RGBA{102, 173, 0, 255}
	for i := rect.Min.X; i < rect.Max.X; i++ {
		for j := rect.Min.Y; j < rect.Max.Y; j++ {
			if i%40 > 20 {
				if j%40 > 20 {
					img.Set(i, j, c1)
				} else {
					img.Set(i, j, c2)
				}
			} else {
				if j%40 > 20 {
					img.Set(i, j, c2)
				} else {
					img.Set(i, j, c1)
				}
			}
		}
	}
	bgTexture := engo.NewImageObject(img)
	field := ecs.NewEntity("RenderSystem")
	fieldRender := engo.NewRenderComponent(engo.NewTexture(bgTexture), engo.Point{1, 1}, "Background1")
	fieldSpace := &engo.SpaceComponent{engo.Point{0, 0}, worldWidth, worldHeight}
	field.AddComponent(fieldRender)
	field.AddComponent(fieldSpace)
	return field
}

func (game *Game) Preload() {}

func (*Game) Hide()        {}
func (*Game) Show()        {}
func (*Game) Exit()        {}
func (*Game) Type() string { return "Game" }

// Setup is called before the main loop is started
func (game *Game) Setup(w *ecs.World) {
	engo.SetBackground(color.White)
	w.AddSystem(&engo.RenderSystem{})

	// Create the background; this way we'll see when we actually scroll
	err := w.AddEntity(generateBackground())
	if err != nil {
		log.Println(err)
	}

	// We issue one camera zoom command at the start, but it takes a while to process because we set a duration
	engo.Mailbox.Dispatch(engo.CameraMessage{
		Axis:        engo.ZAxis,
		Value:       3, // so zooming out a lot
		Incremental: true,
		Duration:    time.Second * 5,
	})
}

func main() {
	opts := engo.RunOptions{
		Title:  "IncrementalCamera Demo",
		Width:  int(worldWidth),
		Height: int(worldHeight),
	}
	engo.Run(opts, &Game{})
}