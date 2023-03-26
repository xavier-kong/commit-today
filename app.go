package main

import (
	"net/http"
	"github.com/gin-gonic/gin"
	"log"
	"io/ioutil"
	"regexp"
	"fmt"
)

func createRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	r.GET("/data", func(c *gin.Context) {
		resp, err := http.Get("https://github.com/users/xavier-kong/contributions")
		if err != nil {
			log.Fatalln(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.Fatalln(err)
		}
		sb := string(body)

		fmt.Println(sb)

		re, err := regexp.Compile(`\d+ contributions? on Monday, March 27, 2023`)

		fmt.Println(re.FindStringSubmatch(sb))
	})

	return r
}

func main() {
	router := createRouter()
	router.Run(":8080")
}
