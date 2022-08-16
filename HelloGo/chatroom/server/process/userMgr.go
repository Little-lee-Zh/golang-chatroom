package process

import (
	"fmt"
	"hellogo/chatroom/common/message"
)

//UserMgr在服务器端有且只有一个 将其定义为全局变量
var (
	userMgr *UserMgr
)

type UserMgr struct {
	onlineUsers  map[int]*UserProcess
	offlineUsers map[int][]*message.Message
}

//完成对userMgr初始化工作
func init() {
	userMgr = &UserMgr{
		onlineUsers:  make(map[int]*UserProcess, 1024),
		offlineUsers: make(map[int][]*message.Message, 1024),
	}
}

//完成对onlineUsers添加
func (this *UserMgr) AddOnlineUsers(up *UserProcess) {
	this.onlineUsers[up.UserId] = up
}

//完成对onlineUsers删除
func (this *UserMgr) DelOnlineUsers(userId int) {
	delete(this.onlineUsers, userId)
}

//完成对offlineUsers添加离线群聊消息
func (this *UserMgr) AddofflineMes(mes *message.Message) {
	fmt.Println(mes)

	for id := range userMgr.offlineUsers {
		userMgr.offlineUsers[id] = append(userMgr.offlineUsers[id], mes)
	}

	fmt.Println(userMgr.offlineUsers)
}

//完成对offlineUsers添加离线私聊消息
func (this *UserMgr) AddofflinePerMes(mes *message.Message, userId int) {
	fmt.Println(mes)
	for id := range userMgr.offlineUsers {
		if id == userId {
			userMgr.offlineUsers[id] = append(userMgr.offlineUsers[id], mes)
			fmt.Println(userMgr.offlineUsers[id])
		}
	}
}

//返回当前所有在线用户
func (this *UserMgr) GetAllOnlineUsers() map[int]*UserProcess {
	return this.onlineUsers
}

//根据iD返回对应的值
func (this *UserMgr) GetOnlineUserById(userId int) (up *UserProcess, err error) {
	up, ok := this.onlineUsers[userId]
	if !ok { //当前不在线
		err = fmt.Errorf("用户%d不存在", userId)
		return
	}
	return
}
