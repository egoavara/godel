// [ reference ] https://github.com/KhronosGroup/glTF-WebGL-PBR
package shader

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/iamGreedy/essence/version"
	"reflect"
)

type Shader struct {
	id        uuid.UUID
	src       string
	shdtype   ShaderType
	arguments map[string]*ShaderArg
}
type ShaderArg struct {
	Count int
	Kind  reflect.Kind
}
type ShaderType uint32

const (
	Vertex   ShaderType = 0x8B31
	Fragment ShaderType = 0x8B30
)

func NewShader(tp ShaderType, src string, args map[string]*ShaderArg) *Shader {
	return &Shader{
		id:        uuid.Must(uuid.NewRandom()),
		src:       src + "\x00",
		shdtype:   tp,
		arguments: args,
	}
}
func (s *Shader) ID() uuid.UUID {
	return s.id
}
func (s *Shader) Source(v version.Version, defines ...Define) string {
	var res string
	res = fmt.Sprintf("#version %d%d0\n", v.Major, v.Minor)
	for _, def := range defines {
		res += string(def)
		res += "\n"
	}
	res += s.src
	return res
}
func (s *Shader) Type() ShaderType {
	return s.shdtype
}
func (s *Shader) Args() []string {
	res := make([]string, 0, len(s.arguments))
	for key := range s.arguments {
		res = append(res, key)
	}
	return res
}
func (s *Shader) Arg(key string) *ShaderArg {
	return s.arguments[key]
}

var (
	Standard *Shader
	Flat     *Shader
	PBR      *Shader
)

func init() {
	Standard = NewShader(Vertex, string(FileStandardVsGlsl), map[string]*ShaderArg{
		//"CameraMatrix": {
		//	Count: 16,
		//	Kind:  reflect.Float32,
		//},
		//"ModelMatrix": {
		//	Count: 16,
		//	Kind:  reflect.Float32,
		//},
		//"NormalMatrix": {
		//	Count: 16,
		//	Kind:  reflect.Float32,
		//},
	})
	Flat = NewShader(Fragment, string(FileFlatFsGlsl), map[string]*ShaderArg{
		//"FlatColor": {
		//	Count: 4,
		//	Kind:  reflect.Float32,
		//},
	})
	PBR = NewShader(Fragment, string(FilePBRFsGlsl), map[string]*ShaderArg{})
}
