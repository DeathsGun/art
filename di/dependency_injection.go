package di

import "sync"

var services sync.Map

func Set[T any](name string, service T) {
	services.Store(name, service)
}

func Instance[T any](name string) T {
	t, ok := services.Load(name)
	if ok {
		return t.(T)
	}
	var result T
	return result
}
