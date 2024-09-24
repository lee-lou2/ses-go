package main

import (
	"ses-go/api"
	"ses-go/cmd/post_send"
	"ses-go/cmd/pre_send"
	"ses-go/cmd/send"
)

func main() {
	// 발송 전 처리
	go func() {
		pre_send.Run()
	}()

	// 발송
	go func() {
		send.Run()
	}()

	// 발송 후 처리
	go func() {
		post_send.Run()
	}()

	// API 서버 실행
	api.Run()
}
