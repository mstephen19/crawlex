package core

func defaultIfZero[T comparable](value *T, zero, def T) {
	if *value == zero {
		*value = def
	}
}

func shift[T any](slice *[]T) (item T) {
	item = (*slice)[0]
	*slice = (*slice)[1:]
	return
}
