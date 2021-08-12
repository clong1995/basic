package logistic

import (
	"testing"
)

func TestTraces(t *testing.T) {
	Server{
		EBusinessID: "",
		ApiKey:      "",
	}.CreateClient()

	traces, err := Logistic.Traces("YTO", "YT9657662191918")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(traces)
}
