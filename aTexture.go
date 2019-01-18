package godel

type Texture struct {
	glptr uint32
}

func (s *Texture) GetTexture() uint32 {
	return s.glptr
}