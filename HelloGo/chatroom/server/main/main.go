package main

import (
	"fmt"
	"hellogo/chatroom/server/model"
	"net"
	"time"
)

//处理与客户端的通讯
func process1(conn net.Conn) {
	//需要延时关闭conn
	defer conn.Close()
	//创建一个总控
	processor := &Precessor{
		Conn: conn,
	}
	err := processor.process2()
	if err != nil {
		fmt.Println("客户端和服务器之间的协程错误=", err)
		return
	}
}

// func Init() {
// 	initPool("localhost:6379", 16, 0, 300*time.Second)
// 	initUserDao()
// }

//编写一个函数，完成对UserDao的初始化任务
func initUserDao() {
	//pool本身就是一个全局的变量
	//初始化顺序先initpool再initUserDao
	model.GlobalUserDao = model.NewUserDao(pool)
}

func main() {
	//当服务器启动时就去初始化redis的连接池
	initPool("127.0.0.1:6379", 16, 0, 300*time.Second)
	initUserDao()
	//提示信息
	fmt.Println("服务器在8889端口监听....")
	listen, err := net.Listen("tcp", "0.0.0.0:8889")
	if err != nil {
		fmt.Println("net.listen err = ", err)
		return
	}
	defer listen.Close()
	//一旦监听成功，等待客户端来链接服务器
	for {
		fmt.Println("等待客户端来链接服务器")
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println("listen.Accept err = ", err)
		}
		//一旦链接成功则启动一个协程与客户端保持通讯。。
		go process1(conn)
	}
}
