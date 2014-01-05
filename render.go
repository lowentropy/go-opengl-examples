package main

import (
	"github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"time"
)

func drawScene() {
	setTime()
	clearFrame()
	drawTopCube()
	drawFloor()
	drawBottomCube()
	update()
}

func drawTopCube() {
	m = mathgl.HomogRotate3DZ(angle)                    // compute rotation matrix
	uniModel.UniformMatrix4f(false, (*[16]float32)(&m)) // store as model transform
	gl.DepthMask(true)                                  // top cube writes to depth buffer
	uniColor.Uniform3f(1, 1, 1)                         // top cube is at full brightness
	gl.Disable(gl.STENCIL_TEST)                         // top cube always drawn regardless of stencils
	gl.DrawArrays(gl.TRIANGLES, 0, 36)                  // draw top cube
}

func drawFloor() {
	gl.Clear(gl.STENCIL_BUFFER_BIT)                         // clear stencil (so top cube doesn't matter)
	gl.Enable(gl.STENCIL_TEST)                              // enable the test
	gl.StencilFunc(gl.ALWAYS, 1, 0xff)                      // the floor will always write 1's
	gl.StencilOp(gl.KEEP, gl.KEEP, gl.REPLACE)              // replace w/ 1 for floor writes
	gl.StencilMask(0xff)                                    // has no effect, but needed 'cause of the 0x00 below
	gl.DepthMask(false)                                     // if this were true, floor would occlude reflection
	uniModel.UniformMatrix4f(false, (*[16]float32)(&ident)) // floor doesn't move
	gl.DrawArrays(gl.TRIANGLES, 36, 6)                      // draw the floor, as well as stencil
}

func drawBottomCube() {
	gl.StencilFunc(gl.EQUAL, 1, 0xff) // now set to only draw on stencil 1's
	gl.StencilMask(0xff)              // but don't affect the stencil itself (not that it matters, since we clear the stencil)

	m = mathgl.Scale3D(1, 1, -1).Mul4(m)                // reflect the transform
	m = mathgl.Translate3D(0, 0, -1).Mul4(m)            // and move it below the floor
	uniModel.UniformMatrix4f(false, (*[16]float32)(&m)) // and set as model transform

	gl.DepthMask(true)                 // bottom cube writes to depth buffer
	uniColor.Uniform3f(0.3, 0.3, 0.3)  // reflected cube is shady
	gl.DrawArrays(gl.TRIANGLES, 0, 36) // draw the reflected cube
}

func clearFrame() {
	gl.Enable(gl.DEPTH_TEST)
	gl.ClearColor(1, 1, 1, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
}

func setTime() {
	uniTime.Uniform1f(float32(time.Now().Sub(t0).Seconds()))
}

func update() {
	angle += speed / 30
	speed /= 1.2
}
