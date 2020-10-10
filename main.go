package main

import (
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/veandco/go-sdl2/sdl"
	"github.com/veandco/go-sdl2/ttf"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(2)
	}
}

const (
	windowWidth  = 800
	windowHeight = 600
)

func run() error {
	err := sdl.Init(sdl.INIT_EVERYTHING)
	if err != nil {
		fmt.Errorf("Could not initialize SDL: %v", err)
	}
	defer sdl.Quit()

	if err := ttf.Init(); err != nil {
		fmt.Errorf("Could not initialize TTF: %v", err)
	}
	defer ttf.Quit()

	window, renderer, err := sdl.CreateWindowAndRenderer(windowWidth, windowHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Errorf("Could not create window: %v", err)
	}
	defer window.Destroy()

	// Show title
	if err := drawTitle(renderer, "Flappy Gopher"); err != nil {
		fmt.Errorf("Could not drawTitle: %v", err)
	}
	time.Sleep(2 * time.Second)

	// Show scene
	scene, err := newScene(renderer)
	if err != nil {
		fmt.Errorf("Could not create scene: %v", err)
	}
	defer scene.destroy()

	// Run scene
	events := make(chan sdl.Event)
	errc := scene.run(events, renderer)
	runtime.LockOSThread()
	for {
		select {
		case events <- sdl.WaitEvent():
		case err := <-errc:
			return err
		}
	}
}

func drawTitle(renderer *sdl.Renderer, text string) error {
	// Renderer management
	renderer.Clear()
	defer renderer.Present()

	// Get font
	font, err := ttf.OpenFont("res/fonts/Flappy.ttf", 14)
	if err != nil {
		return fmt.Errorf("Could not load font: %v", err)
	}
	defer font.Close()

	// Write message
	color := sdl.Color{R: 255, G: 100, B: 0, A: 255}
	surface, err := font.RenderUTF8Solid(text, color)
	if err != nil {
		return fmt.Errorf("Could not render title: %v", err)
	}
	defer surface.Free()

	// Get it ready for renderer
	texture, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("Could not create texture: %v", err)
	}
	defer texture.Destroy()

	// Show it
	if err := renderer.Copy(texture, nil, nil); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	return nil
}
