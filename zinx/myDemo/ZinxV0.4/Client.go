package main

import (
	"fmt"
	"net"
	"time"
)

// 模拟客户端
func main() {
	fmt.Println("client start.....")
	time.Sleep(1 * time.Second)
	//1直接连接远程服务器 得到一个conn连接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err, exit....")
		return
	}
	defer conn.Close()
	//2连接调用write写数据
	for {
		_, err := conn.Write([]byte("Hello Zinx V0.4...."))
		if err != nil {
			fmt.Println("write conn err", err)
			return
		}
		buf := make([]byte, 512)
		cnt, err := conn.Read(buf)
		if err != nil {
			fmt.Println("read buff error: ", err)
			return
		}
		fmt.Printf("server call back: %s, cnt = %d\n", buf, cnt)
		time.Sleep(1 * time.Second)
	}

}
