package utils

func MapFromSlice[K comparable](slice []K) map[K]struct{} {
	dict := make(map[K]struct{}, len(slice))
	for _, item := range slice {
		dict[item] = struct{}{}
	}

	return dict
}

func MapContains[K comparable](dict map[K]struct{}, key K) bool {
	_, exists := dict[key]
	return exists
}
