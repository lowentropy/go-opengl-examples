package main

import (
	// "github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	"time"
)

func drawScene() {
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	gl.DrawArrays(gl.POINTS, 0, 4)
}

func update() {
	angle += speed / 30
	speed /= 1.2
}

func setTime() {
	uniTime.Uniform1f(float32(time.Now().Sub(t0).Seconds()))
}
