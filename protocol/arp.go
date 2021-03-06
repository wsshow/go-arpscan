package protocol

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"go-arpscan/utils"
	"log"
	"net"
	"strings"
	"sync"

	manuf "go-arpscan/manuf"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
)

func ScanARP(ctx context.Context) {
	// 获取系统所有网络接口的切片
	ifaces, err := net.Interfaces()
	if err != nil {
		panic(err)
	}

	// 获取所有设备接口的切片
	devices, err := pcap.FindAllDevs()
	if err != nil {
		panic(err)
	}

	var wg sync.WaitGroup
	for _, iface := range ifaces {
		wg.Add(1)
		// 对获取到的所有接口进行arp扫描
		go func(iface net.Interface) {
			defer wg.Done()
			if err := scan(&iface, &devices, ctx); err != nil {
				log.Printf("interface %v: %v", iface.Name, err)
			}
		}(iface)
	}
	wg.Wait()
}

func scan(iface *net.Interface, devices *[]pcap.Interface, ctx context.Context) error {
	// We just look for IPv4 addresses, so try to find if the interface has one.
	var addr *net.IPNet
	if addrs, err := iface.Addrs(); err != nil {
		return err
	} else {
		for _, a := range addrs {
			if ipnet, ok := a.(*net.IPNet); ok {
				if ip4 := ipnet.IP.To4(); ip4 != nil {
					addr = &net.IPNet{
						IP:   ip4,
						Mask: ipnet.Mask[len(ipnet.Mask)-4:],
					}
					break
				}
			}
		}
	}
	// Sanity-check that the interface has a good address.
	if addr == nil {
		return errors.New("no good IP network found")
	} else if addr.IP[0] == 127 {
		return errors.New("skipping localhost")
	} else if addr.Mask[0] != 0xff || addr.Mask[1] != 0xff {
		return errors.New("mask means network is too large")
	}
	log.Printf("Using network range %v for interface %v", addr, iface.Name)

	// 主要解决windows中设备名称与接口名称不对应的问题
	var deviceName string
	for _, d := range *devices {
		if strings.Contains(fmt.Sprint(d.Addresses), fmt.Sprint(addr.IP)) {
			deviceName = d.Name
		}
	}

	if deviceName == "" {
		return fmt.Errorf("cannot find the corresponding device for the interface %s", iface.Name)
	}

	curDeviceName = deviceName

	// Open up a pcap handle for packet reads/writes.
	handle, err := pcap.OpenLive(deviceName, 65536, true, pcap.BlockForever)
	if err != nil {
		return err
	}
	defer handle.Close()

	go listenMDNS(deviceName, ctx)
	go listenNBNS(deviceName, addr, ctx)

	// Start up a goroutine to read in packet data.
	go readARP(addr, handle, iface, ctx)

	// Write our scan packets out to the handle.
	if err := writeARP(handle, iface, addr); err != nil {
		log.Printf("error writing packets on %v: %v", iface.Name, err)
		return err
	}
	<-ctx.Done()
	return nil
}

// 读取返回的arp数据包
func readARP(addr *net.IPNet, handle *pcap.Handle, iface *net.Interface, ctx context.Context) {
	src := gopacket.NewPacketSource(handle, layers.LayerTypeEthernet)
	in := src.Packets()
	for {
		var packet gopacket.Packet
		select {
		case <-ctx.Done():
			return
		case packet = <-in:
			arpLayer := packet.Layer(layers.LayerTypeARP)
			if arpLayer == nil {
				continue
			}
			arp := arpLayer.(*layers.ARP)
			// 过滤自己发送的数据包
			if arp.Operation != layers.ARPReply || bytes.Equal([]byte(iface.HardwareAddr), arp.SourceHwAddress) {
				continue
			}
			// 从收到的arp数据包中提取物理地址
			mac := net.HardwareAddr(arp.SourceHwAddress)
			// 获取网卡供应商信息
			m := manuf.Search(mac.String())
			PushData(utils.ParseIP(arp.SourceProtAddress).String(), mac, "", m)
			// 发送对应数据包以获取主机名
			if strings.Contains(m, "Apple") {
				go sendMdns(addr, iface.HardwareAddr, utils.ParseIP(arp.SourceProtAddress), mac)
			} else {
				go sendNbns(addr, iface.HardwareAddr, utils.ParseIP(arp.SourceProtAddress), mac)
			}
		}
	}
}

// 准备arp数据包
func writeARP(handle *pcap.Handle, iface *net.Interface, addr *net.IPNet) error {
	// Set up all the layers' fields we can.
	eth := layers.Ethernet{
		SrcMAC:       iface.HardwareAddr,
		DstMAC:       net.HardwareAddr{0xff, 0xff, 0xff, 0xff, 0xff, 0xff},
		EthernetType: layers.EthernetTypeARP,
	}
	arp := layers.ARP{
		AddrType:          layers.LinkTypeEthernet,
		Protocol:          layers.EthernetTypeIPv4,
		HwAddressSize:     6,
		ProtAddressSize:   4,
		Operation:         layers.ARPRequest,
		SourceHwAddress:   []byte(iface.HardwareAddr),
		SourceProtAddress: []byte(addr.IP),
		DstHwAddress:      []byte{0, 0, 0, 0, 0, 0},
	}
	// Set up buffer and options for serialization.
	buf := gopacket.NewSerializeBuffer()
	opts := gopacket.SerializeOptions{
		FixLengths:       true,
		ComputeChecksums: true,
	}
	// Send one packet for every address.
	for _, ip := range ips(addr) {
		arp.DstProtAddress = []byte(ip)
		gopacket.SerializeLayers(buf, opts, &eth, &arp)
		if err := handle.WritePacketData(buf.Bytes()); err != nil {
			return err
		}
	}
	return nil
}

func ips(n *net.IPNet) (out []net.IP) {
	num := binary.BigEndian.Uint32([]byte(n.IP))
	mask := binary.BigEndian.Uint32([]byte(n.Mask))
	network := num & mask
	broadcast := network | ^mask
	for network++; network < broadcast; network++ {
		var buf [4]byte
		binary.BigEndian.PutUint32(buf[:], network)
		out = append(out, net.IP(buf[:]))
	}
	return
}
