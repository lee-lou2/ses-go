package ses

import (
	"context"
	"ses-go/config"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

// SendEmail 함수는 이메일을 발송
func SendEmail(subject, body *string, receivers *[]string) (string, error) {
	// AWS Config 로드
	server := config.EmailServer
	AccessKeyId := config.GetEnv("AWS_ACCESS_KEY_ID")
	SecretAccessKey := config.GetEnv("AWS_SECRET_ACCESS_KEY")
	cfg, err := awsConfig.LoadDefaultConfig(
		context.TODO(),
		awsConfig.WithRegion("ap-northeast-2"),
		awsConfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(AccessKeyId, SecretAccessKey, ""),
		),
	)
	if err != nil {
		return "", err
	}
	client := sesv2.NewFromConfig(cfg)
	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(server.Sender),
		Destination: &types.Destination{
			ToAddresses: *receivers,
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{
					Data: aws.String(*subject),
				},
				Body: &types.Body{
					Html: &types.Content{
						Data: aws.String(*body),
					},
				},
			},
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, err := client.SendEmail(ctx, input)
	if err != nil {
		return "", err
	}
	return *result.MessageId, nil
}
