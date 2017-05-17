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

type SiteApi struct {
  Id                     bson.ObjectId `json:"id" bson:"_id"`
  SiteId                 bson.ObjectId `json:"siteId" bson:"siteId"`
  Owner                  QuotedUser `json:"owner,omitempty"`
  Editor                 QuotedUser `json:"editor,omitempty"`
  Name                   string `json:"name"`
  Url                    string `json:"url"`
  Desc                   string `json:"desc"`
  Method                 string `json:"method"`
  IsMock                 bool `json:"isMock" bson:"isMock"`
  Headers                []ApiParam `json:"headers"`
  ContentType            string `json:"contentType" bson:"contentType"`
  DataType               string `json:"dataType" bson:"dataType"`
  QueryParams            []ApiParam `json:"queryParams" bson:"queryParams"`
  BodyParams             []ApiParam `json:"bodyParams" bson:"bodyParams"`
  ResponseParams         []ApiParam `json:"responseParams" bson:"responseParams"`
  BodyParamsText         string `json:"bodyParamsText" bson:"bodyParamsText"`
  BodyParamsTextDesc     string `json:"bodyParamsTextDesc" bson:"bodyParamsTextDesc"`
  ResponseParamsText     string `json:"responseParamsText" bson:"responseParamsText"`
  ResponseParamsTextDesc string `json:"responseParamsTextDesc" bson:"responseParamsTextDesc"`
  CreatedAt              time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
  UpdatedAt              time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
  DeletedAt              time.Time `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}

type ApiParam struct {
  Id       int    `json:"id" bson:"id,omitempty"`
  Name     string `json:"name"`
  Type     string `json:"type" bson:"type,omitempty"`
  Desc     string `json:"desc"`
  Required string `json:"required"`
  Mock     string `json:"mock" bson:"mock,omitempty"`
  Level    int    `json:"level" bson:"level,omitempty"`
}
