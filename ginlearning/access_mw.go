package ginlearning

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
)

func AccessLogMW(c *gin.Context) {
	reqbody, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(400)
	}
	uri := c.Request.URL
	//
	startTime := time.Now()

	myWrite, ok := c.Writer.(*MyWriter)
	if !ok {
		myWrite = &MyWriter{
			ResponseWriter: c.Writer,
			body:           bytes.NewBuffer(nil),
		}
		c.Writer = myWrite

	}

	//上面读了以后还需要重新设置一下body，不然后面业务代码读不到了
	c.Request.Body = io.NopCloser(bytes.NewReader(reqbody))
	c.Next()

	endTime := time.Now()

	fmt.Printf("%s url:%s, req:%s, ret_status_code:%d, req_body:%s, delay:%dms", startTime, uri, string(reqbody), c.Writer.Status(), myWrite.GetBody(), endTime.Sub(startTime).Milliseconds())

}

type MyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (m MyWriter) Write(data []byte) (int, error) {
	m.body.Write(data)
	return m.ResponseWriter.Write(data)
}

func (m MyWriter) GetBody() string {
	return m.body.String()
}
