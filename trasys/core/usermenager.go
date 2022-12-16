package core

/*
* 这是用户的管理模块，主要是做uid到用户的映射
 */

type UserMgr struct {
	users map[int32]User
}

// //创建一个用户管理模块
// func NewUserMgr() *UserMgr {
// 	return &UserMgr{
// 		users: make(map[int32]user),
// 	}
// }

// 向用户管理模块中增加一个新的连接用户
func (um *UserMgr) AddUser(user User) {
	um.users[user.GetId()] = user
}

// 使用Uid查询一个用户
func (um *UserMgr) GetUserByUid(uid int32) User {
	return um.users[uid]
}

// 全局变量的用户管理模块
var GlobalUserMgr *UserMgr

// 初始化方法
func init() {
	GlobalUserMgr = &UserMgr{
		users: make(map[int32]User),
	}
}
