package google

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"ses-go/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Response 결과
type Response struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// ConfigGoogle 구글 설정
func ConfigGoogle() *oauth2.Config {
	conf := &oauth2.Config{
		ClientID:     config.GetEnv("GOOGLE_CLIENT_ID"),
		ClientSecret: config.GetEnv("GOOGLE_CLIENT_SECRET"),
		RedirectURL:  config.GetEnv("GOOGLE_REDIRECT_URL"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email"},
		Endpoint: google.Endpoint,
	}
	return conf
}

// GetUserInfo 이메일 조회
func GetUserInfo(token string) (*Response, error) {
	reqURL, err := url.Parse("https://www.googleapis.com/oauth2/v1/userinfo")
	if err != nil {
		return nil, err
	}
	res := &http.Request{
		Method: "GET",
		URL:    reqURL,
		Header: map[string][]string{
			"Authorization": {fmt.Sprintf("Bearer %s", token)}},
	}
	req, err := http.DefaultClient.Do(res)
	if err != nil {
		return nil, err
	}
	defer func() { _ = req.Body.Close() }()
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	var data Response
	err = json.Unmarshal(body, &data)
	if err != nil {
		return nil, err
	}
	return &data, nil
}
