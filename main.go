package main

import (
	"./geohash"
	"fmt"
	"math"
)
func main(){

	var box *geohash.Box
	var sgeohash string
	var precision float64 = 0.01
	var hashprecision int = 6

	tollhash, box := geohash.Encode(-37.82496,144.97083,  hashprecision )

	width := box.Width()
	height := box.Height()

	fmt.Println("toll geohash:", tollhash, "width:", width, "height:", height)

	mylat := -37.82496
	mylon := 144.98083
	precision = math.Max(precision, width)
	for i:=mylat-precision;i<=mylat+precision;i+=height{
		for j:=mylon-precision;j<=mylon+precision;j+=width {
			sgeohash, _ = geohash.Encode(i,j,hashprecision)
			fmt.Println(sgeohash, i, j)
		}
	}

	return
}
