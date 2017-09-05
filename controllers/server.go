package controllers

import (
	"encoding/json"
	"time"

	"github.com/alsmile/goApiGateway/models"
	"github.com/alsmile/goApiGateway/services"
	"github.com/alsmile/goApiGateway/services/sites"
	"github.com/kataras/iris/context"
	"gopkg.in/mgo.v2/bson"
)

// ApisSet 批量设置api
func ApisSet(ctx context.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	sdkSite := &models.SdkSite{}
	err := ctx.ReadJSON(sdkSite)
	if err != nil {
		ret["error"] = services.ErrorParam
		return
	}

	if sdkSite.Name == "" {
		sdkSite.Name = sdkSite.APIDomain
	}

	site := &models.Site{
		Name:      sdkSite.Name,
		Gzip:      sdkSite.Gzip,
		HTTPS:     sdkSite.HTTPS,
		APIDomain: sdkSite.APIDomain,
		Group:     sdkSite.Group,
		DstURL:    sdkSite.DstURL,
	}

	// 先删除
	sites.DelSiteBySDK(site)

	// 再保存site信息
	err = sites.Save(site, "sdk")
	if err != nil {
		ret["error"] = err.Error()
		return
	}

	var apis []models.SiteAPI
	json.Unmarshal([]byte(sdkSite.Apis), &apis)

	if len(apis) < 1 {
		return
	}

	docs := make([]interface{}, len(apis))
	for i := 0; i < len(apis); i++ {
		apis[i].Site.ID = site.ID
		apis[i].Site.APIDomain = site.APIDomain
		apis[i].Site.DstURL = site.DstURL
		apis[i].Site.Group = site.Group
		apis[i].Site.HTTPS = site.HTTPS
		apis[i].Site.Gzip = site.Gzip
		apis[i].ID = bson.NewObjectId()
		apis[i].UpdatedAt = time.Now().UTC()
		apis[i].CreatedAt = apis[i].UpdatedAt

		if apis[i].URL == "" {
			apis[i].URL = site.Group + apis[i].ShortURL
		}

		docs[i] = apis[i]
	}
	sites.SaveApis(docs)
}

// ApisDelete 批量删除api
func ApisDelete(ctx context.Context) {
	ret := make(map[string]interface{})
	defer ctx.JSON(ret)

	sdkSite := &models.SdkSite{}
	err := ctx.ReadJSON(sdkSite)
	if err != nil {
		ret["error"] = services.ErrorParam
		return
	}

	site := &models.Site{
		Name:      sdkSite.Name,
		Gzip:      sdkSite.Gzip,
		HTTPS:     sdkSite.HTTPS,
		APIDomain: sdkSite.APIDomain,
		Group:     sdkSite.Group,
		DstURL:    sdkSite.DstURL,
	}

	// 删除site信息
	err = sites.DelSiteBySDK(site)
	if err != nil {
		ret["error"] = err.Error()
		return
	}
}
