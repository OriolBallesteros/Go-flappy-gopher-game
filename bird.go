package main

import (
	"fmt"
	"sync"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	gravity   = 0.25
	jumpSpeed = 5
)

type bird struct {
	mutex    sync.RWMutex //handle different goRoutines accessing bird; at write, access locked for a moment
	time     int
	textures []*sdl.Texture

	x, y  int32
	w, h  int32
	speed float64
	dead  bool
}

func newBird(renderer *sdl.Renderer) (*bird, error) {
	// Load birds images as an animation
	var textures []*sdl.Texture
	for j := 1; j <= 4; j++ {
		path := fmt.Sprintf("res/img/bmp/bird_frame_%d.bmp", j)
		texture, err := img.LoadTexture(renderer, path)
		if err != nil {
			return nil, fmt.Errorf("Could not load bird image: %v", err)
		}
		textures = append(textures, texture)
	}

	return &bird{textures: textures, x: 10, y: 300, w: 50, h: 43}, nil
}

func (bird *bird) paint(renderer *sdl.Renderer) error {
	bird.mutex.RLock()
	defer bird.mutex.RUnlock()

	// Set position and render bird image
	rect := &sdl.Rect{X: bird.x, Y: (windowHeight - bird.y) - bird.h/2, W: bird.w, H: bird.h}

	i := bird.time / 10 % len(bird.textures)
	if err := renderer.Copy(bird.textures[i], nil, rect); err != nil {
		return fmt.Errorf("Could not copy bird: %v", err)
	}

	return nil
}

func (bird *bird) update() {
	bird.mutex.Lock()
	defer bird.mutex.Unlock()

	// Bird physics
	bird.time++
	bird.y -= int32(bird.speed)
	if bird.y < 0 {
		bird.dead = true
	}
	bird.speed += gravity

}

func (bird *bird) jump() {
	bird.mutex.Lock()
	defer bird.mutex.Unlock()

	bird.speed = -jumpSpeed
}

func (bird *bird) isDead() bool {
	bird.mutex.RLock()
	defer bird.mutex.RUnlock()

	return bird.dead
}

func (bird *bird) restart() {
	bird.mutex.Lock()
	defer bird.mutex.Unlock()

	bird.y = 300
	bird.speed = 0
	bird.dead = false
}

func (bird *bird) touch(pipe *pipe) {
	bird.mutex.Lock()
	defer bird.mutex.Unlock()

	if pipe.x > bird.x+bird.w { // not touching, too far right
		return
	}
	if pipe.x+pipe.w < bird.x { //not touching, too far left
		return
	}
	if !pipe.inverted && pipe.h < bird.y-bird.h/2 { // not touching, pipe is too low
		return
	}
	if pipe.inverted && windowHeight-pipe.h > bird.y+bird.h/2 { // inverted pipe is too high
		return
	}
	bird.dead = true
}

func (bird *bird) destroy() {
	bird.mutex.Lock()
	defer bird.mutex.Unlock()

	for _, texture := range bird.textures {
		texture.Destroy()
	}
}
