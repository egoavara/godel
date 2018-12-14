package godel

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
)

type Program struct {
	ptr uint32
	uniformIdx map[string]int32
	uniform map[string]interface{}
}

func (s *Program)GL() uint32 {
	return s.ptr
}
// Do not use just Array float, nor Slice float
func (s *Program)UniformIndex(key string) int32 {
	if v, ok := s.uniformIdx[key];ok{
		return v
	}
	gl.GetUniformLocation(glProgram, gl.Str(key + "\x00"))
}


// Do not use just Array float, nor Slice float
func (s *Program)Uniform(key string, data interface{}) {
	switch data.(type) {
	case int:

	case int32:
	case float32:
	case mgl32.Vec2:
	case mgl32.Vec3:
	case mgl32.Vec4:
	case mgl32.Mat2:
	case mgl32.Mat3:
	case mgl32.Mat4:

	}
}
