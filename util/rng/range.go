package rng

type Range [2]int

func (r Range) Covers(o Range) bool {
	return r[0] <= o[0] && r[1] >= o[1]
}

func (r Range) Intersects(o Range) bool {
	return o.Covers(r) || (r[0] <= o[0] && r[1] >= o[0]) || (r[0] <= o[1] && r[1] >= o[1])
}
