## 基于go的负载均衡代理服务器 ##

**Usage**

    lbpgo -p 8080 -t 8 -f test.json

    -f 负载均衡配置文件 (默认 "proxies.json")

    -p 服务器端口号 (默认 8080)

    -t 服务器使用的线程数，默认为cpu总线程数


/status 负载均衡状态

/proxy 代理ip接口