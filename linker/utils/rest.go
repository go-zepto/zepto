package utils

import (
	"github.com/go-zepto/zepto/linker/filter"
	"github.com/go-zepto/zepto/web"
	"github.com/joncalhoun/qson"
)

func GetFilterFromQueryArgCtx(ctx web.Context) *filter.Filter {
	var res = struct {
		Filter filter.Filter `json:"filter"`
	}{}
	qson.Unmarshal(&res, ctx.Request().URL.Query().Encode())
	return &res.Filter
}
