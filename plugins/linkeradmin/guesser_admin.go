package linkeradmin

import (
	"github.com/go-zepto/zepto/plugins/linker"
)

func NewGuesserAdmin(l *linker.Linker) *Admin {
	guesser := NewDefaultGuesser(DefaultGuesserOptions{
		Linker: l,
	})
	resources := guesser.Resources()
	admin := NewAdmin()
	for _, res := range resources {
		listFields := guesser.ListFields(res.Name)
		for _, f := range listFields {
			res.List.AddField(*f)
		}
		createInputs := guesser.CreateInputs(res.Name)
		for _, i := range createInputs {
			res.Create.AddInput(*i)
		}
		updateInputs := guesser.UpdateInputs(res.Name)
		for _, i := range updateInputs {
			res.Update.AddInput(*i)
		}
		admin.AddResource(res)
	}
	return admin
}
