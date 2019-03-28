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

func TestInrail(t *testing.T) {

	A := point{40.082274, 116.339808}
	B := point{39.922376, 116.949549}
	C := point{39.491324, 116.499109}
	D := point{39.660685, 115.817957}

	E := point{39.795876, 116.466150}

	rail := make([]*point, 0)
	rail = append(rail, &A, &B, &C, &D, &A)
	fmt.Println(inrail(&E, &B, rail))
}
