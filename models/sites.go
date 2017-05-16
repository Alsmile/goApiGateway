package models

import (
  "time"
  "gopkg.in/mgo.v2/bson"
)

type Site struct {
  Id     bson.ObjectId `json:"id" bson:"_id"`
  Owner  QuotedUser `json:"owner,omitempty"`
  Editor QuotedUser `json:"editor,omitempty"`
  Name   string `json:"name"`
  Desc   string `json:"desc"`
  Domain string `json:"domain"`
  Gzip   bool `json:"gzip"`
  Https  string `json:"https"`
  NotFound struct {
    Code int `json:"code"`
    Path string `json:"path"`
  } `json:"notFound" bson:"notFound"`
  Statics   []PathUrl `json:"statics"`
  Proxies   []PathUrl `json:"proxies"`
  CreatedAt time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
  UpdatedAt time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
  DeletedAt time.Time `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}

type PathUrl struct {
  Path string `json:"path"`
  Url  string `json:"url"`
}
