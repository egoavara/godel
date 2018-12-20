package godel

import (
	"github.com/go-gl/mathgl/mgl32"
	"math"
)

type Camera struct {
	mode CameraType
	//
	pargs [4]float32
	vargs [3]mgl32.Vec3
}

func NewCamera(cameraType CameraType, args ...float32) *Camera {
	cam := new(Camera)
	switch cameraType {
	case Orthographic:
		if len(args) == 4 {
			cam.Orthographic(args[0], args[1], args[2], args[3])
		} else {
			cam.Orthographic(AUTO, AUTO, NEAR, FAR)
		}
	case Perspective:
		if len(args) == 4 {
			cam.Perspective(args[0], args[1], args[2], args[3])
		} else {
			cam.Perspective(AUTO, mgl32.DegToRad(60), NEAR, INF)
		}
	}
	cam.vargs[0] = mgl32.Vec3{1, 1, 1}
	cam.vargs[1] = mgl32.Vec3{0, 0, 0}
	cam.vargs[2] = mgl32.Vec3{0, 1, 0}

	return cam
}

// Projection Matrix
func (s *Camera) Projection() (camType CameraType, arg0, arg1, arg2, arg3 float32) {
	return s.mode, s.pargs[0], s.pargs[1], s.pargs[2], s.pargs[3]
}
func (s *Camera) Orthographic(xmag, ymag, znear, zfar float32) {
	s.mode = Orthographic
	s.pargs[0] = xmag
	s.pargs[1] = ymag
	s.pargs[2] = znear
	s.pargs[3] = zfar
}
func (s *Camera) Perspective(aspectRatio, yfov, znear, zfar float32) {
	s.mode = Perspective
	s.pargs[0] = aspectRatio
	s.pargs[1] = yfov
	s.pargs[2] = znear
	s.pargs[3] = zfar
}

// View Matrix
func (s *Camera) LookFrom(from mgl32.Vec3) {
	s.vargs[0] = from
}
func (s *Camera) LookTo(to mgl32.Vec3) {
	s.vargs[1] = to
}
func (s *Camera) Upside(direction mgl32.Vec3) {
	s.vargs[2] = direction
}
func (s *Camera) GetLookFrom() (from mgl32.Vec3) {
	return s.vargs[0]

}
func (s *Camera) GetLookTo() (to mgl32.Vec3) {
	return s.vargs[1]
}
func (s *Camera) GetUpside() (direction mgl32.Vec3) {
	return s.vargs[2]
}

//
func (s *Camera) Matrix(screen mgl32.Vec2) mgl32.Mat4 {
	var p mgl32.Mat4
	var v mgl32.Mat4
	//
	switch s.mode {
	case Orthographic:
		var (
			xmag  = s.pargs[0]
			ymag  = s.pargs[1]
			znear = s.pargs[2]
			zfar  = s.pargs[3]
		)
		if xmag == AUTO {
			xmag = screen.X() / 2
		}
		if ymag == AUTO {
			ymag = screen.Y() / 2
		}
		p = mgl32.Ortho(-xmag, xmag, -ymag, ymag, znear, zfar)
	case Perspective:
		var (
			aspect = s.pargs[0]
			yfov   = s.pargs[1]
			znear  = s.pargs[2]
			zfar   = s.pargs[3]
		)
		if aspect == AUTO {
			aspect = screen.X() / screen.Y()
		}
		if znear < NEAR {
			znear = NEAR
		}
		if zfar == INF {
			const e = 0.000001
			f := float32(1. / math.Tan(float64(yfov)/2.0))
			p = mgl32.Mat4{f / aspect, 0, 0, 0, 0, f, 0, 0, 0, 0, -1 + e, -1, 0, 0, -znear, 0}
		} else {
			p = mgl32.Perspective(yfov, aspect, znear, zfar)
		}
	}
	//
	v = mgl32.LookAtV(s.vargs[0], s.vargs[1], s.vargs[2])
	return p.Mul4(v)
}

type CameraType uint32

const (
	Orthographic CameraType = iota
	Perspective  CameraType = iota
)
const (
	AUTO float32 = 0
	INF  float32 = -1
	NEAR float32 = 0.05
	FAR  float32 = 1000
)
