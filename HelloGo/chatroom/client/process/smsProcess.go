package process

import (
	"encoding/json"
	"fmt"
	"hellogo/chatroom/client/utils"
	"hellogo/chatroom/common/message"
)

type SmsProcess struct {
}

//发送群聊消息
func (this *SmsProcess) SendGroupMes(content string) (err error) {
	//1.创建一个Mes
	var mes message.Message
	mes.Type = message.SmsMesType
	//2.创建一个SmsMes实例
	var smsMes message.SmsMes
	smsMes.Content = content
	smsMes.UserId = CurUser.UserId
	smsMes.UserStatus = CurUser.UserStatus
	//3.序列化smsMes
	data, err := json.Marshal(smsMes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal fail := ", err.Error())
		return
	}
	mes.Data = string(data)
	//4.mes再次序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendGroupMes json.Marshal fail := ", err.Error())
		return
	}
	//将序列化后的mes发送给服务器
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	//6.发送
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendGroupMes err := ", err.Error())
		return
	}
	return
}

func (this *SmsProcess) SendPersonalMes(id int, content string) (err error) {
	//1.创建一个Mes
	var mes message.Message
	mes.Type = message.SmsPerMesType
	//2.创建一个SmsPerMes实例
	var smsPerMes message.SmsPerMes
	smsPerMes.Content = content
	smsPerMes.UserId = CurUser.UserId
	smsPerMes.UserStatus = CurUser.UserStatus
	smsPerMes.SendToUser = id
	//3.序列化smsPerMes
	data, err := json.Marshal(smsPerMes)
	if err != nil {
		fmt.Println("SendPersonalMes json.Marshal fail := ", err.Error())
		return
	}
	mes.Data = string(data)
	//4.mes再次序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("SendPersonalMes json.Marshal fail := ", err.Error())
		return
	}
	//将序列化后的mes发送给服务器
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	//6.发送
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("SendPersonalMes err := ", err.Error())
		return
	}
	return
}

func (this *SmsProcess) NoticeExitMes() (err error) {
	var mes message.Message
	mes.Type = message.NoticeExitType
	var noticeExit message.NoticeExit
	noticeExit.UserId = CurUser.UserId
	data, err := json.Marshal(noticeExit)
	if err != nil {
		fmt.Println("NoticeExit json.Marshal fail := ", err.Error())
		return
	}
	mes.Data = string(data)
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("NoticeExit json.Marshal fail := ", err.Error())
		return
	}
	tf := &utils.Transfer{
		Conn: CurUser.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NoticeExit err := ", err.Error())
		return
	}
	return
}
