package google

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"google.golang.org/api/drive/v3"
	"google.golang.org/api/option"
	"google.golang.org/api/sheets/v4"
	"ses-go/models"
	"time"
)

// GetEmailsFromSheet 스프레드시트에서 이메일 조회
func GetEmailsFromSheet(sheetId string, messages *[]models.Message) error {
	// 서비스 생성
	ctx := context.Background()
	srv, err := sheets.NewService(ctx, option.WithCredentialsFile("config/credentials/client_secret.json"))
	if err != nil {
		return err
	}
	// 스프레드시트 데이터 조회
	resp, err := srv.Spreadsheets.Values.Get(sheetId, "Sheet1").Do()
	if err != nil {
		return err
	}
	// 데이터 조회
	if len(resp.Values) != 0 {
		columns := resp.Values[0]
		if columns[0].(string) != "email" {
			return errors.New("invalid sheet format")
		}
		columns = columns[1:]
		emailMap := make(map[string]bool)
		for i, row := range resp.Values {
			if i == 0 {
				continue
			}
			email := row[0].(string)
			// 이메일 중복 확인
			if _, ok := emailMap[email]; ok {
				continue
			}
			params := make(map[string]string)
			params["email"] = email
			for i := 1; i < len(row); i++ {
				params[columns[i-1].(string)] = row[i].(string)
			}
			paramsStr, _ := json.Marshal(params)
			*messages = append(*messages, models.Message{
				To:     email,
				Params: string(paramsStr),
			})
			emailMap[email] = true
		}
		// 메모리 비우기
		emailMap = nil
	}
	return nil
}

// CreateSheetAndShare 새로운 구글 스프레드시트를 생성하고 특정 사용자에게 쓰기 권한 부여
func CreateSheetAndShare(userEmail string, columns *[]string) (string, error) {
	// 서비스 생성
	ctx := context.Background()
	sheetService, err := sheets.NewService(ctx, option.WithCredentialsFile("config/credentials/client_secret.json"))
	if err != nil {
		return "", err
	}

	// 새 스프레드시트 생성
	now := time.Now()
	today := now.In(time.FixedZone("Asia/Seoul", 9*60*60))
	todayStr := today.Format("2006-01-02")
	spreadsheet := &sheets.Spreadsheet{
		Properties: &sheets.SpreadsheetProperties{
			Title: "Request Send Email " + todayStr,
		},
	}
	resp, err := sheetService.Spreadsheets.Create(spreadsheet).Do()
	if err != nil {
		return "", err
	}

	// 생성된 스프레드시트 ID 반환
	sheetId := resp.SpreadsheetId

	// 첫 행에 columns 값 추가
	valueRange := &sheets.ValueRange{
		Range:  "sheet1!A1", // 첫 번째 시트의 첫 행에 추가
		Values: [][]interface{}{make([]interface{}, len(*columns))},
	}
	for i, column := range *columns {
		valueRange.Values[0][i] = column
	}

	_, err = sheetService.Spreadsheets.Values.Update(sheetId, valueRange.Range, valueRange).ValueInputOption("RAW").Do()
	if err != nil {
		return "", err
	}

	// 드라이브 서비스 생성 (구글 스프레드시트와 연동하기 위해 필요)
	driveService, err := drive.NewService(ctx, option.WithCredentialsFile("config/credentials/client_secret.json"))
	if err != nil {
		return "", err
	}

	// 쓰기 권한 부여
	permission := &drive.Permission{
		Type:         "user",
		Role:         "writer",
		EmailAddress: userEmail,
	}
	_, err = driveService.Permissions.Create(sheetId, permission).Do()
	if err != nil {
		return "", err
	}
	return sheetId, nil
}

// ChangePermissionsToReader 소유주를 제외한 모든 사용자의 권한을 읽기 권한으로 변경
func ChangePermissionsToReader(sheetId string) error {
	// 드라이브 서비스 생성
	ctx := context.Background()
	driveService, err := drive.NewService(ctx, option.WithCredentialsFile("config/credentials/client_secret.json"))
	if err != nil {
		return fmt.Errorf("failed to create drive service: %v", err)
	}

	// 스프레드시트의 현재 권한 목록 조회
	permissionsList, err := driveService.Permissions.List(sheetId).Fields("permissions(id, emailAddress, role)").Do()
	if err != nil {
		return fmt.Errorf("failed to retrieve permissions: %v", err)
	}

	for _, permission := range permissionsList.Permissions {
		// 소유주(Owner)는 권한 변경 제외
		if permission.Role == "owner" {
			continue
		}

		permissionUpdate := &drive.Permission{
			Role: "reader",
		}
		_, err := driveService.Permissions.Update(sheetId, permission.Id, permissionUpdate).Do()
		if err != nil {
			return fmt.Errorf("failed to update permission for %s: %v", permission.EmailAddress, err)
		}
	}
	return nil
}
