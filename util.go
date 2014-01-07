package main

import (
	"fmt"
	"github.com/Jragonmiris/mathgl"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
)

func errorCallback(err glfw.ErrorCode, desc string) {
	fmt.Printf("%v: %v\n", err, desc)
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

func loadTexture(active gl.GLenum, path string) gl.Texture {
	tex := gl.GenTexture()
	gl.ActiveTexture(active)
	tex.Bind(gl.TEXTURE_2D)
	data, w, h := loadImageData(path)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, w, h, 0, gl.RGBA, gl.UNSIGNED_BYTE, data)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	return tex
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

func loadProgram(vertPath, fragPath string, other ...string) (program gl.Program) {
	// load the shaders
	vertexShader := loadShader(gl.VERTEX_SHADER, vertPath)
	fragShader := loadShader(gl.FRAGMENT_SHADER, fragPath)

	// attach to the program
	program = gl.CreateProgram()
	program.AttachShader(vertexShader)
	program.AttachShader(fragShader)

	// load geom shader if present
	if len(other) > 0 {
		geomShader := loadShader(gl.GEOMETRY_SHADER, other[0])
		program.AttachShader(geomShader)
	}

	// bind out color, link, use
	program.BindFragDataLocation(0, "outColor")
	program.Link()
	program.Use()

	return
}

func bindArrays(program gl.Program, args ...interface{}) {
	var i, j, size, n uint
	size, n = 0, uint(len(args)/2)
	vars := make([]string, n)
	lens := make([]uint, n)
	for i, j = 0, 0; i < n; i, j = i+1, j+2 {
		vars[i] = args[j+0].(string)
		lens[i] = uint(args[j+1].(int))
		size += lens[i]
	}
	for i, j = 0, 0; i < n; i++ {
		pos := program.GetAttribLocation(vars[i])
		pos.AttribPointer(lens[i], gl.FLOAT, false, int(size*4), uintptr(j*4))
		pos.EnableArray()
		j += lens[i]
	}
}

func drawStencils(id int) {
	gl.Enable(gl.STENCIL_TEST)                 // enable the test
	gl.StencilFunc(gl.ALWAYS, id, 0xff)        // always write the id
	gl.StencilOp(gl.KEEP, gl.KEEP, gl.REPLACE) // replace mode for matches
	gl.StencilMask(0xff)                       // apply any id < 256
}

func useStencils(id int) {
	gl.StencilFunc(gl.EQUAL, id, 0xff) // now set to only draw on stencil
	gl.StencilMask(0xff)               // but don't affect the stencil itself
}

func noStencils() {
	gl.Clear(gl.STENCIL_BUFFER_BIT) // clear the stencil
	gl.Disable(gl.STENCIL_TEST)     // and then disable it
}

func transform(m *mathgl.Mat4f) {
	uniModel.UniformMatrix4f(false, (*[16]float32)(m))
}
