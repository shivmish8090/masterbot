package helpers

import "sync"

func Contains[T comparable](slice []T, item T) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

func LoadTyped[T any](cache *sync.Map, key string) (T, bool) {
	var zero T
	if val, ok := cache.Load(key); ok {
		if typed, ok := val.(T); ok {
			return typed, true
		}
	}
	return zero, false
}
