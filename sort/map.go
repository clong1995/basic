package sort

import "sort"

type item struct {
	Key   string
	Value int64
}
type itemSliceASC []item

func (s itemSliceASC) Len() int           { return len(s) }
func (s itemSliceASC) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s itemSliceASC) Less(i, j int) bool { return s[i].Value < s[j].Value }

type itemSliceDESC []item

func (s itemSliceDESC) Len() int           { return len(s) }
func (s itemSliceDESC) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }
func (s itemSliceDESC) Less(i, j int) bool { return s[i].Value > s[j].Value }

//mapValueASC 正序
func mapValueASC(m map[string]int64) {
	is := make(itemSliceASC, len(m))
	for s, i := range m {
		is[i] = item{
			Key:   s,
			Value: m[s],
		}
	}
	sort.Stable(is)
}

//mapValueDESC 倒序
func mapValueDESC(m map[string]int64) {
	is := make(itemSliceDESC, len(m))
	for s, i := range m {
		is[i] = item{
			Key:   s,
			Value: m[s],
		}
	}
	sort.Stable(is)
}
