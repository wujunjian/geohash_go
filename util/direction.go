package util

import "fmt"

type point struct {
	lat float64
	lng float64
}

//顺时针 >=0
//OF->OL
func direction(O, F, L *point) float64 {
	return (F.lng-O.lng)*(L.lat-O.lat) - (L.lng-O.lng)*(F.lat-O.lat)
}

//A-B with C-D intersect
func intersect(A, B, C, D *point) bool {
	return (direction(A, C, D)*direction(B, C, D) <= 0) && (direction(C, A, B)*direction(D, A, B) <= 0)
}

// if p in r
func inrail(p *point, bigp *point, r []*point) bool {

	tbigp := &point{p.lat, bigp.lng + 0.01}

	var pr *point
	for idx, tr := range r {
		if idx == 0 {
			pr = tr
			continue
		}
		if intersect(p, tbigp, pr, tr) {
			fmt.Println(p.lng, p.lat, tbigp.lng, tbigp.lat, pr.lng, pr.lat, tr.lng, tr.lat)
			return true
		}
		pr = tr
	}
	return false
}
