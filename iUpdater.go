package godel

type Updater interface {
	deltaT(dt float64)
}