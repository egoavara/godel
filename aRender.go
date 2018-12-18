package godel

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/iamGreedy/gltf2"
	"github.com/iamGreedy/godel/shader"
	"github.com/pkg/errors"
	"image"
	"unsafe"
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
	//vbos []uint32
	mVBO map[*gltf2.BufferView]uint32
	// texture
	//textures []uint32
	mTextures map[*gltf2.Texture]uint32
	//
	t mgl32.Vec3
	r mgl32.Quat
	s mgl32.Vec3
	m *mgl32.Mat4
}

func (s *Renderer) _Setup() error {
	s.mTextures = make(map[*gltf2.Texture]uint32)
	s.mVBO = make(map[*gltf2.BufferView]uint32)
	s.sceneidx = -1
	if err := s._Setup_buffers(); err != nil {
		return err
	}
	if err := s._Setup_textures(); err != nil {
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
			// vs defs
			if _, ok := prim.Attributes[gltf2.POSITION]; !ok {
				return errors.New("Must have POSITION")
			}
			if _, ok := prim.Attributes[gltf2.TEXCOORD_0]; ok {
				defs.Add(shader.HAS_COORD_0)
			}
			if _, ok := prim.Attributes[gltf2.TEXCOORD_1]; ok {
				defs.Add(shader.HAS_COORD_1)
			}
			if _, ok := prim.Attributes[gltf2.NORMAL]; ok {
				defs.Add(shader.HAS_NORMAL)
			}
			if _, ok := prim.Attributes[gltf2.TANGENT]; ok {
				defs.Add(shader.HAS_TANGENT)
			}
			// fs defs
			if prim.Material != nil {
				if prim.Material.PBRMetallicRoughness != nil {
					if prim.Material.PBRMetallicRoughness.BaseColorTexture != nil {
						defs.Add(shader.HAS_BASECOLORTEX)
					}
					if prim.Material.PBRMetallicRoughness.MetallicRoughnessTexture != nil {
						defs.Add(shader.HAS_METALROUGHNESSTEX)
					}
				}
				if prim.Material.NormalTexture != nil{
					defs.Add(shader.HAS_NORMALTEX)
				}
				if prim.Material.OcclusionTexture != nil{
					defs.Add(shader.HAS_OCCLUSIONTEX)
				}
				if prim.Material.EmissiveTexture != nil{
					defs.Add(shader.HAS_EMISSIVETEX)
				}
			}
			//
			s.meshPrimitive[i][j].programIndex = s.app.requireProgram(defs)
			// _Setup Vao
			gl.GenVertexArrays(1, &s.meshPrimitive[i][j].vao)
			gl.BindVertexArray(s.meshPrimitive[i][j].vao)
			// VBO POSITION
			pos := prim.Attributes[gltf2.POSITION]
			gl.BindBuffer(gl.ARRAY_BUFFER, s.mVBO[pos.BufferView])
			gl.EnableVertexAttribArray(0)
			gl.VertexAttribPointer(0, int32(pos.Type.Count()), uint32(pos.ComponentType), pos.Normalized, int32(pos.BufferView.ByteStride), gl.PtrOffset(pos.ByteOffset))
			// VBO TEXCOORD_0
			if coord0, ok := prim.Attributes[gltf2.TEXCOORD_0]; ok {
				gl.BindBuffer(gl.ARRAY_BUFFER, s.mVBO[coord0.BufferView])
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
			// VBO TEXCOORD_1
			if coord1, ok := prim.Attributes[gltf2.TEXCOORD_0]; ok {
				gl.BindBuffer(gl.ARRAY_BUFFER, s.mVBO[coord1.BufferView])
				gl.EnableVertexAttribArray(5)
				gl.VertexAttribPointer(
					5,
					int32(coord1.Type.Count()),
					uint32(coord1.ComponentType),
					coord1.Normalized,
					int32(coord1.BufferView.ByteStride),
					gl.PtrOffset(coord1.ByteOffset),
				)
			}
			// VBO NORMAL
			if norm, ok := prim.Attributes[gltf2.NORMAL]; ok {
				gl.BindBuffer(gl.ARRAY_BUFFER, s.mVBO[norm.BufferView])
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
			// VBO TANGENT
			if tangent, ok := prim.Attributes[gltf2.TANGENT]; ok {
				gl.BindBuffer(gl.ARRAY_BUFFER, s.mVBO[tangent.BufferView])
				gl.EnableVertexAttribArray(2)
				gl.VertexAttribPointer(
					2,
					int32(tangent.Type.Count()),
					uint32(tangent.ComponentType),
					tangent.Normalized,
					int32(tangent.BufferView.ByteStride),
					gl.PtrOffset(tangent.ByteOffset),
				)
			}
			// EBO
			if prim.Indices != nil {
				gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, s.mVBO[prim.Indices.BufferView])
			}
		}
	}
	//
	return nil
}
func (s *Renderer) _Setup_buffers() (err error) {
	if len(s.model.BufferViews) < 1{
		return nil
	}
	vbos := make([]uint32, len(s.model.BufferViews))
	gl.GenBuffers(int32(len(vbos)), &vbos[0])
	defer func() {
		if err != nil {
			gl.DeleteBuffers(int32(len(vbos)), &vbos[0])
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
			gl.BindBuffer(gl.ARRAY_BUFFER, vbos[i])
			gl.BufferData(gl.ARRAY_BUFFER, len(bts), gl.Ptr(&bts[0]), gl.STATIC_DRAW)
		case gltf2.ELEMENT_ARRAY_BUFFER:
			gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, vbos[i])
			gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(bts), gl.Ptr(bts), gl.STATIC_DRAW)
		}
		//
		s.mVBO[bv] = vbos[i]
	}
	return nil
}
func (s *Renderer) _Setup_textures() (err error) {
	if len(s.model.Textures) < 1{
		return nil
	}
	textures := make([]uint32, len(s.model.Textures))
	gl.GenTextures(int32(len(textures)), &textures[0])
	defer func() {
		if err != nil {
			gl.DeleteBuffers(int32(len(textures)), &textures[0])
		}
	}()
	for i, tex := range s.model.Textures {
		var img *image.RGBA
		img, err = tex.Source.Load(false)
		if err != nil {
			return err
		}

		//
		gl.ActiveTexture(gl.TEXTURE0)
		gl.BindTexture(gl.TEXTURE_2D, textures[i])
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, tex.Sampler.MagFilter.GL())
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, tex.Sampler.MinFilter.GL())
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, tex.Sampler.WrapS.GL())
		gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, tex.Sampler.WrapT.GL())
		gl.TexImage2D(gl.TEXTURE_2D,
			0,
			gl.RGBA,
			int32(img.Bounds().Dx()), int32(img.Bounds().Dy()),
			0,
			gl.RGBA,
			gl.UNSIGNED_BYTE,
			unsafe.Pointer(&img.Pix[0]),
		)
		if tex.Sampler.MinFilter.IsMipmap() {
			gl.GenerateMipmap(gl.TEXTURE_2D)
		}
		s.mTextures[tex] = textures[i]
	}

	return nil
}
func (s *Renderer) _Recur_node(node *gltf2.Node, cameraMatrix mgl32.Mat4, modelMatrix mgl32.Mat4, cameraPos mgl32.Vec3) {
	if node == nil {
		return
	}
	// mesh index
	var idx_mesh = -1
	for i, v := range s.model.Meshes {
		if v == node.Mesh {
			idx_mesh = i
			break
		}
	}
	// modelMatrix matrix setup
	modelMatrix = modelMatrix.Mul4(node.Transform())
	//modelMatrix = node.Transform().Mul4(modelMatrix)
	if node.Mesh != nil {
		// render mesh
		for idx_prim, prim := range node.Mesh.Primitives {
			prog := s.app.getProgram(s.meshPrimitive[idx_mesh][idx_prim].programIndex)
			prog.Use(func(p *ProgramContext) {
				// matrix
				p.Uniform("CameraMatrix", cameraMatrix)
				p.Uniform("ModelMatrix", modelMatrix)
				p.Uniform("NormalMatrix", modelMatrix.Inv().Transpose())
				p.Uniform("Camera", cameraPos)
				if s.app.Lighting != nil{
					if s.app.Lighting.Global != nil{
						p.Uniform("LightDir", s.app.Lighting.Global.Direction.Mul(-1))
						p.Uniform("LightColor", s.app.Lighting.Global.Color)
					}
				}
				// material
				if prim.Material != nil {
					if prim.Material.PBRMetallicRoughness != nil {
						p.Uniform("BaseColorFactor", prim.Material.PBRMetallicRoughness.BaseColorFactor)
						p.Uniform("MetalRoughnessFactor", mgl32.Vec2{
							prim.Material.PBRMetallicRoughness.MetallicFactor,
							prim.Material.PBRMetallicRoughness.RoughnessFactor,
						})
						p.Uniform("EmissiveFactor", prim.Material.EmissiveFactor)
						if prim.Material.PBRMetallicRoughness.BaseColorTexture != nil {
							p.Uniform("BaseColorTex", 0)
							p.Uniform("BaseColorTexCoord", int32(prim.Material.PBRMetallicRoughness.BaseColorTexture.TexCoord))
							gl.ActiveTexture(gl.TEXTURE0)
							gl.BindTexture(gl.TEXTURE_2D, s.mTextures[prim.Material.PBRMetallicRoughness.BaseColorTexture.Index])
						}
						if prim.Material.PBRMetallicRoughness.MetallicRoughnessTexture != nil {
							p.Uniform("MetalRoughnessTex", 1)
							p.Uniform("MetalRoughnessTexCoord", int32(prim.Material.PBRMetallicRoughness.MetallicRoughnessTexture.TexCoord))
							gl.ActiveTexture(gl.TEXTURE1)
							gl.BindTexture(gl.TEXTURE_2D, s.mTextures[prim.Material.PBRMetallicRoughness.MetallicRoughnessTexture.Index])
						}
					}
					if prim.Material.NormalTexture != nil {
						p.Uniform("NormalTex", 2)
						p.Uniform("NormalScale", prim.Material.NormalTexture.Scale)
						p.Uniform("NormalTexCoord", int32(prim.Material.NormalTexture.TexCoord))
						gl.ActiveTexture(gl.TEXTURE2)
						gl.BindTexture(gl.TEXTURE_2D, s.mTextures[prim.Material.NormalTexture.Index])
					}
					if prim.Material.OcclusionTexture != nil {
						p.Uniform("OcculusionTex", 3)
						p.Uniform("OcclusionStrength", prim.Material.OcclusionTexture.Strength)
						p.Uniform("OcculusionTexCoord", int32(prim.Material.OcclusionTexture.TexCoord))
						gl.ActiveTexture(gl.TEXTURE3)
						gl.BindTexture(gl.TEXTURE_2D, s.mTextures[prim.Material.OcclusionTexture.Index])
					}
					if prim.Material.EmissiveTexture != nil {
						p.Uniform("EmissiveTex", 4)
						p.Uniform("EmissiveTexCoord", int32(prim.Material.EmissiveTexture.TexCoord))
						gl.ActiveTexture(gl.TEXTURE4)
						gl.BindTexture(gl.TEXTURE_2D, s.mTextures[prim.Material.EmissiveTexture.Index])
					}
				}
				//
				if prim.Indices == nil {
					gl.BindVertexArray(s.meshPrimitive[idx_mesh][idx_prim].vao)
					gl.DrawArrays(uint32(prim.Mode.GL()), int32(prim.Indices.ByteOffset), int32(prim.Indices.Count))
				} else {
					gl.BindVertexArray(s.meshPrimitive[idx_mesh][idx_prim].vao)
					gl.DrawElements(uint32(prim.Mode.GL()), int32(prim.Indices.Count), uint32(prim.Indices.ComponentType.GL()), gl.PtrOffset(prim.Indices.ByteOffset))
				}
			})
		}
	}
	// render node child
	for _, child := range node.Children {
		s._Recur_node(child, cameraMatrix, modelMatrix, cameraPos)
	}

}
// Scene
func (s *Renderer) Scene(i int) {
	if i < 0 {
		s.sceneidx = -1
	}
	if i >= len(s.model.Scenes) {
		i = len(s.model.Scenes) - 1
	}
	s.sceneidx = i
}
func (s *Renderer) SceneCount() int {
	return len(s.model.Scenes)
}
func (s *Renderer) SceneName(id string) bool {
	if len(id) < 1 {
		return false
	}
	for i, scene := range s.model.Scenes {
		if scene.Name == id {
			s.Scene(i)
			return true
		}
	}
	return false
}
func (s *Renderer) SceneNames() []string {
	var res = make([]string, 0, len(s.model.Scenes))
	for _, v := range s.model.Scenes {
		res = append(res, v.Name)
	}
	return res
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
		s._Recur_node(n, s.app.Camera.Matrix(s.app.screen), s.Matrix(), s.app.Camera.GetLookFrom())
	}
}
