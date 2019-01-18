package godel

import (
	"github.com/iamGreedy/essence/must"
	"image"
)

type Driver interface {
	UseContext()
	Viewport() image.Rectangle
}

type Application struct {
	driver Driver
	//
	skel map[string]*Skeleton
	skin map[string]*Skin
	cam  map[string]*Camera
}

func NewApplication(driver Driver) *Application {
	must.NotNil(driver)
	return &Application{driver: driver}
}
