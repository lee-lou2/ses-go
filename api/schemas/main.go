package schemas

// ReqCreatePlan 플랜 생성 요청
type ReqCreatePlan struct {
	Title       string `json:"title"`
	TemplateId  uint   `json:"template_id"`
	RecipientId uint   `json:"recipient_id"`
	ScheduledAt string `json:"scheduled_at"`
}

// RespCreatePlan 플랜 생성 응답
type RespCreatePlan struct {
	Id uint `json:"id"`
}

// ReqCreateTemplate 템플릿 생성 요청
type ReqCreateTemplate struct {
	Subject string `json:"subject"`
}

// RespCreateTemplate 템플릿 생성 응답
type RespCreateTemplate struct {
	Id uint `json:"id"`
}

// ReqUpdateTemplate 템플릿 업데이트 요청
type ReqUpdateTemplate struct {
	Body string `json:"body"`
}

// RespUpdateTemplate 템플릿 업데이트 응답
type RespUpdateTemplate struct {
	Id uint `json:"id"`
}

// ReqAddSendEvent 발송 이벤트 추가 요청
type ReqAddSendEvent struct {
	Type         string `json:"Type"`
	Message      string `json:"Message"`
	SubscribeURL string `json:"SubscribeURL"`
}

// ReqCreateRecipients 수신자 생성 요청
type ReqCreateRecipients struct {
	Data [][]string `json:"data"`
}

// RespGetTemplateFields 템플릿 필드 조회 응답
type RespGetTemplateFields struct {
	Fields []string `json:"fields"`
}
