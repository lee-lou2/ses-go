package models

import (
	"regexp"
	"ses-go/config"

	"gorm.io/gorm"
)

type Template struct {
	gorm.Model
	Subject   string `json:"subject" gorm:"not null;type:varchar(255)"`
	Body      string `json:"body" gorm:"not null;type:text"`
	CreatorId uint   `json:"creator_id" gorm:"not null"`
	Creator   User   `json:"creator" gorm:"foreignKey:CreatorId"`
}

func (t *Template) TableName() string {
	return "templates"
}

// GetParams 템플릿의 body 에서 파라미터를 추출하여 반환
func (t *Template) GetParams() *[]string {
	var columns []string
	re := regexp.MustCompile(`{{\s*(\w+)\s*}}`)
	matches := re.FindAllStringSubmatch(t.Body, -1)
	for _, match := range matches {
		if len(match) > 1 {
			columns = append(columns, match[1])
		}
	}
	return &columns
}

func init() {
	db := config.GetDB()
	_ = db.AutoMigrate(&Template{})
}
