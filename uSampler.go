package godel

import (
	"github.com/go-gl/mathgl/mgl32"
	"github.com/iamGreedy/gltf2"
	"github.com/pkg/errors"
	"math"
)

type PlayerSampler interface {
	P(current float32) *mgl32.VecN
	Range() (min, max float32)
}

func MakeSampler(sampler *gltf2.AnimationSampler, normalize bool, rotation bool) (PlayerSampler, error) {
	switch sampler.Interpolation {
	case gltf2.STEP:
		if sampler.Input.Type != gltf2.SCALAR {
			return nil, errors.New("Must be scalar type")
		}
		if sampler.Input.ComponentType != gltf2.FLOAT {
			return nil, errors.New("Must be float type")
		}
		return &PlayerStepSampler{
			in:       sampler.Input.MustSliceMapping(new([]float32), true, true).([]float32),
			out:      nFloat32s(sampler.Output, normalize),
			outCount: sampler.Output.Type.Count(),
		}, nil
	case gltf2.LINEAR:
		if sampler.Input.Type != gltf2.SCALAR {
			return nil, errors.New("Must be scalar type")
		}
		if sampler.Input.ComponentType != gltf2.FLOAT {
			return nil, errors.New("Must be float type")
		}
		if rotation{
			return &PlayerSphereLinearSampler{
				in:       sampler.Input.MustSliceMapping(new([]float32), true, true).([]float32),
				out:      nFloat32s(sampler.Output, normalize),
				outCount: sampler.Output.Type.Count(),
			}, nil
		}
		return &PlayerLinearSampler{
			in:       sampler.Input.MustSliceMapping(new([]float32), true, true).([]float32),
			out:      nFloat32s(sampler.Output, normalize),
			outCount: sampler.Output.Type.Count(),
		}, nil
	case gltf2.CUBICSPLINE:
		if sampler.Input.Type != gltf2.SCALAR {
			return nil, errors.New("Must be scalar type")
		}
		if sampler.Input.ComponentType != gltf2.FLOAT {
			return nil, errors.New("Must be float type")
		}
		return &PlayerCubicSampler{
			in:       sampler.Input.MustSliceMapping(new([]float32), true, true).([]float32),
			out:      nFloat32s(sampler.Output, normalize),
			outCount: sampler.Output.Type.Count(),
		}, nil
	}
	panic("Unreachable")
}

type PlayerStepSampler struct {
	in       []float32
	out      []float32
	outCount int
}

func (s *PlayerStepSampler) Range() (min, max float32) {
	return s.in[0], s.in[len(s.in)-1]
}
func (s *PlayerStepSampler) P(current float32) *mgl32.VecN {
	current = mgl32.Clamp(current, s.in[0], s.in[len(s.in)-1])
	for i, v := range s.in {
		if current < v {
			return mgl32.NewVecNFromData(s.out[(i-1)*s.outCount : (i+0)*s.outCount])
		}
	}
	panic("Unreachable")
}

type PlayerLinearSampler struct {
	in       []float32
	out      []float32
	outCount int
}

func (s *PlayerLinearSampler) Range() (min, max float32) {
	return s.in[0], s.in[len(s.in)-1]
}
func (s *PlayerLinearSampler) P(current float32) *mgl32.VecN {
	current = mgl32.Clamp(current, s.in[0], s.in[len(s.in)-1])
	if current >= s.in[len(s.in)-1] {
		i := len(s.in) - 1
		return mgl32.NewVecNFromData(s.out[i*s.outCount : (i+1)*s.outCount])
	}
	for i, k1 := range s.in {
		if current < k1 {
			k0 := s.in[i-1]
			dt := k1 - k0
			t := (current - k0) / dt
			startk0 := (i - 1) * s.outCount
			startk1 := (i) * s.outCount
			p0 := mgl32.NewVecNFromData(s.out[startk0 : startk0+s.outCount])
			p1 := mgl32.NewVecNFromData(s.out[startk1 : startk1+s.outCount])
			return p0.Add(nil, p1.Sub(nil, p0).Mul(nil, t))
		}
	}
	panic("Unreachable")
}

type PlayerSphereLinearSampler struct {
	in       []float32
	out      []float32
	outCount int
}

func (s *PlayerSphereLinearSampler) Range() (min, max float32) {
	return s.in[0], s.in[len(s.in)-1]
}
func (s *PlayerSphereLinearSampler) P(current float32) *mgl32.VecN {
	current = mgl32.Clamp(current, s.in[0], s.in[len(s.in)-1])
	if current >= s.in[len(s.in)-1] {
		i := len(s.in) - 1
		return mgl32.NewVecNFromData(s.out[i*s.outCount : (i+1)*s.outCount])
	}
	for i, k1 := range s.in {
		if current < k1 {
			k0 := s.in[i-1]
			dt := k1 - k0
			t := (current - k0) / dt
			startk0 := (i - 1) * s.outCount
			startk1 := (i) * s.outCount
			p0 := mgl32.NewVecNFromData(s.out[startk0 : startk0+s.outCount])
			p1 := mgl32.NewVecNFromData(s.out[startk1 : startk1+s.outCount])
			//
			q0 := mgl32.Quat{
				W : p0.Get(3),
				V:p0.Vec3(),
			}
			q1 := mgl32.Quat{
				W : p1.Get(3),
				V:p1.Vec3(),
			}
			//res := mgl32.QuatSlerp(q0, q1, t)
			res := slerp(q0, q1, t)
			return mgl32.NewVecNFromData([]float32{
				res.V[0],
				res.V[1],
				res.V[2],
				res.W,
			})
		}
	}
	panic("Unreachable")
}

type PlayerCubicSampler struct {
	in       []float32
	out      []float32
	outCount int
}

func (s *PlayerCubicSampler) Range() (min, max float32) {
	return s.in[0], s.in[len(s.in)-1]
}
func (s *PlayerCubicSampler) P(current float32) *mgl32.VecN {
	current = mgl32.Clamp(current, s.in[0], s.in[len(s.in)-1])
	if current >= s.in[len(s.in)-1] {
		i := len(s.in) - 1
		return mgl32.NewVecNFromData(s.out[i*s.outCount : (i+1)*s.outCount])
	}
	for i, k1 := range s.in {
		if current < k1 {
			k0 := s.in[i-1]
			dt := k1 - k0
			startk0 := (i - 1) * s.outCount * 3
			startk1 := (i) * s.outCount * 3
			//
			t := (current - k0) / dt
			p0 := mgl32.NewVecNFromData(s.out[startk0+s.outCount*1 : startk0+s.outCount*2])
			m0 := mgl32.NewVecNFromData(s.out[startk0+s.outCount*2 : startk0+s.outCount*3])
			m1 := mgl32.NewVecNFromData(s.out[startk1+s.outCount*0 : startk1+s.outCount*1])
			p1 := mgl32.NewVecNFromData(s.out[startk1+s.outCount*1 : startk1+s.outCount*2])
			// P(t)
			a := 2*t*t*t - 3*t*t + 1
			b := t*t*t - 2*t*t + t
			c := -2*t*t*t + 3*t*t
			d := t*t*t - t*t
			return p0.Mul(nil, a).Add(nil, m0.Mul(nil, b)).Add(nil, p1.Mul(nil, c)).Add(nil, m1.Mul(nil, d))
		}
	}
	panic("Unreachable")
}
func nFloat32s(accessor *gltf2.Accessor, normalize bool) []float32 {
	if accessor.ComponentType != gltf2.FLOAT {
		if !normalize {
			return nil
		}
	}
	var res []float32
	switch accessor.ComponentType {
	case gltf2.FLOAT:
		accessor.MustSliceMapping(&res, false, true)
	case gltf2.BYTE:
		data := accessor.MustSliceMapping(new([]int8), false, true).([]int8)
		res = make([]float32, 0, len(data))
		for _, v := range data {
			res = append(res, mgl32.Clamp(float32(v)/127.0, -1, 1))
		}
	case gltf2.UNSIGNED_BYTE:
		data := accessor.MustSliceMapping(new([]uint8), false, true).([]uint8)
		res = make([]float32, 0, len(data))
		for _, v := range data {
			res = append(res, mgl32.Clamp(float32(v)/255.0, 0, 1))
		}
	case gltf2.SHORT:
		data := accessor.MustSliceMapping(new([]int16), false, true).([]int16)
		res = make([]float32, 0, len(data))
		for _, v := range data {
			res = append(res, mgl32.Clamp(float32(v)/32767.0, -1, 1))
		}
	case gltf2.UNSIGNED_SHORT:
		data := accessor.MustSliceMapping(new([]uint16), false, true).([]uint16)
		res = make([]float32, 0, len(data))
		for _, v := range data {
			res = append(res, mgl32.Clamp(float32(v)/65535.0, 0, 1))
		}
	}
	return res
}


func slerp(q0 mgl32.Quat, q1 mgl32.Quat, amount float32) mgl32.Quat {
	q0, q1 = q0.Normalize(), q1.Normalize()
	dot := q0.Dot(q1)
	if dot < 0{
		q1 = mgl32.Quat{
			W: -q1.W,
			V : mgl32.Vec3{
				-q1.V[0],
				-q1.V[1],
				-q1.V[2],
			},
		}
		dot = -dot
	}
	// If the inputs are too close for comfort, linearly interpolate and normalize the result.
	if dot > 0.9995 {
		return mgl32.QuatNlerp(q0, q1, amount)
	}

	// This is here for precision errors, I'm perfectly aware that *technically* the dot is bound [-1,1], but since Acos will freak out if it's not (even if it's just a liiiiitle bit over due to normal error) we need to clamp it
	dot = mgl32.Clamp(dot, 0, 1)

	theta := float32(math.Acos(float64(dot))) * amount
	c, s := float32(math.Cos(float64(theta))), float32(math.Sin(float64(theta)))
	rel := q1.Sub(q0.Scale(dot)).Normalize()
	return q0.Scale(c).Add(rel.Scale(s))
}