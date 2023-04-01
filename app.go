package main

import (
	"net/http"
	"context"
	"log"
	"io/ioutil"
	"regexp"
	"fmt"
	"github.com/nleeper/goment"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/events"
	"encoding/json"
)

func createCurrDateString() {
	d, err := goment.New()
	if err != nil {
		fmt.Println("error with goment")
	}

	fmt.Println(d.Format("YYYY-MM-DD"))
}

func checkContributionFromHtml() {
	resp, err := http.Get("https://github.com/xavier-kong/")
	if err != nil {
		log.Fatalln(err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}
	sb := string(body)

	re, err := regexp.Compile(`\d+ contributions? on Monday, March 27, 2023`)

	fmt.Println(re.FindStringSubmatch(sb))

}

type DbEntry struct {
	Date string
	Repo string
	NumCommits int16
}

type RequestBody struct {
	Test string `json:"test"`
}

func HandleWebhookRequest(ctx context.Context, req events.LambdaFunctionURLRequest) (string, error) {
	var body RequestBody

	err := json.Unmarshal([]byte(req.Body), &body)

	if err != nil {
		fmt.Println(err.Error())
	}

	fmt.Println(body.Test)

	return "success", nil
}

// https://docs.github.com/en/webhooks-and-events/webhooks/creating-webhooks#setting-up-a-webhook
func main() {
	lambda.Start(HandleWebhookRequest)
}
