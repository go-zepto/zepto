package filter

type Filter struct {
	Skip  *int                    `json:"skip"`
	Limit *int                    `json:"limit"`
	Where *map[string]interface{} `json:"where"`
}
