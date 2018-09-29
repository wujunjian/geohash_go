package toll

import (
	"encoding/json"
	"fmt"

	"io/ioutil"
	"path"

	"../geohash"
)

const (
	TOLL_CLOSE_TO_POINT   = 100
	TOLL_CLOSE_TO_STATION = 1000

	//TODO 这俩个参数需要可配置
	precision     float64 = 0.001 //精度的上限,格子不能全部超出该范围 0.01(889-1113m)  0.001(88.9-111.3m)
	hashprecision int     = 8     //精度的下限,一个格子代表的面积    6(1223*488m)  7(150*120m) 8(±0.019km)

	before  = -1
	station = 0
	after   = 1
)

type OriTollData struct {
	Tolls []OriToll `json:"tolls"`
}

type OriToll struct {
	Id                      int     `json:"id,string"`
	ExternalId              string  `json:"externalId"`
	Name                    string  `json:"name"`
	TollLatitude            float64 `json:"tollLatitude,string"`
	TollLongitude           float64 `json:"tollLongitude,string"`
	BeforeLocationLatitude  float64 `json:"beforeLocationLatitude,string"`
	BeforeLocationLongitude float64 `json:"beforeLocationLongitude,string"`
	AfterLocationLatitude   float64 `json:"afterLocationLatitude,string"`
	AfterLocationLongitude  float64 `json:"afterLocationLongitude,string"`
	Heading                 string  `json:"heading"`
	Price                   float64 `json:"price,string"`
	LastUpdate              string  `json:"lastUpdate"`
}

type Toll struct {
	Id           int            `json:"id"`
	BeforPoint   *helpers.Point `json:"befor_location"`
	StationPoint *helpers.Point `json:"toll_location"`
	AfterPoint   *helpers.Point `json:"after_location"`
	Price        float64        `json:"price"`
}

type SignedToll struct {
	Toll *Toll

	//before, station, after
	Site int
	Box  *geohash.Box
}

var globalGeoHashMap map[string][]*SignedToll

func init() {
	globalGeoHashMap = make(map[string][]*SignedToll)
	LoadBrTolls("./conf/tolls")
}

func (s *SignedToll) siteString() string {
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

func Debug(tollid, site int) {
	for geohash, signedtolls := range globalGeoHashMap {
		for _, signedtoll := range signedtolls {
			if (tollid == signedtoll.Toll.Id && site == signedtoll.Site) || tollid == 0 {
				lat, lng := signedtoll.point()
				fmt.Printf("geohash=%s||site=%s||tollId=%03d||lng&lat=%.8f,%.8f||box=[%.8f,%.8f;%.8f,%.8f;%.8f,%.8f;%.8f,%.8f]\n",
					geohash, signedtoll.siteString(), signedtoll.Toll.Id, lng, lat,
					signedtoll.Box.MaxLng, signedtoll.Box.MaxLat,
					signedtoll.Box.MinLng, signedtoll.Box.MaxLat,
					signedtoll.Box.MaxLng, signedtoll.Box.MinLat,
					signedtoll.Box.MinLng, signedtoll.Box.MinLat)
			}
		}
	}
}

func GetTolls(latitude, longitude float64) []*SignedToll {
	hash, _ := geohash.Encode(latitude, longitude, hashprecision)
	v, ok := globalGeoHashMap[hash]
	if !ok {
		return nil
	}
	return v
}

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
		fileBytes, err := ioutil.ReadFile(path.Join(tollsFile, f.Name()))
		if err != nil {
			return fmt.Errorf("readfile failed, file=%v err=%v", f.Name(), err)
		}

		data := OriTollData{}
		if err = json.Unmarshal(fileBytes, &data); err != nil {
			return fmt.Errorf("unmarshal file failed, file=%v err=%v", f.Name(), err)
		}

		for _, oriToll := range data.Tolls {
			Toll := oriToll.ToToll()

			Tomap(oriToll.BeforeLocationLatitude, oriToll.BeforeLocationLongitude, Toll, before)

			Tomap(oriToll.TollLatitude, oriToll.TollLongitude, Toll, station)

			Tomap(oriToll.AfterLocationLatitude, oriToll.AfterLocationLongitude, Toll, after)
		}
	}
	return nil
}

func Tomap(latitude, longitude float64, toll *Toll, site int) {

	hashBox := geohash.GetNearGeoHash(latitude, longitude, precision, hashprecision)

	for _, hb := range hashBox {
		signedToll := &SignedToll{
			Toll: toll,
			Site: site,
			Box:  hb.Box,
		}
		globalGeoHashMap[hb.Hash] = append(globalGeoHashMap[hb.Hash], signedToll)
	}
}

func (oriToll *OriToll) ToToll() *Toll {
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
