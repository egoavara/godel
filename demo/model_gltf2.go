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
	//f := must.MustGet(os.Open("./demo/models/2CylinderEngine/glTF-Draco/2CylinderEngine.gltf")).(*os.File)
	//f := must.MustGet(os.Open("./demo/models/Avocado/glTF-pbrSpecularGlossiness/Avocado.gltf")).(*os.File)
	//f := must.MustGet(os.Open("./demo/models/Avocado/glTF-Draco/Avocado.gltf")).(*os.File)
	//f := must.MustGet(os.Open("./demo/models/RiggedFigure/glTF/RiggedFigure.gltf")).(*os.File)
	f := must.MustGet(os.Open("./demo/models/RiggedSimple/glTF/RiggedSimple.gltf")).(*os.File)
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
	fmt.Println(md.Asset)
	fmt.Println(md.Skins[0])
	fmt.Println(md.Skins[0].Joints)
	fmt.Println(md.Skins[0].Skeleton)
//
}
