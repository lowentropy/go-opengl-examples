package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"io/ioutil"
	"runtime"
	"time"
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

var (
	fullscreen = false
	mainCh     = make(chan func(), 10)
	vertices   = [...]float32{
		0.0, 0.5, 1.0, 0.0, 0.0,
		0.5, -0.5, 0.0, 1.0, 0.0,
		-0.5, -0.5, 0.0, 0.0, 1.0,
	}
	vao          gl.VertexArray
	buf          gl.Buffer
	vertexShader gl.Shader
	fragShader   gl.Shader
	program      gl.Program
	color        gl.AttribLocation
	t0           time.Time
	window       *glfw.Window
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
	}
	drawScene()
	window.SwapBuffers()
	glfw.PollEvents()
}

func initGlfw() {
	glfw.SetErrorCallback(errorCallback)

	if !glfw.Init() {
		panic("Can't init glfw!")
	}

	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 2)
	glfw.WindowHint(glfw.OpenglForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.OpenglProfile, glfw.OpenglCoreProfile)

	var mon *glfw.Monitor
	var err error

	if fullscreen {
		mon, err = glfw.GetPrimaryMonitor()
		if err != nil {
			panic(err)
		}
	}

	window, err = glfw.CreateWindow(640, 480, "Testing", mon, nil)
	if err != nil {
		panic(err)
	}

	window.MakeContextCurrent()
}

func initGl() {
	gl.Init()
	fmt.Println("OpenGL Version:", gl.GetString(gl.VERSION))
	fmt.Println("Shader Version:", gl.GetString(gl.SHADING_LANGUAGE_VERSION))
}

func initScene() {
	t0 = time.Now()
	vao = gl.GenVertexArray()
	vao.Bind()
	loadBuf()
	loadShaders()
	attachShaders()
}

func loadBuf() {
	buf = gl.GenBuffer()
	buf.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, &vertices, gl.STATIC_DRAW)
}

func loadShaders() {
	vertexShader = loadShader(gl.VERTEX_SHADER, "shaders/vertex.vert")
	fragShader = loadShader(gl.FRAGMENT_SHADER, "shaders/fragment.frag")
}

func loadShader(kind gl.GLenum, path string) gl.Shader {
	shader := gl.CreateShader(kind)
	source, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	shader.Source(string(source))
	shader.Compile()
	if shader.Get(gl.COMPILE_STATUS) != gl.TRUE {
		panic(fmt.Sprintf("can't compile shader: %s: %s", path, shader.GetInfoLog()))
	}
	return shader
}

func attachShaders() {
	program = gl.CreateProgram()
	program.AttachShader(vertexShader)
	program.AttachShader(fragShader)
	program.BindFragDataLocation(0, "outColor")
	program.Link()
	program.Use()

	// pos is the first two floats on a five-float line (stride 5, offset 0)
	pos := program.GetAttribLocation("position")
	pos.AttribPointer(2, gl.FLOAT, false, 5*4, uintptr(0))
	pos.EnableArray()

	// color is the last three floats on a five-float line (stride 5, offset 2)
	color = program.GetAttribLocation("color")
	color.AttribPointer(3, gl.FLOAT, false, 5*4, uintptr(2*4))
	color.EnableArray()
}

func drawScene() {
	gl.DrawArrays(gl.TRIANGLES, 0, 3)
}

func do(f func()) {
	done := make(chan bool, 1)
	mainCh <- func() {
		f()
		done <- true
	}
	<-done
}

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
