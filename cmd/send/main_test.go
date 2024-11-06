package send

import (
	"context"
	"testing"
	"time"
)

func TestAddMessage(t *testing.T) {
	tests := []struct {
		name      string
		messageId int
		to        string
		subject   string
		body      string
		params    map[string]string
		ctx       context.Context
		wantErr   bool
	}{
		{
			name:      "정상적인 메시지 추가",
			messageId: 1,
			to:        "test@example.com",
			subject:   "Test Subject",
			body:      "Test Body",
			params:    map[string]string{"name": "Test"},
			ctx:       context.Background(),
			wantErr:   false,
		},
		{
			name:      "취소된 컨텍스트로 메시지 추가",
			messageId: 2,
			to:        "test@example.com",
			subject:   "Test Subject",
			body:      "Test Body",
			params:    map[string]string{"name": "Test"},
			ctx: func() context.Context {
				ctx, cancel := context.WithCancel(context.Background())
				cancel()
				return ctx
			}(),
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			done := make(chan bool)
			go func() {
				AddMessage(tt.messageId, tt.to, tt.subject, tt.body, tt.params, tt.ctx)
				done <- true
			}()

			select {
			case <-done:
				// 메시지가 성공적으로 추가됨
			case <-time.After(time.Second):
				if !tt.wantErr {
					t.Error("AddMessage() timed out")
				}
			}
		})
	}
}
