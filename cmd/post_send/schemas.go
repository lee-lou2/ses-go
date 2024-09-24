package post_send

// message 메세지
type message struct {
	Id        int
	MessageId string
	Status    int
	Error     string
}

// messages 메세지 채널
var messages = make(chan message)
