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
	"github.com/joho/godotenv"
)

func createCurrDateString() {
	d, err := goment.New()
	if err != nil {
		fmt.Println("error with goment")
	}

	fmt.Println(d.Format("YYYY-MM-DD"))
}

type DbEntry struct {
	Date string
	Repo string
	NumCommits int16
}

type RequestBody struct {
	Test string `json:"test"`
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
	}

	dynamoSession := createDynamoSession()

	err = godotenv.Load(".env")

	if err != nil {
		fmt.Println(err.Error())
	}

	bodyMap, err := dynamodbattribute.MarshalMap(body)

	input := &dynamodb.PutItemInput{
		Item: bodyMap,
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	_, err = dynamoSession.PutItem(input)


	return "success", nil
}

// https://docs.github.com/en/webhooks-and-events/webhooks/creating-webhooks#setting-up-a-webhook
func main() {
	lambda.Start(HandleWebhookRequest)
}
