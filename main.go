// Author: ZHU HAIHUA
// Date: 8/10/16
package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/kimiazhu/ginweb"
	"github.com/kimiazhu/ginweb/conf"
	"github.com/kimiazhu/golib/utils"
	"github.com/kimiazhu/log4go"
	"io/ioutil"
	"net/http"
	"strings"
)

var target = conf.ExtString("target", "http://www.google.com")
var host = conf.ExtString("host", "http://localhost:8888")

func main() {
	r := ginweb.New()

	r.Use(func(c *gin.Context) {
		log4go.Debug("handled request...uri is: %s, url is: %s", c.Request.RequestURI, c.Request.URL)
		log4go.Debug("client request is: %s", util.ReflectToString(c.Request))

		url := fmt.Sprintf("%s%s", target, c.Request.RequestURI)
		log4go.Debug("ready to request url: %s", url)
		req, _ := http.NewRequest(c.Request.Method, url, c.Request.Body)
		//req.Host = target
		for k, v := range c.Request.Header {
			for _, vv := range v {
				req.Header.Add(k, vv)
			}
		}
		req.Header.Del("Referer")
		req.Header.Add("Referer", target)

		log4go.Debug("new request is: %s", util.ReflectToString(c.Request))
		resp, err := http.DefaultClient.Do(req)

		if err != nil {
			log4go.Error("error occur while do request: %v", err)
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
			log4go.Error("error occur while read response body: %v", err)
			c.AbortWithError(http.StatusInternalServerError, err)
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
		s := strings.Replace(string(body), target, host, -1)
		c.Writer.Write([]byte(s))
	})

	ginweb.Run(conf.Conf.SERVER.PORT, r)
}