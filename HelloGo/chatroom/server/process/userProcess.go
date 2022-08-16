package process

import (
	"encoding/json"
	"fmt"
	"hellogo/chatroom/common/message"
	"hellogo/chatroom/server/model"
	"hellogo/chatroom/server/utils"
	"net"
)

type UserProcess struct {
	Conn net.Conn
	//增加字段表示该Conn是那个用户
	UserId int
}

//通知所有在线用户的方法
//userID要通知其他的在线
func (this *UserProcess) NotifyOthersOnlineUser(userId int) {
	//遍历OnlineUsers，然后一个个的发送
	for id, up := range userMgr.onlineUsers {
		//过滤自己
		if id == userId {
			continue
		}
		up.NotifyMeOnline(userId)
	}
}

func (this *UserProcess) NotifyMeOnline(userId int) {
	//组装NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType

	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOnline

	//将notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//将序列化后的notifyUserStatusMes赋值给mes.Data
	mes.Data = string(data)
	//对mes再次序列化，准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//发送，创建Transfer示例
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOnline err=", err)
		return
	}
}

func (this *UserProcess) Offline(mes *message.Message) (err error) {
	var noticeExit message.NoticeExit
	err = json.Unmarshal([]byte(mes.Data), &noticeExit)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}
	userId := noticeExit.UserId
	if _, ok := userMgr.offlineUsers[userId]; !ok {
		userMgr.offlineUsers[userId] = make([]*message.Message, 0, 1024)
	}
	userMgr.DelOnlineUsers(userId)

	fmt.Println(userMgr)

	for id, up := range userMgr.onlineUsers {
		//过滤自己
		if id == userId {
			continue
		}
		up.NotifyMeOffline(userId)
	}
	return
}

func (this *UserProcess) NotifyMeOffline(userId int) (err error) {
	//组装NotifyUserStatusMes
	var mes message.Message
	mes.Type = message.NotifyUserStatusMesType
	var notifyUserStatusMes message.NotifyUserStatusMes
	notifyUserStatusMes.UserId = userId
	notifyUserStatusMes.Status = message.UserOffline

	//将notifyUserStatusMes序列化
	data, err := json.Marshal(notifyUserStatusMes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//将序列化后的notifyUserStatusMes赋值给mes.Data
	mes.Data = string(data)
	//对mes再次序列化，准备发送
	data, err = json.Marshal(mes)
	if err != nil {
		fmt.Println("json.Marshal err=", err)
		return
	}
	//发送，创建Transfer示例
	tf := &utils.Transfer{
		Conn: this.Conn,
	}
	err = tf.WritePkg(data)
	if err != nil {
		fmt.Println("NotifyMeOffline err=", err)
		return
	}
	return
}

func (this *UserProcess) SeverProcessRegister(mes *message.Message) (err error) {
	//1.先从mes中取出mes.data直接反序列化成RegisterMes
	var registerMes message.RegisterMes
	err = json.Unmarshal([]byte(mes.Data), &registerMes)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}
	//1.先声明一个resMes
	var resMes message.Message
	resMes.Type = message.RegisterResMesType
	var registerResMes message.RegisterResMes

	//需要到redis数据库去完成注册
	//1.使用model.MyUserDao到redis去验证
	err = model.GlobalUserDao.Register(&registerMes.User)
	if err != nil {
		if err == model.ERROR_USER_EXISTS {
			registerResMes.Code = 400
			registerResMes.Error = model.ERROR_USER_EXISTS.Error()
		} else {
			registerResMes.Code = 509
			registerResMes.Error = "注册发生未知错误..."
		}
	} else {
		registerResMes.Code = 200
	}

	data, err := json.Marshal(registerResMes)
	if err != nil {
		fmt.Println("json.Marshal fail", err)
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

//编写一个severProcessLogin函数专门处理登录
func (this *UserProcess) SeverProcessLogin(mes *message.Message) (err error) {
	//1.先从mes中取出mes.data直接反序列化成LoginMes
	var loginMes message.LoginMes
	err = json.Unmarshal([]byte(mes.Data), &loginMes)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}

	//1.先声明一个resMes
	var resMes message.Message
	resMes.Type = message.LoginResMesType

	//2.再声明一个LoginResMes并完成赋值
	var loginResMes message.LoginResMes

	//需要到redis数据库去完成验证
	//1.使用model.MyUserDao到redis去验证
	user, err := model.GlobalUserDao.Login(loginMes.UserId, loginMes.UserPwd)

	if err != nil {
		if err == model.ERROR_USER_NOTEXISTS {
			loginResMes.Code = 500
			loginResMes.Error = err.Error()
		} else if err == model.ERROR_USER_PWD {
			loginResMes.Code = 403
			loginResMes.Error = err.Error()
		} else {
			loginResMes.Code = 505
			loginResMes.Error = "服务器内部错误..."
		}
	} else {
		loginResMes.Code = 200
		//登录成功把登录成功的用户放到userMgr
		//将登录成功的用户userId赋给this
		this.UserId = loginMes.UserId
		userMgr.AddOnlineUsers(this)

		//通知其他的在线用户，上线了
		this.NotifyOthersOnlineUser(loginMes.UserId)
		//将当前在线用户id放入loginResMes.UserIds中
		for id := range userMgr.onlineUsers {
			loginResMes.UsersId = append(loginResMes.UsersId, id)
		}

		defer func() {
			smsProcess := &SmsProcess{}
			if _, ok := userMgr.offlineUsers[this.UserId]; ok {
				for _, offlineMes := range userMgr.offlineUsers[this.UserId] {
					smsProcess.SendofflineMes(offlineMes, this.Conn)
				}
				delete(userMgr.offlineUsers, this.UserId)
			}
		}()

		fmt.Println(user, "登录成功")
	}

	//3.将loginResMes序列化
	data, err := json.Marshal(loginResMes)
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
