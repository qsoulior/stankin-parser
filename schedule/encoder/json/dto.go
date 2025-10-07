package json_encoder

const (
	layoutTime = "15:04"
	layoutDate = "2006-01-02"
)

type EventTime struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type EventDate struct {
	Start    string `json:"start"`
	End      string `json:"end"`
	Interval int    `json:"interval"`
}

type Event struct {
	Title    string      `json:"title"`
	Teacher  string      `json:"teacher"`
	Type     string      `json:"type"`
	Subgroup string      `json:"subgroup"`
	Location string      `json:"location"`
	Time     EventTime   `json:"time"`
	Dates    []EventDate `json:"dates"`
}
