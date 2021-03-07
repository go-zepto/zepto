package where

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEqualSQLOperator(t *testing.T) {
	eq, exists := TYPES["eq"]
	assert.True(t, exists)
	assert.Equal(t, eq.Key, "eq")
	assert.NotNil(t, eq.Operator.SQL)
	assert.Equal(t, eq.Operator.SQL, "%s = ?")
}

func TestEqualSQLOperatorApply(t *testing.T) {
	eq := TYPES["eq"]
	res, err := eq.ApplySQL("name")
	assert.NoError(t, err)
	assert.Equal(t, res, "name = ?")
}
