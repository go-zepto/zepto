package web

import (
	"github.com/gorilla/sessions"
	"net/http"
)

type Session struct {
	gSession *sessions.Session
	req      *http.Request
	res      http.ResponseWriter
}

func (s *Session) Save() error {
	return s.gSession.Save(s.req, s.res)
}

func (s *Session) Get(key interface{}) interface{} {
	return s.gSession.Values[key]
}

func (s *Session) Set(key, value interface{}) {
	s.gSession.Values[key] = value
}

func (s *Session) Delete(key interface{}) {
	delete(s.gSession.Values, key)
}

func (s *Session) Clear() {
	for k := range s.gSession.Values {
		s.Delete(k)
	}
}
