package storage

import "net"

type BaseInfo struct {
	// 物理地址
	Mac net.HardwareAddr
	// 主机名
	Hostname string
	// 厂商信息
	Manuf string
}

// 存放最终的数据
var ResultData = make(map[string]BaseInfo)
