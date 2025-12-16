package numrange

import "golang.org/x/exp/constraints"

type NumRange[T constraints.Integer | constraints.Float] struct {
	mn, mx T
}

func New[T constraints.Integer | constraints.Float](a, b T) NumRange[T] {
	return NumRange[T]{
		mn: min(a, b),
		mx: max(a, b),
	}
}

func (r NumRange[T]) Min() T {
	return r.mn
}

func (r NumRange[T]) Max() T {
	return r.mx
}

func (r NumRange[T]) Intersect(other NumRange[T]) (NumRange[T], bool) {
	if r.mn > other.mx || r.mx < other.mn {
		return NumRange[T]{}, false
	}
	return New(max(r.mn, other.mn), min(r.mx, other.mx)), true
}
