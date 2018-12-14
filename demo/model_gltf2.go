package main

import (
	"fmt"
	"github.com/iamGreedy/essence/align"
	"github.com/iamGreedy/essence/axis"
	"github.com/iamGreedy/essence/meter"
	"github.com/iamGreedy/essence/must"
	"github.com/iamGreedy/essence/prefix"
	"github.com/iamGreedy/gltf2"

	"os"
)

func main() {
	//f, err := os.Open("./demo/RubiksCube/RubiksCube_01.gltf")
	//f, err := os.Open("./demo/dice/dice.gltf")
	f := must.MustGet(os.Open("./demo/damagedHelmet/damagedHelmet.gltf")).(*os.File)
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
			gltf2.Tasks.ModelScale(axis.Y, meter.New(prefix.No, 13)),
			gltf2.Tasks.ModelAlign(align.Center, align.Center, align.Center),
			gltf2.Tasks.ByeWorld,
		).
		Strictness(gltf2.LEVEL1).
		Parse()).(*gltf2.GLTF)
	//
	fmt.Println(md.Asset)
	fmt.Println(md.Materials[0])
	fmt.Println(md.Materials[0].EmissiveTexture)
}
