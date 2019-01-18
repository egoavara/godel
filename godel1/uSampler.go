package godel1

import (
	"gonum.org/v1/gonum/floats"
	"math"
	"sort"
)

type Sampler interface {
	P(current float64) []float64
	Range() (min, max float64)
}

func inputIsValid(in []float64) bool {
	return sort.Float64sAreSorted(in)
}
func outputValidate(in []float64, out []float64, outcount int, mode int) bool {
	l := len(in)
	if l*outcount*mode != len(out) {
		return false
	}
	return true
}
func neighbor(in []float64, a float64) (lidx, ridx int, diff float64) {
	if a <= in[0] {
		return 0, 0, 0
	}
	for i := 1; i < len(in); i++ {
		if a < in[i] {
			return i - 1, i, (a - in[i-1]) / (in[i] - in[i-1])
		}
	}
	return len(in) - 1, len(in) - 1, 0
}
func lerp(dst, a, b []float64, diff float64) []float64 {
	floats.SubTo(dst, b, a)
	floats.Scale(diff, dst)
	floats.Add(dst, a)
	return dst
}
func norm(dst []float64) []float64 {
	floats.Scale(math.Sqrt(floats.Dot(dst, dst)), dst)
	return dst
}
func normTo(dst, a []float64) []float64 {
	copy(dst, a)
	floats.Scale(math.Sqrt(floats.Dot(dst, dst)), dst)
	return dst
}

type StepSampler struct {
	in    []float64
	out   []float64
	count int
}

func NewStepSampler(in, out []float64, count int) *StepSampler {
	if !inputIsValid(in) {
		return nil
	}
	if !outputValidate(in, out, count, 1) {
		return nil
	}
	return &StepSampler{
		in:    in,
		out:   out,
		count: count,
	}
}
func (s *StepSampler) P(current float64) []float64 {
	var l, _, _ = neighbor(s.in, current)
	return s.out[s.count*l : s.count*(l+1)]
}
func (s *StepSampler) Range() (min, max float64) {
	return s.in[0], s.in[len(s.in)-1]
}

type LinearSampler struct {
	in    []float64
	out   []float64
	count int
}

func NewLinearSampler(in, out []float64, count int) *LinearSampler {
	if !inputIsValid(in) {
		return nil
	}
	if !outputValidate(in, out, count, 1) {
		return nil
	}
	return &LinearSampler{
		in:    in,
		out:   out,
		count: count,
	}
}
func (s *LinearSampler) P(current float64) []float64 {
	var l, r, diff = neighbor(s.in, current)
	var lo, ro = s.out[s.count*l : s.count*(l+1)], s.out[s.count*r : s.count*(r+1)]
	return lerp(make([]float64, s.count), lo, ro, diff)
}
func (s *LinearSampler) Range() (min, max float64) {
	return s.in[0], s.in[len(s.in)-1]
}

type SphericalLinearSampler struct {
	in    []float64
	out   []float64
	count int
}

func NewSphericalLinearSampler(in, out []float64, count int) *SphericalLinearSampler {
	if !inputIsValid(in) {
		return nil
	}
	if !outputValidate(in, out, count, 1) {
		return nil
	}
	return &SphericalLinearSampler{
		in:    in,
		out:   out,
		count: count,
	}
}
func (s *SphericalLinearSampler) P(current float64) []float64 {
	var l, r, diff = neighbor(s.in, current)
	var lo, ro = s.out[s.count*l : s.count*(l+1)], s.out[s.count*r : s.count*(r+1)]
	//
	var f0, f1 = normTo(make([]float64, s.count), lo), normTo(make([]float64, s.count), ro)
	dot := floats.Dot(f0, f1)
	if dot < 0 {
		floats.Scale(-1, f1)
		dot = -dot
	}
	if dot > .9995 {
		return normTo(make([]float64, s.count), lerp(make([]float64, s.count), f0, f1, diff))
	}

	theta := math.Acos(dot) * diff
	cos, sin := math.Cos(theta), math.Sin(theta)
	//
	relin := floats.ScaleTo(make([]float64, s.count), dot, f0)
	floats.Sub(f1, relin)
	norm(f1)
	//
	floats.Scale(cos, f0)
	floats.Scale(sin, f1)
	floats.Add(f0, f1)
	return f0
}
func (s *SphericalLinearSampler) Range() (min, max float64) {
	return s.in[0], s.in[len(s.in)-1]
}

type CubicSampler struct {
	in    []float64
	out   []float64
	count int
}
func NewCubicSampler(in, out []float64, count int) *CubicSampler {
	if !inputIsValid(in) {
		return nil
	}
	if !outputValidate(in, out, count, 3) {
		return nil
	}
	return &CubicSampler{
		in:    in,
		out:   out,
		count: count,
	}
}
func (s *CubicSampler) P(current float64) []float64 {
	var l, r, diff = neighbor(s.in, current)
	var lstart = l * s.count * 3
	var rstart = r * s.count * 3
	var (
		p0 = s.out[lstart+1*s.count : lstart+2*s.count]
		m0 = s.out[lstart+2*s.count : lstart+3*s.count]
	)
	var (
		m1 = s.out[rstart+0*s.count : rstart+1*s.count]
		p1 = s.out[rstart+1*s.count : rstart+2*s.count]
	)
	//
	diff3 := diff * diff * diff
	diff2 := diff * diff
	//
	a := 2*diff3 - 3*diff2 + 1
	b := diff3 - 2*diff2 + diff
	c := -2*diff3 + 3*diff2
	d := diff3 - diff2
	//
	ea := floats.ScaleTo(make([]float64, s.count), a, p0)
	eb := floats.ScaleTo(make([]float64, s.count), b, m0)
	ec := floats.ScaleTo(make([]float64, s.count), c, m1)
	ed := floats.ScaleTo(make([]float64, s.count), d, p1)
	floats.Add(ea, eb)
	floats.Add(ea, ec)
	floats.Add(ea, ed)
	return ea
}
func (s *CubicSampler) Range() (min, max float64) {
	return s.in[0], s.in[len(s.in)-1]
}
