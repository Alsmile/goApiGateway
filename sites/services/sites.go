package services

import (
  "github.com/alsmile/goMicroServer/sites/models"
  "github.com/alsmile/goMicroServer/services"
  "github.com/alsmile/goMicroServer/db/mongo"
  "errors"
  "log"
  "gopkg.in/mgo.v2/bson"
  "github.com/alsmile/goMicroServer/utils"
  "time"
)

func Save(site *models.Site) (err error) {
  if site.Id == "" {
    site.Id = bson.NewObjectId()
    site.CreatedAt = time.Now().UTC()
  } else {
    site.UpdatedAt = time.Now().UTC()
  }

  mongoSession := mongo.MgoSession.Clone()
  defer mongoSession.Close()

  _, err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionSites).Upsert(bson.M{"_id": site.Id}, site)
  if err != nil {
    log.Printf("[error]sites.serivces.sites.Save: err=%v\r\n", err)
    err = errors.New(services.ErrorSave)
  }

  return
}
