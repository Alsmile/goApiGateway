package services

import "gopkg.in/mgo.v2/bson"

const (
	ErrorPhoneCode     = "手机验证码错误"
	ErrorPhoneEmpty    = "手机号不能为空"
	ErrorPhoneExists   = "该手机号已经被注册"
	ErrorEmailEmpty    = "邮箱不能为空"
	ErrorEmailExists   = "该邮箱已经被注册"
	ErrorPasswordEmpty = "密码不能为空"
	ErrorCaptchaCode   = "验证码错误"
	ErrorLogin         = "邮箱或密码错误"
	ErrorParam         = "参数错误"
	ErrorParamPage     = "分页参数错误"
	ErrorPermission    = "权限错误"
	ErrorSave          = "保存数据错误，请稍后重试"
	ErrorFile          = "读取上传文件错误"
	ErrorFileMaxSize   = "上传文件大小不能超过10M"
	ErrorFileInfo      = "无法读取文件"
	ErrorFileSave      = "保存文件错误"
	ErrorPassword      = "密码错误"
	ErrorUserNoExists  = "用户不存在"
	ErrorEmail         = "发送邮件错误"
	ErrorActiveCode    = "激活码错误"
	ErrorNoActive      = "账号没有激活"
	ErrorRead          = "数据不存在；或读取数据错误，请稍后重试"
	ErrorNeedSign      = "请先登录"
	ErrorURL           = "无效的url"
	ErrorAPIPause      = "API已经暂停访问"
	ErrorProxyNotFound = "访问您的开发接口错误：404"
)

const (
	TokenValidHours    = 16
	TokenValidRemember = 8760  // 24 * 365 = 1 years
	TokenValidMobile   = 61320 // 24 * 365 * 7 = 7 years
)

const (
	HeaderTrim = "Bearer "
)

const (
	PageIndex = "pageIndex"
	PageCount = "pageCount"
)

const (
	FilePathTemp          = "/temp"
	FilePathImage         = "/images"
	FilePathImageNotFound = "img/notFound.png"
	FilePathImageUser     = "img/user.png"
)

var SelectHide = bson.M{"editor": false, "deletedAt": false}
