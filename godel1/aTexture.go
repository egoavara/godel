package godel1

type Texture struct {
	glptr uint32
}

func (s *Texture) GetTexture() uint32 {
	return s.glptr
}