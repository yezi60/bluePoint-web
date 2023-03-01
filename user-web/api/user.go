package api

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"bluepoint_api/user-web/forms"
	"bluepoint_api/user-web/global"
	"bluepoint_api/user-web/global/response"
	"bluepoint_api/user-web/middlewares"
	"bluepoint_api/user-web/models"
	"bluepoint_api/user-web/proto"
)

// removeTopStruct 去掉错误提示中结构体的定位
func removeTopStruct(fileds map[string]string) map[string]string {
	rsp := map[string]string{}

	for filed, err := range fileds {
		// 找到.所在的位置，取出剩下的字符串
		rsp[filed[strings.Index(filed, ".")+1:]] = err
	}

	return rsp
}

// HandleGrpcErrorToHttp 将grpc的code转换成http的状态码 并以json的格式进行返回
func HandleGrpcErrorToHttp(c *gin.Context, err error) {

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound:
				c.JSON(http.StatusNotFound, gin.H{
					"msg": e.Message(),
				})
			case codes.Internal:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "内部错误",
				})
			case codes.InvalidArgument:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "参数错误",
				})
			case codes.Unavailable:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "用户服务不可用",
				})
			case codes.AlreadyExists:
				c.JSON(http.StatusBadRequest, gin.H{
					"msg": "用户已存在",
				})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "其它错误",
				})
			}
		}
	}
}

// HandleValidatorError 用于封装错误处理
func HandleValidatorError(c *gin.Context, err error) {
	if errs, ok := err.(validator.ValidationErrors); !ok {
		c.JSON(http.StatusBadRequest, gin.H{
			"msg": err.Error(),
		})
	} else {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": removeTopStruct(errs.Translate(global.Trans)),
		})
	}
}

// GetUserList 调用grpc获取用户列表
func GetUserList(ctx *gin.Context) {

	claims, _ := ctx.Get("claims")
	currentUser := claims.(*models.CustomClaims)
	zap.S().Infof("访问的用户：%v", currentUser.ID)

	pn := ctx.DefaultQuery("pn", "0")
	pSize := ctx.DefaultQuery("pnum", "10")
	pnInt, _ := strconv.Atoi(pn)
	pSizeInt, _ := strconv.Atoi(pSize)

	rsp, err := global.UserSrvClient.GetUserList(context.Background(), &proto.PageInfo{
		Pn:    uint32(pnInt),
		PSize: uint32(pSizeInt),
	})

	if err != nil {
		zap.S().Errorw("[GetUserList] 查询 【用户列表】失败")
		HandleGrpcErrorToHttp(ctx, err)
		return
	}

	result := make([]interface{}, 0)

	for _, value := range rsp.Data {

		users := response.UserResponse{
			Id:       value.Id,
			NickName: value.NickName,
			Gender:   value.Gender,
			Mobile:   value.Mobile,
			Birthday: response.JsonTime(time.Unix(int64(value.Birthday), 0)),
			//Birthday: time.Unix(int64(value.Birthday), 0),
		}

		result = append(result, users)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"data":  result,
		"total": rsp.Total,
	})
}

// PassWordLogin 使用验证码登陆 表单验证
func PassWordLogin(c *gin.Context) {

	passwordLoginForm := forms.PasswordLoginForm{}

	if err := c.ShouldBindJSON(&passwordLoginForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	// 输入进行验证
	if !store.Verify(passwordLoginForm.CaptchaId, passwordLoginForm.Captcha, true) {
		c.JSON(http.StatusBadRequest, gin.H{
			"captcha": "验证码错误",
		})
		return
	}

	// 登陆的逻辑
	rspUser, err := global.UserSrvClient.GetUserByMobile(context.Background(), &proto.MobileRequest{
		Mobile: passwordLoginForm.Mobile,
	})

	if err != nil {
		if e, ok := status.FromError(err); ok {
			switch e.Code() {
			case codes.NotFound: //找不到用户
				c.JSON(http.StatusBadRequest, map[string]string{
					"msg": "用户不存在",
				})
			default:
				zap.S().Errorw("[GetUserByMobile] 查询 【用户】失败",
					"msg", err.Error(),
				)
				c.JSON(http.StatusInternalServerError, map[string]string{
					"msg": "登陆失败",
				})

			}
		}
		return
	} else {
		// 知识查询到用户了而已，但是没有检查密码
		passRsp, err := global.UserSrvClient.CheckPassWord(context.Background(), &proto.PasswordCheckInfo{
			Password:          passwordLoginForm.PassWord,
			EncryptedPassword: rspUser.Password,
		})

		if err != nil {
			zap.S().Errorw("[CheckPassWord] 检查 【用户密码服务】失败",
				"msg", err.Error(),
			)
			c.JSON(http.StatusInternalServerError, map[string]string{
				"msg": "登陆失败",
			})
			return
		}

		if passRsp.Success {
			// 生成token
			newJwt := middlewares.NewJWT()
			claims := models.CustomClaims{
				ID:          uint(rspUser.Id),
				NickName:    rspUser.NickName,
				AuthorityID: uint(rspUser.Role),
				StandardClaims: jwt.StandardClaims{
					NotBefore: time.Now().Unix(),                                      // 签名的生效时间
					ExpiresAt: time.Now().Add(time.Second * 60 * 60 * 24 * 30).Unix(), // 过期时间30天
					Issuer:    "soleaf",
				},
			}

			token, err := newJwt.CreateToken(claims)
			if err != nil {
				zap.S().Errorf("CreateToken failed,err :%v", err.Error())
				c.JSON(http.StatusInternalServerError, gin.H{
					"msg": "生成token失败",
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"id":         rspUser.Id,
				"nick_name":  rspUser.NickName,
				"token":      token,
				"expired_at": time.Now().Add(time.Second*60*60*24*30).Unix() * 1000,
			})
		} else {
			c.JSON(http.StatusBadRequest, map[string]string{
				"msg": "密码错误",
			})
			return
		}
	}
}

// Register 用户注册的handleFunc
func Register(c *gin.Context) {

	registerForm := forms.RegisterForm{}

	if err := c.ShouldBindJSON(&registerForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	// 验证码校验

	// 连接redis
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d", global.ServerConfig.RedisInfo.Host, global.ServerConfig.RedisInfo.Port),
	})

	value, err := rdb.Get(context.Background(), registerForm.Mobile).Result()
	if err != nil {
		if err == redis.Nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"code": "验证码错误",
			})
		} else {
			zap.S().Errorf("redis query failed:%v", err.Error())
		}
		return
	}

	if value != registerForm.Code {
		c.JSON(http.StatusBadRequest, gin.H{
			"code": "验证码错误",
		})
		return
	}

	rspUser, err := global.UserSrvClient.CreateUser(context.Background(), &proto.CreateUserInfo{
		NickName: registerForm.Mobile,
		Password: registerForm.PassWord,
		Mobile:   registerForm.Mobile,
	})

	if err != nil {
		zap.S().Errorf("[Register] 查询 【新建用户失败】 err:%s", err.Error())
		HandleGrpcErrorToHttp(c, err)
		return
	}

	// 注册即登陆
	// 生成token
	newJwt := middlewares.NewJWT()
	claims := models.CustomClaims{
		ID:          uint(rspUser.Id),
		NickName:    rspUser.NickName,
		AuthorityID: uint(rspUser.Role),
		StandardClaims: jwt.StandardClaims{
			NotBefore: time.Now().Unix(),                                      // 签名的生效时间
			ExpiresAt: time.Now().Add(time.Second * 60 * 60 * 24 * 30).Unix(), // 过期时间30天
			Issuer:    "soleaf",
		},
	}

	token, err := newJwt.CreateToken(claims)
	if err != nil {
		zap.S().Errorf("CreateToken failed,err :%v", err.Error())
		c.JSON(http.StatusInternalServerError, gin.H{
			"msg": "生成token失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"id":         rspUser.Id,
		"nick_name":  rspUser.NickName,
		"token":      token,
		"expired_at": time.Now().Add(time.Second*60*60*24*30).Unix() * 1000,
	})
}

// 补充，用户信息页中的获取用户信息
func GetUserDetail(c *gin.Context) {
	userId, _ := c.Get("userId")

	rsp, err := global.UserSrvClient.GetUserById(context.Background(), &proto.IdRequest{Id: int32(userId.(uint))})
	if err != nil {
		zap.S().Errorf("[GetUserDetail] 查询 【用户详细信息失败】 err:%s", err.Error())
		HandleGrpcErrorToHttp(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"name":     rsp.NickName,
		"birthday": time.Unix(int64(rsp.Birthday), 0).Format("2006-01-02"),
		"gender":   rsp.Gender,
		"mobile":   rsp.Mobile,
	})
}

func UpdateUser(c *gin.Context) {
	var updateForm forms.UpdateUserForm

	if err := c.ShouldBindJSON(&updateForm); err != nil {
		HandleValidatorError(c, err)
		return
	}

	userId, _ := c.Get("userId")

	// 将birthday转uint
	// 前端传什么格式
	loc, _ := time.LoadLocation("Local") //首字母必须要大写！！！
	birthDay, _ := time.ParseInLocation("2006-01-02", updateForm.Birthday, loc)

	_, err := global.UserSrvClient.UpdateUser(context.Background(), &proto.UpdateUserInfo{
		Id:       int32(userId.(uint)),
		NickName: updateForm.Name,
		Gender:   updateForm.Gender,
		Birthday: uint64(birthDay.Unix()),
	})

	if err != nil {
		zap.S().Errorf("[UpdateUser] 更新 【用户信息失败】 err:%s", err.Error())
		HandleGrpcErrorToHttp(c, err)
		return
	}

	c.Status(http.StatusOK)
}
