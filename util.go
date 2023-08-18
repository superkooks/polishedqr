package main

import (
	"image"
	"math"
)

func iterateRect(x, y int, callback func(x, y int)) {
	for z := 0; z < y; z++ {
		for w := 0; w < x; w++ {
			callback(w, z)
		}
	}
}

func vecLen(x image.Point) float64 {
	return math.Sqrt(math.Pow(float64(x.X), 2) + math.Pow(float64(x.Y), 2))
}

func reverseMap[K, V comparable](in map[K]V) map[V]K {
	out := make(map[V]K)
	for k, v := range in {
		out[v] = k
	}
	return out
}
