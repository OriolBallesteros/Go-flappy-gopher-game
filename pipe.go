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
func newPipes(rend *sdl.Renderer) (*pipes, error) {
	tex, err := img.LoadTexture(rend, "res/img/bmp/pipe.bmp")
	if err != nil {
		return nil, fmt.Errorf("Could not load pipe image: %v", err)
	}

	pipes := &pipes{
		texture: tex,
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

func (ps *pipes) paint(rend *sdl.Renderer) error {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	for _, p := range ps.pipes {
		if err := p.paint(rend, ps.texture); err != nil {
			return err
		}
	}

	return nil
}

func (ps *pipes) touch(b *bird) {
	ps.mutex.RLock()
	defer ps.mutex.RUnlock()

	for _, p := range ps.pipes {
		p.touch(b)
	}
}

func (ps *pipes) restart() {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	ps.pipes = nil
}

func (ps *pipes) update() {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	var rem []*pipe
	for _, p := range ps.pipes {
		p.mutex.Lock()
		p.x -= ps.speed
		p.mutex.Unlock()
		if p.x+p.w > 0 {
			rem = append(rem, p)
		}
	}
	ps.pipes = rem
}

func (ps *pipes) destroy() {
	ps.mutex.Lock()
	defer ps.mutex.Unlock()

	ps.texture.Destroy()
}

//--PIPE--
func newPipe() *pipe {
	return &pipe{
		x:        800,
		h:        100 + int32(rand.Intn(300)),
		w:        50,
		inverted: rand.Float32() > 0.5,
	}

}

func (p *pipe) paint(rend *sdl.Renderer, tex *sdl.Texture) error {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	rect := &sdl.Rect{X: p.x, Y: (600 - p.h), W: p.w, H: p.h}

	flip := sdl.FLIP_NONE
	if p.inverted {
		rect.Y = 0
		flip = sdl.FLIP_VERTICAL
	}
	if err := rend.CopyEx(tex, nil, rect, 0, nil, flip); err != nil {
		return fmt.Errorf("Could not copy pipe: %v", err)
	}

	return nil
}

func (p *pipe) touch(b *bird) {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	b.touch(p)
}
