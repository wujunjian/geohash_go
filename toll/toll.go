package toll

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"io/ioutil"
	"path"

	"../geohash"
	"go.intra.xiaojukeji.com/gulfstream/plutus/common/helpers"
)

const (
	//TOLL_CLOSE_TO_POINT   = 100
	//TOLL_CLOSE_TO_STATION = 1000

	//需要可配置的参数
	precision          float64 = 0.0012 //精度的上限,格子不能全部超出该范围 0.01(889-1113m)  0.001(88.9-111.3m)
	hashprecision      int     = 8      //精度的下限,一个格子代表的面积    6(1223*488m)  7(150*120m) 8(±0.019km)
	reloadtimeinterval         = 1800   //配置文件reload的时间间隔,单位秒

	//不需要配置的参数
	before  = -1
	station = 0
	after   = 1
)

type matched struct {
	before  int
	station int
	after   int
}

//TollsMatched key:Toll id
type TollsMatched struct {
	toll map[int]*matched
}

//DebugPrint message
func (t *TollsMatched) DebugPrint() {
	for id, m := range t.toll {
		fmt.Printf("id=%d||beforematched=%d||stationmatched=%d||aftermatched=%d\n", id, m.before, m.station, m.after)
	}
}

//Match all posible tolls
func (t *TollsMatched) Match(points []*helpers.Point, country string) {
	t.toll = make(map[int]*matched)
	for _, p := range points {
		signedtolls := GetTolls(p.Lat, p.Lng, country)
		if signedtolls != nil {
			for _, signedtoll := range signedtolls {

				id := signedtoll.Toll.Id
				m, ok := t.toll[id]
				if !ok {
					m = &matched{}
					t.toll[id] = m
				}
				switch signedtoll.Site {
				case before:
					m.before++
				case station:
					m.station++
				case after:
					m.after++
				}
			}
		}
	}
}

type oriTollData struct {
	Tolls []oriToll `json:"tolls"`
}

type tPoint struct {
	Lat float64 `json:"lat,string"`
	Lng float64 `json:"lng,string"`
}

type oriToll struct {
	Id                      int      `json:"id,string"`
	ExternalId              string   `json:"externalId"`
	Name                    string   `json:"name"`
	TollLatitude            float64  `json:"tollLatitude,string"`
	TollLongitude           float64  `json:"tollLongitude,string"`
	ExtendStation           []tPoint `json:"extendStation,[]interface{}"`
	DelStation              []tPoint `json:"delStation,[]interface{}"`
	BeforeLocationLatitude  float64  `json:"beforeLocationLatitude,string"`
	BeforeLocationLongitude float64  `json:"beforeLocationLongitude,string"`
	ExtendBefore            []tPoint `json:"extendBefore,[]interface{}"`
	DelBefore               []tPoint `json:"delBefore,[]interface{}"`
	AfterLocationLatitude   float64  `json:"afterLocationLatitude,string"`
	AfterLocationLongitude  float64  `json:"afterLocationLongitude,string"`
	ExtendAfter             []tPoint `json:"extendAfter,[]interface{}"`
	DelAfter                []tPoint `json:"delAfter,[]interface{}"`
	Heading                 string   `json:"heading"`
	Price                   float64  `json:"price,string"`
	LastUpdate              string   `json:"lastUpdate"`
}

// Toll info
type Toll struct {
	Id           int            `json:"id"`
	BeforPoint   *helpers.Point `json:"befor_location"`
	StationPoint *helpers.Point `json:"toll_location"`
	AfterPoint   *helpers.Point `json:"after_location"`
	Price        float64        `json:"price"`
}

// SignedToll 命中的收费站
type SignedToll struct {
	Toll *Toll

	//before, station, after
	Site int
	Box  *geohash.Box
}
type geoHashMap map[string][]*SignedToll

var globalGeolock *sync.RWMutex
var countryMap map[string]geoHashMap

func init() {
	globalGeolock = &sync.RWMutex{}

	updateGlobalGeoHashMap := func() {
		for {
			globalGeolock.Lock()
			countryMap = make(map[string]geoHashMap)
			LoadBrTolls("./conf/tolls")
			globalGeolock.Unlock()

			time.Sleep(reloadtimeinterval * time.Second)
		}

	}
	go updateGlobalGeoHashMap()
}

func (s *SignedToll) SiteString() string {
	switch s.Site {
	case before:
		return "b"
	case station:
		return "s"
	case after:
		return "a"
	}
	return ""
}

func (s *SignedToll) point() (float64, float64) {
	switch s.Site {
	case before:
		return s.Toll.BeforPoint.Lat, s.Toll.BeforPoint.Lng
	case station:
		return s.Toll.StationPoint.Lat, s.Toll.StationPoint.Lng
	case after:
		return s.Toll.AfterPoint.Lat, s.Toll.AfterPoint.Lng
	}
	return 0, 0
}

//Debug Print globalGeoHashMap
func Debug(tollid, site int, where string) {
	globalGeolock.RLock()
	defer globalGeolock.RUnlock()

	for country, geohashmap := range countryMap {
		for geohash, signedtolls := range geohashmap {
			for _, signedtoll := range signedtolls {
				if (tollid == signedtoll.Toll.Id && site == signedtoll.Site && country == where) || tollid == 0 {
					lat, lng := signedtoll.point()
					fmt.Printf("country=%s||geohash=%s||site=%s||tollId=%03d||lng&lat=%.8f,%.8f||box=[%.8f,%.8f;%.8f,%.8f;%.8f,%.8f;%.8f,%.8f]\n",
						country, geohash, signedtoll.SiteString(), signedtoll.Toll.Id, lng, lat,
						signedtoll.Box.MaxLng, signedtoll.Box.MaxLat,
						signedtoll.Box.MinLng, signedtoll.Box.MaxLat,
						signedtoll.Box.MaxLng, signedtoll.Box.MinLat,
						signedtoll.Box.MinLng, signedtoll.Box.MinLat)
				}
			}
		}
	}

}

//GetTolls Get hit Tolls
func GetTolls(latitude, longitude float64, country string) []*SignedToll {
	hash, _ := geohash.Encode(latitude, longitude, hashprecision)

	globalGeolock.RLock()
	defer globalGeolock.RUnlock()

	globalGeoHashMap, ok := countryMap[country]
	if !ok {
		return nil
	}

	v, ok := globalGeoHashMap[hash]
	if !ok {
		return nil
	}
	return v
}

// LoadBrTolls read files and loadBrTolls
func LoadBrTolls(tollsFile string) error {
	if tollsFile == "" {
		return nil
	}
	fileList, err := ioutil.ReadDir(tollsFile)
	if err != nil {
		return err
	}
	for _, f := range fileList {
		if f.IsDir() {
			continue
		}

		filename := f.Name()
		country := strings.Split(f.Name(), ".")[0]
		fileBytes, err := ioutil.ReadFile(path.Join(tollsFile, filename))
		if err != nil {
			return fmt.Errorf("readfile failed, file=%v err=%v", filename, err)
		}

		data := oriTollData{}
		if err = json.Unmarshal(fileBytes, &data); err != nil {
			return fmt.Errorf("unmarshal file failed, file=%v err=%v", filename, err)
		}

		for _, oriToll := range data.Tolls {
			Toll := oriToll.ToToll()

			var tp tPoint
			for _, tp = range oriToll.ExtendBefore {
				narrowtomap(tp.Lat, tp.Lng, Toll, before, country)
			}
			tomap(oriToll.BeforeLocationLatitude, oriToll.BeforeLocationLongitude, Toll, before, country)
			for _, tp = range oriToll.DelBefore {
				narrowdelmap(tp.Lat, tp.Lng, Toll, before, country)
			}

			for _, tp = range oriToll.ExtendStation {
				narrowtomap(tp.Lat, tp.Lng, Toll, station, country)
			}
			tomap(oriToll.TollLatitude, oriToll.TollLongitude, Toll, station, country)
			for _, tp = range oriToll.DelStation {
				narrowdelmap(tp.Lat, tp.Lng, Toll, station, country)
			}

			for _, tp = range oriToll.ExtendAfter {
				narrowtomap(tp.Lat, tp.Lng, Toll, after, country)
			}
			tomap(oriToll.AfterLocationLatitude, oriToll.AfterLocationLongitude, Toll, after, country)
			for _, tp = range oriToll.DelAfter {
				narrowdelmap(tp.Lat, tp.Lng, Toll, after, country)
			}
		}
	}
	return nil
}

func narrowdelmap(latitude, longitude float64, toll *Toll, site int, country string) {
	fmt.Println("delete:", latitude, longitude, site, country)
	h, _ := geohash.Encode(latitude, longitude, hashprecision)
	hashmap, ok := countryMap[country]
	if !ok {
		return
	}

	signedTolls, ok := hashmap[h]
	if !ok {
		return
	}
	var newSignedToll []*SignedToll
	for _, st := range signedTolls {
		if st.Toll.Id == toll.Id && st.Site == site {
			continue
		}
		fmt.Println(st.Toll.Id, toll.Id, st.Site, site)
		newSignedToll = append(newSignedToll, st)
	}
	hashmap[h] = newSignedToll
}

func narrowtomap(latitude, longitude float64, toll *Toll, site int, country string) {
	h, box := geohash.Encode(latitude, longitude, hashprecision)
	hashmap, ok := countryMap[country]
	if !ok {
		countryMap[country] = make(geoHashMap)
		hashmap, _ = countryMap[country]
	}
	signedToll := &SignedToll{
		Toll: toll,
		Site: site,
		Box:  box,
	}
	hashmap[h] = append(hashmap[h], signedToll)
}

func tomap(latitude, longitude float64, toll *Toll, site int, country string) {

	hashBox := geohash.GetNearGeoHash(latitude, longitude, precision, hashprecision)
	hashmap, ok := countryMap[country]
	if !ok {
		countryMap[country] = make(geoHashMap)
		hashmap, _ = countryMap[country]
	}

	for _, hb := range hashBox {
		signedToll := &SignedToll{
			Toll: toll,
			Site: site,
			Box:  hb.Box,
		}
		hashmap[hb.Hash] = append(hashmap[hb.Hash], signedToll)
	}
}

func (oriToll *oriToll) ToToll() *Toll {
	res := &Toll{
		Id:    oriToll.Id,
		Price: oriToll.Price,
		BeforPoint: &helpers.Point{
			Lat: oriToll.BeforeLocationLatitude,
			Lng: oriToll.BeforeLocationLongitude,
		},
		StationPoint: &helpers.Point{
			Lat: oriToll.TollLatitude,
			Lng: oriToll.TollLongitude,
		},
		AfterPoint: &helpers.Point{
			Lat: oriToll.AfterLocationLatitude,
			Lng: oriToll.AfterLocationLongitude,
		},
	}
	return res
}
