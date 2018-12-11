package ogl

import "github.com/go-gl/gl/v4.1-core/gl"

var (
	_OGL_INIT     bool
	_OGL41_single *_OGL41 = nil
	//_OGL45_single *_OGL45= nil
)

const (
	OGL4 = OGL41

	OGL41 = 41
	OGL43 = 43
	OGL45 = 45
)

func Get(version int) OpenGL4 {

	switch version {
	case OGL41:
		if !_OGL_INIT {
			if err := gl.Init(); err != nil {
				panic(err)
			}
			_OGL41_single = new(_OGL41)
		}
		if _OGL41_single == nil {
			panic("OGL4 panic")
		}
		return _OGL41_single
	case OGL43:
		fallthrough
	case OGL45:
		fallthrough
	default:
		panic("Error OGL Version")
	}
}
