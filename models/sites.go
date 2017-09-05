package models

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

// Site 网站数据结构
type Site struct {
	ID             bson.ObjectId `json:"id" bson:"_id"`
	OwnerID        string        `json:"ownerId,omitempty" bson:"ownerId,omitempty"`
	EditorID       string        `json:"editorId,omitempty" bson:"editorId,omitempty"`
	Name           string        `json:"name"`
	Desc           string        `json:"desc"`
	Gzip           bool          `json:"gzip"`
	HTTPS          string        `json:"https"`
	Subdomain      string        `json:"subdomain" `
	IsCustomDomain bool          `json:"isCustomDomain" bson:"isCustomDomain"`
	APIDomain      string        `json:"apiDomain" bson:"apiDomain"`
	Group          string        `json:"group" `
	DstURL         string        `json:"dstUrl" bson:"dstUrl"`
	CreatedAt      time.Time     `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt      time.Time     `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	DeletedAt      time.Time     `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}

// SiteAPI api数据结构
type SiteAPI struct {
	ID                     bson.ObjectId `json:"id" bson:"_id"`
	Site                   SiteParam     `json:"site" bson:"site"`
	OwnerID                string        `json:"ownerId,omitempty" bson:"ownerId,omitempty"`
	EditorID               string        `json:"editorId,omitempty" bson:"editorId,omitempty"`
	Name                   string        `json:"name"`
	ShortURL               string        `json:"shortUrl" bson:"shortUrl"`
	URL                    string        `json:"url"`
	Desc                   string        `json:"desc"`
	Method                 string        `json:"method"`
	IsMock                 bool          `json:"isMock" bson:"isMock"`
	Headers                []APIParam    `json:"headers"`
	ContentType            string        `json:"contentType" bson:"contentType"`
	DataType               string        `json:"dataType" bson:"dataType"`
	QueryParams            []APIParam    `json:"queryParams" bson:"queryParams"`
	BodyParams             []APIParam    `json:"bodyParams" bson:"bodyParams"`
	ResponseParams         []APIParam    `json:"responseParams" bson:"responseParams"`
	BodyParamsText         string        `json:"bodyParamsText" bson:"bodyParamsText"`
	BodyParamsTextDesc     string        `json:"bodyParamsTextDesc" bson:"bodyParamsTextDesc"`
	ResponseParamsText     string        `json:"responseParamsText" bson:"responseParamsText"`
	ResponseParamsTextDesc string        `json:"responseParamsTextDesc" bson:"responseParamsTextDesc"`
	AutoReg                bool          `json:"autoReg" bson:"autoReg,omitempty"`
	Visited                uint64        `json:"visited" bson:"visited,omitempty"`
	StatusCode             int           `json:"statusCode" bson:"statusCode,omitempty"`
	CreatedAt              time.Time     `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt              time.Time     `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	DeletedAt              time.Time     `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}

// APIParam 接口调用时，用到的api数据结构
type APIParam struct {
	ID       int    `json:"id" bson:"id,omitempty"`
	ParentID int    `json:"parentId" bson:"parentId,omitempty"`
	Name     string `json:"name"`
	Type     string `json:"type" bson:"type,omitempty"`
	Desc     string `json:"desc"`
	Required string `json:"required"`
	Mock     string `json:"mock" bson:"mock,omitempty"`
	Level    int    `json:"level" bson:"level,omitempty"`
	HasChild bool   `json:"hasChild" bson:"hasChild,omitempty"`
}

// SiteParam 接口调用时，用到的网站数据
type SiteParam struct {
	ID             bson.ObjectId `json:"id" bson:"_id"`
	Gzip           bool          `json:"gzip"`
	HTTPS          string        `json:"https"`
	Subdomain      string        `json:"subdomain" `
	IsCustomDomain bool          `json:"isCustomDomain"  bson:"isCustomDomain"`
	APIDomain      string        `json:"apiDomain" bson:"apiDomain"`
	Group          string        `json:"group" `
	DstURL         string        `json:"dstUrl" bson:"dstUrl"`
}
