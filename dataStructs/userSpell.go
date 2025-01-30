package dataStructs

type UserSpell struct {
	Id          int     `json:"id" form:"id"`
	Name        string  `json:"name" form:"name"`
	Level       *int    `json:"level" form:"level"`
	School      *int    `json:"school" form:"school"`
	IsRitual    *int    `json:"isRitual" form:"isRitual"`
	CastingTime *string `json:"castingTime" form:"castingTime"`
	SpellRange  *string `json:"spellRange" form:"spellRange"`
	Components  *string `json:"components" form:"components"`
	Duration    *string `json:"duration" form:"duration"`
	Description *string `json:"description" form:"description"`
	Upcast      *string `json:"upcast" form:"upcast"`
	User_id     int
	IsPublic    *int `json:"isPublic" form:"isPublic"`
}
