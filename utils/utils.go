package utils

import "fmt"

type any interface{}

func Filter[T any](slice []T, tester func(t T) bool) []T {
	var len int
	for _, v := range slice {
		if tester(v) {
			len++
		}
	}
	res := make([]T, 0, len)
	for _, v := range slice {
		if tester(v) {
			res = append(res, v)
		}
	}
	return res
}

func Reduce[T any, M any](s []T, f func(M, T) M, initValue M) M {
	acc := initValue
	for _, v := range s {
		acc = f(acc, v)
	}
	return acc
}

func Map[T any, M any](s []T, f func(T) M) []M {
	var a []M = make([]M, len(s))
	for i, v := range s {
		a[i] = f(v)
	}
	return a
}

type Queue[T any] []T

func (r *Queue[T]) Push(a T) {
	*r = append([]T(*r), a)
}

func (r *Queue[T]) Pop() {
	*r = (*r)[1:]
}

type ByteValue int64

func (r ByteValue) String() string {
	num := float64(r)
	if num < 1000 {
		return fmt.Sprint(int64(num), "B")
	} else if num < 1e6 {
		return fmt.Sprintf("%.1f%s", num/1e3, "KB")
	} else if num < 1e9 {
		return fmt.Sprintf("%.1f%s", num/1e6, "MB")
	} else {
		return fmt.Sprintf("%.1f%s", num/1e9, "GB")
	}
}

// transform map to array
func Flatten[T comparable, V any, M any](m map[T]V, flatterner func(V) M) []M {
	var a []M = make([]M, 0, len(m))
	for _, v := range m {
		a = append(a, flatterner(v))
	}
	return a
}
