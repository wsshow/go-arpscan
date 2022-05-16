package protocol

import (
	"goscan/global"
	"goscan/storage"
	"net"
	"sync"
)

var curDeviceName string

// 将抓到的数据集加入到data中，同时重置计时器
func PushData(ip string, mac net.HardwareAddr, hostname, manuf string) {
	global.Ticker.Stop()
	var mu sync.RWMutex
	mu.RLock()
	defer func() {
		global.Ticker.Reset()
		mu.RUnlock()
	}()
	if _, ok := storage.ResultData[ip]; !ok {
		storage.ResultData[ip] = storage.BaseInfo{Mac: mac, Hostname: hostname, Manuf: manuf}
		return
	}
	info := storage.ResultData[ip]
	if len(hostname) > 0 && len(info.Hostname) == 0 {
		info.Hostname = hostname
	}
	if len(manuf) > 0 && len(info.Manuf) == 0 {
		info.Manuf = manuf
	}
	if mac != nil {
		info.Mac = mac
	}
	storage.ResultData[ip] = info
}
