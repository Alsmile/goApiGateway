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
	Desc           string        `json:"desc" bson:"desc,omitempty"`
	Gzip           bool          `json:"gzip" bson:"gzip,omitempty"`
	HTTPS          string        `json:"https" bson:"https,omitempty"`
	Subdomain      string        `json:"subdomain" bson:"subdomain,omitempty"`
	IsCustomDomain bool          `json:"isCustomDomain" bson:"isCustomDomain"`
	APIDomain      string        `json:"apiDomain" bson:"apiDomain"`
	Group          string        `json:"group" bson:"group,omitempty"`
	DstURL         string        `json:"dstUrl" bson:"dstUrl"`
	Pause          bool          `json:"pause" bson:"pause,omitempty"`
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
	Desc                   string        `json:"desc" bson:"desc,omitempty"`
	Method                 string        `json:"method" bson:"method,omitempty"`
	IsMock                 bool          `json:"isMock" bson:"isMock,omitempty"`
	URLReg                 string        `json:"urlReg"  bson:"urlReg,omitempty"`
	URLParams              []*APIParam   `json:"urlParams" bson:"urlParams,omitempty"`
	Headers                []*APIParam   `json:"headers" bson:"headers,omitempty"`
	ContentType            string        `json:"contentType" bson:"contentType,omitempty"`
	DataType               string        `json:"dataType" bson:"dataType,omitempty"`
	QueryParams            []*APIParam   `json:"queryParams" bson:"queryParams,omitempty"`
	BodyParams             []*APIParam   `json:"bodyParams" bson:"bodyParams,omitempty"`
	ResponseParams         []*APIParam   `json:"responseParams" bson:"responseParams,omitempty"`
	BodyParamsText         string        `json:"bodyParamsText" bson:"bodyParamsText,omitempty"`
	BodyParamsTextDesc     string        `json:"bodyParamsTextDesc" bson:"bodyParamsTextDesc,omitempty"`
	ResponseParamsText     string        `json:"responseParamsText" bson:"responseParamsText,omitempty"`
	ResponseParamsTextDesc string        `json:"responseParamsTextDesc" bson:"responseParamsTextDesc,omitempty"`
	AutoReg                bool          `json:"autoReg" bson:"autoReg,omitempty"`
	Visited                uint64        `json:"visited" bson:"visited,omitempty"`
	StatusCode             int           `json:"statusCode" bson:"statusCode,omitempty"`
	Pause                  bool          `json:"pause" bson:"pause,omitempty"`
	CreatedAt              time.Time     `json:"createdAt,omitempty" bson:"createdAt,omitempty"`
	UpdatedAt              time.Time     `json:"updatedAt,omitempty" bson:"updatedAt,omitempty"`
	DeletedAt              time.Time     `json:"deletedAt,omitempty" bson:"deletedAt,omitempty"`
}

// APIParam 接口调用时，用到的api数据结构
type APIParam struct {
	ID       int    `json:"id" bson:"id,omitempty"`
	ParentID int    `json:"parentId" bson:"parentId,omitempty"`
	Name     string `json:"name" bson:"name,omitempty"`
	Type     string `json:"type" bson:"type,omitempty"`
	Desc     string `json:"desc" bson:"desc,omitempty"`
	Required string `json:"required" bson:"required,omitempty"`
	Mock     string `json:"mock" bson:"mock,omitempty"`
	Level    int    `json:"level" bson:"level,omitempty"`
	HasChild bool   `json:"hasChild" bson:"hasChild,omitempty"`
}

// SiteParam 接口调用时，用到的网站数据
type SiteParam struct {
	ID             bson.ObjectId `json:"id" bson:"_id"`
	Gzip           bool          `json:"gzip" bson:"gzip,omitempty"`
	HTTPS          string        `json:"https" bson:"https,omitempty"`
	Subdomain      string        `json:"subdomain" bson:"subdomain,omitempty"`
	IsCustomDomain bool          `json:"isCustomDomain" bson:"isCustomDomain,omitempty"`
	APIDomain      string        `json:"apiDomain" bson:"apiDomain,omitempty"`
	Group          string        `json:"group" bson:"group,omitempty"`
	DstURL         string        `json:"dstUrl" bson:"dstUrl,omitempty"`
	Pause          bool          `json:"pause" bson:"pause,omitempty"`
}
