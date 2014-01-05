package main

import (
	"fmt"
	"github.com/go-gl/gl"
	glfw "github.com/go-gl/glfw3"
)

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
