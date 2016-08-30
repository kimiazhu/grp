// Author: ZHU HAIHUA
// Date: 8/10/16
package midware

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kimiazhu/grp/filter"
	"github.com/kimiazhu/grp/model"
	"github.com/kimiazhu/grp/util/io"
	"github.com/kimiazhu/log4go"
	"net/http"
	"net/http/httputil"
)

// Router 返回一个中间件函数, 将本地请求重定向至远端服务器,
// 在拿到远端服务器应答之后, 替换应答中的远程服务器域名后将
// 其回写到本地。
func Route() func(*gin.Context) {
	return func(c *gin.Context) {
		local := c.Request.Host
		remote := model.Proxies[local]
		if remote == "" {
			log4go.Error("Local host [%s] cannot be found", local)
			msg := fmt.Errorf("no proxy for %s", local)
			c.AbortWithError(http.StatusNotFound, msg)
			c.String(http.StatusNotFound, msg.Error())
			return
		}

		requestURI := ioutils.ReplaceHost(c.Request.RequestURI, true)
		reqUrl := fmt.Sprintf("%s://%s%s", model.SvrCnf[remote].Schema, remote, requestURI)
		method := c.Request.Method
		log4go.Fine("ready to request url: %s, method: %s", reqUrl, method)

		req := filter.CreateRequest(c, method, reqUrl, remote)
		resp, err := http.DefaultTransport.RoundTrip(req)

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

		log4go.Debug("continue to contruct response of url: %s, httpStatus: %v", reqUrl, resp.Status)
		defer resp.Body.Close()

		body, unzipped, err := filter.SmartRead(resp, true)
		filter.SmartWrite(c, resp, body, unzipped)
	}

}
