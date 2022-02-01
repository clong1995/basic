package random

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

func Number(width int) string {
	numeric := [10]byte{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
	r := len(numeric)
	rand.Seed(time.Now().UnixNano())
	var sb strings.Builder
	for i := 0; i < width; i++ {
		_, _ = fmt.Fprintf(&sb, "%d", numeric[rand.Intn(r)])
	}
	return sb.String()
}

func NumberNotZeroStart(width int) string {
	codeStr := Number(width)
	if strings.HasPrefix(codeStr, "0") {
		r := rand.New(rand.NewSource(time.Now().UnixNano()))
		n := r.Intn(9) + 1
		codeStr = fmt.Sprintf("%d%s", n, codeStr[1:])
	}
	return codeStr
}

func String(n int) string {
	letterRunes := []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
