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
	query, err := w.ToSQL()
	assert.NoError(t, err)
	assert.Len(t, query.Vars, 2)
	assert.Equal(t, query.Vars[0].(string), "Carlos")
	assert.Equal(t, query.Vars[1].(string), "Strand")
	assert.Equal(t, "first_name = ? AND last_name <> ?", query.Text)
}

func TestWhereWithAllowedFields_1(t *testing.T) {
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
	w := NewFromMapWithAllowedFields(jsonToMap(filterJson), []string{"first_name"})
	_, err := w.ToSQL()
	assert.EqualError(t, err, "last_name field is not allowed")
}

func TestWhereWithAllowedFields_2(t *testing.T) {
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
	w := NewFromMapWithAllowedFields(jsonToMap(filterJson), []string{"first_name", "last_name"})
	query, err := w.ToSQL()
	assert.NoError(t, err)
	assert.Len(t, query.Vars, 2)
	assert.Equal(t, query.Vars[0].(string), "Carlos")
	assert.Equal(t, query.Vars[1].(string), "Strand")
	assert.Equal(t, "first_name = ? AND last_name <> ?", query.Text)
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
	query, err := w.ToSQL()
	assert.NoError(t, err)
	assert.Len(t, query.Vars, 2)
	assert.Equal(t, query.Vars[0].(string), "Carlos")
	assert.Equal(t, query.Vars[1].(string), "Strand")
	assert.Equal(t, "(first_name = ? OR last_name <> ?)", query.Text)
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
	query, err := w.ToSQL()
	assert.NoError(t, err)
	assert.Len(t, query.Vars, 2)
	assert.Equal(t, query.Vars[0].(string), "Carlos")
	assert.Equal(t, query.Vars[1].(string), "Strand")
	assert.Equal(t, "(first_name = ? AND last_name <> ?)", query.Text)
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
	query, err := w.ToSQL()
	assert.NoError(t, err)
	assert.Len(t, query.Vars, 4)
	assert.Equal(t, query.Vars[0].(string), "Carlos")
	assert.Equal(t, query.Vars[1].(string), "Strand")
	assert.Equal(t, query.Vars[2].(string), "Bill")
	assert.Equal(t, query.Vars[3].(string), "Gates")
	assert.Equal(t, "((first_name = ? AND last_name = ?) OR (first_name = ? AND last_name = ?))", query.Text)
}

func TestGetWhereType(t *testing.T) {
	filterJson := `
	{
		"name": {
			"or": [
				{
					"test": "or can't be inside field"
				}
			]
		}
	}
`
	w := NewFromMap(jsonToMap(filterJson))
	assert.Equal(t, "__field__", w.GetWhereType("name").Key)
}

func TestWhereWithInvalidFormat(t *testing.T) {
	filterJson := `
		{
			"name": {
				"or": [
					{
						"test": "or can't be inside field"
					}
				]
			}
		}
	`
	w := NewFromMap(jsonToMap(filterJson))
	_, err := w.ToSQL()
	assert.EqualError(t, err, "OR operator in unsupported parent")
	filterJson = `
		{
			"name": {
				"and": [
					{
						"test": "and can't be inside field"
					}
				]
			}
		}
	`
	w = NewFromMap(jsonToMap(filterJson))
	_, err = w.ToSQL()
	assert.EqualError(t, err, "AND operator in unsupported parent")
}
