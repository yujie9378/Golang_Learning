package main

import (
	"errors"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"testing"
)

func TestWork_01(t *testing.T) {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})
	r.Run(":8080")
}
func TestWork_02(t *testing.T) {
	r := gin.New()
	r.Use(ValidHeader())
	r.GET("/ping", func(context *gin.Context) {

		context.JSON(http.StatusOK, gin.H{
			"mess": "success",
		})
	})
	r.Run(":8081")
}
func TestWork_03(t *testing.T) {
	r := gin.New()
	r.Use(Logger())
	r.GET("/logtest", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "success",
		})

	})
	r.Run(":8082")

}

func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Print("请求体", c.Request)

	}
}
func ValidHeader() gin.HandlerFunc {
	return func(c *gin.Context) {
		xtoken := c.Request.Header.Get("X-Tokenw")
		if !strings.HasPrefix(xtoken, "123") {
			//fmt.Println("响应", strconv.Itoa(c.Writer.Status()))
			c.AbortWithError(http.StatusBadRequest, errors.New("wrong header"))
			c.Writer.Write([]byte("header wrong"))

			return
		}
		c.JSON(http.StatusOK, gin.H{
			"mess": "hello",
		})
	}
}
