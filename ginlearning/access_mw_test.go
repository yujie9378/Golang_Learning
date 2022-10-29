package ginlearning

import (
	"io/ioutil"
	"testing"

	"github.com/gin-gonic/gin"
)

func world(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(400)
	}
	if string(body) == "hello" {
		c.String(200, "world")
	} else {
		c.String(400, "i want to hello")
	}

}
func TestAccessMW(t *testing.T) {
	r := gin.Default()
	r.POST("/hello", AccessLogMW, world)
	r.Run(":8080")
}
