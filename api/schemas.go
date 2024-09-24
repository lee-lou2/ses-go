package api

// reqCreatePlan 플랜 생성 요청
type reqCreatePlan struct {
	Title       string `json:"title"`
	TemplateId  uint   `json:"template_id"`
	SheetId     string `json:"sheet_id"`
	ScheduledAt string `json:"scheduled_at"`
}

// respCreatePlan 플랜 생성 응답
type respCreatePlan struct {
	Id uint `json:"id"`
}

// reqCreateSheetAndShare 시트 생성 및 공유 요청
type reqCreateSheetAndShare struct {
	TemplateId uint `json:"template_id"`
}

// respCreateSheetAndShare 시트 생성 및 공유 응답
type respCreateSheetAndShare struct {
	SheetId string `json:"sheet_id"`
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
	Type      string `json:"Type"`
	MessageId string `json:"MessageId"`
	Message   string `json:"Message"`
}
