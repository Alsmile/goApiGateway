package user

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
	"gopkg.in/mgo.v2/bson"

	"github.com/alsmile/goApiGateway/db/mongo"
	"github.com/alsmile/goApiGateway/db/redis"
	"github.com/alsmile/goApiGateway/models"
	"github.com/alsmile/goApiGateway/services"
	"github.com/alsmile/goApiGateway/services/email"
	"github.com/alsmile/goApiGateway/utils"
)

// EncodePassword 密码加盐
func EncodePassword(pwd string) string {
	mac := hmac.New(sha256.New, []byte(utils.GlobalConfig.Secret))
	mac.Write([]byte(pwd))
	str := hex.EncodeToString(mac.Sum([]byte("goApiGateway")))

	return str
}

// AddUser 添加用户
func AddUser(u *models.User) (err error) {
	if u.Profile.Email == "" && u.Profile.Phone == "" {
		err = errors.New(services.ErrorEmailEmpty)
		return
	}

	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionUsers).Find(
		bson.M{"profile.email": u.Profile.Email}).Select(bson.M{"_id": 1}).One(&u)

	if u.ID != "" {
		err = errors.New(services.ErrorEmailExists)
		return
	}

	u.ID = bson.NewObjectId()
	u.Active.Code = bson.NewObjectId().Hex()
	u.Profile.GetUsername()
	u.Password = EncodePassword(u.Password)
	u.CreatedAt = time.Now().UTC()

	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionUsers).Insert(u)
	u.Password = ""
	if err != nil {
		if strings.Contains(err.Error(), "unique_email") {
			err = errors.New(services.ErrorEmailExists)
		} else {
			log.Printf("services.user.user.AddUser error: %v\r\n", err)
			err = errors.New(services.ErrorSave)
		}
	} else {
		email.SendSignUpEmail(u)
	}

	return
}

// GetUserByPassword 通过邮箱和密码查找用户
func GetUserByPassword(u *models.User) (err error) {
	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionUsers).Find(
		bson.M{
			"profile.email": u.Profile.Email,
			"password":      EncodePassword(u.Password),
		}).Select(bson.M{"_id": 1, "profile": 1, "roles": 1, "active": 1}).One(&u)
	u.Password = ""
	if err != nil {
		log.Printf("[error]services.user.GetUserByPassword: err=%v\r\n", err)
		err = errors.New("用户名或密码错误")
	} else if u.Active.Code != "" {
		err = errors.New(services.ErrorNoActive)
	}

	return
}

// Active 用户激活
func Active(u *models.User) (err error) {
	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionUsers).Find(
		bson.M{"active.code": u.Active.Code}).Select(bson.M{"_id": 1, "profile": 1, "roles": 1, "active": 1}).One(&u)
	if err != nil {
		err = errors.New(services.ErrorActiveCode)
		return
	}

	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionUsers).
		Update(bson.M{"active.code": u.Active.Code}, bson.M{"$set": bson.M{"active": bson.M{"time": time.Now().UTC()}}})
	if err != nil {
		err = errors.New(services.ErrorSave)
		return
	}

	return
}

// ForgetPassword 忘记密码请求
func ForgetPassword(u *models.User) (err error) {
	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionUsers).Find(
		bson.M{
			"profile.email": u.Profile.Email,
		}).Select(bson.M{"_id": 1, "profile": 1}).One(&u)
	if err != nil || u.ID == "" {
		err = errors.New(services.ErrorUserNoExists)
		return
	}

	redisConn := redis.RedisPool.Get()
	defer redisConn.Close()

	u.PasswordCode = bson.NewObjectId().Hex()
	_, err = redisConn.Do("SETEX", u.PasswordCode, services.TokenValidHours*3600, u.ID.Hex())
	if err != nil {
		log.Printf("services.user.user.ForgetPassword: redis error=%v\r\n", err)
		err = errors.New(services.ErrorSave)
		return
	}
	email.SendForgetPasswordEmail(u)

	return
}

// NewPassword 忘记密码时设置新密码
func NewPassword(u *models.User) (err error) {
	redisConn := redis.RedisPool.Get()
	defer redisConn.Close()

	val, err := redisConn.Do("GET", u.PasswordCode)
	if err != nil {
		err = errors.New(services.ErrorCaptchaCode)
		return
	}
	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	u.ID = bson.ObjectIdHex(utils.String(val))
	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionUsers).
		Update(bson.M{"_id": u.ID}, bson.M{"$set": bson.M{"password": EncodePassword(u.Password)}})
	u.Password = ""
	if err != nil {
		err = errors.New(services.ErrorUserNoExists)
	}

	return
}

// GetUserByID 根据用户id查询用户信息
func GetUserByID(u *models.User) (err error) {
	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	err = mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionUsers).Find(bson.M{"_id": u.ID}).
		Select(bson.M{"_id": 1, "profile": 1, "roles": 1}).One(&u)
	if err != nil {
		err = errors.New(services.ErrorUserNoExists)
	}

	return
}

// GetToken 根据小时数，计算一个有效时长的token
func GetToken(u *models.User, hours int) (data string) {
	if u == nil {
		return
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"uid": u.ID,
		"exp": time.Now().Add(time.Hour * time.Duration(hours)).Unix(),
	})
	data, err := token.SignedString([]byte(utils.GlobalConfig.Jwt))
	if err != nil {
		log.Printf("service.user.GetToken: error= %s", err)
	}
	return
}

// ValidToken 通过token校验用户身份
func ValidToken(ctx iris.Context) (uid string) {
	data := ctx.GetHeader("Authorization")
	if data == "" {
		data = ctx.GetHeader("token")
	}
	if data == "" {
		ctx.StatusCode(401)
		return
	}

	token, err := jwt.Parse(data, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("签名方法错误: %v", token.Header["alg"])
		}
		return []byte(utils.GlobalConfig.Jwt), nil
	})

	if err != nil {
		// java平台采用的是base64编码的secret
		// secret base64 encoded
		token, err = jwt.Parse(data, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("签名方法错误: %v", token.Header["alg"])
			}
			return base64.StdEncoding.DecodeString(utils.GlobalConfig.Jwt)
		})
	}

	if err != nil {
		log.Printf("service.user.ValidToken.Parse: error= %s", err)
		return
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return
	}

	uid = utils.String(claims["uid"])

	return
}

// GetUserByTokenID 通过token中的id（此id可能时系统内部id，也可能时外部id）获取用户信息
func GetUserByTokenID(token, uid string) (*models.User, error) {
	if uid == "" {
		return nil, errors.New(services.ErrorNeedSign)
	}

	mongoSession := mongo.MgoSession.Clone()
	defer mongoSession.Close()

	user := new(models.User)
	if bson.IsObjectIdHex(uid) {
		user.ID = bson.ObjectIdHex(uid)
		err := mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionUsers).Find(
			bson.M{"_id": user.ID}).Select(bson.M{"_id": 1, "profile": 1, "roles": 1}).One(&user)

		if err == nil {
			return user, nil
		}
	}

	err := mongoSession.DB(utils.GlobalConfig.Mongo.Database).C(mongo.CollectionUsers).Find(
		bson.M{"profile.userId": uid}).Select(bson.M{"_id": 1, "profile": 1, "roles": 1}).One(&user)

	if err == nil {
		return user, nil
	}

	client := &http.Client{}
	clientReq, err := http.NewRequest("GET", utils.GlobalConfig.User.ProfileURL, nil)
	if err != nil {
		return nil, errors.New(services.ErrorRead)
	}

	clientReq.Header.Set("Content-Type", "application/json")
	clientReq.Header.Set("Authorization", token)
	clientReq.Header.Set("token", token)
	clientResp, err := client.Do(clientReq)
	if err != nil {
		return nil, errors.New(services.ErrorRead)
	}

	defer clientResp.Body.Close()

	if clientResp.StatusCode != 200 {
		return nil, errors.New(services.ErrorProxyNotFound)
	}

	body, err := ioutil.ReadAll(clientResp.Body)
	if err != nil {
		return nil, errors.New(services.ErrorRead)
	}

	if utils.GlobalConfig.User.ProfileType == 0 {
		tempUser := struct {
			Code string `json:"code"`
			Data struct {
				Nickname string `json:"nickname"`
				UserID   string `json:"userId"`
				Phone    string `json:"userMobile"`
				Email    string `json:"userEmail"`
			} `json:"data"`
			Message string `json:"message"`
		}{}

		if err = json.Unmarshal(body, &tempUser); err == nil {
			if tempUser.Code == "0" || tempUser.Code == "" {
				user.Profile.Username = tempUser.Data.Nickname
				user.Profile.Nickname = tempUser.Data.Nickname
				user.Profile.Email = tempUser.Data.Email
				user.Profile.Phone = tempUser.Data.Phone
				user.Profile.UserID = tempUser.Data.UserID
			} else {
				return user, errors.New(tempUser.Message)
			}

			return user, nil
		}
	}

	return nil, errors.New(services.ErrorUserNoExists)
}
