package process

import (
	"encoding/json"
	"fmt"
	"hellogo/chatroom/common/message"
)

func outputGroupMes(mes *message.Message) {
	//显示出来即可
	//1.反序列化mes.Data
	var smsMes message.SmsMes
	err := json.Unmarshal([]byte(mes.Data), &smsMes)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err.Error())
		return
	}
	//显示信息
	info := fmt.Sprintf("用户id:%d对大家说:%s", smsMes.UserId, smsMes.Content)
	fmt.Println(info)
	fmt.Println()
}

func outputPersonalMes(mes *message.Message) {
	var smsPerMes message.SmsPerMes
	err := json.Unmarshal([]byte(mes.Data), &smsPerMes)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err.Error())
		return
	}
	//显示信息
	info := fmt.Sprintf("用户id:%d对你说:%s", smsPerMes.UserId, smsPerMes.Content)
	fmt.Println(info)
	fmt.Println()
}
