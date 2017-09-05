package models

import (
	"gopkg.in/mgo.v2/bson"
)

// SdkSite 平台内部sdk调用
type SdkSite struct {
	ID        bson.ObjectId `json:"id" bson:"_id"`
	Name      string        `json:"name"`
	Gzip      bool          `json:"gzip"`
	HTTPS     string        `json:"https"`
	APIDomain string        `json:"apiDomain" `
	Group     string        `json:"group" `
	DstURL    string        `json:"dstUrl" bson:"dstUrl"`
	Apis      string        `json:"apis"`
}
