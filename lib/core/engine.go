// engine is inspired by https://www.youtube.com/watch?v=pBK-lb-k-rs&list=PLEETnX-uPtBXP_B2yupUKlflXBznWIlL5&index=4
package core

import (
	"fmt"
	"time"

	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/input"
	"github.com/stojg/graphics/lib/rendering"
)

const maxFps time.Duration = 300

func NewEngine(width, height int, title string) (*Engine, error) {

	window, err := NewWindow(width, height, title)
	if err != nil {
		return nil, err
	}

	input.SetWindow(window.Instance())

	engine := &Engine{
		game:            NewGame(),
		window:          window,
		renderingEngine: rendering.NewEngine(),
	}
	engine.game.SetEngine(engine)

	return engine, nil

}

type Engine struct {
	window          *Window
	game            *Game
	renderingEngine *rendering.Engine
	isRunning       bool
}

func (m *Engine) Start() {
	if m.isRunning {
		return
	}
	m.run()
}

func (m *Engine) Game() *Game {
	return m.game
}

func (m *Engine) Stop() {
	if !m.isRunning {
		return
	}
	m.isRunning = false
}

func (m *Engine) run() {
	m.isRunning = true

	var frames int
	var frameCount time.Duration
	const frameTime = time.Second / maxFps

	lastTime := time.Now()
	var unProcessedTime time.Duration

	for m.isRunning {

		render := false

		startTime := time.Now()
		elapsed := startTime.Sub(lastTime)
		lastTime = startTime

		unProcessedTime += elapsed
		frameCount += elapsed

		for unProcessedTime > frameTime {

			render = true

			unProcessedTime -= frameTime

			if m.window.ShouldClose() {
				m.Stop()
			}

			input.Update()

			m.game.Input(frameTime)
			m.game.Update(frameTime)

			if frameCount >= time.Second {
				fmt.Printf("%s, %d fps\n", time.Second/time.Duration(frames), frames)
				frames = 0
				frameCount = 0
			}
		}

		if render {
			m.render()
			frames++
		} else {
			time.Sleep(time.Millisecond)
		}
	}
	m.cleanup()
}

func (m *Engine) render() {
	m.game.Render(m.renderingEngine)
	m.window.Render()
}

func (m *Engine) cleanup() {
	m.window.Close()
}

func (m *Engine) GetRenderingEngine() components.RenderingEngine {
	return m.renderingEngine
}
