package utils

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hellogo/chatroom/common/message"
	"net"
)

//这里将这些方法关联到结构体中
type Transfer struct {
	Conn net.Conn
	Buf  [8096]byte
}

func (this *Transfer) ReadPkg() (mes message.Message, err error) {
	//读取数据包，直接封装成一个函数
	// buf := make([]byte, 8096)
	fmt.Println("读取客户端发送的数据")
	//conn.Read在conn没有被关闭的情况下才会阻塞
	_, err = this.Conn.Read(this.Buf[:4])
	if err != nil {
		//err = errors.New("read pkg header error")//自定义error
		return
	}
	//根据buf[:4]转成一个uint32类型
	var pkgLen uint32
	pkgLen = binary.BigEndian.Uint32(this.Buf[:4])
	//根据pkgLen读取消息内容
	n, err := this.Conn.Read(this.Buf[:pkgLen]) //从conn中读pkgLen个字节到buf
	if n != int(pkgLen) || err != nil {
		//err = errors.New("read pkg body error")
		return
	}
	//把pkgLen反序列化message.Message
	err = json.Unmarshal(this.Buf[:pkgLen], &mes) //注意&
	if err != nil {
		fmt.Println("json.Unmarshal err", err)
		return
	}
	return
}

func (this *Transfer) WritePkg(data []byte) (err error) {
	//先发送长度给对方
	var pkgLen uint32
	pkgLen = uint32(len(data))
	// var buf [4]byte
	binary.BigEndian.PutUint32(this.Buf[0:4], pkgLen)
	//发送长度
	n, err := this.Conn.Write(this.Buf[0:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) err = ", err)
		return
	}
	//再发送data本身
	n, err = this.Conn.Write(data)
	if n != int(pkgLen) || err != nil {
		fmt.Println("conn.Write(bytes) err = ", err)
		return
	}
	return
}
