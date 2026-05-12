package mysql

import (
	"bluebell/models"
	"crypto/md5"
	"database/sql"
	"encoding/hex"
	"errors"
)

// 把每一步数据库操作封装成函数
// 待logic层根据业务需求调用

var (
	ErrorUserExist       = errors.New("用户名已存在")
	ErrorUserNotExist    = errors.New("用户名不存在")
	ErrorInvalidPassword = errors.New("用户名或密码错误")
	ErrorInvalidID       = errors.New("无效的ID")
)

const secret = "42"

// ChackUserExist 检查指定用户的用户名是否存在
func ChackUserExist(username string) (err error) {
	sqlstr := `select count(user_id) from user where username = ?`
	var count int
	if err := db.Get(&count, sqlstr, username); err != nil {
		return err
	}
	if count > 0 {
		return ErrorUserExist
	}
	return
}

// InsertUser 向数据库中插入一条新的用户记录
func InsertUser(user *models.User) (err error) {
	//对密码进行加密
	user.Password = encryptPassword(user.Password)
	//执行SQL入库
	sqlstr := `insert into user(user_id,username,password) values (?,?,?)`
	_, err = db.Exec(sqlstr, user.UserID, user.Username, user.Password)
	return
}

// encryptPassword 将数据库得密码加密
func encryptPassword(oPassword string) string {
	h := md5.New()
	h.Write([]byte(secret))
	return hex.EncodeToString(h.Sum([]byte(oPassword)))
}

// Login登录
func Login(user *models.User) (err error) {
	oPassword := user.Password // 用户登录得密码
	sqlstr := `select user_id,username,password from user where username = ?`
	if err = db.Get(user, sqlstr, user.Username); err != nil {
		//查询数据库失败
		return err
	}
	if err == sql.ErrNoRows {
		return ErrorUserNotExist
	}
	password := encryptPassword(oPassword) // 数据库里得到得密码
	if password != user.Password {
		return ErrorInvalidPassword
	}
	return
}

// GetUserBuId 根据用户ID查询用户信息
func GetUserBuId(uid int64) (user *models.User, err error) {
	user = new(models.User)
	sqlstr := `select user_id, username from user where user_id = ?`
	db.Get(user, sqlstr, uid)
	return
}
