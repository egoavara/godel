package main

import (
	"fmt"
	"github.com/iamGreedy/essence/version"
	"github.com/iamGreedy/gltf2"

	"os"
)

func main() {
	//f, err := os.Open("godel/demo/RubiksCube/RubiksCube_01.gltf")
	f, err := os.Open("godel/demo/dice/dice.gltf")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	md, err := gltf2.Parser().
		Reader(f).
		Logger(os.Stdout).
		Plugin(
			gltf2.Tasks.HelloWorld,
			gltf2.Tasks.Caching,
			gltf2.Tasks.BufferCaching,
			gltf2.Tasks.ImageCaching,
			gltf2.Tasks.AccessorMinMax,
			//gltf2.Tasks.ModelAlign(align.Center, align.Center, align.Center),
			gltf2.Tasks.AutoBufferTarget,
			gltf2.Tasks.ByeWorld,
		).
		Strictness(gltf2.LEVEL1).
		Parse()
	if err != nil {
		panic(err)
	}
	//
	fmt.Println(md.Asset)
	v := version.New(4, 1)
	fmt.Println(fmt.Sprintf("# version %d%d0", v.Major, v.Minor))
}
