package model

import (
	"encoding/json"
	"fmt"
	"hellogo/chatroom/common/message"

	"github.com/garyburd/redigo/redis"
)

// 在服务器启动后，就初始化全局UserDao实例，在需要和redis交互时直接使用
// 避免每次和redis交互都创建UserDao实例，影响性能
var (
	GlobalUserDao *UserDao
)

// 定义 UserDao结构体，完成对 user的各种操作
type UserDao struct {
	pool *redis.Pool
}

// 工厂函数：使用工厂模式，创建UserDao实例
func NewUserDao(pool *redis.Pool) (userDao *UserDao) {
	// main包初始化的时候，就把pool创建好
	userDao = &UserDao{
		pool: pool,
	}
	return
}

// 提供一些方法
// 1. 根据用户Id返回一个user实例，err
func (this *UserDao) getUserById(conn redis.Conn, id int) (user *User, err error) {

	// 通过给定id，去redis查询用户
	res, err := redis.String(conn.Do("HGET", "users", id)) //butong
	if err != nil {
		if err == redis.ErrNil {
			// 这个错误表示在users的hash中没有找到对应的id
			err = ERROR_USER_NOTEXISTS
		}
		return
	}

	user = &User{}

	// 将res反序列化为User实例对象
	err = json.Unmarshal([]byte(res), &user)
	if err != nil {
		fmt.Println("json.Unmarshal err=", err)
		return
	}

	return
}

// 2. 登录的校验 Login
// 如果用户的Id和密码都正确，则返回user实例
// 如果用户的Id和密码有错误，则返回错误信息
func (this *UserDao) Login(userId int, userPwd string) (user *User, err error) {

	// 先从UserDao的连接池取出连接
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.getUserById(conn, userId)
	if err != nil {
		return
	}

	// 这里表示用户已经找到了，接下来验证密码
	if user.UserPwd != userPwd {
		err = ERROR_USER_PWD
		return
	}
	return
}

// 注册用户
func (this *UserDao) Register(user *message.User) (err error) {

	// 先从UserDao的连接池取出连接
	conn := this.pool.Get()
	defer conn.Close()
	_, err = this.getUserById(conn, user.UserId)
	if err == nil { //bug location
		err = ERROR_USER_EXISTS
		return
	}

	// 这里表示用户不存在，所以可以进行注册
	// 序列化
	data, err := json.Marshal(user)
	if err != nil {
		return
	}
	// 入库
	_, err = conn.Do("hset", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存用户信息错误, err=", err)
		return
	}
	return
}

//检查密码
func (this *UserDao) CheckPass(userId int, olduserPwd string) (user *User, err error) {
	// 先从UserDao的连接池取出连接
	conn := this.pool.Get()
	defer conn.Close()
	user, err = this.getUserById(conn, userId)
	if err != nil {
		return
	}

	// 这里表示用户已经找到了，接下来验证密码
	if user.UserPwd != olduserPwd {
		err = ERROR_CHECK_PWD
		return
	}
	return
}

//修改密码
func (this *UserDao) ChangePass(userid int, newpassport string) (err error) {
	// 先从UserDao的连接池取出连接
	conn := this.pool.Get()
	defer conn.Close()
	user, err := this.getUserById(conn, userid)
	if err != nil {
		return
	}

	user.UserPwd = newpassport

	// 修改密码
	// 序列化
	data, err := json.Marshal(user)
	if err != nil {
		return
	}
	// 入库
	_, err = conn.Do("hset", "users", user.UserId, string(data))
	if err != nil {
		fmt.Println("保存用户信息错误, err=", err)
		return
	}
	return
}
