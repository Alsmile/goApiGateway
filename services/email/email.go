package email

import (
	"log"
	"strings"

	"github.com/alsmile/goApiGateway/models"
	"github.com/alsmile/goApiGateway/utils"
	"gopkg.in/gomail.v2"
)

// SendEmail 发送邮件
func SendEmail(email, title, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", utils.GlobalConfig.Email.User)
	m.SetHeader("To", email)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(utils.GlobalConfig.Email.Address, utils.GlobalConfig.Email.Port, utils.GlobalConfig.Email.User, utils.GlobalConfig.Email.Password)
	err := d.DialAndSend(m)
	if err != nil {
		log.Printf("[email]error: %v; config=%v, email=%s, title=%s\r\n", err.Error(), utils.GlobalConfig.Email, email, title)
	}
	return err
}

// SendSignUpEmail 发送注册邮件
func SendSignUpEmail(u *models.User) error {
	bodyStr, err := utils.ReadFile("./assets/signup.html")
	if err != nil {
		return err
	}

	bodyStr = strings.Replace(bodyStr, "{{name}}", u.Profile.Username, -1)
	bodyStr = strings.Replace(bodyStr, "{{signupUrl}}", utils.GlobalConfig.Website+"/?active="+u.Active.Code, -1)
	go SendEmail(u.Profile.Email, utils.GlobalConfig.Name, bodyStr)

	return nil
}

// SendForgetPasswordEmail 发送忘记密码邮件
func SendForgetPasswordEmail(u *models.User) error {
	bodyStr, err := utils.ReadFile("./assets/forgetPassword.html")
	if err != nil {
		return err
	}

	bodyStr = strings.Replace(bodyStr, "{{name}}", u.Profile.Username, -1)
	bodyStr = strings.Replace(bodyStr, "{{signupUrl}}", utils.GlobalConfig.Website+"/?forgetPassword="+u.PasswordCode+
		"&email="+u.Profile.Email, -1)
	go SendEmail(u.Profile.Email, utils.GlobalConfig.Name, bodyStr)

	return nil
}
