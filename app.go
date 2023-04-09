package main

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"os"
	"time"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/nleeper/goment"
	"strconv"
)

func createCurrDateString() string {
	d, err := goment.New()
	if err != nil {
		fmt.Println("error with goment")
	}

	return d.Format("YYYY-MM-DD")
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

func ComputeExpectedSHA256Hash(data []byte) string {
	secret := os.Getenv("GITHUB_WEBHOOK_SECRET")

	if secret == "" {
		fmt.Println("Error: no secret found")
		return ""
	}

	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(data)

	return "sha256=" + hex.EncodeToString(mac.Sum(nil))
}

func verifyOrigin(req *events.LambdaFunctionURLRequest) (isVerified bool) {
	isVerified = false

	var signature string

	signature, ok := req.Headers["x-hub-signature-256"]

	if !ok || len(signature) == 0 {
		fmt.Println("Error: no signature found in header")
		return
	}

	expectedHash := ComputeExpectedSHA256Hash([]byte(req.Body))

	if expectedHash == "" {
		fmt.Println("Error: no hash calculated")
		return
	}

	isVerified = hmac.Equal([]byte(signature), []byte(expectedHash))

	if !isVerified {
		fmt.Println("hashes are not equal")
	} else {
		fmt.Println("signature has been verified")
	}

	return
}

type DbEntry struct {
	Date string `json:"date"`
	Repo string `json:"repo"`
}

func WriteToDb(repoName string) {
	dynamoSession := createDynamoSession()

	dbEntry := DbEntry{
		Date: strconv.FormatInt(time.Now().UnixMilli(), 10),
		Repo: repoName,
	}

	bodyMap, err := dynamodbattribute.MarshalMap(dbEntry)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	input := &dynamodb.PutItemInput{
		Item: bodyMap,
		TableName: aws.String(os.Getenv("TABLE_NAME")),
	}

	_, err = dynamoSession.PutItem(input)

	if err != nil {
		fmt.Println(err.Error())
	}
}

func HandleWebhookRequest(ctx context.Context, req events.LambdaFunctionURLRequest) (string, error) {
	isVerified := verifyOrigin(&req)

	if isVerified == false {
		return "verification error", nil
	}

	var body RequestBody

	err := json.Unmarshal([]byte(req.Body), &body)

	if err != nil {
		fmt.Println(err.Error())
		return err.Error(), err
	}


	if body.Repository.Name == "" {
		fmt.Println("no repo name")
		return "no repo name", nil
	}

	WriteToDb(body.Repository.Name)

	return "done", nil
}

func GenerateHighEntropyString(length int) (string, error) {
	// Determine the required number of bytes
	byteLength := length * 3 / 4
	if length%4 != 0 {
		byteLength++
	}

	// Generate a random byte slice
	bytes := make([]byte, byteLength)
	_, err := rand.Read(bytes)
	if err != nil {
		return "", err
	}

	// Encode the byte slice as a base64 string
	encoded := base64.RawURLEncoding.EncodeToString(bytes)

	// Truncate the string to the desired length
	fmt.Println(encoded[:length])
	return encoded[:length], nil
}

//https://docs.github.com/en/webhooks-and-events/webhooks/creating-webhooks#setting-up-a-webhook
func main() {
	lambda.Start(HandleWebhookRequest)
}
