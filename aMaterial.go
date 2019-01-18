package godel

import (
	"github.com/go-gl/gl/v3.2-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"github.com/go-gl/mathgl/mgl64"
	"unsafe"
)

type Material interface {
	BindBlock() string
	BufferPointer() uint32
	//
	DeallocateBuffer() uint32
	AllocateBuffer(bufptr uint32) bool
	UploadBuffer()
}
//type DebugMaterial struct {
//}
type (
	EmissiveMaterial struct {
		EmissiveFactor               mgl64.Vec3
		EmissiveTexure               *Texture
		EmissiveTexureCoord          int
		//
		buf uint32
	}
	std140EmissiveMaterial struct {
		EmissiveFactor               mgl32.Vec3
		EmissiveTexure               int32
		EmissiveTexureCoord          int32
	}
)
func (s *EmissiveMaterial) BindBlock() string {
	return "EmissiveMaterial"
}
func (s *EmissiveMaterial) BufferPointer() uint32 {
	return s.buf
}
func (s *EmissiveMaterial) DeallocateBuffer() uint32{
	temp := s.buf
	s.buf = gl.INVALID_INDEX
	return temp
}
func (s *EmissiveMaterial) AllocateBuffer(bufptr uint32) bool{
	if s.buf == gl.INVALID_INDEX{
		s.buf = bufptr
		return true
	}
	return false
}
func (s *EmissiveMaterial) UploadBuffer() {
	if s.buf == gl.INVALID_INDEX{
		return
	}
	var temp = std140EmissiveMaterial{
		EmissiveFactor:dVec3ToVec3(s.EmissiveFactor),
		EmissiveTexure: 0,
		EmissiveTexureCoord:int32(s.EmissiveTexureCoord),
	}
	gl.BindBuffer(gl.ARRAY_BUFFER, s.buf)
	gl.BufferData(gl.ARRAY_BUFFER, int(unsafe.Sizeof(temp)), gl.Ptr(&temp), gl.STATIC_DRAW)
}

//type PhongMaterial struct {
//}

type PBRStandardMaterial struct {
	BaseColorFactor              mgl64.Vec4
	MetallicRoughnessFactor      mgl64.Vec2
	EmissiveFactor               mgl64.Vec3
	BaseColorTexure              *Texture
	BaseColorTexureCoord         int
	MetallicRoughnessTexure      *Texture
	MetallicRoughnessTexureCoord int
	NormalTexure                 *Texture
	NormalTexureCoord            int
	NormalScale                  float64
	OcculusionTexure             *Texture
	OcculusionTexureCoord        int
	OcculusionStrength           float64
	EmissiveTexure               *Texture
	EmissiveTexureCoord          int
}