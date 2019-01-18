package godel

import "github.com/go-gl/mathgl/mgl64"

type Camera struct {
	Projection Projection
	LookFrom mgl64.Vec3
	LookTo mgl64.Vec3
	Upside mgl64.Vec3
}

func (s *Camera) View() mgl64.Mat4 {
	if s.Upside.ApproxEqual(mgl64.Vec3{}){
		return mgl64.LookAtV(s.LookFrom, s.LookTo, mgl64.Vec3{0, 1, 0})
	}
	return mgl64.LookAtV(s.LookFrom, s.LookTo, s.Upside)
}

func (s *Camera) PV(width, height float64) mgl64.Mat4 {
	return s.Projection.Projection(width, height).Mul4(s.View())
}