package utils

import (
	"io"
	"net/http"
)

type MailerWriter struct {
	io.Writer
	value []byte
}

func (mw *MailerWriter) Header() http.Header {
	return http.Header{}
}

func (mw *MailerWriter) WriteHeader(int) {}

func (mw *MailerWriter) Write(value []byte) (int, error) {
	mw.value = value
	return len(value), nil
}

func (mw *MailerWriter) Value() []byte {
	return mw.value
}
