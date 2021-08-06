package ip

import "testing"

func TestBoundInternetIP(t *testing.T) {
	ips, err := BoundInternetIP()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(ips)
}
