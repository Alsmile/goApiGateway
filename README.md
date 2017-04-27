#goMicroServer  
goMicroServer - 简单、快捷、可扩展的WebServer管理工具；具备反向代理、api网关、在线mock等功能；网页式后台管理和网络日志。  


#window下编译linux    
1.下载 [tdm-gcc](http://tdm-gcc.tdragon.net/download)  
2.运行build.bat  
3. 命令行执行：
  set GOOS=linux  
  set GOPACH=amd64  
  go build

当前目录下新增的le5le即为编译程序  


#运行  
当前目录下的config/default.json为缺省配置（主要用于本地开发）。需要修改配置时，在config下自定义json文件即可覆盖合并缺省配置，如my.json。  


