package geohash

import (
	"bytes"
	"fmt"
	"math"
)

const (
	BASE32                = "0123456789bcdefghjkmnpqrstuvwxyz"
	MAX_LATITUDE  float64 = 90
	MIN_LATITUDE  float64 = -90
	MAX_LONGITUDE float64 = 180
	MIN_LONGITUDE float64 = -180
)

var (
	bits   = []int{16, 8, 4, 2, 1}
	base32 = []byte(BASE32)
)

type Box struct {
	MinLat, MaxLat float64 // 纬度
	MinLng, MaxLng float64 // 经度
}

type HashBox struct {
	Box  *Box
	Hash string
}

func (this *Box) Width() float64 {
	return this.MaxLng - this.MinLng
}

func (this *Box) Height() float64 {
	return this.MaxLat - this.MinLat
}

func Debug(latitude, longitude float64, hashprecision int) {
	hash, box := Encode(latitude, longitude, hashprecision)
	fmt.Printf("hash=%s||lng&lat=%.8f,%.8f||box=[%.8f,%.8f;%.8f,%.8f;%.8f,%.8f;%.8f,%.8f]\n",
		hash, longitude, latitude,
		box.MaxLng, box.MaxLat,
		box.MinLng, box.MaxLat,
		box.MaxLng, box.MinLat,
		box.MinLng, box.MinLat)
}

// 输入值: 纬度, 经度, 精度(经纬度,最终取max(width, precision)), 精度(geohash的长度)
// 返回geohash, 以及该点精度范围内的geohash
func GetNearGeoHash(latitude, longitude, precision float64, hashprecision int) []*HashBox {
	var hashboxs []*HashBox

	_, box := Encode(latitude, longitude, hashprecision)
	width := box.Width()
	height := box.Height()

	precision = math.Max(precision, width)
	for i := latitude - precision; i < latitude+precision+height; i += height {
		for j := longitude - precision; j < longitude+precision+width; j += width {
			TmpSgeohash, b := Encode(i, j, hashprecision)
			hashboxs = append(hashboxs, &HashBox{Box: b, Hash: TmpSgeohash})
		}
	}
	return hashboxs
}

// 输入值：纬度，经度，精度(geohash的长度)
// 返回geohash, 以及该点所在的区域
func Encode(latitude, longitude float64, precision int) (string, *Box) {
	var geohash bytes.Buffer
	var minLat, maxLat float64 = MIN_LATITUDE, MAX_LATITUDE
	var minLng, maxLng float64 = MIN_LONGITUDE, MAX_LONGITUDE
	var mid float64 = 0

	bit, ch, length, isEven := 0, 0, 0, true
	for length < precision {
		if isEven {
			if mid = (minLng + maxLng) / 2; mid < longitude {
				ch |= bits[bit]
				minLng = mid
			} else {
				maxLng = mid
			}
		} else {
			if mid = (minLat + maxLat) / 2; mid < latitude {
				ch |= bits[bit]
				minLat = mid
			} else {
				maxLat = mid
			}
		}

		isEven = !isEven
		if bit < 4 {
			bit++
		} else {
			geohash.WriteByte(base32[ch])
			length, bit, ch = length+1, 0, 0
		}
	}

	b := &Box{
		MinLat: minLat,
		MaxLat: maxLat,
		MinLng: minLng,
		MaxLng: maxLng,
	}

	return geohash.String(), b
}
