package main

import (
	"flag"
)

// 返回服务器端口号,线程数，配置文件路径
func clArgs() (int, int, string) {
	p := flag.Int("p", 8080, "服务器端口号")
	t := flag.Int("t", 1, "是否启用多进程模式，0为关闭，1为启用")
	f := flag.String("f", "proxies.json", "负载均衡配置文件")
	flag.Parse()
	return *p, *t, *f
}
