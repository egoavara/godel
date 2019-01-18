package godel1

type Updater interface {
	deltaT(dt float64)
}