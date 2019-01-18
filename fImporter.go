package godel

import (
	"github.com/iamGreedy/essence/must"
	"io"
	"os"
)

type Importer interface {
	// To Importer
	BeginStream(name string)
	Stream(data []byte)
	EndStream(name string)
	// From Importer
	Request() string
	Consume() Asset
	IsError() error
}
type RequestHandler interface {
	Available(target string) bool
	Use(target string) error
	io.Reader
}
type ImportManager struct {
	app *Application
	src io.Reader
	importer Importer
	//
	Buffer   []byte
	Total    int
	Name     string
	nameHint string
	//
	flushed []Asset
}
func (s *Application) Import(name string, src io.Reader, importer Importer) error {
	im := s.NewImportManager(name, src, importer)
	switch t := src.(type) {
	case *os.File:
		info := must.MustGet(t.Stat()).(os.FileInfo)
		im.nameHint = info.Name()
		im.Total = int(info.Size())
	}
	return im.Import()
}
func (s *Application) NewImportManager(name string, src io.Reader, importer Importer) *ImportManager{
	must.NotNil(src)
	must.NotNil(importer)
	return &ImportManager{
		app:s,
		importer:importer,
		src:src,
		//
		Buffer:   make([]byte, 2048),
		Total:    0,
		Name:     name,
		nameHint: "",
	}
}
func (s *ImportManager) Import() error {
	//var (
	//	n int
	//	err error
	//	curr int
	//)
	//var (
	//	buf = s.Buffer
	//	total = s.Total
	//	name = s.Name
	//)
	//if len(name) == 0{
	//	name = s.nameHint
	//}
	return nil
}
func ReadBeforeCheckError(test error, reader io.Reader, buf []byte) (n int, err error) {
	if test != nil{
		return 0, test
	}
	return reader.Read(buf)
}