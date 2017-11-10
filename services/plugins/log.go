package plugins

import (
	"time"

	"github.com/alsmile/goApiGateway/models"

	"github.com/alsmile/goApiGateway/db/mongo"
	"github.com/alsmile/goApiGateway/utils"

	"gopkg.in/mgo.v2/bson"
)

// Log 记录日志
func Log(host, method, url, remoteAddr, dstURL string, difference time.Duration, apiID, siteID bson.ObjectId) {
	log := &models.APILog{
		ID:         bson.NewObjectId(),
		APIID:      apiID,
		SiteID:     siteID,
		Host:       host,
		Method:     method,
		URL:        url,
		DstURL:     dstURL,
		IP:         remoteAddr,
		Difference: difference.String(),
		CreatedAt:  time.Now().UTC(),
	}

	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()
	mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionApiLogs).Insert(log)
}
