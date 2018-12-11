package ogl

import (
	"bytes"
	"fmt"
	"github.com/go-gl/gl/v4.1-core/gl"
	"github.com/go-gl/mathgl/mgl32"
	"image"
	"image/draw"
	"io"
	"io/ioutil"
	"reflect"
	"unsafe"
)

type OpenGL4 interface {
	//
	fmt.Stringer
	Major() int
	Minor() int

	//
	GetError() uint32
	// Program
	CreateProgram() uint32
	AttachShader(program uint32, shader uint32)
	LinkProgram(program uint32)
	GetProgramiv(program uint32, pname uint32) int32
	GetProgramInfoLog(program uint32) string
	GetAttribLocation(program uint32, name string) uint32
	DeleteProgram(program uint32)
	UseProgram(program uint32)
	GetUniformLocation(program uint32, name string) int32
	UniformMatrix4fv(uniformLocation int32, transpose bool, value ...mgl32.Mat4)
	Uniform1iv(uniformLocation int32, value ...int32)
	BindFragDataLocation(program uint32, color uint32, name string)
	// Shader
	CreateShader(shaderType uint32) uint32
	ShaderSource(shader uint32, data io.Reader)
	CompileShader(shader uint32)
	GetShaderiv(shader uint32, pname uint32) int32
	GetShaderInfoLog(shader uint32) string
	DeleteShader(shader uint32)
	// Texture
	GenTextures(n int) []uint32
	ActiveTexture(activeTexture uint32)
	BindTexture(targetTexture uint32, texture uint32)
	TexParameteri(targetTexture uint32, pname uint32, pvalue int32)
	TexImage2D(target uint32, level int32, internalformat int32, width int32, height int32, border int32, format uint32, xtype uint32, img image.Image)
	// Buffer
	GenBuffer(n int) []uint32
	BindBuffer(targetBuffer uint32, buffer uint32)
	BufferData(targetBuffer uint32, size uint32, data io.Reader, usage uint32)

	// Vertex Array
	GenVertexArrays(n int) []uint32
	BindVertexArray(vao uint32)
	EnableVertexAttribArray(attribLocation uint32)
	VertexAttribPointer(attribLocation uint32, size int32, xtype uint32, normalized bool, stride int32, offset int32)
}

type _OGL41 struct {
}

func (_OGL41) String() string {

	return "OpenGL 4.1 Core"
}

func (_OGL41) Major() int {
	return 4
}

func (_OGL41) Minor() int {
	return 1
}

func (_OGL41) GetError() uint32 {
	return gl.GetError()
}
func (_OGL41) CreateProgram() uint32 {
	return gl.CreateProgram()
}
func (_OGL41) AttachShader(program uint32, shader uint32) {
	gl.AttachShader(program, shader)
}
func (_OGL41) LinkProgram(program uint32) {
	gl.LinkProgram(program)
}
func (_OGL41) GetProgramiv(program uint32, pname uint32) int32 {
	var res int32
	gl.GetProgramiv(program, pname, &res)
	return res
}
func (s *_OGL41) GetProgramInfoLog(program uint32) string {
	log := make([]byte, s.GetProgramiv(program, gl.INFO_LOG_LENGTH)+1)
	header := (*reflect.StringHeader)(unsafe.Pointer(&log))
	gl.GetProgramInfoLog(program, int32(header.Len-1), nil, (*uint8)(unsafe.Pointer(header.Data)))
	return string(log)
}
func (s *_OGL41) GetAttribLocation(program uint32, name string) uint32 {
	return uint32(gl.GetAttribLocation(program, gl.Str(name)))
}
func (_OGL41) DeleteProgram(program uint32) {
	gl.DeleteProgram(program)
}
func (_OGL41) UseProgram(program uint32) {
	gl.UseProgram(program)
}
func (_OGL41) GetUniformLocation(program uint32, name string) int32 {
	return gl.GetUniformLocation(program, gl.Str(name))
}
func (_OGL41) UniformMatrix4fv(uniformLocation int32, transpose bool, value ...mgl32.Mat4) {
	if len(value) < 1 {
		return
	}
	gl.UniformMatrix4fv(uniformLocation, int32(len(value)), transpose, &value[0][0])
}
func (_OGL41) Uniform1iv(uniformLocation int32, value ...int32) {
	if len(value) < 0 {
		return
	}
	gl.Uniform1iv(uniformLocation, int32(len(value)), &value[0])
}
func (_OGL41) BindFragDataLocation(program uint32, color uint32, name string) {
	gl.BindFragDataLocation(program, color, gl.Str(name))
}
func (_OGL41) CreateShader(shaderType uint32) uint32 {
	return gl.CreateShader(shaderType)
}
func (_OGL41) ShaderSource(shader uint32, data io.Reader) {
	if data == nil {
		return
	}
	var str string
	if bb, ok := data.(*bytes.Buffer); ok {
		str = bb.String()
	} else {
		bts, err := ioutil.ReadAll(data)
		if err != nil {
			panic(err)
		}
		str = string(bts)
	}
	cstr, free := gl.Strs(str)
	defer free()
	gl.ShaderSource(shader, 1, cstr, nil)
}
func (_OGL41) CompileShader(shader uint32) {
	gl.CompileShader(shader)
}
func (_OGL41) GetShaderiv(shader uint32, pname uint32) int32 {
	var res int32
	gl.GetShaderiv(shader, pname, &res)
	return res
}
func (s *_OGL41) GetShaderInfoLog(shader uint32) string {
	log := make([]byte, s.GetShaderiv(shader, gl.INFO_LOG_LENGTH)+1)
	header := (*reflect.StringHeader)(unsafe.Pointer(&log))
	gl.GetShaderInfoLog(shader, int32(header.Len-1), nil, (*uint8)(unsafe.Pointer(header.Data)))
	return string(log)
}
func (_OGL41) DeleteShader(shader uint32) {
	gl.DeleteShader(shader)
}
func (_OGL41) GenTextures(n int) []uint32 {
	var res = make([]uint32, n)
	gl.GenTextures(int32(n), &res[0])
	return res
}
func (_OGL41) ActiveTexture(activeTexture uint32) {
	gl.ActiveTexture(activeTexture)
}
func (_OGL41) BindTexture(targetTexture uint32, texture uint32) {
	gl.BindTexture(targetTexture, texture)
}
func (_OGL41) TexParameteri(targetTexture uint32, pname uint32, pvalue int32) {
	gl.TexParameteri(targetTexture, pname, pvalue)
}
func (_OGL41) TexImage2D(target uint32, level int32, internalformat int32, width int32, height int32, border int32, format uint32, xtype uint32, img image.Image) {
	var rgba *image.RGBA
	var ok bool
	if rgba, ok = img.(*image.RGBA); !ok {
		rgba = image.NewRGBA(img.Bounds())
		draw.Draw(rgba, rgba.Rect, img, img.Bounds().Min, draw.Src)
	}
	gl.TexImage2D(target, level, internalformat, width, height, border, format, xtype, gl.Ptr(&rgba.Pix[0]))
}
func (_OGL41) GenBuffer(n int) []uint32 {
	var res = make([]uint32, n)
	gl.GenBuffers(int32(n), &res[0])
	return res
}
func (_OGL41) BindBuffer(targetBuffer uint32, buffer uint32) {
	gl.BindBuffer(targetBuffer, buffer)
}
func (_OGL41) BufferData(targetBuffer uint32, size uint32, data io.Reader, usage uint32) {
	var (
		bts []byte
	)
	if buf, ok := data.(*bytes.Buffer); ok {
		bts = buf.Bytes()
	} else {
		var err error
		bts, err = ioutil.ReadAll(data)
		if err != nil {
			panic(err)
		}
	}
	gl.BufferData(targetBuffer, int(size), gl.Ptr(&bts[0]), usage)
}
func (_OGL41) GenVertexArrays(n int) []uint32 {
	var res = make([]uint32, n)
	gl.GenVertexArrays(int32(n), &res[0])
	return res
}
func (_OGL41) BindVertexArray(vao uint32) {
	gl.BindVertexArray(vao)
}
func (_OGL41) EnableVertexAttribArray(attribLocation uint32) {
	gl.EnableVertexAttribArray(attribLocation)
}
func (_OGL41) VertexAttribPointer(attribLocation uint32, size int32, xtype uint32, normalized bool, stride int32, offset int32) {
	gl.VertexAttribPointer(attribLocation, size, xtype, normalized, stride, gl.PtrOffset(int(offset)))
}
