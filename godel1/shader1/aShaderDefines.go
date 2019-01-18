package shader1

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