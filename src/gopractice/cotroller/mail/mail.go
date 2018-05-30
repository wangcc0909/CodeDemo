package mail

import (
	"gopractice/config"
	"fmt"
	"net/smtp"
	"crypto/tls"
	"net"
)

//发送邮件
func SendEmail(toEmail, title, content string) error {

	host := config.ServerConfig.MailHost
	port := config.ServerConfig.MailPort
	email := config.ServerConfig.MailUser
	password := config.ServerConfig.MailPass
	emailFrom := config.ServerConfig.MailFrom

	headers := make(map[string]interface{})
	headers["From"] = emailFrom + "<" + email + ">"
	headers["to"] = toEmail
	headers["Subject"] = title
	headers["content-type"] = "text/html; charset=UTF-8"

	message := ""

	for key, value := range headers {
		message = fmt.Sprintf("%s: %s\n", key, value)
	}

	message += "\r\n" + content

	auth := smtp.PlainAuth("", email, password, host)

	err := SendEmailUsingTLS(fmt.Sprintf("%s:%d", host, port),
		auth,
		email,
		[]string{toEmail},
		message)

	return err
}

//参考net/smtp的func SendMail()
//使用net.Dial连接tls(ssl)端口时, smtp.NewClient()会卡住且不提示err
//len(to) > 1 时, to[1]开始提示是密送
func SendEmailUsingTLS(addr string, auth smtp.Auth, from string, tos []string, msg string) error {
	client,err := createSMTLClient(addr)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	defer client.Close()

	if auth != nil {
		if ok,_ := client.Extension("AUTH");ok {
			if err := client.Auth(auth); err != nil {
				fmt.Println(err.Error())
				return err
			}
		}
	}

	if err := client.Mail(from); err != nil {
		return err
	}

	for _,to := range tos {
		if err := client.Rcpt(to); err != nil {
			return err
		}
	}

	writeCloser,err := client.Data()
	if err != nil {
		return err
	}

	_,err = writeCloser.Write([]byte(msg))
	if err != nil {
		return err
	}

	err = writeCloser.Close()
	if err != nil {
		return err
	}

	return client.Quit()
}

func createSMTLClient(addr string) (*smtp.Client, error) {
	conn,err := tls.Dial("tcp",addr,nil)
	if err != nil {
		fmt.Println(err.Error())
		return nil,err
	}

	host,_,_ := net.SplitHostPort(addr)

	return smtp.NewClient(conn,host)


}
