package controllers

import (
	"strings"

	"github.com/alsmile/goApiGateway/services/sites"
	"github.com/kataras/iris"
	"github.com/kataras/iris/websocket"
	goWebsocket "golang.org/x/net/websocket"
)

// SetupWebsocket 安装websocket
func SetupWebsocket(app *iris.Application) {
	ws := websocket.New(websocket.Config{
		ReadBufferSize:  2048,
		WriteBufferSize: 2048,
	})
	ws.OnConnection(wsConnection)

	app.Get("/{url:path}", ProxyRequest, ws.Handler())
}

func wsConnection(conn websocket.Connection) {
	host := conn.Context().Host()
	url := conn.Context().Request().URL

	// 查找域名
	site, err := sites.GetSiteByDomain(host)
	if err != nil {
		conn.EmitMessage([]byte("Error: Can not find the host."))
		conn.Disconnect()
		return
	}

	// 替换ws协议
	dstURL := strings.Replace(site.DstURL, "http", "ws", 1)
	ws, err := goWebsocket.Dial(dstURL+url.String(), "", "http://localhost")
	if err != nil {
		conn.EmitMessage([]byte("Error: " + err.Error()))
		conn.Disconnect()
		return
	}

	// 收到消息发往目标服务器
	conn.OnMessage(func(data []byte) {
		ws.Write(data)
	})

	// 客户端连接端口
	conn.OnDisconnect(func() {
		ws.Close()
		ws = nil
	})

	go receiveLoop(&conn, ws)
}

// receiveLoop 循环接收消息
func receiveLoop(conn *websocket.Connection, ws *goWebsocket.Conn) {
	msg := make([]byte, 2048)
	var n int
	var err error
	for {
		if n, err = ws.Read(msg); err != nil {
			return
		}
		(*conn).EmitMessage(msg[:n])
	}
}
