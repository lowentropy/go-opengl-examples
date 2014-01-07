package main

import (
	"github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"time"
)

func drawScene() {
	setTime()
	drawCube()
	drawQuad()
	update()
}

func drawCube() {
	// bind framebuffer and scene
	fb.Bind()
	vaoCube.Bind()
	sceneProg.Use()

	clearFrame()
	drawTopCube()
	drawFloor()
	drawBottomCube()
}

func clearFrame() {
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func drawTopCube() {
	activateTextures()
	m = mathgl.HomogRotate3DZ(angle)   // compute rotation matrix
	transform(&m)                      // store
	uniColor.Uniform3f(1, 1, 1)        // top cube is at full brightness
	gl.DrawArrays(gl.TRIANGLES, 0, 36) // draw top cube
}

func drawFloor() {
	drawStencils(1)
	gl.DepthMask(false)                // if this were true, floor would occlude reflection
	transform(&ident)                  // floor doesn't move
	gl.DrawArrays(gl.TRIANGLES, 36, 6) // draw the floor, as well as stencil
	gl.DepthMask(true)                 // bottom cube writes to depth buffer
}

func drawBottomCube() {
	useStencils(1)
	m = mathgl.Scale3D(1, 1, -1).Mul4(m)                // reflect the transform
	m = mathgl.Translate3D(0, 0, -1).Mul4(m)            // and move it below the floor
	uniModel.UniformMatrix4f(false, (*[16]float32)(&m)) // and set as model transform
	uniColor.Uniform3f(0.3, 0.3, 0.3)                   // reflected cube is shady
	gl.DrawArrays(gl.TRIANGLES, 0, 36)                  // draw the reflected cube
	noStencils()
}

func activateTextures() {
	gl.ActiveTexture(gl.TEXTURE0)
	tex0.Bind(gl.TEXTURE_2D)
	gl.ActiveTexture(gl.TEXTURE1)
	tex1.Bind(gl.TEXTURE_2D)
}

func drawQuad() {
	gl.Framebuffer(0).Bind()
	vaoQuad.Bind()
	gl.Disable(gl.DEPTH_TEST)
	screenProg.Use()
	gl.ActiveTexture(gl.TEXTURE0)
	tb.Bind(gl.TEXTURE_2D)
	gl.DrawArrays(gl.TRIANGLES, 0, 6)
}

func update() {
	angle += speed / 30
	speed /= 1.2
}

func setTime() {
	uniTime.Uniform1f(float32(time.Now().Sub(t0).Seconds()))
}
