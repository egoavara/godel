package godel

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/iamGreedy/gltf2"
	"io"
)

type Player struct {
	ref *Instance
	//
	targets  []*playerTarget
	playtime float32
	//
	curr      float32
	PlaySpeed float32
	Loop      bool
	fPlay     bool
}
type playerTarget struct {
	target *Node
	path   gltf2.Path
	sample PlayerSampler
}

func (s *Instance) NewPlayer(i int, callback func(a *Player)) *Player {
	if s.anim != nil {
		s.anim.Close()
	}
	if i < 0 || i >= len(s.model.gltf.Animations) {
		s.anim = nil
		return nil
	}
	s.anim = new(Player)
	if err := s.setupAnimation(s.anim, s.model.gltf.Animations[i]); err != nil {
		s.anim = nil
	} else {
		if callback != nil {
			callback(s.anim)
		}
		s.model.app.append(s.anim)
	}
	return s.anim
}
func (s *Player) Close() error {
	for _, v := range s.targets {
		v.target.clearAnim()
	}
	s.ref.model.app.delete(s)
	return nil
}

func (s *Player) dt(t float32) {
	if s.fPlay {
		s.curr += s.PlaySpeed * t

		if s.Loop {
			for s.curr < 0 {
				s.curr = s.playtime + s.curr
			}
			for s.curr > s.playtime {
				s.curr = s.curr - s.playtime
			}
		}

		for _, t := range s.targets {
			v := t.sample.P(s.curr)
			switch t.path {
			case gltf2.Translation:
				t.target.setT(v.Vec3())
			case gltf2.Rotation:
				t.target.setR(mgl32.Quat{
					W: v.Get(3),
					V: v.Vec3(),
				})
			case gltf2.Scale:
				t.target.setS(v.Vec3())
			case gltf2.Weight:
				// TODO
			}
		}
	}
}

func (s *Player) Play() {
	s.fPlay = true
}
func (s *Player) Pause() {
	s.fPlay = false
}
func (s *Player) Stop() {
	s.fPlay = false
	s.Seek(0, io.SeekStart)
}
func (s *Player) Seek(x float32, whence int) float32 {
	switch whence {
	default:
		fallthrough
	case io.SeekStart:
		s.curr = mgl32.Clamp(x, 0, s.playtime)
	case io.SeekCurrent:
		s.curr = mgl32.Clamp(s.curr+x, 0, s.playtime)
	case io.SeekEnd:
		s.curr = mgl32.Clamp(s.playtime-x, 0, s.playtime)
	}
	temp := s.fPlay
	s.fPlay = true
	s.dt(0)
	s.fPlay = temp
	return s.curr
}
func (s *Player) Playtime() float32 {
	return s.playtime
}
func (s *Player) Current() float32 {
	return s.curr
}
