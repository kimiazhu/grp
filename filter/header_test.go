// Author: ZHU HAIHUA
// Date: 8/28/16
package filter

import "testing"

func TestDealCookie(t *testing.T) {
	cookie := `NID=85=Sa5qdO8_-EADQ54Wz-DMgJLQ; expires=Mon, 27-Feb-2017 05:16:04 GMT; path=/; domain=.localhost.com; HttpOnly`
	t.Log(DealCookie("www.google.com", cookie, true))
}
