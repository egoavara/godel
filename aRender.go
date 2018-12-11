package godel

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/iamGreedy/gltf2"
	"github.com/iamGreedy/godel/shader"
	"github.com/pkg/errors"
)

type primitive struct {
	programIndex int
	vao          uint32
}
type Renderer struct {
	app *Application
	//
	model    *gltf2.GLTF
	sceneidx int
	// program, vao
	meshPrimitive [][]*primitive
	// vbo
	bufferViews []uint32
	//
	t mgl32.Vec3
	r mgl32.Quat
	s mgl32.Vec3
	m *mgl32.Mat4
}

func (s *Renderer) _Setup() error {
	s.sceneidx = -1
	if err := s._Setup_buffers(); err != nil {
		return err
	}
	if err := s._Setup_programs(); err != nil {
		return err
	}
	s.t = mgl32.Vec3{0, 0, 0}
	s.r = mgl32.QuatIdent()
	s.s = mgl32.Vec3{1, 1, 1}
	return nil
}

// privates
func (s *Renderer) _Setup_programs() (err error) {
	var base = shader.NewDefineList()
	// material defs

	//
	for i, mesh := range s.model.Meshes {
		s.meshPrimitive = append(s.meshPrimitive, make([]*primitive, len(mesh.Primitives)))
		for j, prim := range mesh.Primitives {
			s.meshPrimitive[i][j] = new(primitive)
			defs := base.Copy()
			if _, ok := prim.Attributes[gltf2.POSITION]; !ok {
				return errors.New("Must have POSITION")
			}
			if _, ok := prim.Attributes[gltf2.TEXCOORD_0]; ok {
				defs.Add(shader.HAS_NORMAL)
			}
			if _, ok := prim.Attributes[gltf2.NORMAL]; ok {
				defs.Add(shader.HAS_COORD_0)
			}
			//
			s.meshPrimitive[i][j].programIndex = s.app.requireProgram(defs)
			// _Setup Vao
			gl.GenVertexArrays(1, &s.meshPrimitive[i][j].vao)
			gl.BindVertexArray(s.meshPrimitive[i][j].vao)
			// VBO POSITION
			pos := prim.Attributes[gltf2.POSITION]
			gl.BindBuffer(gl.ARRAY_BUFFER, s.bufferViews[s.bufferViewIDX(pos.BufferView)])
			gl.EnableVertexAttribArray(0)
			gl.VertexAttribPointer(0, int32(pos.Type.Count()), uint32(pos.ComponentType), pos.Normalized, int32(pos.BufferView.ByteStride), gl.PtrOffset(pos.ByteOffset))
			// VBO TEXCOORD_0
			if coord0, ok := prim.Attributes[gltf2.TEXCOORD_0]; ok {
				gl.BindBuffer(gl.ARRAY_BUFFER, s.bufferViews[s.bufferViewIDX(coord0.BufferView)])
				gl.EnableVertexAttribArray(4)
				gl.VertexAttribPointer(
					4,
					int32(coord0.Type.Count()),
					uint32(coord0.ComponentType),
					coord0.Normalized,
					int32(coord0.BufferView.ByteStride),
					gl.PtrOffset(coord0.ByteOffset),
				)
			}
			// VBO NORMAL
			if norm, ok := prim.Attributes[gltf2.NORMAL]; ok {
				gl.BindBuffer(gl.ARRAY_BUFFER, s.bufferViews[s.bufferViewIDX(norm.BufferView)])
				gl.EnableVertexAttribArray(1)
				gl.VertexAttribPointer(
					1,
					int32(norm.Type.Count()),
					uint32(norm.ComponentType),
					norm.Normalized,
					int32(norm.BufferView.ByteStride),
					gl.PtrOffset(norm.ByteOffset),
				)
			}
			// EBO
			if prim.Indices != nil {
				gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, s.bufferViews[s.bufferViewIDX(prim.Indices.BufferView)])
			}
		}
	}
	//
	return nil
}
func (s *Renderer) _Setup_buffers() (err error) {
	s.bufferViews = make([]uint32, len(s.model.BufferViews))
	gl.GenBuffers(int32(len(s.bufferViews)), &s.bufferViews[0])
	defer func() {
		if err != nil {
			gl.DeleteBuffers(int32(len(s.bufferViews)), &s.bufferViews[0])
		}
	}()
	for i, bv := range s.model.BufferViews {
		var bts []byte
		bts, err = bv.Load()
		if err != nil {
			return err
		}
		if len(bts) <= 0 {
			continue
		}
		//
		switch bv.Target {
		case gltf2.NEED_TO_DEFINE_BUFFER:
			// TODO : logging
			fallthrough
		case gltf2.ARRAY_BUFFER:
			gl.BindBuffer(gl.ARRAY_BUFFER, s.bufferViews[i])
			gl.BufferData(gl.ARRAY_BUFFER, len(bts), gl.Ptr(&bts[0]), gl.STATIC_DRAW)
		case gltf2.ELEMENT_ARRAY_BUFFER:
			gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, s.bufferViews[i])
			gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(bts), gl.Ptr(bts), gl.STATIC_DRAW)
		}
	}
	return nil
}
func (s *Renderer) bufferViewIDX(bv *gltf2.BufferView) int {
	for i, v := range s.model.BufferViews {
		if v == bv {
			return i
		}
	}
	return -1
}
func (s *Renderer) recur_node(node *gltf2.Node, camera mgl32.Mat4, model mgl32.Mat4) {
	if node == nil {
		return
	}
	// mesh index
	var idx_mesh = -1
	for i, v := range s.model.Meshes {
		if v == node.Mesh {
			idx_mesh = i
		}
	}
	// model matrix setup
	model = model.Mul4(node.Transform())
	// render mesh
	for idx_prim, prim := range node.Mesh.Primitives {
		glProgram := s.app.getProgram(s.meshPrimitive[idx_mesh][idx_prim].programIndex)
		normal := model.Transpose()
		gl.UseProgram(glProgram)
		// matrix
		gl.UniformMatrix4fv(gl.GetUniformLocation(glProgram, gl.Str("CameraMatrix\x00")), 1, false, &camera[0])
		gl.UniformMatrix4fv(gl.GetUniformLocation(glProgram, gl.Str("ModelMatrix\x00")), 1, false, &model[0])
		gl.UniformMatrix4fv(gl.GetUniformLocation(glProgram, gl.Str("NormalMatrix\x00")), 1, false, &normal[0])
		// material
		gl.Uniform4fv(gl.GetUniformLocation(glProgram, gl.Str("FlatColor\x00")), 1, &prim.Material.PBRMetallicRoughness.BaseColorFactor[0])
		//
		if prim.Indices == nil {
			gl.BindVertexArray(s.meshPrimitive[idx_mesh][idx_prim].vao)
			gl.DrawArrays(uint32(prim.Mode.GL()), int32(prim.Indices.ByteOffset), int32(prim.Indices.Count))
		} else {
			gl.BindVertexArray(s.meshPrimitive[idx_mesh][idx_prim].vao)
			gl.DrawElements(uint32(prim.Mode.GL()), int32(prim.Indices.Count), uint32(prim.Indices.ComponentType.GL()), gl.PtrOffset(prim.Indices.ByteOffset))
		}
	}
	// render node child
	for _, child := range node.Childrun {
		s.recur_node(child, camera, model)
	}

}

// Scene
func (s *Renderer) Scene(i int) {
	if i < 0 {
		s.sceneidx = -1
	}
	s.sceneidx = i
}
func (s *Renderer) CountScene() int {
	return len(s.model.Scenes)
}

// TODO Animation

// Transform
func (s *Renderer) Matrix() mgl32.Mat4 {
	if s.m == nil {
		s.m = new(mgl32.Mat4)
		*s.m = mgl32.Translate3D(s.t[0], s.t[1], s.t[2]).
			Mul4(s.r.Mat4()).
			Mul4(mgl32.Scale3D(s.s[0], s.s[1], s.s[2]))
	}

	return *s.m
}
func (s *Renderer) Translate(translate mgl32.Vec3) {
	s.t = translate
	s.m = nil
}
func (s *Renderer) Rotate(rotate mgl32.Quat) {
	s.r = rotate
	s.m = nil
}
func (s *Renderer) Scale(scale mgl32.Vec3) {
	s.s = scale
	s.m = nil
}
func (s *Renderer) GetTranslate() (translate mgl32.Vec3) {
	return s.t
}
func (s *Renderer) GetRotate() (rotate mgl32.Quat) {
	return s.r
}
func (s *Renderer) GetScale() (scale mgl32.Vec3) {
	return s.s
}

// Rendering
func (s *Renderer) Render() {
	var scene *gltf2.Scene
	if s.sceneidx < 0 {
		scene = s.model.Scene
		if scene == nil {
			if len(s.model.Scenes) < 1 {
				return
			}
			scene = s.model.Scenes[0]
		}
	} else {
		if len(s.model.Scenes) < 1 {
			return
		}
		scene = s.model.Scenes[s.sceneidx]
	}

	for _, n := range scene.Nodes {
		s.recur_node(n, s.app.Camera.Matrix(s.app.screen), s.Matrix())
	}
}
