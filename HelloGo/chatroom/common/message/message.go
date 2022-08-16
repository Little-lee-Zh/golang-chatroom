package message

const (
	LoginMesType            = "LoginMes"
	LoginResMesType         = "LoginResMes"
	RegisterMesType         = "RegisterMes"
	RegisterResMesType      = "RegisterResMes"
	NotifyUserStatusMesType = "NotifyUserStatusMes"
	SmsMesType              = "SmsMes"    //群聊
	SmsPerMesType           = "SmsPerMes" //私聊
	CheckPassType           = "CheckPass" //检查密码
	CheckPassResMesType     = "CheckPassResMes"
	ChangePassType          = "ChangePass" //改密码
	ChangePassResMesType    = "ChangePassResMes"
	NoticeExitType          = "NoticeExit"
)

//定义用户状态的常量
const (
	UserOnline = iota
	UserOffline
	UserBusyStatus
)

type Message struct {
	Type string `json:"type"`
	Data string `json:"data"`
}

//定义两个消息
type LoginMes struct {
	UserId   int    `json:"userId"`
	UserPwd  string `json:"userPwd"`
	UserName string `json:"userName"`
}

type LoginResMes struct {
	Code    int    `json:"code"` //状态码 500表示未注册 200表示登录成功
	UsersId []int  //增加字段保存用户id 的切片
	Error   string `json:"error"`
}

type RegisterMes struct {
	User User `json:"user"` //类型就是User结构体
}

type RegisterResMes struct {
	Code  int    `json:"code"` //状态码400用户已经占用 200表示注册成功
	Error string `json:"error"`
}

//为了配合服务器端推送用户状态变化的消息
type NotifyUserStatusMes struct {
	UserId int `json:"userId"`
	Status int `json:"status"`
}

//增加一个SmsMes//发送的消息
type SmsMes struct {
	Content string `json:"content"`
	User           //匿名结构体，继承
}

type SmsPerMes struct {
	Content    string `json:"content"`
	SendToUser int    `json:"sendToUser"`
	User
}

type CheckPass struct {
	OldPassport string `json:"oldpassport"`
	User
}

type CheckPassResMes struct {
	Code  int    `json:"code"` //200表示原密码输入正确
	Error string `json:"error"`
}

type ChangePass struct {
	NewPassport string `json:"newpassport"`
	User
}

type ChangePassResMes struct {
	Code  int    `json:"code"` //200表示密码修改成功
	Error string `json:"error"`
}

type NoticeExit struct {
	User
}
