package controllers

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/alsmile/goApiGateway/models"
	"github.com/alsmile/goApiGateway/services/sites"
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"gopkg.in/mgo.v2/bson"
)

// ProxyDo 收到代理请求并处理
func ProxyDo(ctx context.Context) {
	ret := make(map[string]interface{})

	host := ctx.Host()
	method := string(ctx.Method())
	url := "/" + ctx.Params().Get("url")

	// 查找api级别代理
	siteAPI, err := sites.GetAPIByURL(host, method, url)
	if err == nil {
		if siteAPI.IsMock {
			ctx.ResponseWriter().Header().Set("Content-Type", siteAPI.DataType)
			ctx.WriteWithExpiration([]byte(siteAPI.ResponseParamsText), time.Now())
		} else {
			proxy(ctx, method, siteAPI.Site.DstURL+url)
		}
		return
	}

	// 查找site级别代理
	site, err := sites.GetSiteByDomain(host)
	if err == nil {
		proxy(ctx, method, site.DstURL+url)

		siteAPI = &models.SiteAPI{}

		// 添加到自动发现
		siteAPI.AutoReg = true
		siteAPI.Site.ID = site.ID
		siteAPI.Method = method
		siteAPI.URL = url
		siteAPI.Visited = 1
		siteAPI.StatusCode = ctx.GetStatusCode()
		sites.SaveAPI(siteAPI, site.OwnerID)

		return
	}

	// log
	ret["method"] = method
	ret["url"] = url
	ret["error"] = "Not found."
	ctx.StatusCode(iris.StatusNotFound)
	ctx.JSON(ret)
}

// proxy 执行具体代理
func proxy(ctx context.Context, method, dstURL string) (err error) {
	client := &http.Client{}
	query := "?" + ctx.Request().URL.Query().Encode()
	clientReq, err := http.NewRequest(method, dstURL+query, ctx.Request().Body)
	if err != nil {
		log.Printf("[error]servers.controllers.proxy.proxy: method=%v, url=%v, proxyError=%v [[[[[[[[[[in NewRequest]]]]]]]]]]\r\n",
			method, dstURL, err)

		ctx.StatusCode(iris.StatusNotFound)
		ctx.JSON(bson.M{"error": err.Error()})

		return
	}

	clientReq.Header = ctx.Request().Header
	clientResp, err := client.Do(clientReq)

	if err != nil {
		ctx.StatusCode(iris.StatusBadGateway)
		ctx.JSON(bson.M{"error": err.Error()})
		return
	}
	ctx.StatusCode(clientResp.StatusCode)

	for key, value := range clientResp.Header {
		for _, v := range value {
			ctx.Header(key, v)
		}
	}

	io.Copy(ctx.ResponseWriter(), clientResp.Body)
	clientResp.Body.Close()
	err = nil
	return
}

// ProxyTest web请求代理测试，用于查看代理数据是否正确
func ProxyTest(ctx context.Context) {
	method := string(ctx.Method())
	host := ctx.URLParam("host")
	url := ctx.URLParam("url")

	siteAPI, err := sites.GetAPIByDstURL(host, method, url)
	if err == nil && siteAPI.IsMock {
		ctx.ResponseWriter().Header().Set("Content-Type", siteAPI.DataType)
		ctx.WriteWithExpiration([]byte(siteAPI.ResponseParamsText), time.Now())
		return
	}

	proxy(ctx, method, host+url)
}
