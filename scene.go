package main

import (
	"fmt"
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
	bindArrays(progCube, "position", 3, "color", 3, "texcoord", 2)
	getUniformLocs()
	loadTextures()
	setViewAndProj()
}

func loadTextures() {
	loadTexture(gl.TEXTURE0, "data/kitten.jpg")
	loadTexture(gl.TEXTURE1, "data/gopher.png")
	progCube.GetUniformLocation("tex0").Uniform1i(0)
	progCube.GetUniformLocation("tex1").Uniform1i(1)
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

func bindArrays(program gl.Program, args ...interface{}) {
	var i, j, n, size uint
	size, n = 0, uint(len(args)/2)
	sizes := make([]uint, n)
	names := make([]string, n)
	for i, j = 0, 0; i < n; i, j = i+1, j+2 {
		names[i] = args[j+0].(string)
		sizes[i] = uint(args[j+1].(int))
		size += sizes[i]
		fmt.Println(names[i], sizes[i]) // XXX
	}
	for i, j = 0, 0; i < n; i++ {
		loc := program.GetAttribLocation(names[i])
		fmt.Println(size, j)
		loc.AttribPointer(sizes[i], gl.FLOAT, false, int(size*4), uintptr(j*4))
		loc.EnableArray()
		j += sizes[i]
	}
}

func getUniformLocs() {
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
