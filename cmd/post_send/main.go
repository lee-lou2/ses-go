package post_send

import (
	"log"
	"ses-go/config"
	"ses-go/models"
	"sync/atomic"
	"time"
)

// AddMessage 함수는 메시지를 Messages 채널에 추가
func AddMessage(messageId string, id int, status int, error string) {
	go func() {
		messages <- message{
			MessageId: messageId,
			Id:        id,
			Status:    status,
			Error:     error,
		}
	}()
}

// Run 메시지 후 처리
func Run() {
	var msgCnt int32 = 0

	go func() {
		for {
			time.Sleep(1 * time.Second)
			cnt := atomic.SwapInt32(&msgCnt, 0)
			if cnt != 0 {
				log.Printf("Request per second: %d", cnt)
			}
		}
	}()

	for {
		select {
		case msg := <-messages:
			atomic.AddInt32(&msgCnt, 1)

			// DB 업데이트 작업
			db := config.GetDB()
			var message models.Message
			db.Model(&message).Where("id = ?", msg.Id).Updates(models.Message{
				MessageId: msg.MessageId,
				Status:    msg.Status,
				Error:     msg.Error,
			})
		}
	}
}
