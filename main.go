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

func InitLogConfig() {
	log.SetPrefix("[go-arpscan] ")
	log.SetFlags(log.Lshortfile | log.Ldate | log.Lmicroseconds)
}

func init() {
	InitLogConfig()
	global.Ticker = utils.NewTicker(utils.WithResetTime(5 * time.Second))
}

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
		if !utils.IsUtf8([]byte(d.Hostname)) {
			d.Hostname = utils.ConvertGBK2StrFromStr(d.Hostname)
		}
		log.Printf("%-15s %-17s %-30s %-10s\n", k.String(), mac, d.Hostname, d.Manuf)
	}
}

func main() {
	if runtime.GOOS == "linux" && os.Geteuid() != 0 {
		log.Fatal("goscan must run as root")
	}
	ctx, cancel := context.WithCancel(context.Background())
	go protocol.ScanARP(ctx)
	global.Ticker.Wait()
	cancel()
	PrintData(storage.ResultData)
}
