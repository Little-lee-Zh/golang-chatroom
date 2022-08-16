package process

import (
	"encoding/json"
	"fmt"
	"hellogo/chatroom/common/message"
	"hellogo/chatroom/server/model"
	"hellogo/chatroom/server/utils"
	"net"
)

type PassProcess struct {
	Conn net.Conn
	//增加字段表示该Conn是那个用户
	UserId int
}

//编写一个severProcessCheck函数专门处理登录
func (this *PassProcess) SeverProcessCheckPass(mes *message.Message) (err error) {
	//1.先从mes中取出mes.data直接反序列化成CheckPass
	var checkPass message.CheckPass
	err = json.Unmarshal([]byte(mes.Data), &checkPass)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}

	//1.先声明一个resMes
	var resMes message.Message
	resMes.Type = message.CheckPassResMesType

	//2.再声明一个LoginResMes并完成赋值
	var checkPassResMes message.CheckPassResMes

	//需要到redis数据库去完成验证
	//1.使用model.MyUserDao到redis去验证
	_, err = model.GlobalUserDao.CheckPass(checkPass.User.UserId, checkPass.OldPassport)

	if err != nil {
		if err == model.ERROR_CHECK_PWD {
			checkPassResMes.Code = 500
			checkPassResMes.Error = err.Error()
		} else {
			checkPassResMes.Code = 505
			checkPassResMes.Error = "服务器内部错误..."
		}
	} else {
		checkPassResMes.Code = 200
		fmt.Println("密码校验成功！")
	}

	//3.将checkPassResMes序列化
	data, err := json.Marshal(checkPassResMes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}
	//4.将data赋值给resMes
	resMes.Data = string(data)
	//5.对resMes进行序列化准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}
	//6.发送将其封装到writePkg函数
	//因为使用分层模式（mvc），先创建一个Transfer实例然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}

func (this *PassProcess) SeverProcessChangePass(mes *message.Message) (err error) {
	//1.先从mes中取出mes.data直接反序列化成ChangePass
	var changePass message.ChangePass
	err = json.Unmarshal([]byte(mes.Data), &changePass)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}

	//1.先声明一个resMes
	var resMes message.Message
	resMes.Type = message.ChangePassResMesType

	//2.再声明一个ChangePassResMes并完成赋值
	var changePassResMes message.ChangePassResMes

	//需要到redis数据库去完成验证
	//1.使用model.MyUserDao到redis去验证
	err = model.GlobalUserDao.ChangePass(changePass.User.UserId, changePass.NewPassport)

	if err != nil {
		changePassResMes.Code = 505
		changePassResMes.Error = "服务器内部错误..."
	} else {
		changePassResMes.Code = 200
		fmt.Println("密码修改成功！")
	}

	//3.将changePassResMes序列化
	data, err := json.Marshal(changePassResMes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}
	//4.将data赋值给resMes
	resMes.Data = string(data)
	//5.对resMes进行序列化准备发送
	data, err = json.Marshal(resMes)
	if err != nil {
		fmt.Println("json.Marshal err = ", err)
		return
	}
	//6.发送将其封装到writePkg函数
	//因为使用分层模式（mvc），先创建一个Transfer实例然后读取
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	return
}
