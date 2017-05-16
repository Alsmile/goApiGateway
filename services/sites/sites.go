package sites

import (
  "github.com/alsmile/goMicroServer/models"
  "github.com/alsmile/goMicroServer/services"
  "github.com/alsmile/goMicroServer/db/mongo"
  "errors"
  "log"
  "gopkg.in/mgo.v2/bson"
  "github.com/alsmile/goMicroServer/utils"
  "time"
)

func List(uid bson.ObjectId, pageIndex, pageCount int) (sites []models.Site,err error) {
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

  return
}
