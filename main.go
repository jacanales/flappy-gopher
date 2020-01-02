package main

import (
    "context"
    "fmt"
    "github.com/veandco/go-sdl2/img"
    "github.com/veandco/go-sdl2/sdl"
    "github.com/veandco/go-sdl2/ttf"
    "os"
    "time"
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

    w, r, err := sdl.CreateWindowAndRenderer(800, 600, sdl.WINDOW_SHOWN)

    if err != nil {
        return fmt.Errorf("could not create window: %v", err)
    }
    sdl.PumpEvents()
    defer w.Destroy()

    if err := drawTitle(r); err != nil {
        return fmt.Errorf("could not draw title: %v", err)
    }

    time.Sleep(time.Second)

    s, err := newScene(r)
    if err != nil {
        return fmt.Errorf("could not create scene: %v", err)
    }
    defer s.destroy()

    ctx, cancel := context.WithCancel(context.Background())
    time.AfterFunc(5 * time.Second, cancel)

    return <-s.run(ctx, r)

    /*
    running := true
    for running {
        for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
            switch event.(type) {
            case *sdl.QuitEvent:
                println("Quit")
                running = false
                break
            default:
                err := <-s.run(ctx, r)
                if err != nil {
                    return err
                }
            }
        }
    }

    return nil
    */
}

func drawBackground(r *sdl.Renderer) error {
    r.Clear()

    t, err := img.LoadTexture(r, "res/img/background.png")
    if err != nil {
        return fmt.Errorf("could not load background image: %v", err)
    }
    defer t.Destroy()

    if err := r.Copy(t, nil, nil); err != nil {
        return fmt.Errorf("could not copy background", err)
    }

    r.Present()
    return nil
}

func drawTitle(r *sdl.Renderer) error {
    r.Clear()

    f, err := ttf.OpenFont("res/fonts/FlappyBirdy.ttf", 20)
    if err != nil {
        return fmt.Errorf("could not load Flappy font")
    }
    defer f.Close()

    c := sdl.Color{R: 255, G: 100, B: 0, A: 255}
    s, err := f.RenderUTF8Solid("Flappy Gopher", c)
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

