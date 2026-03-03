package service

func GetPointer[T any](val T) *T {
	return &val
}
