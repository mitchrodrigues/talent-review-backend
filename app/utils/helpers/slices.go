package helpers

func Contains[T comparable](list []T, values ...T) bool {
	for _, item := range list {
		for _, val := range values {
			if item == val {
				return true
			}
		}
	}
	return false
}
