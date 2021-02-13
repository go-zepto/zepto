package where

import (
	"errors"
	"fmt"
)

type WhereOperator struct {
	SQL string
}

type WhereType struct {
	Key      string
	Operator WhereOperator
}

func (wt *WhereType) ApplySQL(vars ...interface{}) (string, error) {
	if wt.Operator.SQL == "" {
		return "", errors.New("sql not supported for operator")
	}
	return fmt.Sprintf(wt.Operator.SQL, vars...), nil
}

var TYPES = map[string]WhereType{
	"__root__": {
		Key: "__root__",
	},
	"__field__": {
		Key: "__field__",
	},
	"and": {
		Key: "and",
		Operator: WhereOperator{
			SQL: "AND",
		},
	},
	"or": {
		Key: "or",
		Operator: WhereOperator{
			SQL: "OR",
		},
	},
	"eq": {
		Key: "eq",
		Operator: WhereOperator{
			SQL: "%s = ?",
		},
	},
	"neq": {
		Key: "neq",
		Operator: WhereOperator{
			SQL: "%s <> ?",
		},
	},
	"gt": {
		Key: "gt",
		Operator: WhereOperator{
			SQL: "%s > ?",
		},
	},
	"gte": {
		Key: "gte",
		Operator: WhereOperator{
			SQL: "%s >= ?",
		},
	},
	"lt": {
		Key: "lt",
		Operator: WhereOperator{
			SQL: "%s < ?",
		},
	},
	"lte": {
		Key: "lte",
		Operator: WhereOperator{
			SQL: "%s <= ?",
		},
	},
	"between": {
		Key: "between",
		Operator: WhereOperator{
			SQL: "%s BETWEEN ? AND ?",
		},
	},
}
