package rating

import "math"

// Wilson returns the lower bound of the Wilson score interval. Callers pass
// only recommend and not-recommend counts, so neutral ratings never reach it.
func Wilson(up, down int64, z float64) float64 {
	n := up + down
	if n <= 0 {
		return 0
	}

	nf := float64(n)
	p := float64(up) / nf
	z2 := z * z
	score := (p + z2/(2*nf) - z*math.Sqrt((p*(1-p)+z2/(4*nf))/nf)) / (1 + z2/nf)

	return math.Min(1, math.Max(0, score))
}
