package service

import (
	"scarlet/common"
)

type UserService struct{}

func NewUserService() *UserService {
	return new(UserService)
}

// TODO:不知道是否能关联地查出server？

func (this *UserService) GetUserByID(id uint) common.User {
	user := common.User{}
	mysqlClient.Find(&user, id)
	return user
}

//func (this * UserService)GetUserEmail()common.User{
//
//}
//
//func (this * UserService)UpdateUser()error{
//
//}
//
//func (this * UserService)GetServers()[]common.Server{
//
//}

func (this *UserService) AddUser(user common.User) {
	// 在业务逻辑层要做好检查
	mysqlClient.Create(&user)
	return
}

func (this *UserService) UpdateUserServers(user common.User) {
	old := common.User{}
	mysqlClient.First(&old, user.ID)
	old.Servers = user.Servers
	mysqlClient.Save(&old)
	return
}

func (this *UserService) Auth(user common.User) (common.User, bool) {
	var num int
	mysqlClient.Table("users").Where("password = ? and email = ?", user.Password, user.Email).First(&user).Count(&num)
	if num >= 1 {
		return user, true
	} else {
		return common.User{}, false
	}
}

func (this *UserService) Exist(user common.User) bool {
	var num int
	mysqlClient.Table("users").Where("email = ?", user.Email).Count(&num)
	if num >= 1 {
		return true
	} else {
		return false
	}
}
