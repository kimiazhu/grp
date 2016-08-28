// Author: ZHU HAIHUA
// Date: 8/28/16
package ioutils

import "testing"

func TestTopDomainName(t *testing.T) {
	t.Log(TopDomainName("g.localhost.com"))
	t.Log(TopDomainName(".localhost.com"))
	t.Log(TopDomainName(".g.localhost.com"))
	t.Log(TopDomainName(".com"))
	t.Log(TopDomainName("gg.localhost.com:8888"))
}
