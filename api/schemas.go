package api

// reqCreatePlan 플랜 생성 요청
type reqCreatePlan struct {
	Title       string `json:"title"`
	TemplateId  uint   `json:"template_id"`
	RecipientId uint   `json:"recipient_id"`
	ScheduledAt string `json:"scheduled_at"`
}

// respCreatePlan 플랜 생성 응답
type respCreatePlan struct {
	Id uint `json:"id"`
}

// reqCreateTemplate 템플릿 생성 요청
type reqCreateTemplate struct {
	Subject string `json:"subject"`
}

// respCreateTemplate 템플릿 생성 응답
type respCreateTemplate struct {
	Id uint `json:"id"`
}

// reqUpdateTemplate 템플릿 업데이트 요청
type reqUpdateTemplate struct {
	Body string `json:"body"`
}

// respUpdateTemplate 템플릿 업데이트 응답
type respUpdateTemplate struct {
	Id uint `json:"id"`
}

// reqAddSendEvent 발송 이벤트 추가 요청
type reqAddSendEvent struct {
	Type         string `json:"Type"`
	Message      string `json:"Message"`
	SubscribeURL string `json:"SubscribeURL"`
}

// reqCreateRecipients 수신자 생성 요청
type reqCreateRecipients struct {
	Data string `json:"data"`
}

// respInitRecipientsData 수신자 초기화 응답
type respInitRecipientsData struct {
	Data [][]string `json:"data"`
}
