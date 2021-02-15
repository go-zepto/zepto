package filter

import "github.com/go-zepto/zepto/linker/filter/include"

type Filter struct {
	Skip    *int64                  `json:"skip"`
	Limit   *int64                  `json:"limit"`
	Where   *map[string]interface{} `json:"where"`
	Include []include.IncludeItem   `json:"include"`
}
