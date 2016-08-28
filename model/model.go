// Author: ZHU HAIHUA
// Date: 8/18/16
package model

import "fmt"

type Server struct {
	Host   string
	Schema string
}

func (s *Server) String() string {
	return fmt.Sprintf("%s://%s", s.Schema, s.Host)
}

// 反向代理是一个map对象,key是需要被代理的远程地址,
// value是服务器本地地址,包括端口号
type ReverseProxy map[string]string

// 代理列表是一个map对象,key和value值和ReverseProxies
// 正好相反
type Proxy map[string]string

// 服务器配置用于存储远程或者本地服务器的配置信息,
// key 是服务器 Host
// value 是Server对象指针
type ServerConfig map[string]*Server

var SvrCnf ServerConfig = make(ServerConfig)
var ReverseProxies ReverseProxy = make(ReverseProxy)
var Proxies Proxy = make(Proxy)
var LocalTopDomain string = ""
