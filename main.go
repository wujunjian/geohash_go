package main

import (
	"./geohash"
	"fmt"
	"time"
)

func main() {

	var box *geohash.Box
	var precision float64 = 0.01 //精度的上限,格子不能全部超出该范围 0.01(889-1113m)  0.001(88.9-111.3m)
	var hashprecision int = 6    //精度的下线,一个格子代表的面积    6(1223*488)  7(150*120)

	mylat := -37.82496
	mylon := 144.98083

	fmt.Println(time.Now())
	tollhash, box := geohash.Encode(mylat, mylon, hashprecision)

	fmt.Println("toll geohash:", tollhash, "width:", box.Width(), "height:", box.Height())

	geohash.GetNearGeoHash(mylat, mylon, precision, hashprecision)
	fmt.Println(time.Now())

	return
}
