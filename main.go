package main

import (
    "fmt"
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/ttf"
    "os"
    "runtime"
    "time"
)

const (
    width = 800
    height = 600
)

func main() {
    if err := run() ; err != nil {
        fmt.Fprintf(os.Stderr, "%v", err)
        os.Exit(2)
    }
}

func run() error {
    if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
        return fmt.Errorf("could not initialize SDL: %v", err)
    }
    defer sdl.Quit()

    if err := ttf.Init(); err != nil {
        return fmt.Errorf("could not initialize TTF: %v", err)
    }
    defer ttf.Quit()

    w, r, err := sdl.CreateWindowAndRenderer(width, height, sdl.WINDOW_SHOWN)

    if err != nil {
        return fmt.Errorf("could not create window: %v", err)
    }
    sdl.PumpEvents()
    defer func() {
        if err := w.Destroy(); err != nil {
            panic(err)
        }
    }()

    if err := drawTitle(r, "Flappy Gopher"); err != nil {
        return fmt.Errorf("could not draw title: %v", err)
    }

    time.Sleep(500 * time.Millisecond)

    s, err := newScene(r)
    if err != nil {
        return fmt.Errorf("could not create scene: %v", err)
    }
    defer s.destroy()

    events := make(chan sdl.Event)
    errc := s.run(events, r)

    runtime.LockOSThread()
    for {
       select {
        case events <- sdl.WaitEvent():
        case err := <-errc:
            return err
       }

    }
}

func drawTitle(r *sdl.Renderer, text string) error {
    r.Clear()

    f, err := ttf.OpenFont("res/fonts/FlappyBirdy.ttf", 20)
    if err != nil {
        return fmt.Errorf("could not load Flappy font")
    }
    defer f.Close()

    c := sdl.Color{R: 255, G: 100, B: 0, A: 255}
    s, err := f.RenderUTF8Solid(text, c)
    if err != nil {
        return fmt.Errorf("could not render title: %v", err)
    }
    defer s.Free()

    t, err := r.CreateTextureFromSurface(s)
    if err != nil {
        return fmt.Errorf("could not create texture: %v", err)
    }
    defer t.Destroy()

    if err := r.Copy(t, nil, nil) ; err != nil {
        return fmt.Errorf("could not copy textyre: %v", err)
    }

    r.Present()

    return nil
}

