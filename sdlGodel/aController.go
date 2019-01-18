package sdlGodel

import "github.com/go-gl/mathgl/mgl64"

type Controller struct {
	proxySetup []IProxy
	proxyNames map[string]uint32
}

type Event interface {
	dummyEvent()
}
type RawEvent interface {
	Event
	dummyRaw()
}
type RawDeltaTime struct {
	Dt float64
}

func (s *RawDeltaTime) dummyEvent() {
	panic("implement me")
}
func (s *RawDeltaTime) dummyRaw() {
	panic("implement me")
}

type RawKey struct {
	Key Key
	IsPress bool
}

func (s *RawKey) dummyEvent() {
	panic("implement me")
}
func (s *RawKey) dummyRaw() {
	panic("implement me")
}

type RawClick struct {
	Click Click
}

func (s *RawClick) dummyEvent() {
	panic("implement me")
}
func (s *RawClick) dummyRaw() {
	panic("implement me")
}

type RawCursor struct {
	X float64
	Y float64
}

func (s *RawCursor) dummyEvent() {
	panic("implement me")
}
func (s *RawCursor) dummyRaw() {
	panic("implement me")
}

type RawScroll struct {
	Dx float64
	Dy float64
}

func (s *RawScroll) dummyEvent() {
	panic("implement me")
}
func (s *RawScroll) dummyRaw() {
	panic("implement me")
}

type ProxyEvent interface {
	Event
	GetID() uint32
}
type Action struct {
	ID uint32
}

func (s *Action) dummyEvent() {
	panic("implement me")
}
func (s *Action) GetID() uint32 {
	return s.ID
}

type Axis1D struct {
	ID   uint32
	Data float64
}

func (s *Axis1D) dummyEvent() {
	panic("implement me")
}
func (s *Axis1D) GetID() uint32 {
	return s.ID
}

type Axis2D struct {
	ID   uint32
	Data mgl64.Vec2
}

func (s *Axis2D) dummyEvent() {
	panic("implement me")
}
func (s *Axis2D) GetID() uint32 {
	return s.ID
}

type IProxy interface {
	IProxyName() string
	RequestNameID(fn func(name string) uint32)
	DoProxy(in []Event) (hijack bool, out []Event)
}

type KeyProxy struct {
	name                string
	keyPressMapping     map[Key]string
	keyReleaseMapping   map[Key]string
	keyPressMappingID   map[Key]uint32
	keyReleaseMappingID map[Key]uint32
}

func (s *KeyProxy) IProxyName() string {
	return s.name
}

func (s *KeyProxy) RequestNameID(fn func(name string) uint32) {
	for k, v := range s.keyPressMapping {
		s.keyPressMappingID[k] = fn(v)
	}
	for k, v := range s.keyReleaseMapping {
		s.keyReleaseMappingID[k] = fn(v)
	}
}

func (s *KeyProxy) DoProxy(in Event) (hijack bool, out []Event) {
	switch i := in.(type) {
	case *RawKey:
		if i.IsPress{
			if id, ok := s.keyPressMappingID[i.Key];ok{
				return false, []Event{&Action{ID:id}}
			}
		}else {
			if id, ok := s.keyReleaseMappingID[i.Key];ok{
				return false, []Event{&Action{ID:id}}
			}
		}
	}
	return false, nil
}
