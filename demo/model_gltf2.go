package main

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/iamGreedy/essence/axis"
	"github.com/iamGreedy/essence/meter"
	"github.com/iamGreedy/essence/must"
	"github.com/iamGreedy/essence/prefix"
	"github.com/iamGreedy/gltf2"

	"os"
)

func main() {
	//f := must.MustGet(os.Open("./demo/models/AnimatedCube/glTF/AnimatedCube.gltf")).(*os.File)
	f := must.MustGet(os.Open("./demo/models/BouncingAnimationTest/Bouncing.gltf")).(*os.File)
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
			gltf2.Tasks.ModelScale(axis.X, meter.New(prefix.No, 30)),
			//gltf2.Tasks.ModelAlign(align.Center, align.Center, align.Center),
			gltf2.Tasks.ByeWorld,
		).
		Strictness(gltf2.LEVEL1).
		Parse()).(*gltf2.GLTF)
	//
	fmt.Println(md.Asset)
	fmt.Println(md.Animations[0])
	fmt.Println(md.Animations[0].Channels[0])
	fmt.Println(md.Animations[0].Channels[0].Target)
	fmt.Println(md.Animations[0].Samplers[0])
	fmt.Println(md.Animations[0].Samplers[0].Input)
	fmt.Println(md.Animations[0].Samplers[0].Input.MustSliceMapping(new([]float32),true, true))
	fmt.Println(md.Animations[0].Samplers[0].Output.MustSliceMapping(new([]mgl32.Vec3),true, true))

}
