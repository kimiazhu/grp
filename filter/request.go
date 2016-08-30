// Author: ZHU HAIHUA
// Date: 8/23/16
package filter

import (
	"github.com/gin-gonic/gin"
	"github.com/kimiazhu/grp/util/io"
	"github.com/kimiazhu/log4go"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

func CreateRequest(c *gin.Context, method, reqUrl, remoteHost string) *http.Request {
	req, _ := http.NewRequest(method, reqUrl, nil)

	contentLength := ""
	if method == "POST" {
		c.Request.ParseMultipartForm(32 << 20) //32M
		form := url.Values{}
		for k, v := range c.Request.PostForm {
			for _, vv := range v {
				form.Add(k, ioutils.ReplaceHost(vv, true))
			}
		}

		encodedForm := form.Encode()
		req.Body = ioutil.NopCloser(strings.NewReader(encodedForm))
		contentLength = strconv.Itoa(len([]byte(encodedForm)))
	} else {
		req.Body = c.Request.Body
	}

	FilterHeader(c.Request.Header, req.Header, remoteHost, true)
	req.Header.Add("Host", remoteHost)

	if cl := req.Header.Get("Content-Length"); len(cl) > 0 {
		log4go.Fine("old content length: %s, new content length: %s", cl, contentLength)
		req.Header.Set("Content-Length", contentLength)
	}

	req.Host = remoteHost
	return req
}
