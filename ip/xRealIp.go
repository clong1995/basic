package ip

import (
	"net/http"
	"strings"
)

func XRealIp(r *http.Request) (ip string) {
	ip = r.Header.Get("X-Real-Ip")
	if ip == "" {
		ip = r.Header.Get("X-Forwarded-For")
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	ip = strings.Split(ip, ":")[0]
	return
}
