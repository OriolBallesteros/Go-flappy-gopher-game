package main

import (
	"context"
	"fmt"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	time int

	background *sdl.Texture
	birds      []*sdl.Texture
}

func newScene(rend *sdl.Renderer) (*scene, error) {
	// Load background of the scene
	bg, err := img.LoadTexture(rend, "res/img/bmp/background.bmp")
	if err != nil {
		return nil, fmt.Errorf("Could not load background image: %v", err)
	}

	// Load birds images as an animation
	var birds []*sdl.Texture
	for n := 1; n <= 4; n++ {
		path := fmt.Sprintf("res/img/bmp/bird_frame_%d.bmp", n)
		bird, err := img.LoadTexture(rend, path)
		if err != nil {
			return nil, fmt.Errorf("Could not load bird image: %v", err)
		}
		birds = append(birds, bird)
	}

	return &scene{background: bg, birds: birds}, nil
}

func (scene *scene) paint(rend *sdl.Renderer) error {
	scene.time++

	// Renderer management
	rend.Clear()
	defer rend.Present()

	// Render background
	if err := rend.Copy(scene.background, nil, nil); err != nil {
		return fmt.Errorf("Could not copy background: %v", err)
	}

	// Set position and render bird image
	rect := &sdl.Rect{X: 10, Y: 300 - 43/2, W: 50, H: 43}
	i := scene.time / 10 % len(scene.birds)
	if err := rend.Copy(scene.birds[i], nil, rect); err != nil {
		return fmt.Errorf("Could not copy bird: %v", err)
	}

	return nil
}

func (scene *scene) destroy() {
	scene.background.Destroy()
}

func (scene *scene) run(ctx context.Context, renderer *sdl.Renderer) chan error {
	errc := make(chan error)

	go func() {
		defer close(errc)
		for range time.Tick(10 * time.Millisecond) {
			select {
			case <-ctx.Done():
				return
			default:
				if err := scene.paint(renderer); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}
