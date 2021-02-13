package where

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func jsonToMap(jsonStr string) map[string]interface{} {
	var filterMap map[string]interface{}
	err := json.Unmarshal([]byte(jsonStr), &filterMap)
	if err != nil {
		panic(err)
	}
	return filterMap
}

func TestWhere(t *testing.T) {
	filterJson := `
		{
			"first_name": {
				"eq": "Carlos"
			},
			"last_name": {
				"neq": "Strand"
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	query := w.ToSQL()
	assert.Len(t, query.Vars, 2)
	assert.Equal(t, query.Vars[0].(string), "Carlos")
	assert.Equal(t, query.Vars[1].(string), "Strand")
	assert.Equal(t, query.Text, "first_name = ? AND last_name <> ?")
}

func TestWhereOR(t *testing.T) {
	filterJson := `
		{
			"or": [
				{
					"first_name": {
						"eq": "Carlos"
					}
				},
				{
					"last_name": {
						"neq": "Strand"
					}
				}
			]
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	query := w.ToSQL()
	assert.Len(t, query.Vars, 2)
	assert.Equal(t, query.Vars[0].(string), "Carlos")
	assert.Equal(t, query.Vars[1].(string), "Strand")
	assert.Equal(t, query.Text, "(first_name = ? OR last_name <> ?)")
}

func TestWhereAND(t *testing.T) {
	filterJson := `
		{
			"and": [
				{
					"first_name": {
						"eq": "Carlos"
					}
				},
				{
					"last_name": {
						"neq": "Strand"
					}
				}
			]
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	query := w.ToSQL()
	assert.Len(t, query.Vars, 2)
	assert.Equal(t, query.Vars[0].(string), "Carlos")
	assert.Equal(t, query.Vars[1].(string), "Strand")
	assert.Equal(t, query.Text, "(first_name = ? AND last_name <> ?)")
}

func TestWhereOR_with_AND(t *testing.T) {
	filterJson := `
		{
			"or": [
				{
					"and": [
						{
							"first_name": {
								"eq": "Carlos"
							}
						},
						{
							"last_name": {
								"eq": "Strand"
							}
						}
					]
				},
				{
					"and": [
						{
							"first_name": {
								"eq": "Bill"
							}
						},
						{
							"last_name": {
								"eq": "Gates"
							}
						}
					]
				}
			]
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	query := w.ToSQL()
	assert.Len(t, query.Vars, 4)
	assert.Equal(t, query.Vars[0].(string), "Carlos")
	assert.Equal(t, query.Vars[1].(string), "Strand")
	assert.Equal(t, query.Vars[2].(string), "Bill")
	assert.Equal(t, query.Vars[3].(string), "Gates")
	assert.Equal(t, query.Text, "((first_name = ? AND last_name = ?) OR (first_name = ? AND last_name = ?))")
}
