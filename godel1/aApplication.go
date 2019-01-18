package godel1

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
}

func NewApplication(driver Driver) *Application {
	must.NotNil(driver)
	return &Application{driver: driver}
}
