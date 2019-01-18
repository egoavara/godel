package godel

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/iamGreedy/godel/shader"
	"image"
)

type Application struct {
	glprograms []*Program
	//
	vs     *shader.Shader
	fs     *shader.Shader
	screen mgl32.Vec2
	// lighting
	//ibl []uint32
	// public
	Camera   *Camera
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
		camera = NewCamera(Perspective)
	}

	if lighting == nil {
		lighting = NewLighting()
	}
	//
	size := viewportSize().Size()

	return &Application{
		vs:       vs,
		fs:       fs,
		screen:   mgl32.Vec2{float32(size.X), float32(size.Y)},
		Camera:   camera,
		Lighting: lighting,
	}
}
func (s *Application) BuildProgram(vs, fs *shader.Shader, defines *shader.DefineList) *Program {
	if temp := s.FindProgram(vs, fs, defines); temp != nil{
		return temp
	}
	temp := NewProgram(vs, fs, defines)
	s.glprograms = append(s.glprograms, temp)
	return temp
}
func (s *Application) FindProgram(vs, fs *shader.Shader, defines *shader.DefineList) *Program {
	for _, v := range s.glprograms {
		if v.vsid == vs.ID() && v.fsid == fs.ID() && v.defines.Condition(defines){
			return v
		}
	}
	return nil
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
	if u == nil {
		return
	}
	s.updaters = append(s.updaters, u)
}
func (s *Application) delete(u Updater) {
	if u == nil {
		return
	}
	for i, updater := range s.updaters {
		if updater == u {
			s.updaters = append(s.updaters[:i], s.updaters[i+1:]...)
			return
		}
	}
}
