package pogo

import "fmt"

func SomeInArray[T any](arr []T, contains func(T) bool) bool {
	for _, v := range arr {
		if contains(v) {
			return true
		}
	}
	return false
}

func EveryInArray[T any](arr []T, contains func(T) bool) bool {
	for _, v := range arr {
		if !contains(v) {
			return false
		}
	}
	return true
}

func FilterArray[T any](arr []T, filter func(T) bool) []T {
	var result []T
	for _, v := range arr {
		if filter(v) {
			result = append(result, v)
		}
	}
	return result
}

func MapArray[T any, R any](arr []T, mapFunc func(T) R) []R {
	var result []R
	for _, v := range arr {
		result = append(result, mapFunc(v))
	}
	return result
}

func ContainsArray[T any](arr []T, value T) bool {
	for _, v := range arr {
		if fmt.Sprintf("%v", v) == fmt.Sprintf("%v", value) {
			return true
		}
	}
	return false
}
