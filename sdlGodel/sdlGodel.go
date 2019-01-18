package sdlGodel

import (
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/iamGreedy/essence/version"
	"image"
)

type Window struct {
	wnd *glfw.Window
}

func (s *Window) Viewport() image.Rectangle {
	w, h := s.wnd.GetSize()
	return image.Rect(0,0,w,h)
}

func (s *Window) UseContext() {
	s.wnd.MakeContextCurrent()
}

func NewWindow(w, h int, gl version.Version) (*Window, error) {

	// glfw3
	if err := glfw.Init(); err != nil {
		return nil, err
	}
	glfw.WindowHint(glfw.Resizable, glfw.False)
	glfw.WindowHint(glfw.ContextVersionMajor, int(gl.Major))
	glfw.WindowHint(glfw.ContextVersionMinor, int(gl.Minor))
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
	glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	glfw.WindowHint(glfw.DoubleBuffer, glfw.True)
	// glfw3.window
	if window, err := glfw.CreateWindow(w, h, "Testing", nil, nil); err != nil {
		return nil, err
	}else {
		return &Window{
			wnd:window,
		}, nil
	}
}
