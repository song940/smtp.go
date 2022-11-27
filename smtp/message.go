package smtp

type Message struct {
	from       string
	recipients []string
	subject    string
	content    string
}

func NewMessage() *Message {
	message := &Message{}
	return message
}

func (msg *Message) AddRecipient(rcpt string) *Message {
	msg.recipients = append(msg.recipients, rcpt)
	return msg
}
