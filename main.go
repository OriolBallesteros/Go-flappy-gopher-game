package main

import (
	"fmt"
	"os"
	"time"

	"github.com/veandco/go-sdl2/img"
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

	w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)
	if err != nil {
		fmt.Errorf("Could not create window: %v", err)
	}
	defer w.Destroy()

	if err := drawTitle(r); err != nil {
		fmt.Errorf("Could not drawTitle: %v", err)
	}

	time.Sleep(5 * time.Second)

	if err := drawBackground(r); err != nil {
		fmt.Errorf("Could not draw background: %v", err)
	}

	time.Sleep(5 * time.Second)

	return nil
}

func drawTitle(renderer *sdl.Renderer) error {
	renderer.Clear()

	font, err := ttf.OpenFont("res/fonts/Flappy.ttf", 20)
	if err != nil {
		return fmt.Errorf("Could not load font: %v", err)
	}
	defer font.Close()

	c := sdl.Color{R: 255, G: 100, B: 0, A: 255}
	surface, err := font.RenderUTF8Solid("Flappy Gopher", c)
	if err != nil {
		return fmt.Errorf("Could not render title: %v", err)
	}

	defer surface.Free()

	tex, err := renderer.CreateTextureFromSurface(surface)
	if err != nil {
		return fmt.Errorf("Could not create texture: %v", err)
	}
	defer tex.Destroy()

	if err := renderer.Copy(tex, nil, nil); err != nil {
		return fmt.Errorf("Could not copy texture: %v", err)
	}

	renderer.Present()

	return nil
}

func drawBackground(renderer *sdl.Renderer) error {
	renderer.Clear()

	tex, err := img.LoadTexture(renderer, "res/img/bmp/background.bmp")
	if err != nil {
		return fmt.Errorf("Could not load background image: %v", err)
	}
	defer tex.Destroy()

	if err := renderer.Copy(tex, nil, nil); err != nil {
		return fmt.Errorf("Could not copy background: %v", err)
	}

	renderer.Present()
	return nil
}
