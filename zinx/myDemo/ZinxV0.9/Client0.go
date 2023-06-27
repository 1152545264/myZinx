package main

import (
	"fmt"
	"io"
	"net"
	"time"
	"zinx/znet"
)

// 模拟客户端
func main() {
	fmt.Println("client0 start.....")
	time.Sleep(1 * time.Second)
	//1直接连接远程服务器 得到一个conn连接
	conn, err := net.Dial("tcp", "0.0.0.0:8999")
	if err != nil {
		fmt.Println("client start err, exit....")
		return
	}
	defer conn.Close()
	//2连接调用write写数据
	for {
		//发送封包消息
		dp := znet.NewDataPack()
		binaryData, err := dp.Pack(znet.NewMessagePackage(0, []byte("ZinxV0.8 client0 Test Message")))
		if err != nil {
			fmt.Println("Pack error ", err)
			return
		}
		if _, err = conn.Write(binaryData); err != nil {
			fmt.Println("Send Data err: ", err)
			return
		}

		//服务器应该回复一个message数据
		//需要处理粘包问题
		//先读取流中的head部分得到ID和dataLen

		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read head error: ", err)
			break
		}
		//将二进制的Head拆包到msg结构体中
		msgHead, err := dp.UnPack(binaryHead)
		if err != nil {
			fmt.Println("client unpack msgHead error ", err)
			break
		}
		if msgHead.GetMsgLen() > 0 {
			//再根据dataLen进行读取将data读出来
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetMsgLen())
			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data error: ", err)
				return
			}
			fmt.Println("------------->Recv Server Msg : ID = ", msg.ID, " len = ", msg.DataLen,
				", data = ", string(msg.Data))
		}

		time.Sleep(1 * time.Second)
	}

}
