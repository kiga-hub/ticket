package email

import (
	"net/mail"
	"Two-Card/models"
	"strconv"

	"github.com/astaxie/beego"
)

func GenerateEmail() (*Message, error) {
	// compose the message
	m := NewMessage("Hi", "this is the body")
	// This is cantained in the message, which has no effect on email sending.
	m.From = mail.Address{Name: "徐睿祺", Address: beego.AppConfig.String("mail::host")}
	m.AddTo(mail.Address{Name: "金泉", Address: "jinquan@aithu.com"})

	// add attachments
	/*
		if err := m.Attach("agent.go"); err != nil {
			return nil, err
		}
	*/

	// add headers
	m.AddHeader("X-CUSTOMER-id", "xxxxx")
	return m, nil
}

func SendMailWithZip() error {
	msg, err := GenerateEmail()
	if err != nil {
		return err
	}
	Host := beego.AppConfig.String("mail::host")
	Port, _ := strconv.Atoi(beego.AppConfig.String("mail::port"))
	User := beego.AppConfig.String("mail::user")
	Pwd := beego.AppConfig.String("mail::password")
	agent := New(User, Pwd, Host, Port, true)
	if err := agent.SendEmail(msg); err != nil {
		return err
	}
	return nil
}
func GenerateEmailMsg(title string, message string, FromUser *models.User, ToUser *models.User) (*Message, error) {
	m := NewMessage(title, message)
	m.From = mail.Address{Name: (FromUser.UserName), Address: beego.AppConfig.String("mail::host")}
	m.AddTo(mail.Address{Name: (ToUser.UserName), Address: ToUser.Email})

	m.AddHeader("X-CUSTOMER-id", "xxxxx")
	return m, nil
}

func SendMailWithMessage(title string, message string, FromUser *models.User, ToUser *models.User) error {
	msg, err := GenerateEmailMsg(title,message,FromUser,ToUser)
	if err != nil {
		return err
	}
	Host := beego.AppConfig.String("mail::host")
	Port, _ := strconv.Atoi(beego.AppConfig.String("mail::port"))
	User := beego.AppConfig.String("mail::user")
	Pwd := beego.AppConfig.String("mail::password")
	agent := New(User, Pwd, Host, Port, true)
	if err := agent.SendEmail(msg); err != nil {
		return err
	}
	return nil
}

func TestPlainAuth() error {
	msg, err := GenerateEmail()
	if err != nil {
		return err
	}
	agent := New("exmaple@aliyun.com", "password", "smtp.aliyun.com", 25, false)
	//agent := emailagent.NewWithIdentify("identify","exmaple@aliyun.com", "password", "smtp.aliyun.com", 25, false)
	if err = agent.SendEmail(msg); err != nil {
		return err
	}
	return nil
}
