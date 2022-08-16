package model

// 定义用户结构体
type User struct {
	// 为了系列化和反序列化成功
	// 用户信息的json字符串的key 要和 结构体的tag名字一致
	UserId int `json:"userId"`
	UserPwd string `json:"userPwd"`
	UserName string `json:"userName"`
}