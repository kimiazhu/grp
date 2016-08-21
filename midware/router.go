// Author: ZHU HAIHUA
// Date: 8/10/16
package midware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kimiazhu/grp/model"
	"github.com/kimiazhu/grp/util/io"
	"github.com/kimiazhu/log4go"
	"net/http"
	"net/http/httputil"
)

// Router 返回一个中间件函数, 将本地请求重定向至远端服务器,
// 在拿到远端服务器应答之后, 替换应答中的远程服务器域名后将
// 其回写到本地。
func Route(r model.ReverseProxies, p model.Proxies) func(*gin.Context) {
	return func(c *gin.Context) {
		local := c.Request.Host
		remote := p[local]
		if remote == "" {
			log4go.Error("Local host [%s] cannot be found", local)
			msg := fmt.Errorf("no proxy for %s", local)
			c.AbortWithError(http.StatusNotFound, msg)
			c.String(http.StatusNotFound, msg.Error())
			return
		}

		url := fmt.Sprintf("https://%s%s", remote, c.Request.RequestURI)
		log4go.Fine("ready to request url: %s", url)
		req, _ := http.NewRequest(c.Request.Method, url, c.Request.Body)
		//req.Host = target
		for k, v := range c.Request.Header {
			for _, vv := range v {
				req.Header.Add(k, vv)
			}
		}

		req.Host = remote
		req.Header.Set("Referer", remote)
		req.Header.Set("X_Forward_For", local)
		req.Header.Set("X-Real-IP", local)

		//log4go.Debug("new request is: %s", util.ReflectToString(c.Request))
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			log4go.Error("error occur while do request, error is: %v", err)
			dat, e := httputil.DumpRequestOut(req, true)
			if e != nil {
				log4go.Debug("dump out requst failed: %v", e)
			} else {
				log4go.Error("dumped request: %s", string(dat))
			}
			c.AbortWithError(http.StatusInternalServerError, err)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		log4go.Debug("continue to contruct response of url: %s", url)
		defer resp.Body.Close()

		//for _, value := range resp.Request.Cookies() {
		//	c.Writer.Header().Add(value.Name, value.Value)
		//}

		body, unzipped, err := ioutils.SmartRead(resp, p, true)
		ioutils.SmartWrite(c, resp, body, unzipped)
	}

}
