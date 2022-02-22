package mock

import (
	"testing"
)

func TestUserInfo(t *testing.T) {
	actual := UserInfo()
	t.Log(actual)
}
