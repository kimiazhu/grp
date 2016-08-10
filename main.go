// Author: ZHU HAIHUA
// Date: 8/10/16
package main

import (
	"github.com/kimiazhu/ginweb"
	"github.com/kimiazhu/ginweb/conf"
	"github.com/kimiazhu/grp/route"
	"github.com/kimiazhu/log4go"
)

//var Proxies []route.ReverseProxy = make([]route.ReverseProxy, 0)
var ReverseProxies route.ReverseProxies = make(route.ReverseProxies)
var Proxies route.Proxies = make(route.Proxies)

func init() {
	proxy := conf.Ext("proxy", map[interface{}]interface{}{})
	for k, v := range proxy.(map[interface{}]interface{}) {
		log4go.Debug("found proxy config: %s", k.(string))
		_v := v.(map[interface{}]interface{})
		local := _v["local"].(string)
		remote := _v["remote"].(string)
		//Proxies = append(Proxies, route.ReverseProxy{Local: local.(string), Remote: remote.(string)})
		ReverseProxies[remote] = local
		Proxies[local] = remote
	}
	log4go.Info("Load proxy list: %v", Proxies)
}

func main() {
	g := ginweb.New()
	g.Use(route.Route(ReverseProxies, Proxies))
	ginweb.Run(conf.Conf.SERVER.PORT, g)
}
