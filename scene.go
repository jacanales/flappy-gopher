package main

import (
    "fmt"
    "github.com/veandco/go-sdl2/img"
    "github.com/veandco/go-sdl2/sdl"
    "log"
    "time"
)

type scene struct {
    bg *sdl.Texture
    bird *bird
}

func newScene(r *sdl.Renderer) (*scene, error) {
    bg, err := img.LoadTexture(r, "res/img/background.png")
    if err != nil {
        return nil, fmt.Errorf("could not load background image: %v", err)
    }

    b, err := newBird(r)
    if err != nil {
        return nil, err
    }

    return &scene{bg: bg, bird: b}, nil
}

func (s *scene) run(events <-chan sdl.Event, r *sdl.Renderer) chan error {
    errc := make(chan error)

    go func() {
        defer close(errc)
        tick := time.Tick(10 * time.Millisecond)
        done := false

        for {
            select {
            case e := <-events:
                if done = s.handleEvent(e); done {
                    return
                }
            case <-tick:
                if err := s.paint(r); err != nil {
                    errc <- err
                }
            }
        }
    }()

    return errc
}

func (s *scene) handleEvent(event sdl.Event) (bool) {
    switch event.(type) {
    case *sdl.QuitEvent:
        return true
    case *sdl.MouseButtonEvent:
        s.bird.jump()
    case *sdl.MouseMotionEvent, *sdl.WindowEvent, *sdl.MouseWheelEvent, *sdl.AudioDeviceEvent,  *sdl.TextInputEvent:
    default:
        log.Printf("unknown event %T", event)
    }

    return false
}

func (s *scene) paint(r *sdl.Renderer) error {
    r.Clear()

    if err := r.Copy(s.bg, nil, nil); err != nil {
        return fmt.Errorf("could not copy background", err)
    }

    if err := s.bird.paint(r); err != nil {
        return nil
    }

    r.Present()
    return nil
}

func (s *scene) destroy() {
    s.bg.Destroy()
    s.bird.destroy()
}