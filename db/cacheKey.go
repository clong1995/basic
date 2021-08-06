package db

import (
	"sort"
	"strings"
)

func keys(ks []string) string {
	sort.Strings(ks)
	return "|" + strings.Join(ks, "|") + "|"
}

func fields(fs []string) string {
	return strings.Join(fs, "")
}
