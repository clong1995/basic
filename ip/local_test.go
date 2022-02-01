package ip

import "testing"

func TestBoundLocalIP(t *testing.T) {
	ips, err := BoundLocalIP()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ips)
}
