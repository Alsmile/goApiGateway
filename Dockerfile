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
EXPOSE 3200
ENTRYPOINT ["./goApiGateway"]

# docker run --name goApiGateway -p 80:80 -p 3200:3200 [-v /etc/goApiGateway.json:/etc/goApiGateway.json] [--net=host] <image name:tag>
