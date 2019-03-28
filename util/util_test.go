package util

import (
	"fmt"
	"testing"
)

func TestDirection(t *testing.T) {

	A := point{-0.5, -1}
	C := point{0.6, -1}
	D := point{0.5, 1}
	fmt.Println(direction(&A, &C, &D))

	E := point{0, 0}
	F := point{1, 0}
	G := point{0, -1}
	fmt.Println(direction(&E, &F, &G))
}

func TestIntersect(t *testing.T) {
	A := point{-0.5, -1}
	B := point{1, 1}
	C := point{-0.5, -1.1}
	D := point{1, 0.5}
	fmt.Println(intersect(&A, &B, &C, &D))
}
