package realip

import (
	"log"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func newRequest(remoteAddr, xRealIP string, xForwardedFor ...string) *http.Request {
	h := http.Header{}
	h.Set("X-Real-IP", xRealIP)
	h.Set("X-Forwarded-For", strings.Join(xForwardedFor, ", "))
	log.Println(h.Get("X-Forwarded-For"))
	return &http.Request{
		RemoteAddr: remoteAddr,
		Header:     h,
	}
}
func TestFromRequest(t *testing.T) {
	testData := []struct {
		Name     string
		Request  *http.Request
		Expected string
	}{
		{"noheader", newRequest("1.1.1.1:1122", ""), "1.1.1.1"},
		{"noheader", newRequest("1.1.1.1", ""), "1.1.1.1"},
		{"x-real-ip", newRequest("1.1.1.1", "2.2.2.2"), "2.2.2.2"},
		{"x-real-ip", newRequest("1.1.1.1", "a.b.c.d"), "1.1.1.1"},
		{"x-forwarded-for", newRequest("1.1.1.1", "", "1.2.3.4"), "1.2.3.4"},
		{"x-forwarded-for", newRequest("1.1.1.1", "", "1.2.3.4", "2.3.4.5"), "1.2.3.4"},
	}
	for _, v := range testData {
		assert.Equal(t, v.Expected, FromRequest(v.Request), v.Name)
	}
}
