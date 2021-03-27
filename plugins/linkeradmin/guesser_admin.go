package linkeradmin

import (
	"github.com/go-zepto/zepto/plugins/linker"
)

func NewGuesserAdmin(l *linker.Linker) *Admin {
	admin := NewAdmin()
	for name, repo := range l.Repositories() {
		res := NewLinkerResource(name)
		for _, f := range repo.Fields() {
			res.List.AddField(Field{
				Name:    f.Name,
				Type:    "text",
				Options: make(map[string]interface{}),
			})
			res.Create.AddField(Field{
				Name:    f.Name,
				Type:    "text",
				Options: make(map[string]interface{}),
			})
		}
		admin.AddResource(res)
	}
	admin.Resources[0].
	return admin
}
