package web

import (
	"net/http"
	"time"
)

type Cookies struct {
	req *http.Request
	res http.ResponseWriter
}

// Get a cookie from name
func (c *Cookies) Get(name string) (string, error) {
	cookie, err := c.req.Cookie(name)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// Set a cookie on the response, which will expire after the given duration.
func (c *Cookies) Set(name, value string, maxAge time.Duration) {
	ma := int(maxAge.Seconds())
	cookie := http.Cookie{
		Name:   name,
		Value:  value,
		MaxAge: ma,
	}

	http.SetCookie(c.res, &cookie)
}

// SetWithExpirationTime sets a cookie that will expire at a specific time.
func (c *Cookies) SetWithExpirationTime(name, value string, expires time.Time) {
	ck := http.Cookie{
		Name:    name,
		Value:   value,
		Expires: expires,
	}

	http.SetCookie(c.res, &ck)
}

// SetWithPath sets a cookie path on the server in which the cookie will be available on.
func (c *Cookies) SetWithPath(name, value, path string) {
	cookie := http.Cookie{
		Name:  name,
		Value: value,
		Path:  path,
	}
	http.SetCookie(c.res, &cookie)
}

// Delete a cookie
func (c *Cookies) Delete(name string) {
	cookie := http.Cookie{
		Name:    name,
		Value:   "-",
		Expires: time.Unix(0, 0),
	}

	http.SetCookie(c.res, &cookie)
}
