package shader

type DefineList []Define

func NewDefineList(defines ...Define) *DefineList {
	res := new(DefineList)
	res.Add(defines...)
	return res
}
func (s *DefineList) Add(defines ...Define) {
	for _, def := range defines {
		if !s.Contain(def) {
			*s = append(*s, def)
		}
	}
}
func (s *DefineList) Contain(defines ...Define) (ok bool) {
	for _, def0 := range defines {
		temp := false
		for _, def1 := range *s {
			temp = temp || def0 == def1
		}
		if !temp {
			return false
		}
	}
	return true
}
func (s *DefineList) Condition(cond *DefineList) (ok bool) {
	return s.Contain([]Define(*cond)...)
}
func (s *DefineList) Copy() *DefineList {
	return NewDefineList(*s...)
}

type Define string

// Standard.vs.glsl, PBR.fs.glsl, Phong.fs.glsl, Debug.fs.glsl
const (
	HAS_NORMAL  Define = "#define HAS_NORMAL"
	HAS_TANGENT Define = "#define HAS_TANGENT"
	HAS_COORD_0 Define = "#define HAS_COORD_0"
	HAS_COORD_1 Define = "#define HAS_COORD_1"
)

// PBR.fs.glsl
const (
	USE_IBL               Define = "#define USE_IBL"
	HAS_BASECOLORTEX      Define = "#define HAS_BASECOLORTEX"
	HAS_NORMALTEX         Define = "#define HAS_NORMALTEX"
	HAS_EMISSIVETEX       Define = "#define HAS_EMISSIVETEX"
	HAS_METALROUGHNESSTEX Define = "#define HAS_METALROUGHNESSTEX"
	HAS_OCCLUSIONTEX      Define = "#define HAS_OCCLUSIONTEX"
	//
	MANUAL_SRGB             Define = "#define MANUAL_SRGB"
	SRGB_FAST_APPROXIMATION Define = "#define SRGB_FAST_APPROXIMATION"
)
