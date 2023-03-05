package main

func iterateRect(x, y int, callback func(x, y int)) {
	for z := 0; z < y; z++ {
		for w := 0; w < x; w++ {
			callback(w, z)
		}
	}
}
