package schedule

type Meta struct {
	Group string
}

type Cell struct {
	Data                     string
	Left, Right, Top, Bottom int
}
