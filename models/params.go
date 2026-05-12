package models

//定义请求参数的结构体

const (
	OrderTime  = "time"
	OrderScore = "score"
)

// ParamSignUp 注册请求参数
type ParamSignUp struct {
	Username   string `json:"username" binding:"required"` //添加这个binding之后就不需要手动校验请求参数业务规则了
	Password   string `json:"password" binding:"required"`
	RePassword string `json:"re_password" binding:"required,eqfield=Password"`
}

// LoginSignUp 登录请求参数
type LoginSignUp struct {
	Username string `json:"username" binding:"required"` //添加这个binding之后就不需要手动校验请求参数业务规则了
	Password string `json:"password" binding:"required"`
}

// ParamVoteData 投票请求参数
type ParamVoteData struct {
	// UserID 从请求中获取当前用户
	PostID    int64 `json:"post_id" binding:"required"`              //帖子ID
	Direction int8  `json:"direction,string" binding:"oneof=1 0 -1"` // 赞成票(1) 反对票(-1) 取消投票(0)
}

type ParamPostList struct {
	Page  int64  `form:"page" json:"page"`
	Size  int64  `form:"size" json:"size"`
	Order string `form:"order" json:"order"`
}
