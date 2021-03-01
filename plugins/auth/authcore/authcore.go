package authcore

type AuthCore struct {
	DS           AuthDatasource
	TokenEncoder AuthEncoder
	Store        AuthStore
	Notifier     AuthNotifier
}

func (ac *AuthCore) AssertConfigured() {
	if ac.DS == nil {
		panic("[auth] you must define a datasource")
	}
	if ac.Store == nil {
		panic("[auth] you must define a store")
	}
	if ac.Notifier == nil {
		panic("[auth] you must define a notifier")
	}
}
