package filter

type Filter struct {
	Skip  *int64                  `json:"skip"`
	Limit *int64                  `json:"limit"`
	Where *map[string]interface{} `json:"where"`
}
