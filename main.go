package main

import (
	"bitbucket.org/zombiezen/math3/mat32"
	"bitbucket.org/zombiezen/math3/vec32"
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"math"
	"os"
	"runtime"
	"time"
	"unsafe"
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
}

var (
	fullscreen = false
	mainCh     = make(chan func(), 10)
	vertices   = [...]float32{
		-0.5, 0.5, 1.0, 0.0, 0.0, 0.0, 0.0, // top left
		0.5, 0.5, 0.0, 1.0, 0.0, 1.0, 0.0, // top right
		0.5, -0.5, 0.0, 0.0, 1.0, 1.0, 1.0, // bot right
		-0.5, -0.5, 1.0, 1.0, 1.0, 0.0, 1.0, // bot left
	}
	elements = [...]uint32{
		0, 1, 2,
		2, 3, 0,
	}
	vao          gl.VertexArray
	vertexShader gl.Shader
	fragShader   gl.Shader
	program      gl.Program
	color        gl.AttribLocation
	t0           time.Time
	window       *glfw.Window
	uniTime      gl.UniformLocation
	trans        gl.UniformLocation
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

	window, err = glfw.CreateWindow(800, 600, "Testing", mon, nil)
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
	loadVertexBuf()
	loadElemBuf()
	loadShaders()
	attachShaders()
	loadTextures()
}

func loadImageData(path string) ([]byte, int, int) {
	f, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	bounds := img.Bounds()
	w := bounds.Max.X - bounds.Min.X
	h := bounds.Max.Y - bounds.Min.Y

	data := make([]byte, w*h*4)

	for i, y := 0, bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; i, x = i+4, x+1 {
			r, g, b, a := img.At(x, y).RGBA()
			data[i+0] = byte(r >> 8)
			data[i+1] = byte(g >> 8)
			data[i+2] = byte(b >> 8)
			data[i+3] = byte(a >> 8)
		}
	}

	return data, w, h
}

func loadTexture(active gl.GLenum, path string) {
	tex := gl.GenTexture()
	gl.ActiveTexture(active)
	tex.Bind(gl.TEXTURE_2D)
	data, w, h := loadImageData(path)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, w, h, 0, gl.RGBA, gl.UNSIGNED_BYTE, data)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
}

func loadTextures() {
	loadTexture(gl.TEXTURE0, "data/kitten.jpg")
	loadTexture(gl.TEXTURE1, "data/gopher.png")
}

func loadVertexBuf() {
	buf := gl.GenBuffer()
	buf.Bind(gl.ARRAY_BUFFER)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, &vertices, gl.STATIC_DRAW)
}

func loadElemBuf() {
	buf := gl.GenBuffer()
	buf.Bind(gl.ELEMENT_ARRAY_BUFFER)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(elements)*4, &elements, gl.STATIC_DRAW)
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
	// bind the shaders and load the program
	program = gl.CreateProgram()
	program.AttachShader(vertexShader)
	program.AttachShader(fragShader)
	program.BindFragDataLocation(0, "outColor")
	program.Link()
	program.Use()

	// position: stride 7, size 2, offset 0
	pos := program.GetAttribLocation("position")
	pos.AttribPointer(2, gl.FLOAT, false, 7*4, uintptr(0))
	pos.EnableArray()

	// color: stride 7, size 3, offset 2
	color = program.GetAttribLocation("color")
	color.AttribPointer(3, gl.FLOAT, false, 7*4, uintptr(2*4))
	color.EnableArray()

	// texcoord: stride 7, size 2, offset 5
	tc := program.GetAttribLocation("texcoord")
	tc.AttribPointer(2, gl.FLOAT, false, 7*4, uintptr(5*4))
	tc.EnableArray()

	// bind texture interpolators
	program.GetUniformLocation("tex0").Uniform1i(0)
	program.GetUniformLocation("tex1").Uniform1i(1)

	// get 'time' and 'trans' variables as global
	uniTime = program.GetUniformLocation("time")
	trans = program.GetUniformLocation("trans")
}

func setMatrix(loc gl.UniformLocation, m *mat32.Matrix) {
	loc.UniformMatrix4f(false, (*[16]float32)((unsafe.Pointer)(m)))
}

func drawScene() {
	// clear the background
	gl.ClearColor(0, 0, 0, 1)
	gl.Clear(gl.COLOR_BUFFER_BIT)

	// set shader 'time' variable
	t := float32(time.Now().Sub(t0).Seconds())
	uniTime.Uniform1f(t)

	// set the shader transform matrix
	m := mat32.Identity.Rotate(t*math.Pi, vec32.Vector{0, 0, 1, 0})
	setMatrix(trans, &m)

	// draw a rectangle
	gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, uintptr(0))
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
