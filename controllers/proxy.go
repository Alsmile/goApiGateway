package controllers

import (
	"io"
	"log"
	"net/http"
	"time"

	"github.com/alsmile/goApiGateway/models"
	"github.com/alsmile/goApiGateway/services/plugins"
	"github.com/alsmile/goApiGateway/services/sites"
	"github.com/alsmile/goApiGateway/utils"
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"
)

// ProxyRequest 收到代理请求并处理
func ProxyRequest(ctx iris.Context) {
	if ctx.GetHeader("Upgrade") == "websocket" {
		ctx.Next()
		return
	}

	host := ctx.Host()
	method := string(ctx.Method())
	url := "/" + ctx.Params().Get("url")
	remoteAddr := ctx.RemoteAddr()

	// api代理
	found, err := proxyAPI(ctx, host, method, url, remoteAddr)
	if !found {
		found, err = proxySite(ctx, host, method, url, remoteAddr)
	}

	ret := make(map[string]interface{})
	if !found {
		ctx.StatusCode(iris.StatusNotFound)
		ret["error"] = "Not found."
	} else if err != nil {
		ret["error"] = err.Error()
	} else {
		return
	}

	ret["method"] = method
	ret["url"] = url
	ctx.JSON(ret)
}

// proxyAPI 查找api代理，返回： bool - 是否找到api定义; error - error
func proxyAPI(ctx iris.Context, host, method, url, remoteAddr string) (bool, error) {
	var found bool
	siteAPI, err := sites.GetAPIByURL(host, method, url)
	if err == nil {
		found = true
		// 插件功能
		if utils.GlobalConfig.Plugins.IP {
			limit := false
			if len(siteAPI.Whitelist) > 0 || len(siteAPI.Blacklist) > 0 {
				limit = plugins.IPLimit(remoteAddr, &siteAPI.Whitelist, &siteAPI.Blacklist)
			} else {
				limit = plugins.IPLimit(remoteAddr, &siteAPI.Site.Whitelist, &siteAPI.Site.Blacklist)
			}
			if limit {
				ctx.StatusCode(iris.StatusNotFound)
				return true, err
			}
		}
		// end 插件功能

		if siteAPI.IsMock {
			ctx.ResponseWriter().Header().Set("Content-Type", siteAPI.DataType)
			ctx.WriteWithExpiration([]byte(siteAPI.ResponseParamsText), time.Now())
		} else {
			err = proxy(ctx, method, siteAPI.Site.DstURL+url)
		}
	}

	return found, err
}

// proxySite 查找site代理，返回： bool - 是否找到site定义; error - error
func proxySite(ctx iris.Context, host, method, url, remoteAddr string) (bool, error) {
	var found bool
	site, err := sites.GetSiteByDomain(host)
	if err == nil {
		found = true
		// 插件功能
		if utils.GlobalConfig.Plugins.IP {
			if plugins.IPLimit(remoteAddr, &site.Whitelist, &site.Blacklist) {
				ctx.StatusCode(iris.StatusNotFound)
				return true, err
			}
		}
		// end 插件功能

		err = proxy(ctx, method, site.DstURL+url)

		// 添加到自动发现
		siteAPI := &models.SiteAPI{}
		siteAPI.AutoReg = true
		siteAPI.Site.ID = site.ID
		siteAPI.Method = method
		siteAPI.URL = url
		siteAPI.Visited = 1
		siteAPI.StatusCode = ctx.GetStatusCode()
		sites.SaveAPI(siteAPI, site.OwnerID)
	}

	return found, err
}

// proxy 执行具体代理
func proxy(ctx iris.Context, method, dstURL string) (err error) {
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
			ctx.ResponseWriter().Header().Set(key, v)
		}
	}
	io.Copy(ctx.ResponseWriter(), clientResp.Body)
	clientResp.Body.Close()
	err = nil
	return
}

// ProxyTest web请求代理测试，用于查看代理数据是否正确
func ProxyTest(ctx iris.Context) {
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
