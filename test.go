package main

import (
	"fmt"
	"go.intra.xiaojukeji.com/gulfstream/plutus/thirdparty/neweta/base"
	"io"
	"io/ioutil"
	"sync/atomic"
	"time"
)
import "net/http"

var tr *http.Transport = &http.Transport{
	MaxIdleConns:       100,
	//IdleConnTimeout:    30 * time.Second,
	DisableCompression: true,
}
func testclient(ms time.Duration) {
	client := &http.Client{
		Transport: tr,
		Timeout:   time.Millisecond * ms,
	}

	//resp, err := client.Get("http://map-mirror-routeplan00.gz01:9990")
	resp, err := client.Get("http://10.179.116.156:9990/hello")


	if err!= nil {
		fmt.Println("err!", err.Error())
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		io.Copy(ioutil.Discard,resp.Body)
		resp.Body.Close()

		fmt.Println(resp.Status, body)
	}
}


func atomictest () {
	var countSocketNum  uint64

	atomic.AddUint64(&countSocketNum, 1)
	atomic.AddUint64(&countSocketNum, 1)
	fmt.Println(countSocketNum)

	atomic.AddUint64(&countSocketNum, ^uint64(0))
	fmt.Println(countSocketNum)


	atomic.AddUint64(&countSocketNum, 10)
	fmt.Println(countSocketNum)


	atomic.AddUint64(&countSocketNum, ^uint64(10-1))
	fmt.Println(countSocketNum)


	fmt.Println(^uint64(10-1))

}

func mclient() {
	for {

		testclient(1500)
		testclient(990)
		fmt.Println(time.Now())

		time.Sleep(2*time.Second)

	}
}


func main() {

	fmt.Println("begin")

	DecodePolyline("f}zeF{wxsZs@{D]wAaAiEc@mBMm@e@aCQw@CSEQEUEOCKAIMu@MwAEQCQ?GEc@I_AA_@KwA?IQeBA[EiAGeBGeBAaAC_AAgAAeAAS?UAm@?o@BiBNkWLaHH_EHwE`Ae[NaDf@oEb@_Df@sCr@}B~@iCzAwD|AaD|@kBpAsC\\u@~AkD`@oAh@kBt@gDHa@d@kCP_ARgAVoBD[@C@M^eCd@}BpBoLLi@^kC@EHu@H_A")

	c := make(chan int)
	<-c
}


// DecodePolyline 解码polyline为geo序列
func DecodePolyline(poly string) []*base.GeoPoint {
	geoList := make([]*base.GeoPoint, 0, 20)

	var nowLat, nowLon int64
	var shift, result uint64
	var latOrLon = true

	for i := 0; i < len(poly); i++ {
		value := uint64(poly[i]) - 63
		isNotLast := value & 0x20
		value &= 0x1F
		result |= (value << shift)
		shift += 5

		if isNotLast == 0 {
			if result&1 != 0 {
				result = ^(result >> 1)
			} else {
				result = result >> 1
			}

			if latOrLon {
				nowLat += int64(result)
			} else {
				nowLon += int64(result)

				geo := base.NewGeoPoint()
				geo.Latitude = float64(nowLat) / 100000.0
				geo.Longitude = float64(nowLon) / 100000.0
				geoList = append(geoList, geo)
			}
			latOrLon = !latOrLon
			result = 0
			shift = 0
		}
	}
	return geoList

}
