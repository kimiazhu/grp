// Author: ZHU HAIHUA
// Date: 8/10/16
package main

import (
	"fmt"
	"github.com/kimiazhu/ginweb"
	"github.com/kimiazhu/ginweb/conf"
	"github.com/kimiazhu/golib/utils"
	"github.com/kimiazhu/grp/midware"
	"github.com/kimiazhu/grp/model"
	"github.com/kimiazhu/log4go"
)

var ReverseProxies model.ReverseProxies = make(model.ReverseProxies)
var Proxies model.Proxies = make(model.Proxies)

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
		local = combinePort(local, localSchema)
		model.SvrCnf[local] = &model.Server{Host: local, Schema: localSchema}

		remote := _v["remote"].(string)
		remoteSchema := "http"
		if rs, ok := _v["remoteSchema"]; ok {
			remoteSchema = rs.(string)
		}
		model.SvrCnf[remote] = &model.Server{Host: remote, Schema: remoteSchema}

		ReverseProxies[remote] = local
		Proxies[local] = remote
	}
	log4go.Info("Load proxy list: %v", util.ReflectToString(Proxies))
}

func combinePort(host, schema string) string {
	if (schema == "http" && conf.Conf.SERVER.PORT != "80") ||
		schema == "https" && conf.Conf.SERVER.PORT != "443" {
		return fmt.Sprintf("%s:%s", host, conf.Conf.SERVER.PORT)
	} else {
		return host
	}
}

func main() {
	g := ginweb.New()
	g.Use(midware.Route(ReverseProxies, Proxies))
	ginweb.Run(conf.Conf.SERVER.PORT, g)
}
