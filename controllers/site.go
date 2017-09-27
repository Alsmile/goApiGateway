package controllers

import (
	"github.com/alsmile/goApiGateway/models"
	"github.com/alsmile/goApiGateway/services"
	"github.com/alsmile/goApiGateway/services/sites"
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"
)

// SiteList 获取用户网站列表
func SiteList(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	pageIndex, err := ctx.URLParamInt(services.PageIndex)
	if err != nil || pageIndex < 1 {
		ret["error"] = services.ErrorParamPage
		return
	}
	pageCount, err := ctx.URLParamInt(services.PageCount)
	if err != nil || pageCount < 1 {
		ret["error"] = services.ErrorParamPage
		return
	}

	list, err := sites.List(ctx.Values().GetString("uid"), pageIndex, pageCount)

	if err != nil {
		ret["error"] = err.Error()
	}

	ret["list"] = list
}

// SiteGet 获取具体的网站信息
func SiteGet(ctx iris.Context) {
	ret := make(map[string]interface{})

	id := ctx.URLParam("id")
	if id == "" {
		ret["error"] = services.ErrorParam
		ctx.JSON(ret)
		return
	}

	site := &models.Site{ID: bson.ObjectIdHex(id)}
	err := sites.Get(site, ctx.Values().GetString("uid"))

	if err != nil {
		ret["error"] = err.Error()
		ctx.JSON(ret)
		return
	}

	ctx.JSON(site)
}

// SiteSave 保存网站信息
func SiteSave(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	site := &models.Site{}
	err := ctx.ReadJSON(site)
	if err != nil {
		ret["error"] = services.ErrorParam
		ret["errorConsole"] = err.Error()
	}

	err = sites.Save(site, ctx.Values().GetString("uid"))

	if err != nil {
		ret["error"] = err.Error()
	}
}

// SiteDel 删除网站
func SiteDel(ctx iris.Context) {
	ret := make(map[string]interface{})

	id := ctx.URLParam("id")
	if id == "" {
		ret["error"] = services.ErrorParam
		ctx.JSON(ret)
		return
	}

	err := sites.DelSite(id, ctx.Values().GetString("uid"))
	if err != nil {
		ret["error"] = err.Error()
		ctx.JSON(ret)
		return
	}

	ctx.JSON(true)
}

// SiteAPISave 保存网站下的api
func SiteAPISave(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	siteAPI := &models.SiteAPI{}
	err := ctx.ReadJSON(siteAPI)
	if err != nil {
		ret["error"] = services.ErrorParam
		ret["errorConsole"] = err.Error()
	}

	err = sites.SaveAPI(siteAPI, ctx.Values().GetString("uid"))
	if err != nil {
		ret["error"] = err.Error()
	}
	ret["id"] = siteAPI.ID
}

// SiteAPIGet 获取网站下的api
func SiteAPIGet(ctx iris.Context) {
	ret := make(map[string]interface{})

	id := ctx.URLParam("id")
	if id == "" {
		ret["error"] = services.ErrorParam
		ctx.JSON(ret)
		return
	}

	siteAPI := &models.SiteAPI{ID: bson.ObjectIdHex(id)}
	err := sites.GetAPI(siteAPI, ctx.Values().GetString("uid"))

	if err != nil {
		ret["error"] = err.Error()
		ctx.JSON(ret)
		return
	}

	ctx.JSON(siteAPI)
}

// SiteAPIDel 删除网站下的api
func SiteAPIDel(ctx iris.Context) {
	ret := make(map[string]interface{})

	id := ctx.URLParam("id")
	if id == "" {
		ret["error"] = services.ErrorParam
		ctx.JSON(ret)
		return
	}

	err := sites.DelAPI(id, ctx.Values().GetString("uid"))
	if err != nil {
		ret["error"] = err.Error()
		ctx.JSON(ret)
		return
	}

	ctx.JSON(true)
}

// SiteAPIList api列表
func SiteAPIList(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	pageIndex, err := ctx.URLParamInt(services.PageIndex)
	if err != nil || pageIndex < 1 {
		ret["error"] = services.ErrorParamPage
		return
	}
	pageCount, err := ctx.URLParamInt(services.PageCount)
	if err != nil || pageCount < 1 {
		ret["error"] = services.ErrorParamPage
		return
	}

	auto := ctx.URLParam("auto")
	fieldType, _ := ctx.URLParamInt("field")

	siteID := bson.ObjectIdHex(ctx.URLParam("siteId"))
	list, err := sites.APIList(ctx.Values().GetString("uid"), siteID, auto, fieldType, pageIndex, pageCount)

	if err != nil {
		ret["error"] = err.Error()
	}

	ret["list"] = list
}

// SiteAPIListByDomains 查找指定域名下的api
func SiteAPIListByDomains(ctx iris.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	pageIndex, err := ctx.URLParamInt(services.PageIndex)
	if err != nil || pageIndex < 1 {
		ret["error"] = services.ErrorParamPage
		return
	}
	pageCount, err := ctx.URLParamInt(services.PageCount)
	if err != nil || pageCount < 1 {
		ret["error"] = services.ErrorParamPage
		return
	}

	auto := ctx.URLParam("auto")
	fieldType, _ := ctx.URLParamInt("field")

	var domains []string
	err = ctx.ReadJSON(&domains)
	if err != nil || len(domains) < 1 {
		ret["error"] = services.ErrorParam
	}

	list, err := sites.APIListByDomains(ctx.Values().GetString("uid"), domains, auto, fieldType, pageIndex, pageCount)

	if err != nil {
		ret["error"] = err.Error()
	}

	ret["list"] = list
}
