package main

import (
	"fmt"
	"hellogo/chatroom/client/process"
	"os"
)

//定义ID和密码
var userId int
var userPwd string
var userName string
var userSex string

func main() {
	var key int
	for true {
		fmt.Println("--------------------欢迎登陆多人聊天系统---------------------")
		fmt.Println("\t\t\t\t 1.登录系统")
		fmt.Println("\t\t\t\t 2.注册用户")
		fmt.Println("\t\t\t\t 3.退出系统")
		fmt.Println("\t\t\t\t 请选择: (1-3)")

		fmt.Scanf("%d\n", &key)
		switch key {
		case 1:
			fmt.Println("登录聊天室")
			fmt.Println("请输入用户id")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码")
			fmt.Scanf("%s\n", &userPwd)
			//1.创建一个UserProcess的实例
			up := &process.UserProcess{}
			up.Login(userId, userPwd)
		case 2:
			fmt.Println("注册用户")
			fmt.Println("请输入用户id:")
			fmt.Scanf("%d\n", &userId)
			fmt.Println("请输入用户密码:")
			fmt.Scanf("%s\n", &userPwd)
			fmt.Println("请选择用户性别: (male/female)")
			fmt.Scanf("%s\n", &userSex)
			fmt.Println("请输入用户昵称:")
			fmt.Scanf("%s\n", &userName)
			//调用Userprocess完成注册的请求
			up := &process.UserProcess{}
			up.Register(userId, userPwd, userSex, userName)
		case 3:
			fmt.Println("退出系统")
			os.Exit(0)
		default:
			fmt.Println("输入有误,请重新输入")
		}
	}
}
