#goApiGateway 
goApiGateway - golang编写的api gateway；具备反向代理、api网关、在线mock等功能；网页式后台管理和网络日志。  


#交叉编译     
3. 命令行执行：
  set GOOS=linux  
  set GOPACH=amd64  (可省略)  
  go build

当前目录下新增的可执行程序即为编译程序  


#运行  
当前目录下的config/default.json为缺省配置（主要用于本地开发）。需要修改配置时，在config下自定义json文件即可覆盖合并缺省配置，如my.json。  


