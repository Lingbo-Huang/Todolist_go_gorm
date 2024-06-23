package serializer

import "todo_list/model"

type User struct {
	ID       uint   `json:"id" form:"id" example:"1"`                      // 用户ID
	UserName string `json:"userName" form:"user_name" example:"CaiXiaobo"` // 用户名
	Status   string `json:"status" form:"status"`                          // 用户状态
	CreateAt int64  `json:"create_at" form:"create_at"`                    // 创建
}

func BuildUser(user model.User) User {
	return User{
		ID:       user.ID,
		UserName: user.UserName,
		CreateAt: user.CreatedAt.Unix(),
	}
}
