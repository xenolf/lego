package sender

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDo_UserAgentOnAllHTTPMethod(t *testing.T) {
	var ua, method string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ua = r.Header.Get("User-Agent")
		method = r.Method
	}))
	defer ts.Close()

	do := NewDo(http.DefaultClient, "")

	testCases := []struct {
		method string
		call   func(u string) (*http.Response, error)
	}{
		{
			method: http.MethodGet,
			call: func(u string) (*http.Response, error) {
				return do.Get(u, nil)
			},
		},
		{
			method: http.MethodHead,
			call:   do.Head,
		},
		{
			method: http.MethodPost,
			call: func(u string) (*http.Response, error) {
				return do.Post(u, strings.NewReader("falalalala"), "text/plain", nil)
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.method, func(t *testing.T) {

			_, err := test.call(ts.URL)
			require.NoError(t, err)

			assert.Equal(t, test.method, method)
			assert.Contains(t, ua, ourUserAgent, "User-Agent")
		})
	}
}

func TestDo_CustomUserAgent(t *testing.T) {
	customUA := "MyApp/1.2.3"
	do := NewDo(http.DefaultClient, customUA)

	ua := do.formatUserAgent()
	assert.Contains(t, ua, ourUserAgent)
	assert.Contains(t, ua, customUA)
	if strings.HasSuffix(ua, " ") {
		t.Errorf("UA should not have trailing spaces; got '%s'", ua)
	}
	assert.Len(t, strings.Split(ua, " "), 5)
}
