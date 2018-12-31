package godel

import (
	"fmt"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/iamGreedy/gltf2"
)

type (
	Node struct {
		src      *gltf2.Node
		parent   *Node
		children []*Node
		//
		skin *Skin
		//
		aT *mgl32.Vec3
		aR *mgl32.Quat
		aS *mgl32.Vec3
	}
	Skin struct {
		src   *gltf2.Skin
		skt   *Node
		joint []*Node
		ibm   []mgl32.Mat4
	}
)

func NewNode(src *gltf2.Node) *Node {
	res := new(Node)
	recurSetupNode(res, src)
	recurSetupSkin(res, res)
	return res
}
func MakeNodes(src ... *gltf2.Node) []*Node {
	res := make([]*Node, len(src))
	for i, v := range src {
		res[i] = NewNode(v)
	}
	return res
}
func recurSetupNode(dst *Node, src *gltf2.Node) {
	dst.src = src
	//
	dst.children = make([]*Node, len(src.Children))
	for i, v := range src.Children {
		dst.children[i] = new(Node)
		dst.children[i].parent = dst
		recurSetupNode(dst.children[i], v)
	}
}
func recurSetupSkin(root *Node, trg *Node) {
	if trg.src.Skin != nil {
		trg.skin = &Skin{
			src:   trg.src.Skin,
			skt:   root.search(trg.src.Skin.Skeleton),
			joint: make([]*Node, len(trg.src.Skin.Joints)),
		}
		for i, v := range trg.src.Skin.Joints {
			trg.skin.joint[i] = root.search(v)
		}
		//
		trg.skin.src.InverseBindMatrices.MustSliceMapping(&trg.skin.ibm, true, true)
	}
	//
	for _, v := range trg.children {
		recurSetupSkin(root, v)
	}
}

func (s *Node) Root() *Node {
	if s.parent == nil {
		return s
	}
	return s.parent.Root()
}
func (s *Node) Parents() *Node {
	return s.parent
}
func (s *Node) Find(id ...Identifier) *Node {
	if len(id) <= 0 {
		return nil
	}
	var found *Node = nil
	switch vid := validateIdentifier(id[0]).(type) {
	case string:
		if len(vid) > 0 {
			for _, v := range s.children {
				if vid == v.Name() {
					found = v
				}
			}
		}
	case int:
		if vid < len(s.children) {
			found = s.children[vid]
		}
	}
	if found == nil {
		return nil
	}
	if len(id) > 1 {
		return found.Find(id[1:]...)
	}
	return found
}
func (s *Node) search(n *gltf2.Node) *Node {
	if s.src == n {
		return s
	}
	for _, v := range s.children {
		if res := v.search(n); res != nil {
			return res
		}
	}
	return nil
}
func (s *Node) Children() []*Node {
	return s.children
}

func (s *Node) Name() string {
	return s.src.Name
}
func (s *Node) Index() int {
	if s.parent == nil {
		return -1
	}
	for i, v := range s.parent.children {
		if v == s {
			return i
		}
	}
	panic("Unreachable")
}
func (s *Node) Length() int {
	return len(s.children)
}
func (s *Node) Source() *gltf2.Node {
	return s.src
}

func (s *Node) setT(vec3 mgl32.Vec3) {
	s.aT = &vec3
}
func (s *Node) setR(quat mgl32.Quat) {
	s.aR = &quat
}
func (s *Node) setS(vec3 mgl32.Vec3) {
	s.aS = &vec3
}
func (s *Node) clearAnim() {
	s.aT = nil
	s.aR = nil
	s.aS = nil
}

func (s *Node) Skin() *Skin {
	return s.skin
}
func (s *Node) Matrix(globalMode bool) mgl32.Mat4 {
	var a = false
	mt := s.src.Translation
	mr := s.src.Rotation
	ms := s.src.Scale
	if s.aT != nil {
		mt = *s.aT
		a = true
	}
	if s.aR != nil {
		mr = *s.aR
		a = true
	}
	if s.aS != nil {
		ms = *s.aS
		a = true
	}
	if !a && s.src.Matrix != mgl32.Ident4() {
		return s.src.Matrix
	}
	model := mgl32.Translate3D(mt[0], mt[1], mt[2]).Mul4(mr.Mat4()).Mul4(mgl32.Scale3D(ms[0], ms[1], ms[2]))
	//
	if globalMode {
		if s.parent == nil {
			return model
		}
		return s.parent.Matrix(globalMode).Mul4(model)
	}
	return model
}


func (s *Skin) update(ctx *ProgramContext, globalTransformOfNodeThatTheMeshIsAttachedTo mgl32.Mat4) {
	for i, j := range s.joint {
		globalTransformOfJointNode := j.Matrix(true)
		inverseBindMatrixForJoint := s.ibm[i]
		// jointMatrix(j) =
		//     globalTransformOfNodeThatTheMeshIsAttachedTo^-1 *
		//     globalTransformOfJointNode(j) *
		//     inverseBindMatrixForJoint(j);
		jointMatrix := globalTransformOfJointNode.
			Mul4(inverseBindMatrixForJoint)
		ctx.Uniform(fmt.Sprintf("JointMatrix[%d]", i), jointMatrix)
	}
}