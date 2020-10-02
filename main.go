package main

import (
	"context"
	"fmt"
	"os"
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

	wind, rend, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Errorf("Could not create window: %v", err)
	}
	defer wind.Destroy()

	// Show title
	if err := drawTitle(rend); err != nil {
		fmt.Errorf("Could not drawTitle: %v", err)
	}
	time.Sleep(2 * time.Second)

	// Show scene; managed with goroutine
	scene, err := newScene(rend)
	if err != nil {
		fmt.Errorf("Could not create scene: %v", err)
	}
	defer scene.destroy()

	// run a scene with context
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	select {
	case err := <-scene.run(ctx, rend):
		return err
	case <-time.After(5 * time.Second):
		return nil
	}
}

func drawTitle(renderer *sdl.Renderer) error {
	// Renderer management
	renderer.Clear()
	defer renderer.Present()

	// Get font
	font, err := ttf.OpenFont("res/fonts/Flappy.ttf", 20)
	if err != nil {
		return fmt.Errorf("Could not load font: %v", err)
	}
	defer font.Close()

	// Write message
	c := sdl.Color{R: 255, G: 100, B: 0, A: 255}
	surface, err := font.RenderUTF8Solid("Flappy Gopher", c)
	if err != nil {
		return fmt.Errorf("Could not render title: %v", err)
	}
	defer surface.Free()

	// Get it ready for renderer
	tex, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("Could not create texture: %v", err)
	}
	defer tex.Destroy()

	// Show it
	if err := renderer.Copy(tex, nil, nil); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	return nil
}
