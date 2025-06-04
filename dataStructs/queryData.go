package dataStructs

type QueryData struct {
	Id          int
	Source      int
	Level       int
	School      int
	IsRitual    string
	CastingTime string // Predefined
	RVSF        int
	RVEF        int
	RVSM        int
	RVEM        int
	RS          string
	Components  string // For reconsideration about selection REGEX FUCK ME
	Duration    string // Predefined
	Upcast      bool
}

type QueryDataStrings struct {
	Id          string
	Source      string
	Level       string
	School      string
	IsRitual    string
	CastingTime string
	RVSF        int
	RVEF        int
	RVSM        int
	RVEM        int
	Specials    [5]string
	Components  string
	Duration    string
	Upcast      string
}
