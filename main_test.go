package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"

	"./geohash"
	"./toll"
	"go.intra.xiaojukeji.com/gulfstream/plutus/common/helpers"
)

func parseJson(file string) []*helpers.Point {
	type Cood struct {
		Lat float64 `json:"lat,string"`
		Lng float64 `json:"lng,string"`
	}
	type OriData struct {
		Data []*Cood `json:"data"`
	}

	b, _ := ioutil.ReadFile(file)
	d := &OriData{}
	json.Unmarshal(b, d)

	points := []*helpers.Point{}
	for _, c := range d.Data {
		p := &helpers.Point{
			Lat: c.Lat,
			Lng: c.Lng,
		}
		points = append(points, p)
	}
	return points
}

func printpoint(points []*helpers.Point) {
	for _, p := range points {
		fmt.Printf("%.8f,%.8f;", p.Lng, p.Lat)
	}
}

func TestMatch(t *testing.T) {
	//init需要一些时间100-200ms左右.
	time.Sleep(500 * time.Millisecond)
	points := parseJson("./track_1.json")
	begin := time.Now()
	tm := &toll.TollsMatched{}
	tm.Match(points, "AU")
	tm.DebugPrint()
	fmt.Println(time.Now().Sub(begin))
	printpoint(points)
}

func TestGetTolls(t *testing.T) {

	//init需要一些时间, 100-200ms左右.
	time.Sleep(500 * time.Millisecond)

	points := parseJson("./track_1.json")
	begin := time.Now()
	for _, p := range points {
		signedtolls := toll.GetTolls(p.Lat, p.Lng, "AU")
		if signedtolls != nil {
			for _, signedtoll := range signedtolls {
				fmt.Printf("site=%s||tollid=%d||box=[%.8f,%.8f;%.8f,%.8f;%.8f,%.8f;%.8f,%.8f]\n",
					signedtoll.SiteString(),
					signedtoll.Toll.Id,
					signedtoll.Box.MaxLng, signedtoll.Box.MaxLat,
					signedtoll.Box.MinLng, signedtoll.Box.MaxLat,
					signedtoll.Box.MaxLng, signedtoll.Box.MinLat,
					signedtoll.Box.MinLng, signedtoll.Box.MinLat)
			}
		}
	}
	fmt.Println(time.Now().Sub(begin))
}

//获取收费站信息
func TestDebug(t *testing.T) {

	begin := time.Now()
	//toll.Debug(111, 1, "AU")
	//toll.Debug(111, 0, "AU")
	toll.Debug(1, -1, "AU")
	fmt.Println(time.Now().Sub(begin))
}

func TestEncode(t *testing.T) {

	var box *geohash.Box
	var precision = 0.01  //精度的上限,格子不能全部超出该范围 0.01(889-1113m)  0.001(88.9-111.3m)
	var hashprecision = 6 //精度的下线,一个格子代表的面积    6(1223*488)  7(150*120)

	mylat := -37.82496
	mylon := 144.98083

	tollhash, box := geohash.Encode(mylat, mylon, hashprecision)

	fmt.Printf("toll geohash=%s||box=[%.8f,%.8f;%.8f,%.8f;%.8f,%.8f;%.8f,%.8f]",
		tollhash,
		box.MaxLng, box.MaxLat,
		box.MinLng, box.MaxLat,
		box.MaxLng, box.MinLat,
		box.MinLng, box.MinLat)

	geohash.GetNearGeoHash(mylat, mylon, precision, hashprecision)

	return
}
