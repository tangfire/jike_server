package utils

// Map 对切片中的每个元素执行转换函数，返回新的切片
func Slice2Slice[T any, R any](slice []T, fn func(T) R) []R {
	result := make([]R, len(slice))
	for i, item := range slice {
		result[i] = fn(item)
	}
	return result
}

// SliceFilter 过滤并映射，只有满足条件的元素才会被转换
func SliceFilter[T any, R any](slice []T, fn func(T) (R, bool)) []R {
	var result []R
	for _, item := range slice {
		if mapped, ok := fn(item); ok {
			result = append(result, mapped)
		}
	}
	return result
}
