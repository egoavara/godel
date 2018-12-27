package godel

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/iamGreedy/essence/version"
	"github.com/iamGreedy/godel/shader"
	"strings"
)

type Program struct {
	ptr          uint32
	defines      *shader.DefineList
	uniformIndex map[string]int32
	uniformData  map[int32]interface{}
	uboPointer   map[string]uint32
}

func NewProgram(vertex, frag *shader.Shader, defines *shader.DefineList) *Program {
	return &Program{
		ptr:          buildProgram(vertex.Source(version.New(4, 1), *defines...), frag.Source(version.New(4, 1), *defines...)),
		defines:      defines,
		uniformIndex: make(map[string]int32),
		uniformData:  make(map[int32]interface{}),
		uboPointer:   make(map[string]uint32),
	}
}
func (s *Program) GL() uint32 {
	return s.ptr
}
func (s *Program) Use(fn func(p *ProgramContext)) {
	gl.UseProgram(s.ptr)
	fn(&ProgramContext{
		ref: s,
	})
}

type ProgramContext struct {
	ref *Program
}

// Do not use just Array float, nor Slice float
func (s *ProgramContext) UniformIndex(key string) int32 {
	key = strings.Trim(key, "\x00")
	if v, ok := s.ref.uniformIndex[key]; ok {
		return v
	}
	lc := gl.GetUniformLocation(s.ref.ptr, gl.Str(key+"\x00"))
	if lc >= 0 {
		s.ref.uniformIndex[key] = lc
		return lc
	}
	return -1
}

// Do not use just Array float, nor Slice float
func (s *ProgramContext) Uniform(key string, data interface{}) bool {
	if idx := s.UniformIndex(key); idx >= 0 {
		if s.ref.uniformData[idx] == data {
			return true
		}
		s.ref.uniformData[idx] = data
		switch d := data.(type) {
		case int:
			gl.Uniform1i(idx, int32(d))
		case int32:
			gl.Uniform1i(idx, d)
		case uint:
			gl.Uniform1ui(idx, uint32(d))
		case uint32:
			gl.Uniform1ui(idx, d)
		case float32:
			gl.Uniform1f(idx, d)
		case mgl32.Vec2:
			gl.Uniform2f(idx, d[0], d[1])
		case mgl32.Vec3:
			gl.Uniform3f(idx, d[0], d[1], d[2])
		case mgl32.Vec4:
			gl.Uniform4f(idx, d[0], d[1], d[2], d[3])
		case mgl32.Mat2:
			gl.UniformMatrix2fv(idx, 1, false, &d[0])
		case mgl32.Mat3:
			gl.UniformMatrix3fv(idx, 1, false, &d[0])
		case mgl32.Mat4:
			gl.UniformMatrix4fv(idx, 1, false, &d[0])
		}
		return true
	}
	return false
}
