package gin

import (
	"bytes"
	"io/ioutil"
	"time"

	"github.com/consolelabs/mochi-typeset/queue/audit-log/typeset"
	"github.com/gin-gonic/gin"
)

type CaptureRequestOptions struct {
	// exclude paths from capture
	ExcludePaths []string
}

func CaptureRequest(c *gin.Context, opts *CaptureRequestOptions) *typeset.AuditLogMessage {
	// check if path is excluded
	for _, path := range opts.ExcludePaths {
		if c.Request.URL.Path == path {
			c.Next()
			return nil
		}
	}

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

	return &typeset.AuditLogMessage{
		Type: typeset.AUDIT_LOG_MESSAGE_TYPE_API,
		ApiLog: &typeset.AuditLogApi{
			Method:       c.Request.Method,
			Uri:          c.Request.URL.String(),
			RequestBody:  body,
			StatusCode:   c.Writer.Status(),
			Latency:      time.Since(start),
			RequestId:    c.Request.Header.Get("X-Request-Id"),
			ResponseBody: w.body.Bytes(),
		},
	}
}
