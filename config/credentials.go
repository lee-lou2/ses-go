package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

func init() {
	privateKey := GetEnv("GOOGLE_SERVICE_ACCOUNT_PRIVATE_KEY")
	privateKey = strings.ReplaceAll(privateKey, `\n`, "\n")
	credentials := map[string]string{
		"type":                        "service_account",
		"project_id":                  GetEnv("GOOGLE_SERVICE_ACCOUNT_PROJECT_ID"),
		"private_key_id":              GetEnv("GOOGLE_SERVICE_ACCOUNT_PRIVATE_KEY_ID"),
		"private_key":                 privateKey,
		"client_email":                GetEnv("GOOGLE_SERVICE_ACCOUNT_CLIENT_EMAIL"),
		"client_id":                   GetEnv("GOOGLE_SERVICE_ACCOUNT_CLIENT_ID"),
		"auth_uri":                    "https://accounts.google.com/o/oauth2/auth",
		"token_uri":                   "https://oauth2.googleapis.com/token",
		"auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
		"client_x509_cert_url":        GetEnv("GOOGLE_SERVICE_ACCOUNT_CLIENT_X509_CERT_URL"),
		"universe_domain":             "googleapis.com",
	}
	jsonData, err := json.MarshalIndent(credentials, "", "  ")
	if err != nil {
		fmt.Printf("Error marshalling JSON: %v\n", err)
		return
	}
	filePath := "config/credentials/client_secret.json"
	if err := os.MkdirAll("config/credentials", os.ModePerm); err != nil {
		fmt.Printf("Error creating directory: %v\n", err)
		return
	}
	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		fmt.Printf("Error writing to file: %v\n", err)
		return
	}
}
