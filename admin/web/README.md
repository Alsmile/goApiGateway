# le5le-ng-start  
le5le-ng-start是le5le的angular+webpack2框架，支持模块延迟加载。  
<a href="https://le5le-com.github.io/le5le-ng-start/" target="_blank">demo</a>

# 基本功能  
登录  
注册  
忘记密码  
账号激活  
延迟加载模块  

其中，在线预览无法调用mock接口。在本地调试时，可以使用在线mock数据演示登录过程。  
在线数据mock：<a href="http://www.i-dengta.com" target="_blank">灯塔互联</a>提供了登录、注册、忘记密码、账号激活等接口服务。

# 框架依赖  
css UI库：<a href="https://github.com/le5le-com/le5le-ui" target="_blank">le5le-ui</a>
angular 组件库：<a href="https://github.com/le5le-com/le5le-components" target="_blank">le5le-components</a>
angular 数据中心服务：<a href="https://github.com/le5le-com/le5le-store" target="_blank">le5le-store</a>


# 开发环境  
### 1.安装nodejs  
https://nodejs.org/en/

### 2.安装依赖库  
执行npm install安装依赖库。国内用户建议使用cnpm install。


### 3.开发调试  
运行npm start命令即可进行本地开发调试。  


# 编译  
生产环境编译：    
npm run build  
编译后的生产文件在docs文件夹下。

# css使用
src/assets文件夹下的index.pcss，并且按照自己的需求修改变量和导入即可。





