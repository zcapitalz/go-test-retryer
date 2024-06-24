package retryer

func filter[T any](ss []T, f ...func(T) bool) (result []T) {
	for _, s := range ss {
		add := true
		for _, f := range f {
			if add = f(s); !add {
				break
			}
		}
		if add {
			result = append(result, s)
		}
	}
	return
}
