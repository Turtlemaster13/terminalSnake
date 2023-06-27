package main

type snake struct {
	direction    int
	length       int
	tail         []location
	headPosition location
}

type location struct {
	x, y int
}
