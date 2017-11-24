package rendering

import (
	"fmt"

	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/stojg/graphics/lib/components"
	"github.com/stojg/graphics/lib/rendering/framebuffer"
)

func NewEngine(width, height int) *Engine {

	gl.ClearColor(0.00, 0.00, 0.00, 1)

	gl.FrontFace(gl.CCW)
	gl.CullFace(gl.BACK)
	gl.Enable(gl.CULL_FACE)
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)

	gl.Disable(gl.MULTISAMPLE)
	gl.Disable(gl.FRAMEBUFFER_SRGB)

	samplerMap := make(map[string]uint32)
	samplerMap["diffuse"] = 0
	samplerMap["x_shadowMap"] = 9

	return &Engine{
		width:      int32(width),
		height:     int32(height),
		samplerMap: samplerMap,

		screenQuad:   NewScreenQuad(),
		screenShader: NewShader("screen_shader"),

		ambientShader: NewShader("forward_ambient"),

		hdrTexture: framebuffer.NewTexture(0, width, height, gl.RGBA32F, gl.RGBA, gl.FLOAT, gl.LINEAR, false),
		hdrShader:  NewShader("screen_hdr"),
	}
}

type Engine struct {
	width, height int32
	mainCamera    *components.Camera
	lights        []components.Light
	activeLight   components.Light

	samplerMap map[string]uint32

	screenQuad   *ScreenQuad
	screenShader *Shader

	ambientShader *Shader

	hdrTexture *framebuffer.Texture
	hdrShader  *Shader
}

func (e *Engine) Render(object components.Renderable) {
	if e.mainCamera == nil {
		panic("mainCamera not found, the game cannot render")
	}
	checkForError("renderer.Engine.Render [start]")

	// shadow map
	gl.Enable(gl.DEPTH_TEST)
	gl.CullFace(gl.FRONT)
	for _, l := range e.lights {
		caster, ok := l.(components.ShadowCaster)
		if !ok {
			continue
		}
		e.activeLight = l
		caster.BindAsRenderTarget()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		object.RenderAll(caster.ShadowShader(), e)
		// debug
		//gl.Viewport(0, 0, e.width, e.height)
		//gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
		//gl.Disable(gl.DEPTH_TEST)
		//caster.BindShadow()
		//e.screenShader.Bind()
		//gl.Clear(gl.COLOR_BUFFER_BIT)
		//e.screenQuad.Draw()
		//return
	}
	gl.CullFace(gl.BACK)

	e.hdrTexture.BindAsRenderTarget()
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
	e.hdrTexture.SetViewPort()

	// ambient pass
	object.RenderAll(e.ambientShader, e)

	// light pass
	gl.Enable(gl.BLEND)
	gl.BlendFunc(gl.ONE, gl.ONE)
	gl.DepthMask(false)
	gl.DepthFunc(gl.EQUAL)

	for _, l := range e.lights {
		e.activeLight = l
		l.Shader().Bind()
		if caster, ok := l.(components.ShadowCaster); ok {
			caster.BindShadowTexture(e.GetSamplerSlot("x_shadowMap"), "x_shadowMap")
		}
		object.RenderAll(l.Shader(), e)
	}
	gl.DepthFunc(gl.LESS)
	gl.DepthMask(true)
	gl.Disable(gl.BLEND)

	// move to default framebuffer buffer
	gl.BindFramebuffer(gl.FRAMEBUFFER, 0)
	gl.Viewport(0, 0, e.width, e.height)
	// disable depth test so screen-space quad isn't discarded due to depth test
	gl.Disable(gl.DEPTH_TEST)
	e.hdrTexture.Bind()
	e.hdrShader.Bind()
	gl.Clear(gl.COLOR_BUFFER_BIT)
	e.screenQuad.Draw()

	checkForError("renderer.Engine.Render [end]")
}

func (e *Engine) GetActiveLight() components.Light {
	return e.activeLight
}

func (e *Engine) AddLight(l components.Light) {
	e.lights = append(e.lights, l)
}

func (e *Engine) AddCamera(c *components.Camera) {
	e.mainCamera = c
}

func (e *Engine) GetMainCamera() *components.Camera {
	return e.mainCamera
}

func (e *Engine) GetSamplerSlot(samplerName string) uint32 {
	slot, exists := e.samplerMap[samplerName]
	if !exists {
		fmt.Printf("rendering.Engine tried finding texture slot for %s, failed\n", samplerName)
	}
	return slot
}
