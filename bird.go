package main

import (
    "fmt"
    "github.com/veandco/go-sdl2/img"
    "github.com/veandco/go-sdl2/sdl"
    "sync"
)

const (
    gravity = 0.10
    jumpSpeed = 5
)


type bird struct {
    mu sync.RWMutex

    time int
    textures []*sdl.Texture
    x, y int32
    w, h int32
    speed float64
    dead bool
}

func newBird(r *sdl.Renderer) (*bird, error) {
    var textures []*sdl.Texture
    for i := 1; i <= 4; i++ {
        path := fmt.Sprintf("res/img/bird_frame_%d.png", i)
        texture, err := img.LoadTexture(r, path)
        if err != nil {
            return nil, fmt.Errorf("could not load bird frame image: %v", err)
        }

        textures = append(textures, texture)
    }

    return &bird{textures: textures, x: 10, y: width/2, w: 50, h: 43}, nil
}

func (b *bird) update() {
    b.mu.Lock()
    defer b.mu.Unlock()

    b.time++
    b.y -= int32(b.speed)
    if b.y < 0 {
        b.dead = true
        b.speed = -b.speed
        b.y = 0
    }
    b.speed += gravity
}

func (b *bird) paint(r *sdl.Renderer) error {
    b.mu.RLock()
    defer b.mu.RUnlock()

    rect := &sdl.Rect{X: b.x, Y: height - b.y - b.h/2, W: b.w, H: b.h}

    i := b.time/10 % len(b.textures)
    if err := r.Copy(b.textures[i], nil, rect); err != nil {
        return fmt.Errorf("could not copy bird", err)
    }

    if b.y > height - b.x * 2 {
        b.y = height - b.h
        b.speed = 0
    }

    return nil
}

func (b *bird) restart() {
    b.mu.Lock()
    defer b.mu.Unlock()

    b.y = width/2
    b.speed = 0
    b.dead = false
}

func (b *bird) destroy() {
    b.mu.Lock()
    defer b.mu.Unlock()

    for _, t := range b.textures {
        t.Destroy()
    }
}

func (b *bird) jump() {
    b.mu.Lock()
    defer b.mu.Unlock()

    b.speed = -jumpSpeed
}

func (b *bird) isDead() bool {
    b.mu.RLock()
    defer b.mu.RUnlock()

    return b.dead
}

func (b *bird) touch (p * pipe) {
    b.mu.Lock()
    defer b.mu.Unlock()
    p.mu.RLock()
    defer p.mu.RUnlock()

    if p.x > b.x + b.w { // too far right
        return
    }

    if p.x + p.w < b.x { // too far lefts
        return
    }

    if !p.inverted && p.h < b.y - (b.h/2) { // pipe is too low
        return
    }

    if p.inverted && height - p.h > b.y+b.h/2 { // pipe is too high
        return
    }
    b.dead = true
}