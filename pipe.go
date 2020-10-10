package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/veandco/go-sdl2/img"
	"github.com/veandco/go-sdl2/sdl"
)

type pipe struct {
	mutex sync.RWMutex

	x        int32
	h        int32
	w        int32
	inverted bool
}

type pipes struct {
	mutex sync.RWMutex

	texture *sdl.Texture

	speed int32

	pipes []*pipe
}

//--PIPES--
func newPipes(renderer *sdl.Renderer) (*pipes, error) {
	texture, err := img.LoadTexture(renderer, "res/img/bmp/pipe.bmp")
	if err != nil {
		return nil, fmt.Errorf("Could not load pipe image: %v", err)
	}

	pipes := &pipes{
		texture: texture,
		speed:   2,
	}

	go func() {
		for {
			pipes.mutex.Lock()
			pipes.pipes = append(pipes.pipes, newPipe())
			pipes.mutex.Unlock()
			time.Sleep(time.Second)
		}
	}()

	return pipes, nil

}

func (pipes *pipes) paint(renderer *sdl.Renderer) error {
	pipes.mutex.RLock()
	defer pipes.mutex.RUnlock()

	for _, pipe := range pipes.pipes {
		if err := pipe.paint(renderer, pipes.texture); err != nil {
			return err
		}
	}

	return nil
}

func (pipes *pipes) touch(bird *bird) {
	pipes.mutex.RLock()
	defer pipes.mutex.RUnlock()

	for _, pipe := range pipes.pipes {
		pipe.touch(bird)
	}
}

func (pipes *pipes) restart() {
	pipes.mutex.Lock()
	defer pipes.mutex.Unlock()

	pipes.pipes = nil
}

func (pipes *pipes) update() {
	pipes.mutex.Lock()
	defer pipes.mutex.Unlock()

	var rem []*pipe
	for _, pipe := range pipes.pipes {
		pipe.mutex.Lock()
		pipe.x -= pipes.speed
		pipe.mutex.Unlock()
		if pipe.x+pipe.w > 0 {
			rem = append(rem, pipe)
		}
	}
	pipes.pipes = rem
}

func (pipes *pipes) destroy() {
	pipes.mutex.Lock()
	defer pipes.mutex.Unlock()

	pipes.texture.Destroy()
}

//--PIPE--
func newPipe() *pipe {
	return &pipe{
		x:        windowWidth,
		h:        100 + int32(rand.Intn(300)),
		w:        50,
		inverted: rand.Float32() > 0.5,
	}

}

func (pipe *pipe) paint(renderer *sdl.Renderer, texture *sdl.Texture) error {
	pipe.mutex.RLock()
	defer pipe.mutex.RUnlock()

	rect := &sdl.Rect{X: pipe.x, Y: (windowHeight - pipe.h), W: pipe.w, H: pipe.h}

	flip := sdl.FLIP_NONE
	if pipe.inverted {
		rect.Y = 0
		flip = sdl.FLIP_VERTICAL
	}
	if err := renderer.CopyEx(texture, nil, rect, 0, nil, flip); err != nil {
		return fmt.Errorf("Could not copy pipe: %v", err)
	}

	return nil
}

func (pipe *pipe) touch(bird *bird) {
	pipe.mutex.RLock()
	defer pipe.mutex.RUnlock()

	bird.touch(pipe)
}
