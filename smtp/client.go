package smtp

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"sort"
	"strings"
	"time"
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

func resolveMX(hostname string) []string {
	records, err := net.LookupMX(hostname)
	checkError(err)
	sort.Slice(records, func(i, j int) bool {
		return records[i].Pref < records[j].Pref
	})
	hosts := []string{}
	for _, record := range records {
		hosts = append(hosts, record.Host)
	}
	return hosts
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type SMTPClient struct {
	Host    string
	Port    uint32
	Timeout time.Duration
	conn    net.Conn
}

func NewClient() *SMTPClient {
	client := &SMTPClient{Port: 25, Timeout: 30 * time.Second}
	return client
}

func (c *SMTPClient) TryConnection(hosts []string) (conn net.Conn, err error) {
	for _, host := range hosts {
		remote := fmt.Sprintf("%s:%d", host, c.Port)
		conn, err := net.DialTimeout("tcp", remote, c.Timeout)
		if err != nil {
			log.Println("try connection", host, c.Port, err)
			continue
		}
		log.Printf("connect %s success\n", remote)
		return conn, nil
	}
	return
}

func (c *SMTPClient) CreateConnection(hostname string) (conn net.Conn, err error) {
	var hosts []string
	if c.Host != "" {
		hosts = []string{c.Host}
	} else {
		hosts = resolveMX(hostname)
		hosts = append(hosts, hostname)
	}
	return c.TryConnection(hosts)
}

func (c *SMTPClient) SetConnection(conn net.Conn) {
	c.conn = conn
}

func (c *SMTPClient) ExecuteCommand(cmd string, args ...interface{}) error {
	var buf bytes.Buffer
	fmt.Fprintf(&buf, cmd+"\n", args...)
	_, err := c.conn.Write(buf.Bytes())
	return err
}

func (c *SMTPClient) Quit() {
	c.ExecuteCommand("QUIT\n")
}

func (c *SMTPClient) PostMessage(hostname string, from string, recipients []string, content string) {
	conn, err := c.CreateConnection(hostname)
	checkError(err)
	c.SetConnection(conn)
	// reply := make([]byte, 1024)
	// _, err = conn.Read(reply)
	// checkError(err)
	// log.Println(string(reply))
	c.ExecuteCommand("EHLO %s", "localhost")
	c.ExecuteCommand("MAIL FROM:<%s>", from)
	for _, rcpt := range recipients {
		c.ExecuteCommand("RCPT TO:<%s>", rcpt)
	}
	c.ExecuteCommand("DATA")
	c.ExecuteCommand(content + "")
	c.ExecuteCommand(".")
	c.ExecuteCommand("")
}

func (c *SMTPClient) Send(msg *Message) {
	hosts := groupByHost(msg.GetRecipients())
	for hostname, repts := range hosts {
		c.PostMessage(hostname, msg.From, repts, msg.ToMime())
	}
}

func (c *SMTPClient) SendMessage() {
	message := NewMessage()
	c.Send(message)
}

func (c *SMTPClient) Close() error {
	return c.conn.Close()
}
