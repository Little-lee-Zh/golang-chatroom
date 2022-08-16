package process

import (
	"encoding/json"
	"fmt"
	"hellogo/chatroom/client/utils"
	"hellogo/chatroom/common/message"
	"net"
	"os"
)

//显示登录成功后的界面

func ShowMenu() {
	fmt.Println("--------恭喜你登录成功------------")
	fmt.Println("--------1.显示在线用户------------")
	fmt.Println("--------2.发送群聊消息------------")
	fmt.Println("--------3.发送私聊消息------------")
	fmt.Println("--------4.修改用户密码------------")
	fmt.Println("--------5.退出聊天系统------------")
	fmt.Println("请选择(1-5):")
	var key int
	var content string
	var id int
	var choice string
	var oldPassport string
	var newPassport string
	//总会使用SmsProcess实例，因此定义在switch外部
	smsProcess := &SmsProcess{}
	fmt.Scanf("%d\n", &key)
	switch key {
	case 1:
		outputOnlineUser()
	case 2:
		fmt.Println("请输入群聊消息:")
		fmt.Scanf("%s\n", &content)
		smsProcess.SendGroupMes(content)
	case 3:
		fmt.Println("请输入私聊对象id:")
		fmt.Scanf("%d\n", &id)
		fmt.Println("请输入私聊消息:")
		fmt.Scanf("%s\n", &content)
		smsProcess.SendPersonalMes(id, content)
	case 4:
		fmt.Println("请确认是否修改密码？YES/NO")
		fmt.Scanf("%s\n", &choice)
		if choice == "YES" {
			fmt.Println("请输入原始密码:")
			fmt.Scanf("%s\n", &oldPassport)
			//判断输入密码是否正确
			passprocess := &PassProcess{}
			result, err := passprocess.CheckPassport(oldPassport)
			if err != nil {
				fmt.Println(err)
			} else {
				if result {
					fmt.Println("请输入新密码:")
					fmt.Scanf("%s\n", &newPassport)
					passprocess.ChangePassport(newPassport)
				}
			}
		}
	case 5:
		fmt.Println("退出系统...")
		smsProcess.NoticeExitMes()
		os.Exit(0)
	default:
		fmt.Println("输入选项不正确")
	}
}

//和服务器保持通讯
func serverProcessMes(conn net.Conn) {
	//创建一个Transfer实例不停的读取
	tf := &utils.Transfer{
		Conn: conn,
	}
	for {
		//fmt.Println("客户端正在等待读取服务器发送的消息")
		mes, err := tf.ReadPkg()
		if err != nil {
			fmt.Println("tf.ReadPkg err = ", err)
			return
		}
		//如果读到消息，下一步处理逻辑
		switch mes.Type {
		case message.NotifyUserStatusMesType: //有人上线
			//1.取出NotifyUserStatusMes
			var notifyUserStatusMes message.NotifyUserStatusMes
			json.Unmarshal([]byte(mes.Data), &notifyUserStatusMes)
			//2.把用户的信息状态保存到客户map
			updateUserStatus(&notifyUserStatusMes)
		case message.SmsMesType: //有人群发消息
			outputGroupMes(&mes)
		case message.SmsPerMesType: //私聊消息
			outputPersonalMes(&mes)
		default:
			fmt.Println("服务器端返回了未知的消息类型")
		}
	}
}
