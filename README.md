# go-arpscan
## 描述
用于内网ARP探测，以获取在线主机名称和对应网卡型号。

## 基本功能
1. 自动识别本机IP和网卡
2. arp探测获取局域网存活主机IP和物理地址
3. 获取主机用户名
4. 获取网卡型号
5. 兼容Windows、Linux、MacOS
6. 主机名兼容中文
7. 集成网卡厂商信息

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
https://github.com/google/gopacket/blob/master/examples/arpscan/arpscan.go<br>
https://github.com/timest/goscan/tree/master/src/main<br>
https://github.com/timest/gomanuf