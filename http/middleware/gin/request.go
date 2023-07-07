package gin

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/consolelabs/mochi-typeset/queue/audit-log/typeset"
	"github.com/gin-gonic/gin"
)

func CaptureRequest(c *gin.Context) typeset.AuditLogApi {
	start := time.Now()
	var body []byte
	if c.Request.Method == "POST" || c.Request.Method == "PUT" {
		// read request body and parse to string
		body, _ = ioutil.ReadAll(c.Request.Body)
	}
	// Restore the io.ReadCloser to its original state
	c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(body))

	w := &ResponseBodyWriter{body: &bytes.Buffer{}, ResponseWriter: c.Writer}
	c.Writer = w
	c.Next()

	return typeset.AuditLogApi{
		Method:       c.Request.Method,
		Uri:          c.Request.URL.String(),
		RequestBody:  body,
		StatusCode:   c.Writer.Status(),
		Latency:      time.Since(start),
		RequestId:    c.Request.Header.Get("X-Request-Id"),
		ResponseBody: w.body.Bytes(),
	}
}
