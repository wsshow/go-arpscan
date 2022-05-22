package main

import (
	"context"
	"go-arpscan/global"
	"go-arpscan/protocol"
	"go-arpscan/storage"
	"go-arpscan/utils"
	"log"
	"os"
	"runtime"
	"sort"
	"time"
)

// 初始化日志配置
func InitLogConfig() {
	log.SetPrefix("[go-arpscan] ")
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
}

func init() {
	InitLogConfig()
	// 初始化全局计时器对象
	global.Ticker = utils.NewTicker(utils.WithResetTime(5 * time.Second))
}

// 输出结果数据
func PrintData(data map[string]storage.BaseInfo) {
	var keys utils.IPSlice
	for k := range data {
		keys = append(keys, utils.ParseIPString(k))
	}
	sort.Sort(keys)
	for _, k := range keys {
		d := data[k.String()]
		mac := ""
		if d.Mac != nil {
			mac = d.Mac.String()
		}
		// 检查主机名是否含有中文
		if !utils.IsUtf8([]byte(d.Hostname)) {
			d.Hostname = utils.ConvertGBK2StrFromStr(d.Hostname)
		}
		log.Printf("%-15s %-17s %-30s %-10s\n", k.String(), mac, d.Hostname, d.Manuf)
	}
	log.Println("Total Count:", len(keys))
}

func main() {
	// 要求必须以root权限运行
	if (runtime.GOOS == "linux" || runtime.GOOS == "darwin") && os.Geteuid() != 0 {
		log.Fatal("go-arpscan must run as root")
	}
	ctx, cancel := context.WithCancel(context.Background())
	// 开启arp扫描
	go protocol.ScanARP(ctx)
	// 等待计时器停止
	global.Ticker.Wait()
	// 取消未完任务
	cancel()
	// 打印扫描数据
	PrintData(storage.ResultData)
}
