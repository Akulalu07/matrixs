package main

import "errors"

var ErrIrregularMatrix = errors.New("not matrix [m*n]x[n*k]")
var ErrNilMatrix = errors.New("nil matrix")
var ErrSquare = errors.New("matrix must be square")

type Script struct {
	s    string
	path string
}

type Matrix struct {
	lenm   int
	lenk   int
	matrix [][]Num
}

type Num struct {
	a int
	b int
}
