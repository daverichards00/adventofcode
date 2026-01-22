package spatial

import (
	"maps"
	"sync"
)

type Space2D[T comparable] struct {
	m  map[Point2D]T
	bb Box2D
	mu sync.RWMutex
}

func NewSpace2D[T comparable]() *Space2D[T] {
	return &Space2D[T]{
		m:  map[Point2D]T{},
		bb: NilBox2D,
	}
}

func (s *Space2D[T]) Set(p Point2D, v T) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.m[p] = v
	s.bb = s.bb.Expand(p)
}

func (s *Space2D[T]) Unset(p Point2D) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.m, p)
	s.bboxReset()
}

func (s *Space2D[T]) Get(p Point2D) (T, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.m[p]
	return v, ok
}

func (s *Space2D[T]) GetOrDefault(p Point2D, d T) T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if v, ok := s.m[p]; ok {
		return v
	}
	return d
}

func (s *Space2D[T]) Exists(p Point2D) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, ok := s.m[p]
	return ok
}

func (s *Space2D[T]) GetAll() map[Point2D]T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := map[Point2D]T{}
	for p, v := range s.m {
		n[p] = v
	}
	return n
}

func (s *Space2D[T]) GetAllWithin(b Box2D) map[Point2D]T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := map[Point2D]T{}
	for p, v := range s.m {
		if b.Contains(p) {
			n[p] = v
		}
	}
	return n
}

func (s *Space2D[T]) GetAllAdjacent(p Point2D) map[Point2D]T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := map[Point2D]T{}
	for _, v := range []Vector2D{North2D, South2D, West2D, East2D} {
		pp := p.Add(v)
		if vv, ok := s.m[pp]; ok {
			n[pp] = vv
		}
	}
	return n
}

func (s *Space2D[T]) GetAllAdjacentAndDiagonal(p Point2D) map[Point2D]T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := map[Point2D]T{}
	for _, v := range []Vector2D{North2D, South2D, West2D, East2D, NorthWest2D, NorthEast2D, SouthWest2D, SouthEast2D} {
		pp := p.Add(v)
		if vv, ok := s.m[pp]; ok {
			n[pp] = vv
		}
	}
	return n
}

func (s *Space2D[T]) Equals(other *Space2D[T]) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return maps.Equal(s.m, other.m)
}

func (s *Space2D[T]) BBox() Box2D {
	return s.bb
}

func (s *Space2D[T]) bboxReset() {
	s.bb = NilBox2D
	for p := range s.m {
		s.bb = s.bb.Expand(p)
	}
}

func (s *Space2D[T]) Find(t T) []Point2D {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var r []Point2D
	for p, v := range s.m {
		if v == t {
			r = append(r, p)
		}
	}
	return r
}

func (s *Space2D[T]) Clone() *Space2D[T] {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := &Space2D[T]{m: map[Point2D]T{}, bb: s.bb}
	for p, v := range s.m {
		n.m[p] = v
	}
	return n
}

func (s *Space2D[T]) String(drawFunc func(Point2D, T) rune, dflt rune) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.bb == NilBox2D {
		return ""
	}
	out := ""
	for y := s.bb.Min().Y(); y <= s.bb.Max().Y(); y++ {
		for x := s.bb.Min().X(); x <= s.bb.Max().X(); x++ {
			p := NewPoint2D(x, y)
			if v, ok := s.m[p]; ok {
				out += string(drawFunc(p, v))
				continue
			}
			out += string(dflt)
		}
		out += "\n"
	}
	return out
}
