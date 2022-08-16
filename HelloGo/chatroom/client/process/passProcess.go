package process

import (
	"encoding/json"
	"fmt"
	"hellogo/chatroom/client/utils"
	"hellogo/chatroom/common/message"
	"net"
)

type PassProcess struct {
}

func (this *PassProcess) CheckPassport(oldPassport string) (result bool, err error) {
	//链接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return
	}
	//延时关闭
	defer conn.Close()

	//2.准备过conn发送信息给服务器
	var mes message.Message
	mes.Type = message.CheckPassType
	//3.创建一个checkPass结体
	var checkPass message.CheckPass
	checkPass.User.UserId = CurUser.UserId
	checkPass.OldPassport = oldPassport
	//4.将checkPass序列化
	data, err := json.Marshal(checkPass)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}
	//5.把data赋给mes.Data字段
	mes.Data = string(data)
	//6.将mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}
	// 创建 一个Tansfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	//发送data给服务器端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("原密码信息发送错误 err = ", err)
		return
	}
	mes, err = tf.ReadPkg() //mes就是checkPass返回的信息
	if err != nil {
		fmt.Println("readPkg err= ", err)
		return
	}

	//将mes的data部分反序列化成RegisterResMes
	var checkPassResMes message.CheckPassResMes
	err = json.Unmarshal([]byte(mes.Data), &checkPassResMes)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}

	if checkPassResMes.Code == 200 {
		fmt.Println("原密码输入正确！")
		result = true
	} else {
		fmt.Println(checkPassResMes.Error)
		result = false
	}
	return
}

func (this *PassProcess) ChangePassport(newPassport string) (err error) {
	//链接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return
	}
	//延时关闭
	defer conn.Close()

	//2.准备过conn发送信息给服务器
	var mes message.Message
	mes.Type = message.ChangePassType
	//3.创建一个changePass结体
	var changePass message.ChangePass
	changePass.User.UserId = CurUser.UserId
	changePass.NewPassport = newPassport
	//4.将changePass序列化
	data, err := json.Marshal(changePass)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}
	//5.把data赋给mes.Data字段
	mes.Data = string(data)
	//6.将mes序列化
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}
	// 创建 一个Tansfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	//发送data给服务器端
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("新密码信息发送错误 err = ", err)
		return
	}
	mes, err = tf.ReadPkg() //mes就是changePass返回的信息
	if err != nil {
		fmt.Println("readPkg err= ", err)
		return
	}

	//将mes的data部分反序列化成ChangePassResMes
	var changePassResMes message.ChangePassResMes
	err = json.Unmarshal([]byte(mes.Data), &changePassResMes)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}

	if changePassResMes.Code == 200 {
		fmt.Println("密码修改成功！")
	} else {
		fmt.Println("密码修改失败！")
	}
	return
}
