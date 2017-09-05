package controllers

import (
	"github.com/alsmile/goApiGateway/models"
	"github.com/alsmile/goApiGateway/services"
	"github.com/alsmile/goApiGateway/services/captcha"
	"github.com/alsmile/goApiGateway/services/user"
	"github.com/alsmile/goApiGateway/session"
	"github.com/alsmile/goApiGateway/utils"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
)

// Login 登录
func Login(ctx context.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	u := &models.User{}
	ctx.ReadJSON(u)
	rememberMe := u.RememberMe

	sid := session.GetSessionId(ctx)
	if captcha.IsNeedSignCaptcha(sid) {
		if captcha.VerifyImage(ctx, u.Captcha) == false {
			ret["error"] = services.ErrorCaptchaCode
			ret["errorTip"] = "captcha"
			return
		}
	}

	err := user.GetUserByPassword(u)
	if err != nil {
		ret["error"] = err.Error()
		captcha.SignError(sid)
		if captcha.IsNeedSignCaptcha(sid) {
			ret["errorTip"] = "captcha"
		}
		return
	}

	ret["id"] = u.ID
	ret["email"] = u.Profile.Email
	ret["username"] = u.Profile.Username
	if rememberMe {
		ret["token"] = user.GetToken(u, services.TokenValidRemember)
	} else {
		ret["token"] = user.GetToken(u, services.TokenValidHours)
	}
}

// SignUp 注册
func SignUp(ctx context.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	u := &models.User{}
	ctx.ReadJSON(u)

	if captcha.VerifyImage(ctx, u.Captcha) == false {
		ret["error"] = services.ErrorCaptchaCode
		return
	}

	err := user.AddUser(u)
	if err != nil {
		ret["error"] = err.Error()
	}
}

// SignActive 用户激活
func SignActive(ctx context.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	u := &models.User{}
	ctx.ReadJSON(u)

	err := user.Active(u)
	if err != nil {
		ret["error"] = err.Error()
	}

	ret["id"] = u.ID
	ret["email"] = u.Profile.Email
	ret["username"] = u.Profile.Username
	ret["token"] = user.GetToken(u, services.TokenValidHours)
}

// ForgetPassword 忘记密码请求
func ForgetPassword(ctx context.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	u := &models.User{}
	ctx.ReadJSON(u)

	if captcha.VerifyImage(ctx, u.Captcha) == false {
		ret["error"] = services.ErrorCaptchaCode
		return
	}

	err := user.ForgetPassword(u)
	if err != nil {
		ret["error"] = err.Error()
	}
}

// NewPassword 忘记密码时，设置新密码
func NewPassword(ctx context.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	u := &models.User{}
	ctx.ReadJSON(u)

	err := user.NewPassword(u)
	if err != nil {
		ret["error"] = err.Error()
		return
	}

	ret["id"] = u.ID
	ret["email"] = u.Profile.Email
	ret["username"] = u.Profile.Username
	ret["token"] = user.GetToken(u, services.TokenValidHours)
}

// UserProfile 获取用户基本信息
func UserProfile(ctx context.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	u, err := user.GetUserByTokenID(ctx.GetHeader("Authorization"), ctx.Values().GetString("uid"))
	if err != nil {
		ctx.StatusCode(iris.StatusUnauthorized)
		ret["error"] = err.Error()
		return
	}

	ret["id"] = u.ID
	if u.ID == "" {
		ret["id"] = u.Profile.UserID
	}
	ret["email"] = u.Profile.Email
	ret["username"] = u.Profile.Username
}

// Auth 身份认证中间件
func Auth(ctx context.Context) {
	uid := user.ValidToken(ctx)
	if uid == "" {
		ctx.StatusCode(iris.StatusUnauthorized)
		ret := make(map[string]interface{})
		ret["error"] = services.ErrorNeedSign
		ret["loginUrl"] = utils.GlobalConfig.User.LoginURL
		ctx.JSON(ret)
		return
	}

	ctx.Values().Set("uid", uid)

	ctx.Next()
}
