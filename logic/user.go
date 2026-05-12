package logic

import (
	"bluebell/dao/mysql"
	"bluebell/models"
	"bluebell/pkg/jwt"
	"bluebell/pkg/snowflake"
)

//存放业务逻辑的代码

func SignUp(p *models.ParamSignUp) (err error) {
	// 1.判断用户是否存在
	err = mysql.ChackUserExist(p.Username)
	if err != nil {
		//数据库查询出错
		return err
	}
	// 2.生成UIO
	userID, _ := snowflake.GenID()
	//构造一个User实例
	user := &models.User{
		UserID:   int64(userID),
		Username: p.Username,
		Password: p.Password,
	}

	// 3.保存进数据库
	err = mysql.InsertUser(user)
	return
}

func Login(p *models.LoginSignUp) (token string, err error) {
	user := &models.User{
		Username: p.Username,
		Password: p.Password,
	}

	//传递的是指针，所以我们可以拿到user.UserID
	if err := mysql.Login(user); err != nil {
		return "", err
	}
	//生成JWT
	return jwt.GenToken(user.UserID, user.Username)
}
