package main

import "sort"

type (
	StatSlice []StatRaw
)

func (s StatSlice) Len() int {
	return len(s)
}

func (s StatSlice) Less(i, j int) bool {
	return s[i].Time.Before(s[j].Time)
}

func (s StatSlice) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s StatSlice) Sort() {
	sort.Sort(s)
}
