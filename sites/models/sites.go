package models

import (
  "time"
  "gopkg.in/mgo.v2/bson"
)

type Site struct {
  Id        bson.ObjectId `json:"id" bson:"_id"`
  UserId    string `json:"userId" bson:"userId"`
  Name      string `json:"name"`
  Desc      string `json:"desc"`
  Domain    string `json:"domain"`
  Gzip      bool `json:"gzip"`
  Https     string `json:"https"`
  NotFound  string `json:"notFound" bson:"notFound"`
  Statics   []PathUrl `json:"statics"`
  Proxies   []PathUrl `json:"proxies"`
  CreatedAt    time.Time `json:"createdAt" bson:"createdAt,omitempty"`
  UpdatedAt    time.Time `json:"updatedAt" bson:"updatedAt,omitempty"`
  DeletedAt    time.Time `json:"deletedAt" bson:"deletedAt,omitempty"`
}

type PathUrl struct {
  Path string `json:"path"`
  Url  string `json:"url"`
}
