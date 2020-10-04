package main

import (
	"fmt"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

const (
	gravity   = 0.25
	jumpSpeed = 5
)

type bird struct {
	time     int
	textures []*sdl.Texture

	y, speed float64
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

	return &bird{textures: textures, y: 300}, nil
}

func (b *bird) paint(rend *sdl.Renderer) error {
	// Bird physics
	b.y -= b.speed
	if b.y < 0 {
		b.y = 0
		b.speed = -b.speed
	}
	b.speed += gravity

	// Set position and render bird image
	rect := &sdl.Rect{X: 10, Y: (600 - int32(b.y)) - 43/2, W: 50, H: 43}

	b.time++
	i := b.time / 10 % len(b.textures)
	if err := rend.Copy(b.textures[i], nil, rect); err != nil {
		return fmt.Errorf("Could not copy bird: %v", err)
	}

	return nil
}

func (b *bird) destroy() {
	for _, tex := range b.textures {
		tex.Destroy()
	}
}

func (b *bird) jump() {
	b.speed = -jumpSpeed
}
