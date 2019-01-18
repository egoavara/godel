package godel

import (
	"github.com/go-gl/mathgl/mgl64"
	"math"
)

type Projection interface {
	Projection(wicth, height float64) mgl64.Mat4
}

type AutoPerspective struct {
	Yfov float64
	ZNear float64
	ZFar float64
}
func NewAutoPerspective(yfov float64, ZNear float64, ZFar float64) *AutoPerspective {
	return &AutoPerspective{Yfov: yfov, ZNear: ZNear, ZFar: ZFar}
}
func (s AutoPerspective) Projection(wicth, height float64) mgl64.Mat4 {
	return mgl64.Perspective(s.Yfov, wicth / height, s.ZNear, s.ZFar)
}

type Perspective struct {
	AspectRatio float64
	Yfov float64
	ZNear float64
	ZFar float64
}
func NewPerspective(aspectRatio float64, yfov float64, ZNear float64, ZFar float64) *Perspective {
	return &Perspective{AspectRatio: aspectRatio, Yfov: yfov, ZNear: ZNear, ZFar: ZFar}
}
func (s Perspective) Projection(wicth, height float64) mgl64.Mat4 {
	return mgl64.Perspective(s.Yfov, s.AspectRatio, s.ZNear, s.ZFar)
}

type AutoInfPerspective struct {
	Yfov float64
	ZNear float64
}
func NewAutoInfPerspective(yfov float64, ZNear float64) *AutoInfPerspective {
	return &AutoInfPerspective{Yfov: yfov, ZNear: ZNear}
}
func (s AutoInfPerspective) Projection(wicth, height float64) mgl64.Mat4 {
	const e = 0.000001
	f := 1. / math.Tan(float64(s.Yfov)/2.0)
	return mgl64.Mat4{f / (wicth / height), 0, 0, 0, 0, f, 0, 0, 0, 0, -1 + e, -1, 0, 0, -s.ZNear, 0}
}

type InfPerspective struct {
	AspectRatio float64
	Yfov float64
	ZNear float64
}
func NewInfPerspective(aspectRatio float64, yfov float64, ZNear float64) *InfPerspective {
	return &InfPerspective{AspectRatio: aspectRatio, Yfov: yfov, ZNear: ZNear}
}
func (s InfPerspective) Projection(wicth, height float64) mgl64.Mat4 {
	const e = 0.000001
	f := 1. / math.Tan(float64(s.Yfov)/2.0)
	return mgl64.Mat4{f / s.AspectRatio, 0, 0, 0, 0, f, 0, 0, 0, 0, -1 + e, -1, 0, 0, -s.ZNear, 0}
}

type Orthographic struct {
	XMag float64
	YMag float64
	ZNear float64
	ZFar float64
}
func NewOrthographic(XMag float64, YMag float64, ZNear float64, ZFar float64) *Orthographic {
	return &Orthographic{XMag: XMag, YMag: YMag, ZNear: ZNear, ZFar: ZFar}
}
func (s Orthographic) Projection(wicth, height float64) mgl64.Mat4 {
	return mgl64.Ortho(-s.XMag, s.XMag, -s.YMag, s.YMag, s.ZNear, s.ZFar)
}

type AutoOrthographic struct {
	ZNear float64
	ZFar float64
}
func (s AutoOrthographic) Projection(wicth, height float64) mgl64.Mat4 {
	xmag := wicth/2
	ymag := height/2
	return mgl64.Ortho(-xmag, xmag, -ymag, ymag, s.ZNear, s.ZFar)
}