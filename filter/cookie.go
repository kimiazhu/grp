// Author: ZHU HAIHUA
// Date: 8/23/16
package filter

import (
	"net/http"
)

func FilterCookie(cookies []*http.Cookie, reverse bool) []*http.Cookie {
	for _, c := range cookies {
		replaceCookieDomain(c, reverse)
	}
	return cookies
}

func replaceCookieDomain(c *http.Cookie, reverse bool) {
	//if strings.HasPrefix(c.Domain, ".") {
	//if len(c.Domain) > 0 {
	if reverse {
		c.Domain = ".localhost.com"
	} else {
		c.Domain = ".google.com"
	}
	//}
}
