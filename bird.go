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

func newBird(rend *sdl.Renderer) (*bird, error) {
	// Load birds images as an animation
	var textures []*sdl.Texture
	for j := 1; j <= 4; j++ {
		path := fmt.Sprintf("res/img/bmp/bird_frame_%d.bmp", j)
		texture, err := img.LoadTexture(rend, path)
		if err != nil {
			return nil, fmt.Errorf("Could not load bird image: %v", err)
		}
		textures = append(textures, texture)
	}

	return &bird{textures: textures, x: 10, y: 300, w: 50, h: 43}, nil
}

func (b *bird) paint(rend *sdl.Renderer) error {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	// Set position and render bird image
	rect := &sdl.Rect{X: b.x, Y: (600 - b.y) - b.h/2, W: b.w, H: b.h}

	i := b.time / 10 % len(b.textures)
	if err := rend.Copy(b.textures[i], nil, rect); err != nil {
		return fmt.Errorf("Could not copy bird: %v", err)
	}

	return nil
}

func (b *bird) update() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	// Bird physics
	b.time++
	b.y -= int32(b.speed)
	if b.y < 0 {
		b.dead = true
	}
	b.speed += gravity

}

func (b *bird) jump() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.speed = -jumpSpeed
}

func (b *bird) isDead() bool {
	b.mutex.RLock()
	defer b.mutex.RUnlock()

	return b.dead
}

func (b *bird) restart() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	b.y = 300
	b.speed = 0
	b.dead = false
}

func (b *bird) touch(pipe *pipe) {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	if pipe.x > b.x+b.w { // not touching, too far right
		return
	}
	if pipe.x+pipe.w < b.x { //not touching, too far left
		return
	}
	if !pipe.inverted && pipe.h < b.y-b.h/2 { // not touching, pipe is too low
		return
	}
	if pipe.inverted && 600-pipe.h > b.y+b.h/2 { // inverted pipe is too high
		return
	}
	b.dead = true
}

func (b *bird) destroy() {
	b.mutex.Lock()
	defer b.mutex.Unlock()

	for _, tex := range b.textures {
		tex.Destroy()
	}
}
