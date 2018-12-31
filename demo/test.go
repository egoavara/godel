package main

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/iamGreedy/essence/align"
	"github.com/iamGreedy/essence/axis"
	"github.com/iamGreedy/essence/meter"
	"github.com/iamGreedy/essence/must"
	"github.com/iamGreedy/essence/prefix"
	"github.com/iamGreedy/essence/version"
	"github.com/iamGreedy/gltf2"
	"github.com/iamGreedy/godel/shader"
	"os"
	"runtime"
)

func main() {
	glv := version.New(4, 1)
	f := must.MustGet(os.Open("./demo/models/0gltfTutorial19/file.gltf")).(*os.File)
	defer f.Close()
	md := must.MustGet(gltf2.Parser().
		Reader(f).
		Logger(os.Stdout).
		Extensions(
			new(gltf2.KHRMaterialsPBRSpecularGlossiness),
			new(gltf2.KHRDracoMeshCompression),
		).
		Tasks(
			gltf2.Tasks.HelloWorld,
			gltf2.Tasks.Caching,
			gltf2.Tasks.BufferCaching,
			gltf2.Tasks.ImageCaching,
			gltf2.Tasks.AutoBufferTarget,
			gltf2.Tasks.AccessorMinMax,
			gltf2.Tasks.ModelScale(axis.X, meter.New(prefix.No, 30)),
			gltf2.Tasks.ModelAlign(align.Center, align.Center, align.Center),
			gltf2.Tasks.ByeWorld,
		).
		Strictness(gltf2.LEVEL1).
		Parse()).(*gltf2.GLTF)
	//
	defs := shader.NewDefineList(shader.HAS_JOINT_0, shader.HAS_WEIGHT_0)
	vs := shader.Standard.Source(glv, *defs...)
	fs := shader.Flat.Source(glv, *defs...)
	//
	var (
		width  = 640
		height = 480
	)
	runtime.LockOSThread()
	must.Must(glfw.Init())
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, 4)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)
	window := must.MustGet(glfw.CreateWindow(width, height, "Testing", nil, nil)).(*glfw.Window)
	window.MakeContextCurrent()
	must.Must(gl.Init())
	defer glfw.Terminate()
	//
	gl.Enable(gl.DEPTH_TEST)
	gl.DepthFunc(gl.LESS)
	gl.Enable(gl.CULL_FACE)
	gl.ClearColor(0.1, 0.1, 0.1, 1.0)
	//
	//
	for prev, curr := float32(0), float32(glfw.GetTime()); !window.ShouldClose(); prev, curr = curr, float32(glfw.GetTime()) {
		glfw.PollEvents()
		gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
		//
		window.SwapBuffers()
	}
}