package forms

// 前端传的Form表单，用于登陆校验 带验证码图片
type PasswordLoginForm struct {
	Mobile    string `form:"mobile" json:"mobile" binding:"required,mobile"` // 手机号码格式有规范可寻，自定义validation
	PassWord  string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Captcha   string `form:"captcha" json:"captcha" binding:"required,min=5,max=5"`
	CaptchaId string `form:"captcha_id" json:"captcha_id" binding:"required"`
}

// 前端传的Form表单，用于注册
type RegisterForm struct {
	Mobile   string `form:"mobile" json:"mobile" binding:"required,mobile"`
	PassWord string `form:"password" json:"password" binding:"required,min=3,max=20"`
	Code     string `form:"code" json:"code" binding:"required,min=6,max=6"` // 短信验证码
}

// 前端传的Form表单，用于注册
type UpdateUserForm struct {
	Name     string `form:"name" json:"name" binding:"required,max=10,min=3"`
	Birthday string `form:"birthday" json:"birthday" binding:"required,datetime=2006-01-02"`
	Gender   string `form:"gender" json:"gender" binding:"required,oneof=female male"`
}
