package db

import "testing"

func Test_keys(t *testing.T) {
	name := keys([]string{"a", "b", "d", "c"})
	t.Log(name)
}
