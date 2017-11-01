package sites

import (
	"errors"
	"log"
	"time"

	"github.com/alsmile/goApiGateway/db/mongo"
	"github.com/alsmile/goApiGateway/models"
	"github.com/alsmile/goApiGateway/services"
	"github.com/alsmile/goApiGateway/utils"
	"gopkg.in/mgo.v2/bson"
)

// List 网站列表
func List(uid string, pageIndex, pageCount int) (sites []models.Site, err error) {
	if uid == "" {
		err = errors.New(services.ErrorPermission)
		return
	}

	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).
		Find(bson.M{"ownerId": uid, "deletedAt": bson.M{"$exists": false}}).
		Select(services.SelectHide).
		Sort("-updatedAt").Skip((pageIndex - 1) * pageCount).Limit(pageCount).
		All(&sites)

	if err != nil {
		log.Printf("[error]serivces.sites.List: err=%v, data=%s\r\n", err, uid)
		err = errors.New(services.ErrorRead)
	}

	return
}

// Get 获取具体网站
func Get(site *models.Site, uid string) (err error) {
	if site.ID == "" {
		err = errors.New(services.ErrorParam)
		return
	}

	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).
		Find(bson.M{"_id": site.ID, "ownerId": uid}).
		Select(services.SelectHide).
		One(&site)

	if err != nil {
		log.Printf("[error]serivces.sites.Get: err=%v, data=%v\r\n", err, site)
		err = errors.New(services.ErrorRead)
	}

	return
}

// Save 保存网站信息
func Save(site *models.Site, uid string) (err error) {
	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	site.UpdatedAt = time.Now().UTC()
	site.EditorID = uid
	if site.ID == "" {
		site.ID = bson.NewObjectId()
		site.CreatedAt = site.UpdatedAt
		site.OwnerID = uid
	} else {
		mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
			Update(bson.M{"site._id": site.ID, "ownerId": uid}, bson.M{"$set": bson.M{"site": models.SiteParam{
				ID:             site.ID,
				Gzip:           site.Gzip,
				HTTPS:          site.HTTPS,
				Group:          site.Group,
				Subdomain:      site.Subdomain,
				IsCustomDomain: site.IsCustomDomain,
				APIDomain:      site.APIDomain,
				DstURL:         site.DstURL,
				Pause:          site.Pause,
			}}})

	}

	_, err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).
		Upsert(bson.M{"_id": site.ID}, site)

	if err != nil {
		log.Printf("[error]serivces.sites.Save: err=%v, data=%v\r\n", err, site)
		err = errors.New(services.ErrorSave)
	}

	return
}

// DelSite 删除网站
func DelSite(id, uid string) (err error) {
	if id == "" {
		err = errors.New(services.ErrorParam)
		return
	}

	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).
		Remove(bson.M{"_id": bson.ObjectIdHex(id), "ownerId": uid})

	if err != nil {
		log.Printf("[error]serivces.sites.DelSite: err=%v, id=%s\r\n", err, id)
		err = errors.New(services.ErrorPermission)
	}

	mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
		RemoveAll(bson.M{"site._id": bson.ObjectIdHex(id)})

	return
}

// DelSiteBySDK 删除网站，仅平台内部使用
func DelSiteBySDK(site *models.Site) (err error) {
	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	_, err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).
		RemoveAll(bson.M{"apiDomain": site.APIDomain, "group": site.Group})

	if err != nil {
		log.Printf("[error]serivces.sites.DelSiteBySDK - CollectionSites: err=%v, site=%v\r\n", err, site)
		err = errors.New(services.ErrorSave)
	}

	_, errApis := mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
		RemoveAll(bson.M{"site.apiDomain": site.APIDomain, "site.group": site.Group})

	if errApis != nil {
		log.Printf("[error]serivces.sites.DelSiteBySDK - CollectionApis: err=%v, site=%v\r\n", err, site)
		err = errors.New(services.ErrorSave)
	}

	return
}

// SaveAPI 保存api
func SaveAPI(siteAPI *models.SiteAPI, uid string) (err error) {
	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	siteAPI.EditorID = uid
	siteAPI.UpdatedAt = time.Now().UTC()

	s := &models.Site{ID: siteAPI.Site.ID}

	//  site id不存在，表示自动根据api信息保存site
	if siteAPI.ID == "" && siteAPI.Site.ID == "" && siteAPI.Site.Subdomain != "" {
		tempSite := bson.M{}
		err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).
			Find(bson.M{"subdomain": siteAPI.Site.Subdomain}).
			Select(bson.M{"_id": true}).
			One(&tempSite)

		// 不存在site，插入保存
		if err != nil {
			s.UpdatedAt = time.Now().UTC()
			s.ID = bson.NewObjectId()
			siteAPI.Site.ID = s.ID
			s.CreatedAt = s.UpdatedAt
			s.Subdomain = siteAPI.Site.Subdomain
			s.IsCustomDomain = siteAPI.Site.IsCustomDomain
			s.APIDomain = siteAPI.Site.APIDomain
			s.Group = siteAPI.Site.Group
			s.DstURL = siteAPI.Site.DstURL
			s.HTTPS = siteAPI.Site.HTTPS
			s.Gzip = siteAPI.Site.Gzip
			s.OwnerID = uid
			s.EditorID = uid
			s.Name = "[system]" + siteAPI.Site.Subdomain
			s.Desc = "[system]根据api自动保存site信息"
			mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).Insert(s)
		}
	} else {
		err = Get(s, uid)
		if err == nil {
			siteAPI.Site.Subdomain = s.Subdomain
			siteAPI.Site.IsCustomDomain = s.IsCustomDomain
			siteAPI.Site.APIDomain = s.APIDomain
			siteAPI.Site.DstURL = s.DstURL
			siteAPI.Site.Group = s.Group
			siteAPI.Site.HTTPS = s.HTTPS
			siteAPI.Site.Gzip = s.Gzip
		} else {
			log.Printf("[error]serivces.sites.SaveAPI: get site err=%v, data=%v\r\n", err, siteAPI)
			err = errors.New(services.ErrorRead)
			return
		}
	}

	if siteAPI.URL == "" {
		siteAPI.URL = siteAPI.Site.Group + siteAPI.ShortURL
	}

	if siteAPI.ID == "" {
		err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
			Find(bson.M{"method": siteAPI.Method, "site.apiDomain": siteAPI.Site.APIDomain, "url": siteAPI.URL}).
			Select(services.SelectHide).
			One(&siteAPI)

		if err != nil {
			siteAPI.ID = bson.NewObjectId()
			siteAPI.CreatedAt = siteAPI.UpdatedAt
		} else if siteAPI.AutoReg {
			siteAPI.Visited = siteAPI.Visited + 1
		}
	}

	if siteAPI.OwnerID == "" {
		siteAPI.OwnerID = uid
	}

	_, err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
		Upsert(bson.M{"_id": siteAPI.ID}, siteAPI)
	if err != nil {
		siteAPI.ID = ""
		log.Printf("[error]serivces.sites.SaveAPI: err=%v, data=%v\r\n", err, siteAPI)
		err = errors.New(services.ErrorSave)
	}

	return
}

// SaveApis 保存多个api
func SaveApis(apis []interface{}) (err error) {
	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).Insert(apis...)
	if err != nil {
		log.Printf("[error]serivces.sites.SaveApis: err=%v, data=%s\r\n", err, apis)
		err = errors.New(services.ErrorSave)
	}

	return
}

// GetAPI 获取api信息
func GetAPI(siteAPI *models.SiteAPI, uid string) (err error) {
	if siteAPI.ID == "" {
		err = errors.New(services.ErrorParam)
		return
	}

	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
		Find(bson.M{"_id": siteAPI.ID, "ownerId": uid}).
		Select(services.SelectHide).
		One(&siteAPI)

	if err != nil {
		log.Printf("[error]serivces.sites.GetAPI: err=%v, data=%v\r\n", err, siteAPI)
		err = errors.New(services.ErrorPermission)
	}

	return
}

// DelAPI 删除api
func DelAPI(id, uid string) (err error) {
	if id == "" {
		err = errors.New(services.ErrorParam)
		return
	}

	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
		Remove(bson.M{"_id": bson.ObjectIdHex(id), "ownerId": uid})

	if err != nil {
		log.Printf("[error]serivces.sites.DelAPI: err=%v, id=%s\r\n", err, id)
		err = errors.New(services.ErrorPermission)
	}

	return
}

// APIList api列表
func APIList(uid string, siteID bson.ObjectId, autoReg string, fieldType, pageIndex, pageCount int) (apis []models.SiteAPI, err error) {
	if siteID == "" {
		err = errors.New(services.ErrorParam)
		return
	}

	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	where := bson.M{"ownerId": uid, "site._id": siteID, "deletedAt": bson.M{"$exists": false}}
	if autoReg == "true" {
		where["autoReg"] = true
	} else if autoReg == "false" {
		where["autoReg"] = bson.M{"$ne": true}
	}

	selected := bson.M{
		"_id":        true,
		"name":       true,
		"site._id":   true,
		"site.group": true,
		"shortUrl":   true,
		"url":        true,
		"method":     true,
		"visited":    true,
		"createdAt":  true,
		"updatedAt":  true,
	}
	if fieldType == 1 {
		selected = bson.M{"_id": true, "name": true, "site._id": true, "site.group": true}
	}

	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
		Find(where).
		Select(selected).
		Sort("-updatedAt").Skip((pageIndex - 1) * pageCount).Limit(pageCount).
		All(&apis)

	if err != nil {
		log.Printf("[error]serivces.sites.APIList: err=%v, siteID=%v\r\n", err, siteID)
		err = errors.New(services.ErrorPermission)
	}

	return
}

// GetAPIByURL 请求代理时，通过url查找api
func GetAPIByURL(apiDomain, method, url string) (siteAPI *models.SiteAPI, err error) {
	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	js := `function(){return this.urlReg && new RegExp(this.urlReg).test("` + url + `")}`

	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
		Find(bson.M{
			"method":         method,
			"site.apiDomain": apiDomain,
			"$or": []bson.M{
				bson.M{"url": url},
				bson.M{
					"$where": bson.JavaScript{
						Code: js,
					},
				},
			},
			"autoReg": bson.M{"$ne": true},
		}).
		Select(services.SelectHide).
		One(&siteAPI)

	if err != nil || siteAPI.URL == "" {
		err = errors.New(services.ErrorRead)
	} else {
		// 计数+1
		mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
			Update(bson.M{"_id": siteAPI.ID}, bson.M{"$inc": bson.M{"visited": 1}})
	}

	if siteAPI != nil && (siteAPI.Pause || siteAPI.Site.Pause) {
		err = errors.New(services.ErrorAPIPause)
	}

	return
}

// GetAPIByDstURL 通过url和dstUrl查找api
func GetAPIByDstURL(dstURL, method, URL string) (siteAPI *models.SiteAPI, err error) {
	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
		Find(bson.M{"method": method, "site.dstUrl": dstURL, "url": URL}).
		Select(services.SelectHide).
		One(&siteAPI)

	if err != nil {
		err = errors.New(services.ErrorRead)
	} else {
		// 计数+1
		mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
			Update(bson.M{"_id": siteAPI.ID}, bson.M{"$inc": bson.M{"visited": 1}})
	}

	return
}

// GetSiteByDomain 通过域名查找网站信息（请求代理时用）
func GetSiteByDomain(apiDomain string) (site *models.Site, err error) {
	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).
		Find(bson.M{"apiDomain": apiDomain}).
		Select(services.SelectHide).
		One(&site)

	if err != nil {
		err = errors.New(services.ErrorRead)
	}

	if site != nil && site.Pause {
		err = errors.New(services.ErrorAPIPause)
	}

	return
}

// APIListByDomains 查找指定域名下的api
func APIListByDomains(uid string, domain []string, autoReg string, fieldType, pageIndex, pageCount int) (apis []models.SiteAPI, err error) {
	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	where := bson.M{"ownerId": uid, "site.apiDomain": bson.M{"$in": domain}, "deletedAt": bson.M{"$exists": false}}
	if autoReg == "true" {
		where["autoReg"] = true
	} else if autoReg == "false" {
		where["autoReg"] = bson.M{"$ne": true}
	}

	selected := bson.M{
		"_id":        true,
		"name":       true,
		"site._id":   true,
		"site.group": true,
		"shortUrl":   true,
		"url":        true,
		"method":     true,
		"visited":    true,
		"createdAt":  true,
		"updatedAt":  true,
	}
	if fieldType == 1 {
		selected = bson.M{"_id": true, "name": true, "site._id": true, "site.group": true, "site.apiDomain": true}
	}

	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
		Find(where).
		Select(selected).
		Sort("-updatedAt").Skip((pageIndex - 1) * pageCount).Limit(pageCount).
		All(&apis)

	if err != nil {
		log.Printf("[error]serivces.sites.APIListByDomains: err=%v, domain=%s\r\n", err, domain)
		err = errors.New(services.ErrorPermission)
	}

	return
}
