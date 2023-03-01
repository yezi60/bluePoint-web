package forms

// 前端的传来的校验码，信息 用于注册和登陆
type SendMmsForm struct {
	Mobile string `form:"mobile" json:"mobile" binding:"required,mobile"` // 手机号码格式有规范可寻，自定义validation
	Type   uint   `form:"type" json:"type" binding:"required,oneof=1 2"`  // 1:注册 2:动态验证码登陆
	// 注册发送短信验证码 和 动态验证码登陆发送验证码（没有用户创建用户）
}
