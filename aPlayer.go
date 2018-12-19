package godel

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/iamGreedy/gltf2"
	"github.com/pkg/errors"
)

type Player struct {
	ref  *Renderer
	anim *gltf2.Animation
	//
	current float32
	//
	t0 float32
	p0 mgl32.VecN
	p1 mgl32.VecN

}

func (s *Renderer) NewPlayer(i int) (player *Player, err error) {
	if i < 0 || i >= len(s.model.Animations) {
		return nil, errors.New("Not found")
	}
	player = new(Player)
	player.ref = s
	player.anim = s.model.Animations[i]
	player.current = 0

	return player, nil
}
func (s *Player) Dt(dt float32) {
	s.current += dt
}

func (s *Player)p(t float32)  {

}

type PlayerStepSampler struct {
	in []float32
	out []float32
	outCount int
}
func (s *PlayerStepSampler) p(current float32) *mgl32.VecN {
	current = mgl32.Clamp(current, s.in[0], s.in[len(s.in) - 1])
	for i, v := range s.in{
		if v <= current{
			return mgl32.NewVecNFromData(s.out[i * s.outCount : (i + 1) * s.outCount])
		}
	}
	panic("Unreachable")
}

type PlayerLinearSampler struct {
	in []float32
	out []float32
	outCount int
}
func (s *PlayerLinearSampler) p(current float32) *mgl32.VecN {
	current = mgl32.Clamp(current, s.in[0], s.in[len(s.in) - 1])
	if current >= s.in[len(s.in) - 1]{
		i := len(s.in) - 1
		return mgl32.NewVecNFromData(s.out[i * s.outCount : (i + 1) * s.outCount])
	}
	for i, v := range s.in{
		if v <= current{
			nv := s.in[i + 1]
			t := (current - v) / (nv - v)
			p0 := mgl32.NewVecNFromData(s.out[i * s.outCount : (i + 1) * s.outCount])
			p1 := mgl32.NewVecNFromData(s.out[(i + 1) * s.outCount : (i + 2) * s.outCount])
			return p0.Add(nil, p1.Sub(nil, p0).Mul(nil, t))
		}
	}
	panic("Unreachable")
}


type PlayerCubicSampler struct {
	in []float32
	out []float32
	outCount int
	intangent []float32
	outtangent []float32
}
func (s *PlayerCubicSampler) p(current float32) *mgl32.VecN {
	current = mgl32.Clamp(current, s.in[0], s.in[len(s.in) - 1])
	if current >= s.in[len(s.in) - 1]{
		i := len(s.in) - 1
		return mgl32.NewVecNFromData(s.out[i * s.outCount : (i + 1) * s.outCount])
	}
	for i, v := range s.in{
		if v <= current{
			nv := s.in[i + 1]
			dt := nv - v
			t := (current - v) / dt
			p0 := mgl32.NewVecNFromData(s.out[i * s.outCount : (i + 1) * s.outCount])
			p1 := mgl32.NewVecNFromData(s.out[(i + 1) * s.outCount : (i + 2) * s.outCount])
			//
			return p0.Add(nil, p1.Sub(nil, p0).Mul(nil, t))
		}
	}
	panic("Unreachable")
}





