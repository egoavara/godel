package godel

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/iamGreedy/godel/shader"
	"image"
)

//type program struct {
//	glptr   uint32
//	defines *shader.DefineList
//}
type Application struct {
	glprograms []*Program
	//
	vs     *shader.Shader
	fs     *shader.Shader
	screen mgl32.Vec2
	// lighting
	//ibl []uint32
	// public
	Camera *Camera
	Lighting *Lighting

	//
	updaters []Updater
}
type Updater interface {
	dt(t float32)
}

func NewApplication(vs *shader.Shader, fs *shader.Shader, camera *Camera, lighting *Lighting) *Application {
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

	return &Application{
		vs: vs,
		fs: fs,
		screen: mgl32.Vec2{float32(size.X), float32(size.Y)},
		Camera: camera,
		Lighting: lighting,
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
		s.glprograms = append(s.glprograms, NewProgram(s.vs, s.fs, defines))
		findIdx = len(s.glprograms) - 1
	}
	return findIdx
}
func (s *Application) getProgram(i int) *Program {
	return s.glprograms[i]
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
func (s *Application) Update(dt float32) {
	for _, v := range s.updaters {
		v.dt(dt)
	}
}
func (s *Application) append(u Updater) {
	if u == nil{
		return
	}
	s.updaters = append(s.updaters, u)
}
func (s *Application) delete(u Updater) {
	if u == nil{
		return
	}
	for i, updater := range s.updaters {
		if updater == u{
			s.updaters = append(s.updaters[:i], s.updaters[i + 1:]...)
			return
		}
	}
}