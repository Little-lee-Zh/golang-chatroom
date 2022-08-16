package process

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"hellogo/chatroom/client/utils"
	"hellogo/chatroom/common/message"
	"net"
	_ "os"
)

type UserProcess struct {
}

func (this *UserProcess) Register(userId int, userPwd string, userSex string, userName string) (err error) {
	//链接到务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return
	}
	//延时关闭
	defer conn.Close()

	//2.准备过conn发送信息给服务器
	var mes message.Message
	mes.Type = message.RegisterMesType
	//3.创建一个Registermes结体
	var registerMes message.RegisterMes
	registerMes.User.UserId = userId
	registerMes.User.UserPwd = userPwd
	registerMes.User.UserName = userName
	registerMes.User.Sex = userSex
	//4.将registerMes 序列化
	data, err := json.Marshal(registerMes)
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
		fmt.Println("注册信息发送错误 err = ", err)
		return
	}
	mes, err = tf.ReadPkg() //mes就是RegisterResMes
	if err != nil {
		fmt.Println("readPkg err= ", err)
		return
	}

	//将mes的data部分反序列化成RegisterResMes
	var registerResMes message.RegisterResMes
	err = json.Unmarshal([]byte(mes.Data), &registerResMes)
	if err != nil {
		fmt.Println("json.Unmarshal err = ", err)
		return
	}

	if registerResMes.Code == 200 {
		fmt.Println("注册成功请重新登录")
		//os.Exit(0)
	} else {
		fmt.Println(registerResMes.Error)
		//os.Exit(0)
	}
	return
}

func (this *UserProcess) Login(userId int, userPwd string) (err error) {
	//定协议
	//1.链接到服务器
	conn, err := net.Dial("tcp", "127.0.0.1:8889")
	if err != nil {
		fmt.Println("net.Dial err = ", err)
		return
	}
	//延时关闭
	defer conn.Close()

	//2.准备通过conn发送信息服务器
	var mes message.Message
	mes.Type = message.LoginMesType
	//3.创建一个Loginmes结构体
	var loginMes message.LoginMes
	loginMes.UserId = userId
	loginMes.UserPwd = userPwd
	//.将loginMes 序列化
	data, err := json.Marshal(loginMes)
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
	//7.此时data就是要发送的消息
	//7.1先把data的长度发送给服务器
	//先获取data长度再转化成一个示长度的byte切片
	var pkgLen uint32
	pkgLen = uint32(len(data))
	var buf [4]byte
	binary.BigEndian.PutUint32(buf[0:4], pkgLen)
	//发送长度
	n, err := conn.Write(buf[:4])
	if n != 4 || err != nil {
		fmt.Println("conn.Write(bytes) err = ", err)
		return
	}

	fmt.Printf("客户端发送消息的长度=%d 内容=%s", len(data), string(data))

	//发送消息本身
	_, err = conn.Write(data)
	if err != nil {
		fmt.Println("conn.Write(data) err = ", err)
		return
	}
	//休眠20
	// time.Sleep(20 * time.Seond)
	// fmt.Printn("休眠了20...")
	//处理服务器返回的消息
	//创建一个transfer实例
	tf := &utils.Transfer{
		Conn: conn,
	}
	mes, err = tf.ReadPkg()
	if err != nil {
		fmt.Println("readPkg err = ", err)
		return
	}
	//将mes的data部分反序列化成LoginResMes
	var loginResMes message.LoginResMes
	err = json.Unmarshal([]byte(mes.Data), &loginResMes)
	if loginResMes.Code == 200 {
		//初始化CurUser
		CurUser.Conn = conn
		CurUser.UserId = userId
		CurUser.UserStatus = message.UserOnline

		// fmt.Println("登录成功")
		//显示当前用户列表
		//fmt.Println("当前用户列表如下：")
		for _, v := range loginResMes.UsersId {

			//如果要求不显示自己在线
			if v == userId {
				continue
			}

			//fmt.Println("用户id:\t", v)
			//完成客户端onlineUsers初始化
			user := &message.User{
				UserId:     v,
				UserStatus: message.UserOnline,
			}
			onlineUsers[v] = user
		}
		fmt.Print("\n\n")
		//还需要客户端启动一个协程
		//该协程保持与服务器的通.如果服务器有数据推送给客户端
		//接受信息并显示在客户端
		go serverProcessMes(conn)
		//1.显登录成功后的菜单
		for {
			ShowMenu()
		}
	} else {
		fmt.Println(loginResMes.Error)
	}
	return
}
