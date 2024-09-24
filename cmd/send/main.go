package send

import (
	"context"
	"fmt"
	"ses-go/cmd/post_send"
	"ses-go/config"
	"ses-go/pkg/ses"
	"strconv"
	"strings"
	"time"
)

// AddMessage 함수는 메시지를 Messages 채널에 추가
func AddMessage(messageId int, to, subject, body string, params map[string]string, ctx context.Context) {
	messages <- message{
		MessageId: messageId,
		To:        to,
		Params:    params,
		Ctx:       ctx,
		Template: &template{
			Subject: subject,
			Body:    body,
		},
	}
}

// Run 메세지 전송
func Run() {
	limit := config.EmailServerLimit
	rateStr := limit.Rate
	rate, _ := strconv.Atoi(rateStr)

	// 메시지 큐
	for {
		select {
		case msg := <-messages:
			<-time.After(time.Second / time.Duration(rate))
			if msg.Ctx.Err() != nil {
				post_send.AddMessage("", msg.MessageId, 3, msg.Ctx.Err().Error())
				continue
			}
			go func(m *message) {
				// params 적용
				body := m.Template.Body
				for k, v := range m.Params {
					body = strings.ReplaceAll(body, fmt.Sprintf("{{%s}}", k), v)
					body = strings.ReplaceAll(body, fmt.Sprintf("{{ %s }}", k), v)
				}
				// body 마지막에 열림 이벤트를 위한 코드 추가
				serverHost := config.GetEnv("SERVER_HOST", "http://localhost:3000")
				body += `<img src="` + serverHost + `/v1/events/open/?message_id=` + strconv.Itoa(m.MessageId) + `" alt="open-event">`
				msgId, err := ses.SendEmail(
					&m.Template.Subject,
					&body,
					&[]string{m.To},
				)
				if err != nil {
					post_send.AddMessage(msgId, m.MessageId, 2, err.Error())
					return
				}
				post_send.AddMessage(msgId, m.MessageId, 1, "")
			}(&msg)
		}
	}
}
