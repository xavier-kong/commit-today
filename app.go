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

func verifySignature() {

}

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

func verifyOrigin(req *events.LambdaFunctionURLRequest) bool {
	isVerified := false

	var signature string

	signature, ok := req.Headers["x-hub-signature-256"]

	if !ok || len(signature) == 0 {
		fmt.Println("no signature found in header")
		return isVerified
	}

	secret := os.Getenv("TABLE_NAME")

	hash := hmac.New(sha256.New, []byte(secret))

	if _, err := hash.Write(b); err != nil {

	}


	return isVerified
}


func HandleWebhookRequest(ctx context.Context, req events.LambdaFunctionURLRequest) (string, error) {
	isVerified := verifyOrigin()


	if isVerified == false {
		return "verification error", nil
	}

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
		fmt.Println("no repo name")
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
