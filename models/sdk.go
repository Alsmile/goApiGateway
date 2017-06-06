package models

import (
  "gopkg.in/mgo.v2/bson"
)

type SdkSite struct {
  Id             bson.ObjectId `json:"id" bson:"_id"`
  Name           string `json:"name"`
  Gzip           bool `json:"gzip"`
  Https          string `json:"https"`
  ApiDomain      string `json:"apiDomain" `
  Group          string `json:"group" `
  DstUrl         string `json:"dstUrl" bson:"dstUrl"`
  Apis           string `json:"apis"`
}
