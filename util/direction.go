package util

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
