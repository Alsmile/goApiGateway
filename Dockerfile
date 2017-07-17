FROM ubuntu:16.04
MAINTAINER Alsmile "bing@cloudtogo.cn"
ENV REFRESHED_AT 2017-07-14

# 程序安装目录
ENV APIGATEWAY_HOME /usr/local/goApiGateway
ENV PATH $APIGATEWAY_HOME/bin:$PATH
WORKDIR $APIGATEWAY_HOME
RUN mkdir -p "$APIGATEWAY_HOME"
RUN mkdir -p "$APIGATEWAY_HOME/admin/web/dist"
RUN mkdir -p "$APIGATEWAY_HOME/assets"
RUN mkdir -p "$APIGATEWAY_HOME/config"
RUN mkdir -p "$APIGATEWAY_HOME/out/log"

# 拷贝主程序
ADD ./goApiGateway $APIGATEWAY_HOME/
ADD ./admin/web/dist/ $APIGATEWAY_HOME/admin/web/dist/
ADD ./assets/ $APIGATEWAY_HOME/assets/
ADD ./config/default.json $APIGATEWAY_HOME/config/

RUN chmod +x $APIGATEWAY_HOME/goApiGateway

EXPOSE 80
ENTRYPOINT ["./goApiGateway"]

# docker run --name goApiGateway -p 80:80 [-v /etc/goApiGateway.json:/etc/goApiGateway.json] 或 [-e configKey=configValue] [--net=host] <image name:tag>
#挂载配置文件和使用环境变量任选一种。其中，环境变量参数将覆盖配置文件中的参数
# 环境变量参数:
#   ApiGateway_Website：网站的后台管理web首页网址，例如：http://admin.apicloudtogo.cn
#   ApiGateway_Cpu：使用CPU核数
#   ApiGateway_Domain_Domain：网关根域名，例如： api.cloudtogo.cn
#   ApiGateway_Domain_AdminDomain：相对于网关根域名的后台管理子域名，例如：“admin.”
#   ApiGateway_Domain_Port：端口，推荐80
#   ApiGateway_Domain_SdkPort：内部服务端口，供行云创新系统内部局域网使用，没有身份认证过程，接口不公开。
#   ApiGateway_User_LoginUrl：第三方用户登录网址
#   ApiGateway_User_SignUpUrl：第三方用户注册网址
#   ApiGateway_User_InfoUrl：获取第三方用户身份信息
#   ApiGateway_Jwt：Jwt密钥。如果和行云创新系统结合，暂时需要使用相同密钥
#   ApiGateway_Secret：内部密钥
#   ApiGateway_Mongo_Address：mongodb数据库连接地址，例如：localhost:27017
#   ApiGateway_Mongo_Database：mongodb数据库名
#   ApiGateway_Mongo_User：mongodb数据库user
#   ApiGateway_Mongo_Password：mongodb数据库密码
#   ApiGateway_Mongo_MaxConnections：连接池数
#   ApiGateway_Mongo_Mechanism：连接方式
#   ApiGateway_Redis_Address：redis连接地址，例如：127.0.0.1:6379
#   ApiGateway_Redis_ConnectNum：连接池数
#   ApiGateway_Redis_Password：密码
#   ApiGateway_Redis_Db：数据库
#   ApiGateway_Log_Filename：输出日志文件名，例如：./out/log/default.log；当为空时，输出到标准输出。
#   ApiGateway_Log_MaxAge：有效天数
#   ApiGateway_Log_MaxSize：单个日志文件最大大小
#   ApiGateway_Log_MaxBackups：最大日志数
#   ApiGateway_Email_Address：邮箱地址，用于发生系统邮件，例如：smtp.apiGateway.com
#   ApiGateway_Email_Port：邮箱端口，例如：80。25已经被国内限制
#   ApiGateway_Email_User：邮箱登录账户
#   ApiGateway_Email_Password：邮箱登录账户密码


