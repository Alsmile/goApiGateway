package sites

import (
  "github.com/alsmile/goApiGateway/models"
  "github.com/alsmile/goApiGateway/services"
  "github.com/alsmile/goApiGateway/db/mongo"
  "errors"
  "log"
  "gopkg.in/mgo.v2/bson"
  "github.com/alsmile/goApiGateway/utils"
  "time"
)

func ListAll() (sites []models.Site, err error) {
  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).
    Find(bson.M{"deletedAt": bson.M{"$exists": false}}).
    Select(bson.M{"apiDomain": true, "gzip" : true,  "https" : true, "notFound": true, "statics":true, "proxies": true}).
    All(&sites)

  if err != nil {
    log.Printf("[error]serivces.sites.ListAll: err=%v\r\n", err)
    err = errors.New(services.ErrorRead)
  }

  return
}

func List(uid bson.ObjectId, pageIndex, pageCount int) (sites []models.Site, err error) {
  if uid == "" {
    err = errors.New(services.ErrorPermission)
    return
  }

  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).
    Find(bson.M{"owner._id": uid, "deletedAt": bson.M{"$exists": false}}).
    Select(services.SelectHide).
    Sort("-updatedAt").Skip((pageIndex-1)*pageCount).Limit(pageCount).
    All(&sites)

  if err != nil {
    log.Printf("[error]serivces.sites.List: err=%v, data=%s\r\n", err, uid)
    err = errors.New(services.ErrorRead)
  }

  return
}

func Get(site *models.Site) (err error) {
  if site.Id == "" {
    err = errors.New(services.ErrorParam)
    return
  }

  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).
    Find(bson.M{"_id": site.Id}).
    Select(services.SelectHide).
    One(&site)

  if err != nil {
    log.Printf("[error]serivces.sites.Get: err=%v, data=%s\r\n", err, site)
    err = errors.New(services.ErrorRead)
  }

  return
}

func GetByGroup(site *models.Site) (err error) {
  if site.Id == "" {
    err = errors.New(services.ErrorParam)
    return
  }

  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).
    Find(bson.M{"apiDomain": site.ApiDomain, "group": site.Group}).
    Select(services.SelectHide).
    One(&site)

  if err != nil {
    log.Printf("[error]serivces.sites.GetByGroup: err=%v, data=%s\r\n", err, site)
    err = errors.New(services.ErrorRead)
  }

  return
}

func Save(site *models.Site) (err error) {
  site.UpdatedAt = time.Now().UTC()
  if site.Id == "" {
    site.Id = bson.NewObjectId()
    site.CreatedAt = site.UpdatedAt
  }

  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  _, err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).Upsert(bson.M{"_id": site.Id}, site)

  if err != nil {
    log.Printf("[error]serivces.sites.Save: err=%v, data=%s\r\n", err, site)
    err = errors.New(services.ErrorSave)
  }

  mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
    UpdateAll(bson.M{"site._id": site.Id}, bson.M{"$set": bson.M{"site": models.SiteParam{
    Id: site.Id,
    Gzip: site.Gzip,
    Https: site.Https,
    Group: site.Group,
    Subdomain: site.Subdomain,
    IsCustomDomain: site.IsCustomDomain,
    ApiDomain: site.ApiDomain,
    DstUrl: site.DstUrl,
  }}})

  return
}

func DelSite(site *models.Site) (err error) {
  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  _, err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).
    RemoveAll(bson.M{"apiDomain": site.ApiDomain, "group": site.Group})

  if err != nil {
    log.Printf("[error]serivces.sites.DelSite - CollectionSites: err=%v, site=%s\r\n", err, site)
    err = errors.New(services.ErrorSave)
  }

  _, errApis := mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
    RemoveAll(bson.M{"site.apiDomain": site.ApiDomain, "site.group": site.Group})

  if errApis != nil {
    log.Printf("[error]serivces.sites.DelSite - CollectionApis: err=%v, site=%s\r\n", err, site)
    err = errors.New(services.ErrorSave)
  }

  return
}

func SaveApi(siteApi *models.SiteApi) (err error) {
  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  siteApi.UpdatedAt = time.Now().UTC()

  s := &models.Site{Id: siteApi.Site.Id}

  //  site id不存在，表示自动根据api信息保存site
  if siteApi.Id == "" && siteApi.Site.Id == "" && siteApi.Site.Subdomain != ""{
    tempSite := bson.M{}
    err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).
      Find(bson.M{"subdomain": siteApi.Site.Subdomain}).
      Select(bson.M{"_id": true}).
      One(&tempSite)


    // 不存在site，插入保存
    if err != nil {
      s.UpdatedAt = time.Now().UTC()
      s.Id = bson.NewObjectId()
      siteApi.Site.Id = s.Id
      s.CreatedAt = s.UpdatedAt
      s.Subdomain = siteApi.Site.Subdomain
      s.IsCustomDomain = siteApi.Site.IsCustomDomain
      s.ApiDomain = siteApi.Site.ApiDomain
      s.Group = siteApi.Site.Group
      s.DstUrl = siteApi.Site.DstUrl
      s.Https = siteApi.Site.Https
      s.Gzip = siteApi.Site.Gzip
      s.Owner = siteApi.Owner
      s.Editor = siteApi.Editor
      s.Name = "[system]" + siteApi.Site.Subdomain
      s.Desc = "[system]根据api自动保存site信息"
      mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).Insert(s)
    }
  } else {
    err = Get(s)
    if err == nil {
      siteApi.Site.Subdomain = s.Subdomain
      siteApi.Site.IsCustomDomain = s.IsCustomDomain
      siteApi.Site.ApiDomain = s.ApiDomain
      siteApi.Site.DstUrl = s.DstUrl
      siteApi.Site.Group = s.Group
      siteApi.Site.Https = s.Https
      siteApi.Site.Gzip = s.Gzip
    } else {
      log.Printf("[error]serivces.sites.SaveApi: get site err=%v, data=%s\r\n", err, siteApi)
      err = errors.New(services.ErrorRead)
      return
    }
  }

  if siteApi.Id == "" {
    siteApi.Id = bson.NewObjectId()
    siteApi.CreatedAt = siteApi.UpdatedAt
  }
  _, err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).Upsert(bson.M{"_id": siteApi.Id}, siteApi)
  if err != nil {
    siteApi.Id = ""
    log.Printf("[error]serivces.sites.SaveApi: err=%v, data=%s\r\n", err, siteApi)
    err = errors.New(services.ErrorSave)
  }

  return
}

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

func GetApi(siteApi *models.SiteApi) (err error) {
  if siteApi.Id == "" {
    err = errors.New(services.ErrorParam)
    return
  }

  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
    Find(bson.M{"_id": siteApi.Id}).
    Select(services.SelectHide).
    One(&siteApi)

  if err != nil {
    log.Printf("[error]serivces.sites.GetApi: err=%v, data=%s\r\n", err, siteApi)
    err = errors.New(services.ErrorRead)
  }

  return
}

func DelApi(id string) (err error) {
  if id == "" {
    err = errors.New(services.ErrorParam)
    return
  }

  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).RemoveId(bson.ObjectIdHex(id))

  if err != nil {
    log.Printf("[error]serivces.sites.DelApi: err=%v, id=%s\r\n", err, id)
    err = errors.New(services.ErrorSave)
  }

  return
}

func ApiList(siteId bson.ObjectId, autoReg string, fieldType, pageIndex, pageCount int) (apis []models.SiteApi,err error) {
  if siteId == "" {
    err = errors.New(services.ErrorPermission)
    return
  }

  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  where := bson.M{"site._id": siteId, "deletedAt": bson.M{"$exists": false}}
  if autoReg == "true" {
    where["autoReg"] = true
  } else if autoReg == "false" {
    where["autoReg"] = bson.M{"$ne": true}
  }

  selected := bson.M{
    "_id": true,
    "name": true,
    "site._id": true,
    "site.group": true,
    "shortUrl": true,
    "url": true,
    "method": true,
    "visited": true,
    "createdAt": true,
    "updatedAt": true,
  }
  if fieldType == 1 {
    selected = bson.M{"_id": true, "name": true, "site._id": true, "site.group": true}
  }

  err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
    Find(where).
    Select(selected).
    Sort("-updatedAt").Skip((pageIndex-1)*pageCount).Limit(pageCount).
    All(&apis)

  if err != nil {
    log.Printf("[error]serivces.sites.ApiList: err=%v, data=%s\r\n", err, siteId)
    err = errors.New(services.ErrorRead)
  }

  return
}

func GetApiByUrl(apiDomain, method, url string) (siteApi *models.SiteApi, err error) {
  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
    Find(bson.M{"method": method, "site.apiDomain": apiDomain, "url": url}).
    Select(services.SelectHide).
    One(&siteApi)

  if err != nil {
    err = errors.New(services.ErrorRead)
  } else {
    // 计数+1
    mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
      Update(bson.M{"_id": siteApi.Id}, bson.M{"$inc": bson.M{"visited": 1}})
  }

  return
}

func GetSiteByGroup(apiDomain, group string) (site *models.Site, err error) {
  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  groups := [2]string{group, ""}
  err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).
    Find(bson.M{"apiDomain": apiDomain, "group": bson.M{"$in": groups}}).
    Select(services.SelectHide).
    One(&site)

  if err != nil {
    err = errors.New(services.ErrorRead)
  }

  return
}

func ApiListByDomains(domain []string, autoReg string, fieldType, pageIndex, pageCount int) (apis []models.SiteApi,err error) {
  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  where := bson.M{"site.apiDomain": bson.M{"$in": domain}, "deletedAt": bson.M{"$exists": false}}
  if autoReg == "true" {
    where["autoReg"] = true
  } else if autoReg == "false" {
    where["autoReg"] = bson.M{"$ne": true}
  }

  selected := bson.M{
    "_id": true,
    "name": true,
    "site._id": true,
    "site.group": true,
    "shortUrl": true,
    "url": true,
    "method": true,
    "visited": true,
    "createdAt": true,
    "updatedAt": true,
  }
  if fieldType == 1 {
    selected = bson.M{"_id": true, "name": true, "site._id": true, "site.group": true}
  }

  err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
    Find(where).
    Select(selected).
    Sort("-updatedAt").Skip((pageIndex-1)*pageCount).Limit(pageCount).
    All(&apis)

  if err != nil {
    log.Printf("[error]serivces.sites.ApiListByDomains: err=%v, data=%s\r\n", err, domain)
    err = errors.New(services.ErrorRead)
  }

  return
}
