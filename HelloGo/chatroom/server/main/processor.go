package main

import (
	"fmt"
	"hellogo/chatroom/common/message"
	"hellogo/chatroom/server/process"
	"hellogo/chatroom/server/utils"
	"io"
	"net"
)

//先创建一个processor结构体
type Precessor struct {
	Conn net.Conn
}

//编写一个ServerProcessMes函数
//功能：根据客户端发送消息种类不同决定调用那个函数来处理
func (this *Precessor) serverProcessMes(mes *message.Message) (err error) {
	//看看能否收到客户端发送的群消息

	switch mes.Type {
	case message.LoginMesType:
		//处理登录
		//创建一个UserProcess实例
		up := &process.UserProcess{
			Conn: this.Conn,
		}
		err = up.SeverProcessLogin(mes)
	case message.RegisterMesType:
		//处理注册
		up := &process.UserProcess{
			Conn: this.Conn,
		}
		err = up.SeverProcessRegister(mes)
	case message.SmsMesType:
		//创建一个SmsProcess实例完成转发群聊消息
		smsProcess := &process.SmsProcess{}
		smsProcess.SendGroupMes(mes)
	case message.SmsPerMesType:
		smsProcess := &process.SmsProcess{}
		smsProcess.SendPersonalMes(mes)
	case message.CheckPassType:
		passProcess := &process.PassProcess{
			Conn: this.Conn,
		}
		err = passProcess.SeverProcessCheckPass(mes)
	case message.ChangePassType:
		passProcess := &process.PassProcess{
			Conn: this.Conn,
		}
		err = passProcess.SeverProcessChangePass(mes)
	case message.NoticeExitType:
		up := &process.UserProcess{}
		up.Offline(mes)
	default:
		fmt.Println("消息类型不存在无法处理")
	}
	return
}

func (this *Precessor) process2() (err error) {

	//循环读客户端发送的信息
	for {
		//创建一个Transfer实例完成读包任务
		tf := &utils.Transfer{
			Conn: this.Conn,
		}
		mes, err := tf.ReadPkg()
		if err != nil {
			if err == io.EOF {
				fmt.Println("客户端退出 服务器也正常关闭")
				return err
			} else {
				fmt.Println("readPkg err = ", err)
				return err
			}
		}
		// fmt.Println("mes = ", mes)
		err = this.serverProcessMes(&mes)
		if err != nil {
			return err
		}
	}
}
