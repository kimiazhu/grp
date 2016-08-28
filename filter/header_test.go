// Author: ZHU HAIHUA
// Date: 8/28/16
package filter

import "testing"

func TestTopDomainName(t *testing.T) {
	t.Log(TopDomainName("g.localhost.com"))
	t.Log(TopDomainName(".localhost.com"))
	t.Log(TopDomainName(".g.localhost.com"))
	t.Log(TopDomainName(".com"))
	t.Log(TopDomainName("gg.localhost.com:8888"))
}

func TestDealCookie(t *testing.T) {
	cookie := `NID=85=Sa5qdO8_-EADQ54Wz-DMgJLQ; expires=Mon, 27-Feb-2017 05:16:04 GMT; path=/; domain=.localhost.com; HttpOnly`
	t.Log(DealCookie("www.google.com", cookie, ))
}