package ip

import (
	"net/http"
	"strings"
)

func XRealIp(r *http.Request) (ip string) {
	//X-Forwarded-For和X-Real-Ip可能是伪造的
	ip = strings.TrimSpace(strings.Split(r.Header.Get("X-Forwarded-For"), ",")[0])
	if ip == "" {
		ip = strings.TrimSpace(r.Header.Get("X-Real-Ip"))
	}
	if ip == "" {
		ip = r.RemoteAddr
	}
	ip = strings.Split(ip, ":")[0]
	return
}
