package send

import (
	"context"
)

// template 메시지 템플릿
type template struct {
	Subject string
	Body    string
}

// message 메시지
type message struct {
	MessageId int
	To        string
	Template  *template
	Params    map[string]string
	Ctx       context.Context
}

const messageBufferSize = 1000

var messages = make(chan message, messageBufferSize)
