package config

// EmailServer 서버 설정
var EmailServer = struct {
	Sender string
}{
	Sender: GetEnv("EMAIL_SENDER", "no-reply@example.com"),
}

// EmailServerLimit 서버 제한 설정
var EmailServerLimit = struct {
	DailyQuota string
	Rate       string
}{
	DailyQuota: GetEnv("EMAIL_DAILY_QUOTA", "10000"),
	Rate:       GetEnv("EMAIL_RATE", "24"),
}
