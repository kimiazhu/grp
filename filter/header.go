// Author: ZHU HAIHUA
// Date: 8/28/16
package filter

import (
	"fmt"
	"github.com/kimiazhu/grp/model"
	"github.com/kimiazhu/grp/util/io"
	"github.com/kimiazhu/log4go"
	"net/http"
	"strings"
)

// 过滤src的Header信息, reverse=true表示将远程Host替换成本地Host。
// 替换结果将直接放置于target中
func FilterHeader(src, target http.Header, remoteHost string, isRequest bool) http.Header {
	for k, v := range src {
		for _, vv := range v {
			vv = ioutils.ReplaceHost(vv, isRequest)
			kk := strings.ToLower(k)
			if kk == "set-cookie" || kk == "cookie" {
				vv = DealCookie(remoteHost, vv, isRequest)
			}
			target.Add(k, vv)
		}
	}
	return target
}

// 处理Cookie中的域名信息, toHostName 表示这个cookie将要发送到的remoteHost
// isRequest=true表示这个是代理服务器转发客户端或者浏览器的请求。反之表示此时是
// 代理服务器转发服务端的应答
func DealCookie(remoteHost, cookieData string, isRequest bool) string {
	cookies, err := ioutils.ParseRawCookies(cookieData)
	//log4go.Debug("@@@@@@@@@@@@@@@@@@@@@@@@@@%v", cookies)
	if err != nil {
		log4go.Error("parse cookie [%s] failed with error: %v\n", cookieData, err)
		panic(fmt.Sprintf("parse cookie failed: %s", cookieData))
	}

	var newCookies string
	for i := 0; i < len(cookies); i++ {
		cookie := cookies[i]
		if topDomain := ioutils.TopDomainName(cookie.Domain); topDomain != "" {
			if !isRequest {
				cookie.Domain = model.LocalTopDomain
			} else {
				cookie.Domain = ioutils.TopDomainName(remoteHost)
			}
		}
		newCookies += (cookie.String() + ";")
	}

	if strings.HasSuffix(newCookies, ";") {
		newCookies = newCookies[:len(newCookies)-2]
	}

	return newCookies
}
