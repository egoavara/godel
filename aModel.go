package godel

import (
	"github.com/go-gl/gl/v4.1-core/gl"
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
type Model struct {
	app *Application
	gltf     *gltf2.GLTF
}

func (s *Application) BuildModel(model *gltf2.GLTF, clearCache bool) (*Model, error) {
	res := &Model{
		app:  s,
		gltf: model,
	}
	if err := res._Setup(); err != nil {
		return nil, err
	}
	if clearCache{
		gltf2.ThrowAllCache(model)
	}
	return res, nil
}
// privates
func (s *Model) _Setup() error {
	if err := s._Setup_buffers(); err != nil {
		return err
	}
	if err := s._Setup_textures(); err != nil {
		return err
	}
	if err := s._Setup_programs(); err != nil {
		return err
	}
	return nil
}
func (s *Model) _Setup_programs() (err error) {
	var base = shader.NewDefineList()
	// material defs

	//
	for _, mesh := range s.gltf.Meshes {
		for _, prim := range mesh.Primitives {
			temp := new(primitive)
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
				if prim.Material.NormalTexture != nil {
					defs.Add(shader.HAS_NORMALTEX)
				}
				if prim.Material.OcclusionTexture != nil {
					defs.Add(shader.HAS_OCCLUSIONTEX)
				}
				if prim.Material.EmissiveTexture != nil {
					defs.Add(shader.HAS_EMISSIVETEX)
				}
			}
			//
			temp.programIndex = s.app.requireProgram(defs)
			// _Setup Vao
			gl.GenVertexArrays(1, &temp.vao)
			gl.BindVertexArray(temp.vao)
			// VBO POSITION
			pos := prim.Attributes[gltf2.POSITION]

			gl.BindBuffer(gl.ARRAY_BUFFER, pos.BufferView.UserData.(uint32))
			gl.EnableVertexAttribArray(0)
			gl.VertexAttribPointer(0, int32(pos.Type.Count()), uint32(pos.ComponentType), pos.Normalized, int32(pos.BufferView.ByteStride), gl.PtrOffset(pos.ByteOffset))
			// VBO TEXCOORD_0
			if coord0, ok := prim.Attributes[gltf2.TEXCOORD_0]; ok {
				gl.BindBuffer(gl.ARRAY_BUFFER, coord0.BufferView.UserData.(uint32))
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
			if coord1, ok := prim.Attributes[gltf2.TEXCOORD_1]; ok {
				gl.BindBuffer(gl.ARRAY_BUFFER, coord1.BufferView.UserData.(uint32))
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
				gl.BindBuffer(gl.ARRAY_BUFFER, norm.BufferView.UserData.(uint32))
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
				gl.BindBuffer(gl.ARRAY_BUFFER, tangent.BufferView.UserData.(uint32))
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
				gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, prim.Indices.BufferView.UserData.(uint32))
			}
			//
			prim.UserData = temp
		}
	}
	//
	return nil
}
func (s *Model) _Setup_buffers() (err error) {
	if len(s.gltf.BufferViews) < 1 {
		return nil
	}
	vbos := make([]uint32, len(s.gltf.BufferViews))
	gl.GenBuffers(int32(len(vbos)), &vbos[0])
	defer func() {
		if err != nil {
			gl.DeleteBuffers(int32(len(vbos)), &vbos[0])
		}
	}()
	for i, bv := range s.gltf.BufferViews {
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
		bv.UserData = vbos[i]
	}
	return nil
}
func (s *Model) _Setup_textures() (err error) {
	if len(s.gltf.Textures) < 1 {
		return nil
	}
	textures := make([]uint32, len(s.gltf.Textures))
	gl.GenTextures(int32(len(textures)), &textures[0])
	defer func() {
		if err != nil {
			gl.DeleteBuffers(int32(len(textures)), &textures[0])
		}
	}()
	for i, tex := range s.gltf.Textures {
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
		tex.UserData = textures[i]
	}

	return nil
}

