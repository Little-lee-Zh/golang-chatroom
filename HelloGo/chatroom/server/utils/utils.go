package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hellogo/chatroom/common/message"
	"net"
)

// 将工具函数关联到结构体的方法中，供其他模块调用
// 传输工具结构体
type Transfer struct {
	Conn net.Conn
	Buf  [8096]byte // 传输时使用的缓冲
}

// 工具函数：读取客户端发送的数据包
func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	// conn.Read 只有在conn没有被关闭的情况下，才会阻塞
	// 如果客户端关闭了conn就不会阻塞了
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		// err = errors.New("read pkg header error")
		return
	}

	// 将Buf[:4]转成 uint32类型，即为数据的长度
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[0:4])
	// 根据pkgLen读取消息内容
	n, err := this.Conn.Read(this.Buf[:pkgLen])
	if n != int(pkgLen) || err != nil {
		// err = errors.New("read pkg body error")
		return
	}

	// 将数据反序列化 -> message.Message
	err = json.Unmarshal(this.Buf[:pkgLen], &mes)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	return
}

// 工具函数：向客户端返回数据包
func (this *Transfer) WritePkg(data []byte) (err error) {

	// 先发送一个长度给客户端
	var pkgLen uint32
	pkgLen = uint32(len(data))
	// var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	// 发送长度
	n, err := this.Conn.Write(this.Buf[:4])
	if n != 4 || err != nil {
		fmt.Println("this.Conn.Write(this.Buf[:4]) fail:", err)
		return
	}

	// 发送data
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("this.Conn.Write(data) fail:", err)
		return
	}

	return
}
