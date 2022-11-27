package cli

import (
	"flag"

	"github.com/song940/smtp/smtp"
)

func Run() {
	flag.Parse()
	client := smtp.NewClient()
	client.Hello()
	client.Auth()

	message := smtp.NewMessage()
	message.AddRecipient("hi@lsong.org")
	message.AddRecipient("test@lsong.org")
	message.AddRecipient("song940@163.com")
	client.Send(message)
	client.Quit()
}
