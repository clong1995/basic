package logistic

import (
	"testing"
)

func TestTraces(t *testing.T) {
	Server{
		EBusinessID: "1724473",
		ApiKey:      "9127d27f-5395-499f-ad1d-96227256212c",
	}.CreateClient()

	traces, err := Logistic.Traces("YTO", "YT9716397794059")
	if err != nil {
		t.Error(err)
		return
	}
	t.Log(traces)
}

func TestSubscribe(t *testing.T) {
	Server{
		EBusinessID: "",
		ApiKey:      "",
	}.CreateClient()

	err := Logistic.Subscribe("YTO", "YT9716562561519", "ABCAuV-UwRM")
	if err != nil {
		t.Error(err)
		return
	}
}
