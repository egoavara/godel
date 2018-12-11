package godel

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/iamGreedy/essence/version"
	"github.com/iamGreedy/gltf2"
	"github.com/iamGreedy/godel/shader"
	"image"
)

type program struct {
	glptr   uint32
	defines *shader.DefineList
}
type Application struct {
	glprograms []*program
	//
	vs     *shader.Shader
	fs     *shader.Shader
	screen mgl32.Vec2
	// public
	Camera *Camera
}

func NewApplication(vs *shader.Shader, fs *shader.Shader, camera *Camera) *Application {
	if vs.Type() != shader.Vertex {
		panic("vs must be Vertex Shader")
	}
	if fs.Type() != shader.Fragment {
		panic("fs must be Fragment Shader")
	}
	if camera == nil {
		panic("Camera not nillable")
	}
	//
	size := viewportSize().Size()

	return &Application{vs: vs, fs: fs, Camera: camera,
		screen: mgl32.Vec2{float32(size.X), float32(size.Y)},
	}
}
func (s *Application) requireProgram(defines *shader.DefineList) int {
	// find matching program index
	var findIdx = -1
	for i, prog := range s.glprograms {
		if prog.defines.Condition(defines) {
			findIdx = i
			break
		}
	}
	if findIdx == -1 {
		// Compile new program
		s.glprograms = append(s.glprograms, &program{
			glptr:   buildProgram(s.vs.Source(version.New(4, 1), *defines...), s.fs.Source(version.New(4, 1), *defines...)),
			defines: defines,
		})
		findIdx = len(s.glprograms) - 1
	}
	return findIdx
}
func (s *Application) getProgram(i int) uint32 {
	return s.glprograms[i].glptr
}

//
// Set openGL viewport size to `size`
// If `size` == image.ZP, it automatically set `size` by using gl context
func (s *Application) Screen(size image.Point) {
	if size == image.ZP {
		size = viewportSize().Size()
	}
	s.screen = mgl32.Vec2{float32(size.X), float32(size.Y)}
}

//
func (s *Application) NewRenderer(model *gltf2.GLTF) (*Renderer, error) {
	res := &Renderer{
		app:   s,
		model: model,
	}
	if err := res._Setup(); err != nil {
		return nil, err
	}
	return res, nil
}
func (s *Application) MustRenderer(model *gltf2.GLTF) *Renderer {
	res, err := s.NewRenderer(model)
	if err != nil {
		panic(err)
	}
	return res
}