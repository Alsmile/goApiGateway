package services

const (
  ErrorPhoneCode = "手机验证码错误"
  ErrorPhoneEmpty = "手机号不能为空"
  ErrorPhoneExists = "该手机号已经被注册"
  ErrorEmailEmpty = "邮箱不能为空"
  ErrorEmailExists = "该邮箱已经被注册"
  ErrorPasswordEmpty = "密码不能为空"
  ErrorCaptchaCode = "验证码错误"
  ErrorLogin = "邮箱或密码错误"
  ErrorParam = "参数错误"
  ErrorPermission = "权限错误"
  ErrorSave = "保存数据错误，请稍后重试"
  ErrorFile = "读取上传文件错误"
  ErrorFileMaxSize = "上传文件大小不能超过10M"
  ErrorFileInfo = "无法读取文件"
  ErrorFileSave = "保存文件错误"
  ErrorPassword = "密码错误"
  ErrorUserNoExists = "用户不存在"
  ErrorEmail = "发送邮件错误"
  ErrorActiveCode = "激活码错误"
  ErrorNoActive = "账号没有激活"
  ErrorRead = "读取数据错误，请稍后重试"
  ErrorNeedSign = "请先登录"
  ErrorLock = "您的小伙伴正在编辑该数据，被锁定。锁定者："
  ErrorMaxCount = "您已经超出最大个数限制"
  ErrorProjectMaxSize = "您存储空间已满，请先增加您的存储空间"
  ErrorProxyEndTime = "您的接口代理已到期，请先续费"
  ErrorReadProxySetting = "读取您的代理设置失败，请稍后重试"
  ErrorReadProxyUrl = "您没有设置您的服务器地址，请在项目设置中设置"
  ErrorProxyNotFound = "访问您的开发接口错误：404"
  ErrorProxyOther = "访问您的开发接口错误，状态码："
)

const (
  TokenValidHours = 16
  TokenValidRemember = 8760  // 24 * 365 = 1 years
  TokenValidMobile = 61320  // 24 * 365 * 7 = 7 years
)

const (
  HeaderTrim = "Bearer "
)

const (
  FilePathTemp = "/temp"
  FilePathImage = "/images"
  FilePathImageNotFound = "img/notFound.png"
  FilePathImageUser = "img/user.png"
)

const (
  ProjectDataType_Index = 1
  ProjectDataType_Api = 2
  ProjectDataType_Api_Comment = 3
)

const (
  Feedback_New = 1
  Feedback_Dealing = 2
  Feedback_Dealt = 3
  Feedback_No_Deal = 4
  Feedback_Future = 5
)
