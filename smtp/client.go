package smtp

import (
	"log"
	"net"
	"strings"
)

func parseAddress(address string) (string, string) {
	s := strings.Split(address, "@")
	return s[0], s[1]
}

func groupByHost(recipients []string) map[string][]string {
	output := make(map[string][]string)
	for _, recipient := range recipients {
		_, hostname := parseAddress(recipient)
		output[hostname] = append(output[hostname], recipient)
	}
	return output
}

func ResolveMX(hostname string) []string {
	return []string{}
}

type SMTPClient struct {
}

func NewClient() *SMTPClient {
	client := &SMTPClient{}
	return client
}

func (c *SMTPClient) Auth() {

}

func (c *SMTPClient) Hello() {

}

func (c *SMTPClient) PostMessage(hostname string, from string, recipients []string, content string) {
	log.Println(hostname, recipients)
	hosts := ResolveMX(hostname)
	for _, domain := range hosts {
		conn, err := net.Dial("tcp", domain)
		log.Println(conn, err)
	}
}

func (c *SMTPClient) Send(msg *Message) {
	hosts := groupByHost(msg.recipients)
	for hostname, repts := range hosts {
		c.PostMessage(hostname, msg.from, repts, msg.content)
	}
}

func (c *SMTPClient) SendMessage() {
	message := NewMessage()
	c.Send(message)
}

func (c *SMTPClient) Quit() {

}
