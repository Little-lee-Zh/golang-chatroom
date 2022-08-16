package process

import (
	"encoding/json"
	"fmt"
	"hellogo/chatroom/common/message"
	"hellogo/chatroom/server/utils"
	"net"
)

type SmsProcess struct {
	//暂时不需要字段
}

//写方法转发消息
func (this *SmsProcess) SendGroupMes(mes *message.Message) {
	userMgr.AddofflineMes(mes)
	//遍历服务器端onlineUsers map[int]*userProcess
	//将消息转发出去
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}

	for id, up := range userMgr.onlineUsers {
		//过滤自己
		if id == smsMes.UserId {
			continue
		}
		this.SendMesToOnlineUsers(data, up.Conn)
	}
}

func (this *SmsProcess) SendMesToOnlineUsers(data []byte, conn net.Conn) {
	//创建一个Transfer实例，发送data
	tf := &utils.Transfer{
		Conn: conn,
	}
	err := tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err", err)
	}
}

func (this *SmsProcess) SendPersonalMes(mes *message.Message) {
	var smsPerMes message.SmsPerMes
	err := json.Unmarshal([]byte(mes.Data), &smsPerMes)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}

	userMgr.AddofflinePerMes(mes, smsPerMes.SendToUser)

	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}

	for id, up := range userMgr.onlineUsers {
		//过滤自己
		if id == smsPerMes.SendToUser {
			this.SendMesToOnlineUsers(data, up.Conn)
		}
	}
}

func (this *SmsProcess) SendofflineMes(mes *message.Message, conn net.Conn) {
	//当前用户，直接可以发送
	data, err := json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}
	//创建一个Transfer实例，发送data
	tf := &utils.Transfer{
		Conn: conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("转发消息失败 err", err)
	}
}
