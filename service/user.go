package service

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"todo_list/model"
	"todo_list/pkg/utils"
	"todo_list/serializer"
)

type UserService struct {
	UserName string `form:"user_name" json:"user_name" binding:"required,min=3,max=15"`
	Password string `form:"password" json:"password" binding:"required,min=5,max=16"`
}

func (service *UserService) Register() serializer.Response {
	var user model.User
	var count int
	model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).
		First(&user).Count(&count)
	if count == 1 {
		return serializer.Response{
			Status: 400,
			Msg:    "该用户已存在，请勿重复注册",
		}
	}
	user.UserName = service.UserName
	// 密码加密
	err := user.SetPassword(service.Password)
	if err != nil {
		return serializer.Response{
			Status: 400,
			Msg:    fmt.Sprintf("密码设置失败, %s\n", err.Error()),
		}
	}

	// 创建用户
	if err := model.DB.Create(&user).Error; err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "数据库创建用户失败",
		}
	}

	return serializer.Response{
		Status: 200,
		Msg:    "用户注册成功",
	}
}
func (service *UserService) Login() serializer.Response {
	var user model.User
	err := model.DB.Model(&model.User{}).Where("user_name=?", service.UserName).
		First(&user).Error
	if err != nil {
		if gorm.IsRecordNotFoundError(err) {
			return serializer.Response{
				Status: 400,
				Msg:    "用户不存在，需要先注册",
			}
		}
		return serializer.Response{
			Status: 500,
			Msg:    fmt.Sprintf("数据库错误,%s\n", err.Error()),
		}
	}

	// 验证密码
	if user.CheckPassword(service.Password) == false {
		return serializer.Response{
			Status: 400,
			Msg:    "密码错误",
		}
	}

	// 发一个token，为了其它功能需要身份验证，给前端存储的
	token, err := utils.GenerateToken(user.ID, service.UserName, service.Password)
	if err != nil {
		return serializer.Response{
			Status: 500,
			Msg:    "Token 生成错误",
		}
	}

	return serializer.Response{
		Status: 200,
		Data: serializer.TokenData{
			User:  serializer.BuildUser(user),
			Token: token,
		},
		Msg: "用户登录成功",
	}
}
