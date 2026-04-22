package lazy_atomic

import (
	"sync/atomic"
)

type Value[T any] struct {
	atom        atomic.Value
	initialized bool
}

func (v *Value[T]) init() {
	if v.initialized {
		return
	}
	v.initialized = true

	var default_value T
	v.atom.Store(default_value)
}

func (v *Value[T]) CompareAndSwap(old T, new T) (swapped bool) {
	v.init()
	return v.atom.CompareAndSwap(old, new)
}

func (v *Value[T]) Load() (val T) {
	v.init()
	return v.atom.Load().(T)
}

func (v *Value[T]) Store(val T) {
	v.init()
	v.atom.Store(val)
}

func (v *Value[T]) Swap(new T) (old T) {
	v.init()
	return v.atom.Swap(new).(T)
}
