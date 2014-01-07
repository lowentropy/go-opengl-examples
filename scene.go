package main

import (
	// "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"time"
)

func initScene() {
	t0 = time.Now()
	prog := loadProgram("shaders/geom.vert", "shaders/geom.frag", "shaders/geom.geom")
	vbo := gl.GenBuffer()
	vbo.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(geom)*4, &geom, gl.STATIC_DRAW)
	vao := gl.GenVertexArray()
	vao.Bind()
	bindArrays(prog, "pos", 2, "color", 3, "sides", 1)
}
