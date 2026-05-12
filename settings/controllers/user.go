package controllers

import (
	"bluebell/dao/mysql"
	"bluebell/logic"
	"bluebell/models"
	"errors"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
)

// SignUpHandler 处理注册请求函数
func SignUpHandler(c *gin.Context) {
	// 1.获取参数和参数校验
	// var p models.ParamSignUp
	//这样创建对性能更友好,且后续传值就行，不需要传指针
	// 等价于p := &models.ParamSignUp{}
	p := new(models.ParamSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("SignUp With Invalid Param", zap.Error(err))
		//判断err是不是validator.ValidationError类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			// c.JSON(http.StatusOK, gin.H{
			// 	"msg": err.Error(),
			// })
			// 通过改进封装代码实现返回错误响应
			ResponseError(c, CodeInvalidParams)
			return
		}
		// c.JSON(http.StatusOK, gin.H{
		// 	//请求参数有误直接返回响应
		// 	"msg": RemoveTopStruct(errs.Translate(trans)), //翻译错误
		// })
		ResponseErrorWithMsg(c, CodeInvalidParams, RemoveTopStruct(errs.Translate(trans)))
		return
	}
	//通过validator这个库可以实现自动业务规则校验的功能（上面）

	//手动对请求参数进行详细的业务规则校验
	// if len(p.Username) == 0 || len(p.Password) == 0 || len(p.Repassword) == 0 || p.Password != p.Repassword {
	// 	zap.L().Error("SignUp With Invalid Param")
	// 	c.JSON(http.StatusOK, gin.H{
	// 		//请求参数有误直接返回响应
	// 		"msg": "请求参数有误",
	// 	})
	// 	return
	// }

	fmt.Println(p)

	// 2.业务处理
	if err := logic.SignUp(p); err != nil {
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": "注册失败噜",
		// })
		if errors.Is(err, mysql.ErrorUserExist) {
			ResponseError(c, CodeUserExist)
			return
		}
		ResponseError(c, CodeServerBusy)
		return
	}
	// 3.返回响应
	// c.JSON(http.StatusOK, gin.H{
	// 	"msg": "success",
	// })
	ResponseSuccess(c, nil)
}

func LoginHandler(c *gin.Context) {
	// 1.获取参数以及参数校验
	p := new(models.LoginSignUp)
	if err := c.ShouldBindJSON(p); err != nil {
		zap.L().Error("Login With Invalid Param", zap.Error(err))
		//判断err是不是validator.ValidationError类型
		errs, ok := err.(validator.ValidationErrors)
		if !ok {
			ResponseError(c, CodeInvalidParams)
			return
		}
		ResponseErrorWithMsg(c, CodeInvalidParams, RemoveTopStruct(errs.Translate(trans)))
		// c.JSON(http.StatusOK, gin.H{
		// 	//请求参数有误直接返回响应
		// 	"msg": RemoveTopStruct(errs.Translate(trans)), //翻译错误
		// })
		return
	}
	// 2.业务逻辑处理
	token, err := logic.Login(p)
	if err != nil {
		zap.L().Error("Login failed", zap.String("username", p.Username), zap.Error(err))
		// c.JSON(http.StatusOK, gin.H{
		// 	"msg": "登录失败噜，用户名或者密码错误",
		// })
		if errors.Is(err, mysql.ErrorUserNotExist) {
			ResponseError(c, CodeUserNotExist)
			return
		}
		ResponseError(c, CodeInvalidPassword)
		return
	}
	// 3. 返回响应
	// c.JSON(http.StatusOK, gin.H{
	// 	"msg": "登陆成功",
	// })
	ResponseSuccess(c, token)
}
