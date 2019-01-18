package godel

import (
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/iamGreedy/gltf2"
)

type Instance struct {
	model *Model
	tree  []*Node
	anim  *Player
	//
	t mgl32.Vec3
	r mgl32.Quat
	s mgl32.Vec3
}

func (s *Model) NewInstance(i int) *Instance {
	res := &Instance{
		model: s,
		tree:  nil,

		t: mgl32.Vec3{0, 0, 0},
		r: mgl32.QuatIdent(),
		s: mgl32.Vec3{1, 1, 1},
	}
	res.Scene(i)
	return res
}
func (s *Instance) Close() error {
	if s.anim != nil {
		s.anim.Close()
	}
	return nil
}
func (s *Instance) Scene(i int) {
	if s.anim != nil {
		s.anim.Close()
		s.anim = nil
	}
	//
	var scn *gltf2.Scene
	if i < 0 || i >= len(s.model.gltf.Scenes) {
		scn = s.model.gltf.Scene
		if scn == nil {
			scn = s.model.gltf.Scenes[0]
		}
	} else {
		scn = s.model.gltf.Scenes[i]
	}
	s.tree = MakeNodes(scn.Nodes...)
}
func (s *Instance) find(src *gltf2.Node) *Node {
	for _, node := range s.tree {
		if n := node.search(src); n != nil {
			return n
		}
	}
	return nil
}
func (s *Instance) Controller() *Player {
	return s.anim
}
func (s *Instance) setupAnimation(dst *Player, src *gltf2.Animation) error {
	dst.ref = s
	dst.curr = 0
	dst.fPlay = true
	dst.PlaySpeed = 1
	dst.Loop = false
	//
	temp := make(map[*gltf2.AnimationSampler]PlayerSampler)
	dst.targets = make([]*playerTarget, len(src.Channels))
	for i, c := range src.Channels {
		if _, ok := temp[c.Sampler]; !ok {
			var err error
			switch c.Target.Path {
			case gltf2.Rotation:
				temp[c.Sampler], err = MakeSampler(c.Sampler, true, true)
			case gltf2.Weights:
				temp[c.Sampler], err = MakeSampler(c.Sampler, true, false)
			default:
				temp[c.Sampler], err = MakeSampler(c.Sampler, false, false)
			}
			if err != nil {
				return err
			}
		}
		dst.targets[i] = &playerTarget{
			sample: temp[c.Sampler],
			path:   c.Target.Path,
			target: s.find(c.Target.Node),
		}
		_, end := dst.targets[i].sample.Range()
		if dst.playtime < end {
			dst.playtime = end
		}
	}

	return nil
}

func (s *Instance) Matrix() mgl32.Mat4 {
	return mgl32.Translate3D(s.t[0], s.t[1], s.t[2]).Mul4(s.r.Mat4()).Mul4(mgl32.Scale3D(s.s[0], s.s[1], s.s[2]))
}
func (s *Instance) Translate(translate mgl32.Vec3) {
	s.t = translate
}
func (s *Instance) Rotate(rotate mgl32.Quat) {
	s.r = rotate
}
func (s *Instance) Scale(scale mgl32.Vec3) {
	s.s = scale
}
func (s *Instance) GetTranslate() (translate mgl32.Vec3) {
	return s.t
}
func (s *Instance) GetRotate() (rotate mgl32.Quat) {
	return s.r
}
func (s *Instance) GetScale() (scale mgl32.Vec3) {
	return s.s
}

func (s *Instance) Render() {
	for _, n := range s.tree {
		s.recurRender(n, n, s.model.app.Camera.Matrix(s.model.app.screen), s.Matrix(), s.model.app.Camera.GetLookFrom())
	}
}

func (s *Instance) recurRender(root *Node, node *Node, cameraMatrix mgl32.Mat4, modelMatrix mgl32.Mat4, cameraPos mgl32.Vec3) {
	if node == nil {
		return
	}
	// modelMatrix matrix setup
	modelMatrix = modelMatrix.Mul4(node.Matrix(false))
	//
	if node.src.Mesh != nil {
		// render mesh
		for _, prim := range node.src.Mesh.Primitives {
			primUser := prim.UserData.(*primitive)
			primUser.prog.Use(func(p *ProgramContext) {
				// matrix
				p.Uniform("CameraMatrix", cameraMatrix)
				p.Uniform("ModelMatrix", modelMatrix)
				p.Uniform("NormalMatrix", modelMatrix.Inv().Transpose())
				p.Uniform("Camera", cameraPos)
				// lighting
				if s.model.app.Lighting != nil {
					if s.model.app.Lighting.Global != nil {
						p.Uniform("LightDir", s.model.app.Lighting.Global.Direction.Mul(-1))
						p.Uniform("LightColor", s.model.app.Lighting.Global.Color)
					}
				}
				// morphing
				if weight := node.Weight(); weight != nil{
					for i, v := range weight {
						p.Uniform(fmt.Sprintf("MorphWeight[%d]", i), v)
					}
				}
				// skinning
				if node.skin != nil {
					node.Skin().update(p, modelMatrix)
				}
				// material
				var mt *gltf2.Material
				if prim.Material != nil {
					mt = prim.Material
				} else {
					mt = gltf2.DefaultMaterial()
				}
				if mt.PBRMetallicRoughness != nil {
					p.Uniform("BaseColorFactor", mt.PBRMetallicRoughness.BaseColorFactor)
					p.Uniform("MetalRoughnessFactor", mgl32.Vec2{
						mt.PBRMetallicRoughness.MetallicFactor,
						mt.PBRMetallicRoughness.RoughnessFactor,
					})
					p.Uniform("EmissiveFactor", mt.EmissiveFactor)
					if mt.PBRMetallicRoughness.BaseColorTexture != nil {
						p.Uniform("BaseColorTex", 0)
						p.Uniform("BaseColorTexCoord", int32(mt.PBRMetallicRoughness.BaseColorTexture.TexCoord))
						gl.ActiveTexture(gl.TEXTURE0)
						gl.BindTexture(gl.TEXTURE_2D, mt.PBRMetallicRoughness.BaseColorTexture.Index.UserData.(uint32))
					}
					if mt.PBRMetallicRoughness.MetallicRoughnessTexture != nil {
						p.Uniform("MetalRoughnessTex", 1)
						p.Uniform("MetalRoughnessTexCoord", int32(mt.PBRMetallicRoughness.MetallicRoughnessTexture.TexCoord))
						gl.ActiveTexture(gl.TEXTURE1)
						gl.BindTexture(gl.TEXTURE_2D, mt.PBRMetallicRoughness.MetallicRoughnessTexture.Index.UserData.(uint32))
					}
				}
				if mt.NormalTexture != nil {
					p.Uniform("NormalTex", 2)
					p.Uniform("NormalScale", mt.NormalTexture.Scale)
					p.Uniform("NormalTexCoord", int32(mt.NormalTexture.TexCoord))
					gl.ActiveTexture(gl.TEXTURE2)
					gl.BindTexture(gl.TEXTURE_2D, mt.NormalTexture.Index.UserData.(uint32))
				}
				if mt.OcclusionTexture != nil {
					p.Uniform("OcculusionTex", 3)
					p.Uniform("OcclusionStrength", mt.OcclusionTexture.Strength)
					p.Uniform("OcculusionTexCoord", int32(mt.OcclusionTexture.TexCoord))
					gl.ActiveTexture(gl.TEXTURE3)
					gl.BindTexture(gl.TEXTURE_2D, mt.OcclusionTexture.Index.UserData.(uint32))
				}
				if mt.EmissiveTexture != nil {
					p.Uniform("EmissiveTex", 4)
					p.Uniform("EmissiveTexCoord", int32(mt.EmissiveTexture.TexCoord))
					gl.ActiveTexture(gl.TEXTURE4)
					gl.BindTexture(gl.TEXTURE_2D, mt.EmissiveTexture.Index.UserData.(uint32))
				}
				// rendering
				gl.BindVertexArray(primUser.vao)
				if prim.Indices == nil {
					gl.DrawArrays(uint32(prim.Mode.GL()), int32(0), int32(prim.Indices.Count))
				} else {
					gl.DrawElements(uint32(prim.Mode.GL()), int32(prim.Indices.Count), uint32(prim.Indices.ComponentType.GL()), gl.PtrOffset(0))
				}
			})
		}
	}
	for _, child := range node.children {
		s.recurRender(root, child, cameraMatrix, modelMatrix, cameraPos)
	}
}
