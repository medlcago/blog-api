package storage

type KeyFunc func(namespace, key string) string

func DefaultKeyFunc(namespace, key string) string {
	return namespace + "_" + key
}
