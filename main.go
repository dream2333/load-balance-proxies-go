package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/monitor"
	"github.com/shirou/gopsutil/v3/cpu"
	"github.com/shirou/gopsutil/v3/mem"
)

var proxies Proxies

func Index(c *fiber.Ctx) error {
	return c.SendString("/status 负载均衡状态\n\n/dashboard 仪表盘 \n\n/proxy 从上游服务器抽取一条全新代理ip\n\n")
}

func Status(c *fiber.Ctx) error {
	v, _ := mem.VirtualMemory()
	totalCpuPercent, _ := cpu.Percent(1*time.Second, false)
	perCpuPercents, _ := cpu.Percent(1*time.Second, true)
	c.WriteString(fmt.Sprintf("CPU总使用率: %.2f%%\n", totalCpuPercent[0]))
	c.WriteString(fmt.Sprintf("CPU核心使用率: %.2v\n\n", perCpuPercents))
	memUsageStr := fmt.Sprintf("内存总计: %vM, 空闲:%vM, 使用率:%.2f%%\n\n", v.Total/1024/1024, v.Free/1024/1024, v.UsedPercent)
	c.WriteString(memUsageStr)
	c.WriteString("代理 | 权重 | 代理地址\n")
	for _, v := range proxies {
		c.WriteString(v.Name + ": " + strconv.Itoa(v.Weight) + " | " + v.Url + "\n")
	}
	return nil
}

func Redirect(c *fiber.Ctx) error {
	//重定向
	proxy := WeightedRandomProxy(proxies)
	c.Redirect(proxy.Url)
	return nil
}

func RedirectQpsTest(c *fiber.Ctx) error {
	//重定向
	proxy := WeightedRandomProxy(proxies)
	c.Redirect(proxy.Url, 200)
	return nil
}

func main() {
	// 读取配置文件
	port, multiprocess, filename := clArgs()
	proxies = LoadJson("./" + filename)
	fiberConfig := fiber.Config{Prefork: multiprocess == 1, AppName: "ProxyBalancer", DisableDefaultDate: true, DisableStartupMessage: true}
	monitorConfig := monitor.Config{
		Title:      "服务器状态",
		Refresh:    3 * time.Second,
		APIOnly:    false,
		Next:       nil,
		FontURL:    "https://fonts.googleapis.com/css2?family=Roboto:wght@400;900&display=swap",
		ChartJsURL: "https://cdn.jsdelivr.net/npm/chart.js@2.9/dist/Chart.bundle.min.js",
	}
	// 启动服务
	app := fiber.New(fiberConfig)
	app.Get("/", Index)
	app.Get("/dashboard", monitor.New(monitorConfig))
	app.Get("/status", Status)
	app.Get("/proxy", Redirect)
	proxies = LoadJson("./" + filename)
	if !fiber.IsChild() {
		fmt.Print("服务启动，使用-h命令查看帮助\n")
		url := "http://0.0.0.0:" + strconv.Itoa(port)
		fmt.Print(url + "/dashboard\n")
		fmt.Print(url + "/status\n")
		fmt.Print(url + "/proxy\n")
		worker := 1
		if multiprocess == 1 {
			worker = runtime.NumCPU()
		}
		// runtime.GOMAXPROCS(worker)
		log.Printf("服务器端口%d", port)
		log.Printf("进程数：%d", worker)
		log.Printf("配置文件：%s", filename)
	}
	app.Listen("0.0.0.0:" + strconv.Itoa(port))
}
