package utils

func Map[T any, R any](source []T, cb func(T) R) []R {
	list := make([]R, 0, len(source))
	for _, i := range source {
		list = append(list, cb(i))
	}

	return list
}

func Filter[T any](s []T, keep func(T) bool) []T {
	d := make([]T, 0, len(s))
	for _, n := range s {
		if keep(n) {
			d = append(d, n)
		}
	}
	return d
}
