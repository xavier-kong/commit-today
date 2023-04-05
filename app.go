package main

import (
	"context"
	"os"
	"fmt"
	"github.com/nleeper/goment"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"encoding/json"
)

func createCurrDateString() string {
	d, err := goment.New()
	if err != nil {
		fmt.Println("error with goment")
	}

	return d.Format("YYYY-MM-DD")
}

type DbEntry struct {
	Date string `json:"date"`
	Repo string
}

type Repository struct {
	Name string
}

type RequestBody struct {
	Repository Repository `json:"repository"`
}

func createDynamoSession() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return dynamodb.New(sess)
}

func HandleWebhookRequest(ctx context.Context, req events.LambdaFunctionURLRequest) (string, error) {
	var body RequestBody

	err := json.Unmarshal([]byte(req.Body), &body)

	if err != nil {
		fmt.Println(err.Error())
		return err.Error(), err
	}

	dynamoSession := createDynamoSession()

	if err != nil {
		fmt.Println(err.Error())
		return err.Error(), err
	}

	if body.Repository.Name == "" {
		return "no repo name", nil
	}

	dbEntry := DbEntry{
		Date: createCurrDateString(),
		Repo: body.Repository.Name,
	}

	bodyMap, err := dynamodbattribute.MarshalMap(dbEntry)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(bodyMap)

	input := &dynamodb.PutItemInput{
		Item: bodyMap,
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	_, err = dynamoSession.PutItem(input)

	if err != nil {
		fmt.Println(err.Error())
	}

	return "success", nil
}

//https://docs.github.com/en/webhooks-and-events/webhooks/creating-webhooks#setting-up-a-webhook
func main() {
	lambda.Start(HandleWebhookRequest)
}
