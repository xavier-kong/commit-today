package main

import (
	"context"
	"os"
	"fmt"
	"github.com/nleeper/goment"
	//"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"encoding/json"
	"log"
	"net/http"
	"io/ioutil"
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

	if err != nil {
		fmt.Println(err.Error())
	}

	dbEntry := DbEntry{
		Date: createCurrDateString(),
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

// https://docs.github.com/en/webhooks-and-events/webhooks/creating-webhooks#setting-up-a-webhook
//func main() {
//lambda.Start(HandleWebhookRequest)
//}


func handler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	responseBody := body
	var data map[string]interface{}
	err = json.Unmarshal([]byte(responseBody), &data)
	if err != nil {
		panic(err)
	}
	fmt.Println(data)
}

func main() {
	http.HandleFunc("/", handler)
	fmt.Println("started")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
