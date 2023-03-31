package main

import (
	"net/http"
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"io/ioutil"
	"regexp"
	"fmt"
	"github.com/nleeper/goment"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-lambda-go/lambdacontext"
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

func createRouter() *gin.Engine {
	r := gin.Default()

	r.POST("/commit", func(c *gin.Context) {

		// handle webhook

	})

	return r
}

type DbEntry struct {
	Date string
	Repo string
	NumCommits int16
}

type RequestBody struct {
	idk string
}

func HandleWebhookRequest(ctx context.Context, requestBody RequestBody) {
	lc, _:= lambdacontext.FromContext(ctx)
	body, err := json.MarshalIndent(requestBody)
}

// https://docs.github.com/en/webhooks-and-events/webhooks/creating-webhooks#setting-up-a-webhook
func main() {
	router := createRouter()
	router.Run(":8080")
	lambda.Start(HandleWebhookRequest)
}
