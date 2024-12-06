package dataStructs

type QueryData struct {
	Id              int
	Source          int
	Level           int
	School          int
	IsRitual        bool
	CastingTime     string // Predefined
	RangeValueStart int
	RangeValueStop  int
	RangeType       string // Predefined OR feet OR mile
	Components      string // For reconsideration about selection REGEX FUCK ME
	Duration        string // Predefined
	Upcast          bool
}

type QueryDataStrings struct {
	Id              string
	Source          string
	Level           string
	School          string
	IsRitual        string
	CastingTime     string
	RangeValueStart string
	RangeValueStop  string
	RangeType       string
	Components      string
	Duration        string
	Upcast          string
}
