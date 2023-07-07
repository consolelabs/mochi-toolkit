package gin

import (
	"bytes"

	"github.com/gin-gonic/gin"
)

type ResponseBodyWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (r ResponseBodyWriter) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}
