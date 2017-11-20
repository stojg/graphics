package components

import (
	"time"

	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/input"
)

func NewFreelook(width, height float32) *FreeLook {
	return &FreeLook{
		centerPosition: mgl32.Vec2{width / 2, height / 2},
	}
}

type FreeLook struct {
	GameComponent

	centerPosition mgl32.Vec2
	locked         bool
}

func (c *FreeLook) Input(elapsed time.Duration) {

	if input.ButtonDown(glfw.MouseButton1) {
		c.locked = true
		input.HideCursor()
		c.centerCamera()
	}

	if input.KeyDown(glfw.KeySpace) {
		c.locked = false
		c.centerCamera()
		input.ShowCursor()
	}

	if !c.locked {
		return
	}

	delta := mgl32.Vec2(input.CursorPosition()).Sub(c.centerPosition)
	if delta.Len() == 0 {
		return
	}
	c.centerCamera()

	const sensitivity float32 = 0.5

	yaw := mgl32.DegToRad(-delta[0]) * sensitivity
	pitch := mgl32.DegToRad(delta[1]) * sensitivity

	rotation := c.Transform().Rot().Mul(mgl32.QuatRotate(pitch, left()))
	rotation = mgl32.QuatRotate(yaw, up()).Mul(rotation)
	c.Transform().SetRot(rotation)
}

func (c *FreeLook) centerCamera() {
	input.SetCursorPosition(c.centerPosition[0], c.centerPosition[1])
}