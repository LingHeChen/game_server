package email_utils

import (
	"game_server/src/frame"
	"fmt"
	"log"
  "net/smtp"
)


var EmailSender *EmailClient


type EmailClient struct {
  SmtpServer string
  SmtpPort int
  SenderEmail string
  SenderAuth string
  auth smtp.Auth
}


func (client *EmailClient) Send(
  reciver string,
  subject string,
  msg string,
) error {
  	// 组装邮件内容
	message := []byte("To: " + reciver + "\r\n" +
		"Subject: " + subject + "\r\n" +
		"\r\n" +
		msg + "\r\n")

	// 发送邮件
	err := smtp.SendMail(
    fmt.Sprintf("%s:%d", client.SmtpServer, client.SmtpPort),
    client.auth,
    client.SenderEmail,
    []string{reciver},
    message,
  )
	if err != nil {
    return err
	} else {
    return nil
	}
}


func (client *EmailClient) SayHelloTo(reciver string)  {
  if err := client.Send(reciver, "Hello", "Hello World"); err != nil {
    log.Println("Error Sending Email: ", err.Error())
  }
}


func (client *EmailClient) InitAuth() {
  client.auth = smtp.PlainAuth(
    "",
    client.SenderEmail,
    client.SenderAuth,
    client.SmtpServer,
  )
}


func init()  {
  stmpServer := fmt.Sprintf("smtp.%s.com", frame.Config.Email.Type)
  var stmpPort int
  switch frame.Config.Email.Type {
  case "163":
    stmpPort = 25
  default:
    log.Fatal(fmt.Sprintf("Unknown Email Type: %s", frame.Config.Email.Type))
  }

  senderAddr := fmt.Sprintf(
    "%s@%s.com",
    frame.Config.Email.Username,
    frame.Config.Email.Type,
  )

  EmailSender = &EmailClient{
    SmtpServer: stmpServer,
    SmtpPort: stmpPort,
    SenderEmail: senderAddr,
    SenderAuth: frame.Config.Email.Auth,
  }

  EmailSender.InitAuth()
}
