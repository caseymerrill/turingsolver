package optional

type Optional[T any] struct {
	value    T
	hasValue bool
}

func (o *Optional[T]) Value(defaultValue ...T) T {
	if o.hasValue {
		return o.value
	} else if len(defaultValue) > 0 {
		return defaultValue[0]
	} else {
		return o.zero()
	}
}

func (o *Optional[T]) HasValue() bool {
	return o.hasValue
}

func (o *Optional[T]) Set(value T) {
	o.value = value
	o.hasValue = true
}

func (o *Optional[T]) Unset() {
	o.hasValue = false
	o.value = o.zero()
}

func (o *Optional[T]) zero() (zero T) {
	return zero
}

