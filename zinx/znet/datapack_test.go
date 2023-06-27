package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

// 负责测试datapack拆包封包的单元测试
func TestDataPack(t *testing.T) {
	//模拟的服务器
	//1 创建socketTCP Server
	listener, err := net.Listen("tcp", "localhost:8999")
	if err != nil {
		fmt.Println("server listen err : ", err)
		return
	}

	//创建一个go协程 负责从客户端处理业务
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Println("server accept error", err)
				continue
			}
			go func(conn net.Conn) {
				//处理客户端的请求
				/*
					-------------->拆包过程<-----------------------
					定义一个拆包的对象dp
				*/
				dp := NewDataPack()
				for {
					//1 第一次从conn中读取
					headData := make([]byte, dp.GetHeadLen())
					_, err := io.ReadFull(conn, headData)
					if err != nil {
						fmt.Println("read head error")
						return
					}

					msgHead, err := dp.UnPack(headData)
					if err != nil {
						fmt.Println("server unpack err ", err)
						return
					}
					if msgHead.GetMsgLen() > 0 {
						//2 第二次从conn中读取,根据head中的dataLen，再读取dataLen内容
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetMsgLen())
						//根据dataLen的长度再次从io流中读取

						_, err = io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("server unpack data err: ", err)
							return
						}

						//完整的一个小希已经读取完毕
						fmt.Println("------>Recv MsgID: ", msg.ID, " dataLen= ", msg.DataLen, "data=", msg.Data)
					}
				}

			}(conn)
		}
	}()

	//模拟客户端
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client dial err: ", err)
		return
	}

	//创建一个封包对象 dp
	dp := NewDataPack()
	//模拟粘包过程,封装两个msg一同发送
	//封装第一个msg1包
	msg1 := &Message{
		ID:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("cleint pack msg1 error: ", err)
		return
	}
	//封装第二个包msg2
	msg2 := &Message{
		ID:      2,
		DataLen: 7,
		Data:    []byte{'n', 'i', 'h', 'a', 'o', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 error: ", err)
		return
	}

	//将两个包黏在一起
	sendData := append(sendData1, sendData2...)

	//一次性发送给服务器
	conn.Write(sendData)

	//客户端阻塞
	select {}

}
