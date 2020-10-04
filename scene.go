package main

import (
	"fmt"
	"log"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type scene struct {
	background *sdl.Texture
	bird       *bird
}

func newScene(rend *sdl.Renderer) (*scene, error) {
	// Load background of the scene
	bg, err := img.LoadTexture(rend, "res/img/bmp/background.bmp")
	if err != nil {
		return nil, fmt.Errorf("Could not load background image: %v", err)
	}

	// Load bird for the scene
	bird, err := newBird(rend)
	if err != nil {
		return nil, err
	}

	return &scene{background: bg, bird: bird}, nil
}

func (scene *scene) paint(rend *sdl.Renderer) error {
	// Renderer management
	rend.Clear()
	defer rend.Present()

	// Render background
	if err := rend.Copy(scene.background, nil, nil); err != nil {
		return fmt.Errorf("Could not copy background: %v", err)
	}

	// Render animated bird
	if err := scene.bird.paint(rend); err != nil {
		return err
	}

	return nil
}

func (scene *scene) run(events <-chan sdl.Event, renderer *sdl.Renderer) <-chan error {
	errc := make(chan error)

	// Check event and/or paint
	go func() {
		defer close(errc)
		tick := time.Tick(10 * time.Millisecond)
		for {
			select {
			case event := <-events:
				// Check events, close down if it is Quit
				if done := scene.handleEvent(event); done {
					return
				}
			case <-tick:
				if err := scene.paint(renderer); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (scene *scene) destroy() {
	scene.background.Destroy()
	scene.bird.destroy()
}

// Play game interactions
func (scene *scene) handleEvent(event sdl.Event) bool {
	switch event.(type) {

	case *sdl.QuitEvent:
		return true

	case *sdl.MouseButtonEvent:
		scene.bird.jump()

	case *sdl.MouseMotionEvent, *sdl.AudioDeviceEvent, *sdl.WindowEvent:
		// ignored events
	default:
		log.Printf("Unkown event %T", event)

	}
	return false
}
