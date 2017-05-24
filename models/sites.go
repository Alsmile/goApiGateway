package models

import (
  "time"
  "gopkg.in/mgo.v2/bson"
)

type Site struct {
  Id         bson.ObjectId `json:"id" bson:"_id"`
  Owner      QuotedUser `json:"owner,omitempty"`
  Editor     QuotedUser `json:"editor,omitempty"`
  Name       string `json:"name"`
  Desc       string `json:"desc"`
  Gzip       bool `json:"gzip"`
  Https      string `json:"https"`
  Subdomain  string `json:"subdomain" `
  ProxyKey   string `json:"proxyKey" bson:"proxyKey"`
  ProxyValue string `json:"proxyValue" bson:"proxyValue"`
  CreatedAt  time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
  UpdatedAt  time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
  DeletedAt  time.Time `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}

type SiteApi struct {
  Id                     bson.ObjectId `json:"id" bson:"_id"`
  Site                   SiteParam `json:"site" bson:"site"`
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
  AutoReg                bool `json:"autoReg" bson:"autoReg,omitempty"`
  Visited                uint64 `json:"visited" bson:"visited,omitempty"`
  StatusCode             int `json:"statusCode" bson:"statusCode,omitempty"`
  CreatedAt              time.Time `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
  UpdatedAt              time.Time `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
  DeletedAt              time.Time `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}

type ApiParam struct {
  Id       int    `json:"id" bson:"id,omitempty"`
  ParentId int    `json:"parentId" bson:"parentId,omitempty"`
  Name     string `json:"name"`
  Type     string `json:"type" bson:"type,omitempty"`
  Desc     string `json:"desc"`
  Required string `json:"required"`
  Mock     string `json:"mock" bson:"mock,omitempty"`
  Level    int    `json:"level" bson:"level,omitempty"`
  HasChild bool   `json:"hasChild" bson:"hasChild,omitempty"`
}

type SiteParam struct {
  Id         bson.ObjectId `json:"id" bson:"_id"`
  Gzip       bool `json:"gzip"`
  Https      string `json:"https"`
  Subdomain  string `json:"subdomain" `
  ProxyKey   string `json:"proxyKey" bson:"proxyKey"`
  ProxyValue string `json:"proxyValue" bson:"proxyValue"`
}

func (api *SiteApi) GetMockData() {
  if api.DataType == "application/json" || api.DataType == "multipart/form-data" || api.DataType == "application/x-www-form-urlencoded" {

  }
}
