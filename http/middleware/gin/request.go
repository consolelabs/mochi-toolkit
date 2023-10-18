package gin

import (
	"bytes"
	"io"
	"net/http"
	"time"

	"github.com/consolelabs/mochi-typeset/queue/audit-log/typeset"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type CaptureRequestOptions struct {
	// exclude paths from capture
	ExcludePaths []string
}

type MiddlewareOption struct {
	PrivateIP     bool
	WhitelistPath []string
	DebugMode     bool
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
	if c.Request.Method == "POST" || c.Request.Method == "PUT" || c.Request.Method == "DELETE" {
		// read request body and parse to string
		body, _ = io.ReadAll(c.Request.Body)
	}
	// Restore the io.ReadCloser to its original state
	c.Request.Body = io.NopCloser(bytes.NewBuffer(body))

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

func WithAuth(opts *MiddlewareOption) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := logrus.WithFields(logrus.Fields{
			"component": "gin",
		})

		if opts == nil {
			c.Next()
			return
		}

		// check Debug mode
		if opts.DebugMode {
			logger.WithFields(logrus.Fields{"Client Ip": c.ClientIP(), "Request": c.Request}).Info("process request")
			c.Next()
			return
		}

		// check whitelist path
		for _, pApi := range opts.WhitelistPath {
			if validatePublicApi(c.Request.URL.RequestURI(), pApi) {
				c.Next()
				return
			}
		}

		// check Private IP
		if opts.PrivateIP && !isPrivateIP(c.ClientIP()) {
			c.AbortWithStatusJSON(http.StatusForbidden, gin.H{
				"status":  http.StatusForbidden,
				"message": "Permission denied",
			})
			return
		}

		c.Next()
	}
}
