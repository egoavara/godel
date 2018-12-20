package godel

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/iamGreedy/gltf2"
)

type Object struct {
	model *Model
	tree  []*node
	anim  *Player
	//
	t mgl32.Vec3
	r mgl32.Quat
	s mgl32.Vec3
}

type node struct {
	children []*node
	//
	t  mgl32.Vec3
	r  mgl32.Quat
	s  mgl32.Vec3
	aT *mgl32.Vec3
	aR *mgl32.Quat
	aS *mgl32.Vec3
	m  *mgl32.Mat4

	//
	src  *gltf2.Node
	mesh *gltf2.Mesh
}

func (s *Model) NewObject(i int) (player *Object) {
	res := &Object{
		model: s,
		tree:  nil,

		t: mgl32.Vec3{0, 0, 0},
		r: mgl32.QuatIdent(),
		s: mgl32.Vec3{1, 1, 1},
	}
	res.Scene(i)
	return res
}
func (s *Object) Close() error {
	if s.anim != nil{
		s.anim.Close()
	}
	return nil
}
func (s *Object) Scene(i int) {
	if s.anim != nil{
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
	s.tree = make([]*node, len(scn.Nodes))
	for i, v := range scn.Nodes {
		s.tree[i] = new(node)
		s.recurSetupNode(s.tree[i], v)
	}
}
func (s *Object) recurSetupNode(dst *node, src *gltf2.Node) {
	dst.src = src
	dst.mesh = src.Mesh
	//
	dst.t = src.Translation
	dst.r = src.Rotation
	dst.s = src.Scale
	if src.Matrix != mgl32.Ident4() {
		dst.m = new(mgl32.Mat4)
		*dst.m = src.Matrix
	}
	//
	dst.children = make([]*node, len(src.Children))
	for i, v := range src.Children {
		dst.children[i] = new(node)
		s.recurSetupNode(dst.children[i], v)
	}
}
func (s *Object) find(src *gltf2.Node) *node {
	for _, tree := range s.tree {
		if n := s.recurFind(tree, src); n != nil {
			return n
		}
	}
	return nil
}
func (s *Object) recurFind(from *node, src *gltf2.Node) *node {
	if from == nil {
		return nil
	}
	if from.src == src {
		return from
	}
	for _, c := range from.children {
		if n := s.recurFind(c, src); n != nil {
			return n
		}
	}
	return nil
}

func (s *Object) Controller() *Player {
	return s.anim
}
func (s *Object) setupAnimation(dst *Player, src *gltf2.Animation) error {
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
			case gltf2.Weight:
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

func (s *Object) Transform() mgl32.Mat4 {
	return mgl32.Translate3D(s.t[0], s.t[1], s.t[2]).
		Mul4(s.r.Mat4()).
		Mul4(mgl32.Scale3D(s.s[0], s.s[1], s.s[2]))
}
func (s *Object) Translate(translate mgl32.Vec3) {
	s.t = translate
}
func (s *Object) Rotate(rotate mgl32.Quat) {
	s.r = rotate
}
func (s *Object) Scale(scale mgl32.Vec3) {
	s.s = scale
}
func (s *Object) GetTranslate() (translate mgl32.Vec3) {
	return s.t
}
func (s *Object) GetRotate() (rotate mgl32.Quat) {
	return s.r
}
func (s *Object) GetScale() (scale mgl32.Vec3) {
	return s.s
}

func (s *Object) Render() {
	for _, n := range s.tree {
		s.recurRender(n, s.model.app.Camera.Matrix(s.model.app.screen), s.Transform(), s.model.app.Camera.GetLookFrom())
	}
}
func (s *Object) recurRender(node *node, cameraMatrix mgl32.Mat4, modelMatrix mgl32.Mat4, cameraPos mgl32.Vec3) {
	if node == nil {
		return
	}
	// modelMatrix matrix setup
	animon := false
	transform := mgl32.Ident4()
	if !animon && node.m != nil {
		transform = *node.m
	}else {
		if node.aT != nil {
			transform = transform.Mul4(mgl32.Translate3D(node.aT[0], node.aT[1], node.aT[2]))
			animon = true
		} else {
			transform = transform.Mul4(mgl32.Translate3D(node.t[0], node.t[1], node.t[2]))
		}
		if node.aR != nil {
			transform = transform.Mul4(node.aR.Mat4())
			animon = true
		} else {
			transform = transform.Mul4(node.r.Mat4())
		}
		if node.aS != nil {
			transform = transform.Mul4(mgl32.Translate3D(node.aS[0], node.aS[1], node.aS[2]))
			animon = true
		} else {
			transform = transform.Mul4(mgl32.Scale3D(node.s[0], node.s[1], node.s[2]))
		}
	}
	modelMatrix = modelMatrix.Mul4(transform)
	//
	if node.mesh != nil {
		// render mesh
		for _, prim := range node.mesh.Primitives {
			primUser := prim.UserData.(*primitive)
			prog := s.model.app.getProgram(primUser.programIndex)
			prog.Use(func(p *ProgramContext) {
				// matrix
				p.Uniform("CameraMatrix", cameraMatrix)
				p.Uniform("ModelMatrix", modelMatrix)
				p.Uniform("NormalMatrix", modelMatrix.Inv().Transpose())
				p.Uniform("Camera", cameraPos)
				if s.model.app.Lighting != nil {
					if s.model.app.Lighting.Global != nil {
						p.Uniform("LightDir", s.model.app.Lighting.Global.Direction.Mul(-1))
						p.Uniform("LightColor", s.model.app.Lighting.Global.Color)
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
							gl.BindTexture(gl.TEXTURE_2D, prim.Material.PBRMetallicRoughness.BaseColorTexture.Index.UserData.(uint32))
						}
						if prim.Material.PBRMetallicRoughness.MetallicRoughnessTexture != nil {
							p.Uniform("MetalRoughnessTex", 1)
							p.Uniform("MetalRoughnessTexCoord", int32(prim.Material.PBRMetallicRoughness.MetallicRoughnessTexture.TexCoord))
							gl.ActiveTexture(gl.TEXTURE1)
							gl.BindTexture(gl.TEXTURE_2D, prim.Material.PBRMetallicRoughness.MetallicRoughnessTexture.Index.UserData.(uint32))
						}
					}
					if prim.Material.NormalTexture != nil {
						p.Uniform("NormalTex", 2)
						p.Uniform("NormalScale", prim.Material.NormalTexture.Scale)
						p.Uniform("NormalTexCoord", int32(prim.Material.NormalTexture.TexCoord))
						gl.ActiveTexture(gl.TEXTURE2)
						gl.BindTexture(gl.TEXTURE_2D, prim.Material.NormalTexture.Index.UserData.(uint32))
					}
					if prim.Material.OcclusionTexture != nil {
						p.Uniform("OcculusionTex", 3)
						p.Uniform("OcclusionStrength", prim.Material.OcclusionTexture.Strength)
						p.Uniform("OcculusionTexCoord", int32(prim.Material.OcclusionTexture.TexCoord))
						gl.ActiveTexture(gl.TEXTURE3)
						gl.BindTexture(gl.TEXTURE_2D, prim.Material.OcclusionTexture.Index.UserData.(uint32))
					}
					if prim.Material.EmissiveTexture != nil {
						p.Uniform("EmissiveTex", 4)
						p.Uniform("EmissiveTexCoord", int32(prim.Material.EmissiveTexture.TexCoord))
						gl.ActiveTexture(gl.TEXTURE4)
						gl.BindTexture(gl.TEXTURE_2D, prim.Material.EmissiveTexture.Index.UserData.(uint32))
					}
				}
				//
				gl.BindVertexArray(primUser.vao)
				if prim.Indices == nil {
					gl.DrawArrays(uint32(prim.Mode.GL()), int32(prim.Indices.ByteOffset), int32(prim.Indices.Count))
				} else {
					gl.DrawElements(uint32(prim.Mode.GL()), int32(prim.Indices.Count), uint32(prim.Indices.ComponentType.GL()), gl.PtrOffset(prim.Indices.ByteOffset))
				}
			})
		}
	}
	for _, child := range node.children {
		s.recurRender(child, cameraMatrix, modelMatrix, cameraPos)
	}
}
