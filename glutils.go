package godel

import (
	"github.com/go-gl/gl/v4.1-core/gl"
	"image"
	"strings"
)

func buildProgram(srcvs string, srcfs string) uint32 {
	// shader
	vs := buildShader(srcvs, gl.VERTEX_SHADER)
	fs := buildShader(srcfs, gl.FRAGMENT_SHADER)
	//
	prog := gl.CreateProgram()
	gl.AttachShader(prog, vs)
	gl.AttachShader(prog, fs)
	gl.LinkProgram(prog)
	gl.DeleteShader(vs)
	gl.DeleteShader(fs)
	//
	var status int32
	gl.GetProgramiv(prog, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(prog, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(prog, logLength, nil, gl.Str(log))
		panic(log)
	}
	gl.BindFragDataLocation(prog, 0, gl.Str("outputColor\x00"))
	return prog
}
func buildShader(src string, shdtype uint32) uint32 {
	shd := gl.CreateShader(shdtype)
	cstr, free := gl.Strs(src)
	gl.ShaderSource(shd, 1, cstr, nil)
	gl.CompileShader(shd)
	free()
	//
	var status int32
	gl.GetShaderiv(shd, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shd, gl.INFO_LOG_LENGTH, &logLength)
		log := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shd, logLength, nil, gl.Str(log))
		panic(log)
	}
	return shd
}

func viewportSize() image.Rectangle {
	var bound [4]int32
	gl.GetIntegerv(gl.VIEWPORT, &bound[0])
	return image.Rectangle{
		Min: image.Point{
			X: int(bound[0]),
			Y: int(bound[1]),
		},
		Max: image.Point{
			X: int(bound[2]),
			Y: int(bound[3]),
		},
	}
}
