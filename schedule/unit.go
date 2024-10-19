package schedule

// Meta represents schedule metadata encoded in input data.
// It typically contains group name.
type Meta struct {
	Group string
}

// Unit represents schedule unit encoded in input data.
type Unit struct {
	Data                     string
	Left, Right, Top, Bottom int
}
