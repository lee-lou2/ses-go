package models

import (
	"reflect"
	"testing"
)

func TestTemplate_GetParams(t *testing.T) {
	tests := []struct {
		name     string
		template Template
		want     []string
	}{
		{
			name: "기본 변수 추출",
			template: Template{
				Body: "Hello {{name}}, your email is {{email}}",
			},
			want: []string{"name", "email"},
		},
		{
			name: "공백이 있는 변수 추출",
			template: Template{
				Body: "Hello {{ name }}, your email is {{ email }}",
			},
			want: []string{"name", "email"},
		},
		{
			name: "중복 변수 처리",
			template: Template{
				Body: "Hello {{name}}, bye {{name}}",
			},
			want: []string{"name", "name"},
		},
		{
			name: "변수가 없는 경우",
			template: Template{
				Body: "Hello World",
			},
			want: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.template.GetParams()
			if !reflect.DeepEqual(*got, tt.want) {
				t.Errorf("Template.GetParams() = %v, want %v", *got, tt.want)
			}
		})
	}
}
