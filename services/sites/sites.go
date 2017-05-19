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
    Select(bson.M{"domain": true, "gzip" : true,  "https" : true, "notFound": true, "statics":true, "proxies": true}).
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
    ProxyKey: site.ProxyKey,
    ProxyValue: site.ProxyValue,
  }}})

  return
}

func SaveApi(siteApi *models.SiteApi) (err error) {
  siteApi.UpdatedAt = time.Now().UTC()
  if siteApi.Id == "" {
    siteApi.Id = bson.NewObjectId()
    siteApi.CreatedAt = siteApi.UpdatedAt
  }

  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  s := &models.Site{Id: siteApi.Site.Id}
  //  site id不存在，表示自动根据api信息保存site
  if siteApi.Site.Id == "" {
    s.ProxyValue = siteApi.Site.ProxyValue
    s.ProxyKey = siteApi.Site.ProxyKey
    s.Https = siteApi.Site.Https
    s.Gzip = siteApi.Site.Gzip
    s.Owner = siteApi.Owner
    s.Editor = siteApi.Editor
    s.Desc = "[system]根据api自动保存site信息"
    Save(s)
  }

  err = Get(s)
  if err == nil {
    siteApi.Site.ProxyValue = s.ProxyValue
    siteApi.Site.ProxyKey = s.ProxyKey
    siteApi.Site.Https = s.Https
    siteApi.Site.Gzip = s.Gzip
  } else {
    log.Printf("[error]serivces.sites.SaveApi: get site err=%v, data=%s\r\n", err, siteApi)
    err = errors.New(services.ErrorRead)
    return
  }

  _, err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).Upsert(bson.M{"_id": siteApi.Id}, siteApi)
  if err != nil {
    log.Printf("[error]serivces.sites.SaveApi: err=%v, data=%s\r\n", err, siteApi)
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

func ApiList(siteId bson.ObjectId, pageIndex, pageCount int) (apis []models.SiteApi,err error) {
  if siteId == "" {
    err = errors.New(services.ErrorPermission)
    return
  }

  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  selected := bson.M{"_id": true, "name": true, "site._id": true, "site.proxyKey": true}
  err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
    Find(bson.M{"site._id": siteId, "deletedAt": bson.M{"$exists": false}}).
    Select(selected).
    Sort("-updatedAt").Skip((pageIndex-1)*pageCount).Limit(pageCount).
    All(&apis)

  if err != nil {
    log.Printf("[error]serivces.sites.ApiList: err=%v, data=%s\r\n", err, siteId)
    err = errors.New(services.ErrorRead)
  }

  return
}

func GetApiByUrl(subdomain, method, key, url string) (siteApi *models.SiteApi, err error) {
  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApis).
    Find(bson.M{"method": method, "site.subdomain": subdomain, "site.proxyKey": key, "url": url}).
    Select(services.SelectHide).
    One(&siteApi)

  if err != nil {
    //log.Printf("[error]serivces.sites.GetApiByUrl: err=%v, method=%s, key=%s, url=%s\r\n", err, method, key, url)
    err = errors.New(services.ErrorRead)
  }

  return
}

func GetSiteByProxyKey(subdomain, proxyKey string) (site *models.Site, err error) {
  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).
    Find(bson.M{"subdomain": subdomain, "proxyKey": proxyKey}).
    Select(services.SelectHide).
    One(&site)

  if err != nil {
    //log.Printf("[error]serivces.sites.GetSiteByProxyKey: err=%v,key=%s\r\n", err, site)
    err = errors.New(services.ErrorRead)
  }

  return
}
