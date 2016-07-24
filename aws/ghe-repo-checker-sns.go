package aws

import (
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/aws"
	"fmt"
	"github.com/aws/aws-sdk-go/aws/session"
)

func SendNotification(subject string, message string) {
	svc := sns.New(session.New(&aws.Config{Region: aws.String("eu-west-1")}))

	params := &sns.PublishInput{
		Message: aws.String(message), // Required
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"Key": { // Required
				DataType:    aws.String("String"), // Required
				StringValue: aws.String("String"),
			},
		},
		MessageStructure: aws.String("messageStructure"),
		Subject:          aws.String(subject),
		TopicArn:         aws.String("arn:aws:sns:eu-west-1:291656607665:ms-repo-notifier"),
	}
	resp, err := svc.Publish(params)

	if err != nil {
		panic(err)
	}

	fmt.Println(resp)
}
