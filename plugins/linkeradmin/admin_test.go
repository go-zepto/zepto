package linkeradmin_test

import (
	"testing"

	"github.com/go-zepto/zepto/plugins/linkeradmin"
	"github.com/go-zepto/zepto/plugins/linkeradmin/fields"
	"github.com/stretchr/testify/assert"
)

func TestAdmin(t *testing.T) {
	admin := linkeradmin.NewAdmin()

	// Post Resource
	post := linkeradmin.NewLinkerResource("Post")

	post.List.
		AddField(fields.NewTextField("id", nil)).
		AddField(fields.NewTextField("title", nil))

	admin.AddResource(post)

	assert.Equal(t, "posts", admin.Resources[0].Endpoint)
	assert.Len(t, admin.Resources, 1)
	assert.Len(t, admin.Resources[0].List.Fields, 2)
}
