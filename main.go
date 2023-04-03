package main

import (
	"fmt"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
	"io"
	"log"
	"net/http"
	"runtime"
	"strconv"
	"time"
)

var proxies Proxies

func Index(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "/status 负载均衡状态\n\n/proxy 代理ip接口")
}

func Status(w http.ResponseWriter, r *http.Request) {
	v, _ := mem.VirtualMemory()
	totalCpuPercent, _ := cpu.Percent(1*time.Second, false)
	perCpuPercents, _ := cpu.Percent(1*time.Second, true)
	io.WriteString(w, fmt.Sprintf("CPU总使用率: %.2f%%\n", totalCpuPercent[0]))
	io.WriteString(w, fmt.Sprintf("CPU核心使用率: %.2v\n\n", perCpuPercents))
	memUsageStr := fmt.Sprintf("内存总计: %vM, 空闲:%vM, 使用率:%.2f%%\n\n", v.Total/1024/1024, v.Free/1024/1024, v.UsedPercent)
	io.WriteString(w, memUsageStr)
	io.WriteString(w, "代理 | 权重 | 代理地址\n")
	for _, v := range proxies {
		io.WriteString(w, v.Name+": "+strconv.Itoa(v.Weight)+" | "+v.Url+"\n")
	}
}

func Redirect(writer http.ResponseWriter, request *http.Request) {
	//重定向
	proxy := WeightedRandomProxy(proxies)
	http.Redirect(writer, request, proxy.Url, http.StatusTemporaryRedirect)
}

func main() {
	port, worker, filename := clArgs()
	fmt.Print("服务启动，使用-h命令查看帮助\n")
	fmt.Print("http://127.0.0.1:"+strconv.Itoa(port)+"\n")
	proxies = LoadJson("./" + filename)
	if worker < 1 {
		worker = runtime.NumCPU()
	}
	runtime.GOMAXPROCS(worker)
	log.Printf("服务器端口%d", port)
	log.Printf("线程数：%d", worker)
	log.Printf("配置文件：%s", filename)
	http.HandleFunc("/", Index)
	http.HandleFunc("/status", Status)
	http.HandleFunc("/proxy", Redirect)
	http.ListenAndServe("0.0.0.0:"+strconv.Itoa(port), nil)
}
