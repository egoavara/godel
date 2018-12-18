package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/iamGreedy/essence/axis"
	"github.com/iamGreedy/essence/meter"
	"github.com/iamGreedy/essence/must"
	"github.com/iamGreedy/essence/prefix"
	"github.com/iamGreedy/gltf2"
	"github.com/iamGreedy/godel"
	"github.com/iamGreedy/godel/shader"
	"os"
	"runtime"
)


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
	//f := must.MustGet(os.Open("./demo/models/damagedHelmet/damagedHelmet.gltf")).(*os.File)
	//f := must.MustGet(os.Open("./demo/models/2CylinderEngine/glTF/2CylinderEngine.gltf")).(*os.File)
	f := must.MustGet(os.Open("./demo/models/AnimatedCube/glTF/AnimatedCube.gltf")).(*os.File)
	defer f.Close()
	md := must.MustGet(gltf2.Parser().
		Reader(f).
		Logger(os.Stdout).
		Plugin(
			gltf2.Tasks.HelloWorld,
			gltf2.Tasks.Caching,
			gltf2.Tasks.BufferCaching,
			gltf2.Tasks.ImageCaching,
			gltf2.Tasks.AutoBufferTarget,
			gltf2.Tasks.AccessorMinMax,
			gltf2.Tasks.ModelScale(axis.X, meter.New(prefix.No, 15)),
			//gltf2.Tasks.ModelAlign(align.Center, align.Center, align.Center),
			//gltf2.Tasks.ModelAlign(align.No, align.No, align.Center),
			gltf2.Tasks.ByeWorld,
		).
		Strictness(gltf2.LEVEL1).
		Parse()).(*gltf2.GLTF)

	// godel Renderer
	app := godel.NewApplication(shader.Standard, shader.PBR, godel.NewCamera(godel.Perspective), godel.NewLighting())
	app.Camera.LookFrom(mgl32.Vec3{0, 10, -50})
	app.Lighting.Global.Direction = mgl32.Vec3{
		0, -1, 1,
	}
	rd := must.MustGet(app.NewRenderer(md)).(*godel.Renderer)
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
		rd.Rotate(mgl32.QuatRotate(mgl32.DegToRad(curr * 50), mgl32.Vec3{0, 1, 0}))
		//rd.Rotate(mgl32.QuatRotate(45, mgl32.Vec3{0, 1, 0}))
		rd.Render()
		//
		wnd.SwapBuffers()
		glfw.PollEvents()
	}

}
func window(w, h int) *glfw.Window {
	// glfw3
	must.Must(glfw.Init())
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)
	// glfw3.window
	window := must.MustGet(glfw.CreateWindow(w, h, "Testing", nil, nil)).(*glfw.Window)
	window.MakeContextCurrent()
	// OpenGL
	must.Must(gl.Init())
	return window
}
