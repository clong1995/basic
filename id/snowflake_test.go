package id

import (
	"testing"
)

func TestInfo(t *testing.T) {
	Server{
		Node: 1,
	}.CreateNode()

	t.Log(SId.String())
}
