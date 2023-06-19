package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"strings"
)

var (
	svcaddr = flag.String("svcaddr",
		"jaeger-agent-cluster.inner-udp.efficiency.ww5sawfyut0k.bitsvc.io:6831", "server addr")
)

func init() {
	flag.Parse()
}

func main() {
	addr, err := net.ResolveUDPAddr("udp", *svcaddr)
	if err != nil {
		log.Fatal(err)
	}

	socket, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		log.Fatal(err)
	}
	localAddr := socket.LocalAddr().String()
	var ss []string
	if localAddr[0] == '[' {
		localAddr = localAddr[1:]
		ss = strings.SplitN(localAddr, "]", 2)
	} else {
		ss = strings.SplitN(localAddr, ":", 2)
	}
	fmt.Printf("local host=%v port=%v\n", ss[0], ss[1])
	defer socket.Close()
	// sendData := []byte("Hello server")
	// _, err = socket.Write(sendData) // 发送数据
	// if err != nil {
	// 	fmt.Println("发送数据失败, err:", err)
	// 	return
	// }
	// data := make([]byte, 4096)
	// n, remoteAddr, err := socket.ReadFromUDP(data) // 接收数据
	// if err != nil {
	// 	fmt.Println("接收数据失败, err:", err)
	// 	return
	// }
	// fmt.Printf("recv:%v addr:%v count:%v\n", string(data[:n]), remoteAddr, n)
}
