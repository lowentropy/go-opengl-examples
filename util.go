package main

import (
	"fmt"
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
