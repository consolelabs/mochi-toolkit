package gin

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func TestWithAuth(t *testing.T) {
	testCases := []struct {
		ClientIP        string
		requestPath     string
		applicationPath string
		opts            *MiddlewareOption
		expected        interface{}
	}{
		{
			"51.0.0.8",
			"/api/v1/1111/transaction",
			"/api/v1/1111/transaction",
			&MiddlewareOption{
				PrivateIP:     false,
				WhitelistPath: []string{},
				DebugMode:     false,
			},
			http.StatusOK,
		},
		{
			"51.0.0.8",
			"/api/v1/1111/transaction?query=1",
			"/api/v1/1111/transaction",
			&MiddlewareOption{
				PrivateIP:     true,
				WhitelistPath: []string{"/api/v1/:id/transaction"},
				DebugMode:     false,
			},
			http.StatusOK,
		},
		{
			"51.0.0.8",
			"/api/v1/1111/transaction",
			"/api/v1/1111/transaction",
			&MiddlewareOption{
				PrivateIP:     true,
				WhitelistPath: []string{"/api/v1/:id/transaction"},
				DebugMode:     false,
			},
			http.StatusOK,
		},
		{
			"51.0.0.8",
			"/api/v1/1111/transaction",
			"/api/v1/1111/transaction",
			&MiddlewareOption{
				PrivateIP:     true,
				WhitelistPath: []string{},
				DebugMode:     false,
			},
			http.StatusForbidden,
		},
		{
			"10.0.0.8",
			"/api/v1/1111/transaction",
			"/api/v1/1111/transaction",
			&MiddlewareOption{
				PrivateIP:     true,
				WhitelistPath: []string{},
				DebugMode:     false,
			},
			http.StatusOK,
		},
		{
			"51.0.0.8",
			"/api/v1/1111/transaction",
			"/api/v1/1111/transaction",
			&MiddlewareOption{
				PrivateIP:     true,
				WhitelistPath: []string{},
				DebugMode:     true,
			},
			http.StatusOK,
		},
	}

	for _, tc := range testCases {
		tc := tc
		t.Run(fmt.Sprintf("WithAuth(%v) => %v ", tc.opts, tc.expected), func(t *testing.T) {
			gin.SetMode(gin.TestMode)

			w := httptest.NewRecorder()
			router := gin.Default()
			router.TrustedPlatform = "X-Appengine-Remote-Addr"
			router.Use(WithAuth(tc.opts))
			router.GET(tc.applicationPath, MockHandler)

			req, _ := http.NewRequest(http.MethodGet, tc.requestPath, nil)
			req.Header.Add("X-Appengine-Remote-Addr", tc.ClientIP)

			router.ServeHTTP(w, req)
			require := require.New(t)

			actual := w.Result().StatusCode
			require.Equal(tc.expected, actual)
		})
	}

}

func MockHandler(c *gin.Context) {
	c.JSONP(http.StatusOK, "abc")
}
