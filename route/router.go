// Author: ZHU HAIHUA
// Date: 8/10/16
package route

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kimiazhu/log4go"
	"io/ioutil"
	"net/http"
	"strings"
	"net/http/httputil"
)

type Server struct {
	Host string
	Schema string
}

// 反向代理是一个map对象,key是需要被代理的远程地址,
// value是服务器本地地址,包括端口号
type ReverseProxies map[string]string

// 代理列表是一个map对象,key和value值和ReverseProxies
// 正好相反
type Proxies map[string]string

// 服务器配置用于存储远程或者本地服务器的配置信息,
// key 是服务器 Host
// value 是Server对象指针
type ServerConfig map[string]*Server

var SvrCnf ServerConfig = make(ServerConfig)

// Router 返回一个中间件函数, 将本地请求重定向至远端服务器,
// 在拿到远端服务器应答之后, 替换应答中的远程服务器域名后将
// 其回写到本地。
func Route(r ReverseProxies, p Proxies) func(*gin.Context) {
	return func(c *gin.Context) {
		local := c.Request.Host
		remote := p[local]
		if remote == "" {
			log4go.Error("Local host [%s] cannot be found", local)
			msg := fmt.Errorf("no proxy for %s", local)
			c.AbortWithError(http.StatusNotFound, msg);
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
		req.Header.Del("Referer")
		req.Header.Add("Referer", remote)
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

		for k, v := range resp.Header {
			for _, vv := range v {
				c.Writer.Header().Add(k, vv)
			}
		}

		for _, value := range resp.Request.Cookies() {
			c.Writer.Header().Add(value.Name, value.Value)
		}

		c.Writer.WriteHeader(resp.StatusCode)
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log4go.Error("error occur while read response, error is: %v", err)
			dat, e := httputil.DumpResponse(resp, true)
			if e != nil {
				log4go.Error("dump response failed: %v", e)
			} else {
				log4go.Error("dumped response body: %s", string(dat))
			}
			c.AbortWithError(http.StatusInternalServerError, err)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}

		result := string(body)
		for _remote, _local := range r {
			result = strings.Replace(result, _remote, _local, -1)
		}
		c.Writer.Write([]byte(result))
	}

}
