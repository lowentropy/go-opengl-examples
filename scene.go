package main

import (
	"github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"time"
)

func initScene() {
	t0 = time.Now()
	initCube()
}

func initCube() {
	vaoCube = gl.GenVertexArray()
	vaoCube.Bind()
	gl.Enable(gl.DEPTH_TEST)
	loadVertexBuf()
	progCube = loadProgram("shaders/vertex.vert", "shaders/fragment.frag")
	attachShaders()
	loadTextures()
	setViewAndProj()
}

func loadTextures() {
	loadTexture(gl.TEXTURE0, "data/kitten.jpg")
	loadTexture(gl.TEXTURE1, "data/gopher.png")
}

func loadVertexBuf() {
	buf := gl.GenBuffer()
	buf.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(cubeVertices)*4, &cubeVertices, gl.STATIC_DRAW)
}

func loadProgram(vertPath, fragPath string) gl.Program {
	vertexShader = loadShader(gl.VERTEX_SHADER, vertPath)
	fragShader = loadShader(gl.FRAGMENT_SHADER, fragPath)
	program := gl.CreateProgram()
	program.AttachShader(vertexShader)
	program.AttachShader(fragShader)
	program.BindFragDataLocation(0, "outColor")
	program.Link()
	program.Use()
	return program
}

func attachShaders() {
	// position: stride 8, size 3, offset 0
	pos := progCube.GetAttribLocation("position")
	pos.AttribPointer(3, gl.FLOAT, false, 8*4, uintptr(0))
	pos.EnableArray()

	// color: stride 8, size 3, offset 3
	color = progCube.GetAttribLocation("color")
	color.AttribPointer(3, gl.FLOAT, false, 8*4, uintptr(3*4))
	color.EnableArray()

	// texcoord: stride 8, size 2, offset 6
	tc := progCube.GetAttribLocation("texcoord")
	tc.AttribPointer(2, gl.FLOAT, false, 8*4, uintptr(6*4))
	tc.EnableArray()

	// bind texture interpolators
	progCube.GetUniformLocation("tex0").Uniform1i(0)
	progCube.GetUniformLocation("tex1").Uniform1i(1)

	// get 'time', 'model', and 'overrideColor' variables as global
	uniTime = progCube.GetUniformLocation("time")
	uniModel = progCube.GetUniformLocation("model")
	uniColor = progCube.GetUniformLocation("overrideColor")
}

func setViewAndProj() {
	loc := progCube.GetUniformLocation("view")
	m := mathgl.LookAt(2.5, 2.5, 2, 0, 0, 0, 0, 0, 1)
	loc.UniformMatrix4f(false, (*[16]float32)(&m))

	loc = progCube.GetUniformLocation("proj")
	m = mathgl.Perspective(45, 800/600, 1, 10)
	loc.UniformMatrix4f(false, (*[16]float32)(&m))
}
