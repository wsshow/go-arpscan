# go-arpscan
## 描述
用于内网ARP探测，以获取在线主机名称和对应网卡型号。

## 使用
```shell
sudo ./go-arpscan
```
## 效果图
![go-arpscan演示图](./screenshot/go-arpscan%E6%95%88%E6%9E%9C%E5%9B%BE.png)

## 依赖
linux/macos需要安装libpcap<br>
源码链接：https://github.com/the-tcpdump-group/libpcap

windows需要安装npcap<br>
下载链接：https://npcap.com/#download

## 参考链接
https://github.com/google/gopacket/blob/master/examples/arpscan/arpscan.go
https://github.com/timest/goscan/tree/master/src/main
https://github.com/timest/gomanuf