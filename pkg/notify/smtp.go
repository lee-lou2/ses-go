package notify

import (
	"errors"
	"fmt"
	"ses-go/config"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"gopkg.in/mail.v2"
)

// SendSMTPEmail 이메일 발송
func SendSMTPEmail(subject, body *string, receivers *[]string) (string, error) {
	// SMTP 서버 설정
	smtpHost := config.GetEnv("EMAIL_SMTP_HOST")
	smtpPortString := config.GetEnv("EMAIL_SMTP_PORT")
	username := config.GetEnv("EMAIL_USERNAME")
	password := config.GetEnv("EMAIL_PASSWORD")
	if smtpHost == "" || smtpPortString == "" || username == "" || password == "" {
		return "", errors.New("SMTP 서버 설정이 누락되었습니다")
	}
	smtpPort, err := strconv.Atoi(smtpPortString)
	if err != nil {
		return "", err
	}

	// 메일 내용 설정
	from := username

	// 메시지 생성
	to := strings.Join(*receivers, ",")
	m := mail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", *subject)
	m.SetBody("text/html", *body)

	// 다이얼러 설정
	d := mail.NewDialer(smtpHost, smtpPort, username, password)
	d.StartTLSPolicy = mail.MandatoryStartTLS

	// 이메일 발송 시도
	if err := d.DialAndSend(m); err != nil {
		// 발송 실패 시 에러 반환
		return "", fmt.Errorf("이메일 발송 실패: %v", err)
	}

	// 발송 성공 시 메시지 ID 반환
	messageId := uuid.Must(uuid.NewV7()).String()
	return messageId, nil
}
