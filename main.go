// Author: ZHU HAIHUA
// Date: 8/10/16
package main

import (
	"github.com/kimiazhu/ginweb"
	"github.com/kimiazhu/ginweb/conf"
)

var remote = conf.ExtString("target", "http://www.google.com")
var local = conf.ExtString("host", "http://localhost:8888")

func main() {
	r := ginweb.New()

	r.Use(Route(local, remote))

	ginweb.Run(conf.Conf.SERVER.PORT, r)
}
