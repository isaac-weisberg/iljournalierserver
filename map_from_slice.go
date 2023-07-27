package main

func mapFromSlice[K comparable](slice []K) map[K]struct{} {
	dict := make(map[K]struct{}, len(slice))
	for _, item := range slice {
		dict[item] = struct{}{}
	}

	return dict
}

func mapContains[K comparable](dict map[K]struct{}, key K) bool {
	_, exists := dict[key]
	return exists
}
