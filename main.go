package main

import (
	"./geohash"
	"./toll"
	"fmt"
	"time"
)

func main3() {
	signedtolls := toll.GetTolls(-37.82613500, 144.96554400)
	if signedtolls != nil {
		for _, signedtoll := range signedtolls {
			fmt.Printf("site=%d||tollid=%d||box=[%.8f,%.8f;%.8f,%.8f;%.8f,%.8f;%.8f,%.8f]\n",
				signedtoll.Site,
				signedtoll.Toll.Id,
				signedtoll.Box.MaxLng, signedtoll.Box.MaxLat,
				signedtoll.Box.MinLng, signedtoll.Box.MaxLat,
				signedtoll.Box.MaxLng, signedtoll.Box.MinLat,
				signedtoll.Box.MinLng, signedtoll.Box.MinLat)
		}
	}
}

func main() {

	//geohash.Debug(-37.82613500,144.96554400, 8)
	//return
	starttime := time.Now()

	spend := time.Now().Sub(starttime)
	fmt.Println(spend)

	toll.Debug(110, 0)
}

func main1() {

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
