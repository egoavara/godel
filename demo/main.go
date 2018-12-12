package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/iamGreedy/gltf2"
	"github.com/iamGreedy/godel"
	"github.com/iamGreedy/godel/shader"
	"os"
	"runtime"
)

func Must(err error) {
	if err != nil {
		panic(err)
	}
}
func MustGet(i interface{}, err error) interface{} {
	if err != nil {
		panic(err)
	}
	return i
}

func main() {
	var (
		width  = 640
		height = 480
	)
	runtime.LockOSThread()
	// GLFW, GL Init
	wnd := window(width, height)
	defer glfw.Terminate()
	version := gl.GoStr(gl.GetString(gl.VERSION))
	fmt.Println("OpenGL version : ", version)
	// GLTF
	//f := MustGet(os.Open("./demo/RubiksCube/RubiksCube_01.gltf")).(*os.File)
	//f := MustGet(os.Open("./demo/dice/dice.gltf")).(*os.File)
	f := MustGet(os.Open("./demo/damagedHelmet/damagedHelmet.gltf")).(*os.File)
	defer f.Close()
	md := MustGet(gltf2.Parser().
		Reader(f).
		Logger(os.Stdout).
		Plugin(
			gltf2.Tasks.HelloWorld,
			gltf2.Tasks.Caching,
			gltf2.Tasks.BufferCaching,
			gltf2.Tasks.ImageCaching,
			//gltf2.Tasks.AutoBufferTarget,
			//gltf2.Tasks.AccessorMinMax,
			//gltf2.Tasks.ModelScale(axis.Y, meter.New(prefix.No, 10)),
			//gltf2.Tasks.ModelAlign(align.Center, align.Center, align.Center),
			gltf2.Tasks.ByeWorld,
		).
		Strictness(gltf2.LEVEL1).
		Parse()).(*gltf2.GLTF)

	// godel Renderer
	app := godel.NewApplication(shader.Standard, shader.Flat, godel.NewCamera(godel.Perspective))
	app.Camera.LookFrom(mgl32.Vec3{-3, 0, 0})
	rd := app.MustRenderer(md)
	//rd.Scale(mgl32.Vec3{0.001, 0.001, 0.001})
	//
	// Configure global settings
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.CULL_FACE)
	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	//
	for !wnd.ShouldClose() {
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		curr := float32(glfw.GetTime())
		rd.Rotate(mgl32.QuatRotate(mgl32.DegToRad(curr*50), mgl32.Vec3{0, 1, 0}))
		rd.Render()
		//
		wnd.SwapBuffers()
		glfw.PollEvents()
	}

}
func window(w, h int) *glfw.Window {
	// glfw3
	Must(glfw.Init())
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)
	// glfw3.window
	window := MustGet(glfw.CreateWindow(w, h, "Testing", nil, nil)).(*glfw.Window)
	window.MakeContextCurrent()
	// OpenGL
	Must(gl.Init())
	return window
}
