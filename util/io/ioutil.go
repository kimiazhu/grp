// Author: ZHU HAIHUA
// Date: 8/23/16
package ioutils

import (
	"bufio"
	"fmt"
	"github.com/kimiazhu/grp/model"
	"net/http"
	"strings"
)

// 替换请求中的Host, 在转发请求的时候, isRequest应该设置为true, 在转发应答的时候, reverse应该设置为false
func ReplaceHost(data string, isRequest bool) string {
	for _local, _remote := range model.Proxies {
		if !isRequest {
			data = strings.Replace(data, model.SvrCnf[_remote].String(), model.SvrCnf[_local].String(), -1)
		} else {
			data = strings.Replace(data, model.SvrCnf[_local].String(), model.SvrCnf[_remote].String(), -1)
		}
	}
	return data
}

func ParseRawCookies(cookie string) ([]*http.Cookie, error) {
	req, err := http.ReadRequest(bufio.NewReader(strings.NewReader(fmt.Sprintf("GET / HTTP/1.0\r\nCookie: %s\r\n\r\n", cookie))))
	if err != nil {
		return nil, err
	}
	cookies := req.Cookies()
	if len(cookies) == 0 {
		return nil, fmt.Errorf("no cookies")
	}
	return cookies, nil
}

func IsTopDomain(domain string) bool {
	if strings.Count(domain, ".") == 1 || (strings.HasPrefix(domain, ".") && strings.Count(domain, ".") == 2) {
		return true
	}
	return false
}

func TopDomainName(domain string) string {
	if strings.HasPrefix(domain, ".") {
		domain = domain[1:]
	}

	if ci := strings.LastIndex(domain, ":"); ci > 0 {
		domain = domain[:ci]
	}

	if strings.Count(domain, ".") == 1 {
		return domain
	}

	ss := strings.Split(domain, ".")
	length := len(ss)
	if length < 2 {
		return ""
	} else {
		return strings.Join(ss[length-2:], ".")
	}
}
