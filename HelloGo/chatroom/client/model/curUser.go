package model

import (
	"hellogo/chatroom/common/message"
	"net"
)

//在客户端很多地方用到CurUser，将其设为全局

type CurUser struct {
	Conn net.Conn
	message.User
}
