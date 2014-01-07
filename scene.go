package main

import (
	"github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"time"
)

func initScene() {
	t0 = time.Now()
	initCube()
	initQuad()
}

func initCube() {
	vaoCube = gl.GenVertexArray()
	vaoCube.Bind()
	vboCube := gl.GenBuffer()
	vboCube.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeArray)*4, &cubeArray, gl.STATIC_DRAW)
	loadTextures()
	sceneProg = loadSceneProgram()
	initFramebuffer()
	setViewAndProj(sceneProg)
}

func initQuad() {
	vaoQuad = gl.GenVertexArray()
	vaoQuad.Bind()
	vboQuad := gl.GenBuffer()
	vboQuad.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(quadArray)*4, &quadArray, gl.STATIC_DRAW)
	screenProg = loadScreenProgram()
}

func initFramebuffer() {
	fb = gl.GenFramebuffer()
	fb.Bind()

	tb = gl.GenTexture()
	tb.Bind(gl.TEXTURE_2D)

	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGB, 800, 600, 0, gl.RGB, gl.UNSIGNED_BYTE, nil)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)

	gl.FramebufferTexture2D(gl.FRAMEBUFFER, gl.COLOR_ATTACHMENT0, gl.TEXTURE_2D, tb, 0)

	rb = gl.GenRenderbuffer()
	rb.Bind()
	gl.RenderbufferStorage(gl.RENDERBUFFER, gl.DEPTH24_STENCIL8, 800, 600)
	rb.FramebufferRenderbuffer(gl.FRAMEBUFFER, gl.DEPTH_STENCIL_ATTACHMENT, gl.RENDERBUFFER)
}

func loadTextures() {
	tex0 = loadTexture(gl.TEXTURE0, "data/kitten.jpg")
	tex1 = loadTexture(gl.TEXTURE1, "data/gopher.png")
}

func loadSceneProgram() (program gl.Program) {
	program = loadProgram("shaders/scene.vert", "shaders/scene.frag")
	bindArrays(program, "position", 3, "color", 3, "texcoord", 2)
	program.GetUniformLocation("tex0").Uniform1i(0)
	program.GetUniformLocation("tex1").Uniform1i(1)
	uniTime = program.GetUniformLocation("time")
	uniModel = program.GetUniformLocation("model")
	uniColor = program.GetUniformLocation("overrideColor")
	return
}

func loadScreenProgram() (program gl.Program) {
	program = loadProgram("shaders/screen.vert", "shaders/screen.frag")
	bindArrays(program, "position", 2, "texcoord", 2)
	program.GetUniformLocation("texFramebuffer").Uniform1i(0)
	return
}

func setViewAndProj(program gl.Program) {
	loc := program.GetUniformLocation("view")
	m := mathgl.LookAt(2.5, 2.5, 2, 0, 0, 0, 0, 0, 1)
	loc.UniformMatrix4f(false, (*[16]float32)(&m))

	loc = program.GetUniformLocation("proj")
	m = mathgl.Perspective(45, 800/600, 1, 10)
	loc.UniformMatrix4f(false, (*[16]float32)(&m))
}
