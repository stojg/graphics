package lights

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering"
	"github.com/stojg/graphics/lib/rendering/framebuffer"
)

func NewDirectional(r, g, b, intensity float32) *Directional {
	return &Directional{
		BaseLight: BaseLight{
			color:  mgl32.Vec3{r, g, b}.Mul(intensity),
			shader: rendering.NewShader("forward_directional"),
		},
		shadowBuffer: framebuffer.NewShadow(1024, 1024),
		shadowShader: rendering.NewShader("shadow"),
	}
}

type Directional struct {
	BaseLight

	shadowBuffer *framebuffer.FBO
	shadowShader components.Shader
}

func (b *Directional) AddToEngine(e components.Engine) {
	e.GetRenderingEngine().AddLight(b)
}

func (b *Directional) Direction() mgl32.Vec3 {
	return b.BaseLight.Position().Normalize()
}

func (b *Directional) ViewProjection() mgl32.Mat4 {
	const nearPlane float32 = 0.1
	const farPlane float32 = 10
	lightProjection := mgl32.Ortho(-8, 8, -8, 8, nearPlane, farPlane)
	lightView := mgl32.LookAt(b.Position().X(), b.Position().Y(), b.Position().Z(), 0, 0, 0, 0, 1, 0)
	return lightProjection.Mul4(lightView)
}

func (b *Directional) BindShadowBuffer() {
	b.shadowBuffer.Bind()
	b.shadowBuffer.Texture().SetViewPort()
	gl.Enable(gl.DEPTH_TEST)
	gl.Clear(gl.DEPTH_BUFFER_BIT)
}

func (b *Directional) ShadowShader() components.Shader {
	return b.shadowShader
}

func (b *Directional) BindShadowTexture(samplerSlot uint32, samplerName string) {
	gl.ActiveTexture(gl.TEXTURE0 + uint32(samplerSlot))
	b.Shader().SetUniform(samplerName, int32(samplerSlot))
	gl.BindTexture(gl.TEXTURE_2D, b.shadowBuffer.Texture().ID())
}
