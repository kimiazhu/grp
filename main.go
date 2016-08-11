// Author: ZHU HAIHUA
// Date: 8/10/16
package main

import (
	"github.com/kimiazhu/ginweb"
	"github.com/kimiazhu/ginweb/conf"
	"github.com/kimiazhu/grp/route"
	"github.com/kimiazhu/log4go"
	"github.com/kimiazhu/golib/utils"
)

var ReverseProxies route.ReverseProxies = make(route.ReverseProxies)
var Proxies route.Proxies = make(route.Proxies)
//var Servers []route.Server = make([]route.Server, 0)

func init() {
	proxy := conf.Ext("proxy", map[interface{}]interface{}{})
	for k, v := range proxy.(map[interface{}]interface{}) {
		log4go.Debug("found proxy config: %s", k.(string))
		_v := v.(map[interface{}]interface{})

		local := _v["local"].(string)
		localSchema := "http"
		if ls, ok := _v["localSchema"]; ok {
			localSchema = ls.(string)
		}
		route.SvrCnf[local] = &route.Server{Host: local, Schema: localSchema}

		remote := _v["remote"].(string)
		remoteSchema := "http"
		if rs, ok := _v["remoteSchema"]; ok {
			remoteSchema = rs.(string)
		}
		route.SvrCnf[remote] = &route.Server{Host: remote, Schema: remoteSchema}

		ReverseProxies[remote] = local
		Proxies[local] = remote
	}
	log4go.Info("Load proxy list: %v", util.ReflectToString(Proxies))
}

func main() {
	g := ginweb.New()
	g.Use(route.Route(ReverseProxies, Proxies))
	ginweb.Run(conf.Conf.SERVER.PORT, g)
}
