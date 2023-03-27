package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"log"
	"io/ioutil"
	"regexp"
	"fmt"
)

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

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.POST("/commit", func(c *gin.Context) {
		// handle webhook
		// https://docs.github.com/en/webhooks-and-events/webhooks/creating-webhooks#setting-up-a-webhook

	})

	return r
}

func main() {
	router := createRouter()
	router.Run(":8080")
}
