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
	pipes      *pipes
}

func newScene(renderer *sdl.Renderer) (*scene, error) {
	// Load background of the scene
	background, err := img.LoadTexture(renderer, "res/img/bmp/background.bmp")
	if err != nil {
		return nil, fmt.Errorf("Could not load background image: %v", err)
	}

	// Load bird for the scene
	bird, err := newBird(renderer)
	if err != nil {
		return nil, err
	}

	// Load pipes for the scene
	pipes, err := newPipes(renderer)
	if err != nil {
		return nil, err
	}

	return &scene{background: background, bird: bird, pipes: pipes}, nil
}

func (scene *scene) paint(renderer *sdl.Renderer) error {
	// Renderer management
	renderer.Clear()
	defer renderer.Present()

	// Render background
	if err := renderer.Copy(scene.background, nil, nil); err != nil {
		return fmt.Errorf("Could not copy background: %v", err)
	}

	// Render animated bird
	if err := scene.bird.paint(renderer); err != nil {
		return err
	}

	// Render pipes
	if err := scene.pipes.paint(renderer); err != nil {
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
				scene.update()
				if scene.bird.isDead() {
					drawTitle(renderer, "Game Over")
					time.Sleep(time.Second)
					scene.restart()
				}
				if err := scene.paint(renderer); err != nil {
					errc <- err
				}
			}
		}
	}()

	return errc
}

func (scene *scene) restart() {
	scene.bird.restart()
	scene.pipes.restart()
}

func (scene *scene) update() {
	scene.bird.update()
	scene.pipes.update()
	scene.pipes.touch(scene.bird)
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

func (scene *scene) destroy() {
	scene.background.Destroy()
	scene.bird.destroy()
	scene.pipes.destroy()
}
