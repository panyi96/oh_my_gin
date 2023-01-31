package pojo

import (
	"ohmygin/dbconfig"
)

type User struct {
	Id       uint    `json:"userId" gorm:"primaryKey" binding:"required"`
	Name     string  `json:"userName" binding:"required,max=20,min=1"`
	Email    *string `json:"userEmail" binding:"required,email"`
	Password *string `json:"password" binding:"required,max=20,min=1,validPw"`
}

type Users struct {
	UserList     []User `json:"userList" binding:"required,gte=2,lte=5"`
	UserListSize int    `json:"userListSize"`
}

func FindAllUser() ([]User, error) {
	var userList []User
	result := dbconfig.DBConnect.Table("user").Find(&userList)
	return userList, result.Error
}

func FindAllUserById(userId uint64) (User, error) {
	var user User
	result := dbconfig.DBConnect.Table("user").Where("id = ? ", userId).Find(&user)
	return user, result.Error
}

func AddUser(user User) (User, error) {
	result := dbconfig.DBConnect.Table("user").Save(&user)
	return user, result.Error
}

func AddUsers(users Users) (Users, error) {
	result := dbconfig.DBConnect.Table("user").CreateInBatches(&(users.UserList), users.UserListSize)
	return users, result.Error
}

func UpdateUser(user User) (User, error) {
	result := dbconfig.DBConnect.Table("user").Where("id = ?", user.Id).Updates(&user)
	return user, result.Error
}

func DelUser(id int) (int64, error) {
	result := dbconfig.DBConnect.Table("user").Delete(&User{}, id)
	return result.RowsAffected, result.Error
}

func GetUserByPw(userName string, password string) (User, error) {
	var user User
	result := dbconfig.DBConnect.Table("user").
		Where("name = ? and password = ?", userName, password).
		Find(&user)
	return user, result.Error
}
