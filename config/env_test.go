package config

import (
	"os"
	"testing"
)

func TestGetEnv(t *testing.T) {
	tests := []struct {
		name         string
		key          string
		defaultValue string
		envValue     string
		want         string
	}{
		{
			name:         "환경변수가 설정된 경우",
			key:          "TEST_KEY",
			defaultValue: "default",
			envValue:     "test_value",
			want:         "test_value",
		},
		{
			name:         "환경변수가 없는 경우 기본값 반환",
			key:          "MISSING_KEY",
			defaultValue: "default",
			envValue:     "",
			want:         "default",
		},
		{
			name:         "기본값이 없고 환경변수도 없는 경우",
			key:          "MISSING_KEY",
			defaultValue: "",
			envValue:     "",
			want:         "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				os.Setenv(tt.key, tt.envValue)
				defer os.Unsetenv(tt.key)
			}

			got := GetEnv(tt.key, tt.defaultValue)
			if got != tt.want {
				t.Errorf("GetEnv() = %v, want %v", got, tt.want)
			}
		})
	}
}
