package main

import "github.com/yantao1995/ds/skiplist"

func main() {
	//var cmp *skiplist.CmpInstanceInt
	skiplist.New(func(a, b interface{}) int {
		return -1
	})
}
