package utils

import (
	"encoding/json"
	"errors"
	"fmt"

	"github.com/go-zepto/zepto/linker/filter"
	"github.com/go-zepto/zepto/web"
	"github.com/joncalhoun/qson"
)

func GetFilterFromQueryArgCtx(ctx web.Context) (*filter.Filter, error) {
	var res = struct {
		Filter filter.Filter `json:"filter"`
	}{}
	query := ctx.Request().URL.Query()
	filterJson := query.Get("filter")
	if filterJson != "" {
		var f filter.Filter
		err := json.Unmarshal([]byte(filterJson), &f)
		if err != nil {
			return nil, errors.New("could not parse json in filter query param")
		}
		return &f, nil
	}
	rawQuery := ctx.Request().URL.RawQuery
	if rawQuery != "" {
		err := qson.Unmarshal(&res, rawQuery)
		if err != nil {
			fmt.Println(err)
			return nil, errors.New("could not parse URI filter")
		}
		return &res.Filter, nil
	}
	return nil, nil
}
