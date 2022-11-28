package cli

import (
	"flag"

	"github.com/song940/smtp/smtp"
)

func Run() {
	flag.Parse()
	client := smtp.NewClient()
	client.Host = "localhost"
	client.Port = 8989

	message := smtp.NewMessage()
	message.From = "hi@lsong.org"
	message.To = "song940@gmail.com"
	message.Subject = "Test Email"
	message.Content = "This is a test message"

	client.Send(message)
	client.Quit()
}
