package aws

import (
	"log"
	"fmt"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/google/go-github/github"
	"strings"
	//"reflect"
)

type Item struct {
	Name       string `json:"Name"`
	UpdatedAt   string `json:"UpdatedAt",omitempty`
}

func CreateTable() {
	svc := dynamodb.New(session.New(&aws.Config{Region: aws.String("eu-west-1")}))
	result, err := svc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Tables:")
	for _, table := range result.TableNames {
		log.Println(*table)
	}
}

func ListAllTables () {
	svc := dynamodb.New(session.New(&aws.Config{Region: aws.String("eu-west-1")}))
	result, err := svc.ListTables(&dynamodb.ListTablesInput{})
	if err != nil {
		log.Println(err)
		return
	}

	log.Println("Tables:")
	for _, table := range result.TableNames {
		log.Println(*table)
	}
}

func ScanTable (dynamoDbTable string) map[string]map[string]*dynamodb.AttributeValue {
	svc := dynamodb.New(session.New(&aws.Config{Region: aws.String("eu-west-1")}))

	params := &dynamodb.ScanInput{
		TableName: aws.String(dynamoDbTable), // Required
	}
	resp, err := svc.Scan(params)

	if err != nil {
		panic(err)
	}

	m := make(map[string]map[string]*dynamodb.AttributeValue)

	for _, value := range resp.Items {
		m[*value["Name"].S] = value
	}

	return m
}

func GetAccessTokenFromDynamo() string {
	svc := dynamodb.New(session.New(&aws.Config{Region: aws.String("eu-west-1")}))
	 getItemInput := &dynamodb.GetItemInput{
		Key: map[string]*dynamodb.AttributeValue{
			"key" : {
				S: aws.String("AccessToken"),
			},
		},
		TableName: aws.String("credentials"),
	}

	item, err := svc.GetItem(getItemInput)

	if err != nil {
		panic(err)
	}

	return *item.Item["value"].S
}

func CreateItems(items map[string]github.Repository) {

	svc := dynamodb.New(session.New(&aws.Config{Region: aws.String("eu-west-1")}))

	var putreq []*dynamodb.WriteRequest

	i:=0
	for _, repo := range items {

		if !strings.HasPrefix(*repo.Name, "ms-"){
			continue
		}
		it, err := dynamodbattribute.ConvertToMap(Item{
			Name:      string(*repo.Name),
			UpdatedAt:    repo.PushedAt.String(),
		})

		if *repo.Name == "ms-test--test" {
			fmt.Println(*repo.PushedAt)
		}

		if err != nil {
			panic(err)
		}



		putreq = append(putreq, &dynamodb.WriteRequest{
			PutRequest: &dynamodb.PutRequest{
				Item: it,
			},
		})
		i++
		if i==20 {
			params := &dynamodb.BatchWriteItemInput{
				RequestItems: map[string][]*dynamodb.WriteRequest{
					"ghe-repositories": putreq,
				},
			}

			_, err :=  svc.BatchWriteItem(params)

			fmt.Println(err);

			i=0
			putreq = nil
		}
	}

	if putreq != nil {
		params := &dynamodb.BatchWriteItemInput{
			RequestItems: map[string][]*dynamodb.WriteRequest{
				"ghe-repositories": putreq,
			},
		}

		_, err :=  svc.BatchWriteItem(params)

		fmt.Println(err);

		putreq = nil
	}
}
