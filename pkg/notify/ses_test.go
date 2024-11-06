package notify

import (
	"os"
	"testing"
)

func TestSendSESEmail(t *testing.T) {
	tests := []struct {
		name      string
		subject   string
		body      string
		receivers []string
		setupEnv  func()
		wantErr   bool
	}{
		{
			name:      "AWS 자격증명이 없을 때 SMTP로 폴백",
			subject:   "테스트 제목",
			body:      "<p>테스트 내용</p>",
			receivers: []string{"test@example.com"},
			setupEnv: func() {
				os.Setenv("AWS_ACCESS_KEY_ID", "")
				os.Setenv("AWS_SECRET_ACCESS_KEY", "")
			},
			wantErr: true,
		},
		{
			name:      "잘못된 수신자 이메일",
			subject:   "테스트 제목",
			body:      "<p>테스트 내용</p>",
			receivers: []string{"invalid-email"},
			setupEnv: func() {
				os.Setenv("AWS_ACCESS_KEY_ID", "test")
				os.Setenv("AWS_SECRET_ACCESS_KEY", "test")
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupEnv()
			_, err := SendSESEmail(&tt.subject, &tt.body, &tt.receivers)
			if (err != nil) != tt.wantErr {
				t.Errorf("SendSESEmail() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
