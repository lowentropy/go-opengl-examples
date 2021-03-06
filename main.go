package main

import (
	"github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"runtime"
	"time"
)

var (
	ident      = mathgl.Ident4f()
	fullscreen = false
	mainCh     = make(chan func(), 10)

	angle = float32(-45)
	speed = float32(0)

	m            mathgl.Mat4f
	vertexShader gl.Shader
	fragShader   gl.Shader
	color        gl.AttribLocation
	t0           time.Time
	window       *glfw.Window

	uniColor, uniTime, uniModel gl.UniformLocation
	vaoCube, vaoQuad            gl.VertexArray
	progCube, progQuad          gl.Program
)

func main() {
	runtime.LockOSThread()
	initGlfw()
	initGl()
	initScene()
	for !window.ShouldClose() {
		loop()
	}
	glfw.Terminate()
}

func loop() {
	flushMain()
	if window.GetKey(glfw.KeyEscape) == glfw.Press {
		window.SetShouldClose(true)
		return
	} else if window.GetKey(glfw.KeySpace) == glfw.Press {
		speed = 180
	}
	drawScene()
	window.SwapBuffers()
	glfw.PollEvents()
}

// pass a function here to have it executed in
// the main thread
func doMain(f func()) {
	done := make(chan bool, 1)
	mainCh <- func() {
		f()
		done <- true
	}
	<-done
}

// actually run the functions passed to doMain,
// until the buffer is empty
func flushMain() {
	for {
		select {
		case f := <-mainCh:
			f()
		default:
			return
		}
	}
}
